package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GitHub is a minimal OAuth 2.0 web-flow client for "Sign in with GitHub".
type GitHub struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string // APP_URL + /api/v1/auth/github/callback
	// Endpoints default to github.com; tests point them at a fake.
	AuthorizeURL string
	TokenURL     string
	APIURL       string
	HTTP         *http.Client
}

// Profile is what sign-in needs from GitHub.
type Profile struct {
	ID     string
	Login  string
	Name   string
	Avatar string
	Email  string // primary, verified
}

// Name is "github".
func (g *GitHub) Name() string { return "github" }

func (g *GitHub) client() *http.Client {
	if g.HTTP != nil {
		return g.HTTP
	}
	return &http.Client{Timeout: 10 * time.Second}
}

func or(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// AuthCodeURL is where to send the browser; state must be checked on return.
func (g *GitHub) AuthCodeURL(state string) string {
	v := url.Values{
		"client_id":    {g.ClientID},
		"redirect_uri": {g.RedirectURL},
		"scope":        {"read:user user:email"},
		"state":        {state},
		"allow_signup": {"true"},
	}
	return or(g.AuthorizeURL, "https://github.com/login/oauth/authorize") + "?" + v.Encode()
}

// Exchange trades the callback code for an access token.
func (g *GitHub) Exchange(ctx context.Context, code string) (string, error) {
	form := url.Values{"client_id": {g.ClientID}, "client_secret": {g.ClientSecret}, "code": {code}, "redirect_uri": {g.RedirectURL}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, or(g.TokenURL, "https://github.com/login/oauth/access_token"), strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	var out struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := g.do(req, &out); err != nil {
		return "", err
	}
	if out.AccessToken == "" {
		return "", fmt.Errorf("github token: %s", or(out.Error, "no access token"))
	}
	return out.AccessToken, nil
}

// Profile fetches the user and their primary verified email.
func (g *GitHub) Profile(ctx context.Context, token string) (Profile, error) {
	api := strings.TrimRight(or(g.APIURL, "https://api.github.com"), "/")
	var u struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := g.get(ctx, api+"/user", token, &u); err != nil {
		return Profile{}, err
	}
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := g.get(ctx, api+"/user/emails", token, &emails); err != nil {
		return Profile{}, err
	}
	p := Profile{ID: strconv.FormatInt(u.ID, 10), Login: u.Login, Name: or(u.Name, u.Login), Avatar: u.AvatarURL}
	for _, e := range emails {
		if e.Primary && e.Verified {
			p.Email = e.Email
		}
	}
	if u.ID == 0 {
		return Profile{}, errors.New("github: no user id")
	}
	if p.Email == "" {
		return Profile{}, errors.New("github: no verified primary email")
	}
	return p, nil
}

func (g *GitHub) get(ctx context.Context, u, token string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("github request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	return g.do(req, v)
}

func (g *GitHub) do(req *http.Request, v any) error {
	resp, err := g.client().Do(req)
	if err != nil {
		return fmt.Errorf("github: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("github: read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github: %s %s: status %d", req.Method, req.URL.Path, resp.StatusCode)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("github: decode %s: %w", req.URL.Path, err)
	}
	return nil
}
