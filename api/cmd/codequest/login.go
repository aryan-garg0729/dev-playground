package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func (a *app) login(ctx context.Context) error {
	c := a.cfg
	if err := checkServiceURL(c.apiURL()); err != nil {
		return err
	}
	loginURL := c.appURL() + "/cli/login"
	_, _ = fmt.Fprintf(a.out, "Welcome to the CodeQuest CLI!\n\nOpen this page, signed in, to get a login code:\n  %s\n\n", loginURL)
	if a.open != nil {
		_ = a.open(loginURL) // printing the URL is enough when there is no browser
	}
	code, err := a.ask("Paste your login code: ")
	if err != nil {
		return fmt.Errorf("read code: %w", err)
	}
	if code == "" {
		return usageError{"no code entered"}
	}
	host, _ := os.Hostname()
	var tok tokenResponse
	err = newAPIClient(c.apiURL(), "").do(ctx, "POST", "/cli/token", map[string]string{"code": code, "label": host}, &tok)
	if isAPIError(err, "invalid_code") {
		return errors.New("that code is unknown, used or expired; refresh the page for a new one")
	}
	if err != nil {
		return err
	}
	if tok.Token == "" {
		return errors.New("the API returned no token")
	}
	c.Token, c.Name, c.Email = tok.Token, tok.Name, tok.Email
	if err := c.save(); err != nil {
		return err
	}
	who := tok.Email
	if tok.Name != "" {
		who = tok.Name + " <" + tok.Email + ">"
	}
	_, _ = fmt.Fprintf(a.out, "Signed in as %s.\n", who)
	return nil
}

// logout revokes the token on the server if it can, and always forgets it here.
func (a *app) logout(ctx context.Context) error {
	c := a.cfg
	if c.Token == "" {
		_, _ = fmt.Fprintln(a.out, "You are not signed in.")
		return nil
	}
	rctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	revokeErr := newAPIClient(c.apiURL(), c.Token).do(rctx, "POST", "/cli/logout", struct{}{}, nil)
	c.Token, c.Name, c.Email = "", "", ""
	if err := c.save(); err != nil {
		return err
	}
	if revokeErr != nil {
		var ae *apiError
		if !errors.As(revokeErr, &ae) || ae.Status != 401 {
			_, _ = fmt.Fprintf(a.err, "warning: could not revoke the token on the server (%v); it expires on its own\n", revokeErr)
		}
	}
	_, _ = fmt.Fprintln(a.out, "Signed out.")
	return nil
}

func (a *app) client() (*apiClient, error) {
	if a.cfg.Token == "" {
		return nil, errors.New("you are not signed in; run `codequest login`")
	}
	if err := checkServiceURL(a.cfg.apiURL()); err != nil {
		return nil, err
	}
	return newAPIClient(a.cfg.apiURL(), a.cfg.Token), nil
}

func (a *app) status(ctx context.Context) error {
	cl, err := a.client()
	if err != nil {
		return err
	}
	var me meResponse
	if err := cl.do(ctx, "GET", "/me", nil, &me); err != nil {
		return err
	}
	s := me.Stats
	name := me.User.Name
	if name == "" {
		name = me.User.Email
	}
	_, _ = fmt.Fprintf(a.out, "%s (%s)\nLevel %d · %d XP (next level at %d)\nStreak: %d days · today %d/%d lessons\nLeague: %s · %d Bytes\n",
		name, me.User.Email, s.Level, s.TotalXP, s.NextLevelXP, s.StreakDays, s.TodayLessons, s.DailyGoal, s.LeagueTier, s.Bytes)
	return nil
}
