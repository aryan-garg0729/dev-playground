package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Google is a minimal OAuth 2.0 web-flow client for "Sign in with Google"
// (OpenID Connect scopes; the profile comes from the userinfo endpoint).
type Google struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string // APP_URL + /api/v1/auth/google/callback
	// Endpoints default to Google's; tests point them at a fake.
	AuthorizeURL string
	TokenURL     string
	UserinfoURL  string
	HTTP         *http.Client
}

// Name is "google".
func (g *Google) Name() string { return "google" }

func (g *Google) client() *http.Client {
	if g.HTTP != nil {
		return g.HTTP
	}
	return &http.Client{Timeout: 10 * time.Second}
}

// AuthCodeURL is where to send the browser; state must be checked on return.
func (g *Google) AuthCodeURL(state string) string {
	v := url.Values{
		"client_id":     {g.ClientID},
		"redirect_uri":  {g.RedirectURL},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
		"prompt":        {"select_account"},
	}
	return or(g.AuthorizeURL, "https://accounts.google.com/o/oauth2/v2/auth") + "?" + v.Encode()
}

// Exchange trades the callback code for an access token.
func (g *Google) Exchange(ctx context.Context, code string) (string, error) {
	form := url.Values{
		"client_id": {g.ClientID}, "client_secret": {g.ClientSecret}, "code": {code},
		"redirect_uri": {g.RedirectURL}, "grant_type": {"authorization_code"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, or(g.TokenURL, "https://oauth2.googleapis.com/token"), strings.NewReader(form.Encode()))
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
		return "", fmt.Errorf("google token: %s", or(out.Error, "no access token"))
	}
	return out.AccessToken, nil
}

// Profile fetches the user; Google must report the email as verified.
func (g *Google) Profile(ctx context.Context, token string) (Profile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, or(g.UserinfoURL, "https://openidconnect.googleapis.com/v1/userinfo"), nil)
	if err != nil {
		return Profile{}, fmt.Errorf("google request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var u struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := g.do(req, &u); err != nil {
		return Profile{}, err
	}
	if u.Sub == "" {
		return Profile{}, errors.New("google: no user id")
	}
	if u.Email == "" || !u.EmailVerified {
		return Profile{}, errors.New("google: no verified email")
	}
	name := u.Name
	if name == "" {
		name, _, _ = strings.Cut(u.Email, "@")
	}
	return Profile{ID: u.Sub, Login: u.Email, Name: name, Avatar: u.Picture, Email: u.Email}, nil
}

func (g *Google) do(req *http.Request, v any) error {
	resp, err := g.client().Do(req)
	if err != nil {
		return fmt.Errorf("google: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("google: read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("google: %s %s: status %d", req.Method, req.URL.Path, resp.StatusCode)
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("google: decode %s: %w", req.URL.Path, err)
	}
	return nil
}
