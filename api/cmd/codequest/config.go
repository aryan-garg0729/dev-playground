package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Defaults for a local `make up`; release builds override them with
// -ldflags "-X main.defaultAPIURL=... -X main.defaultAppURL=...".
var (
	defaultAPIURL = "http://localhost:8080"
	defaultAppURL = "http://localhost:5173"
)

// config is ~/.config/codequest/config.json (os.UserConfigDir), mode 0600.
type config struct {
	APIURL  string `json:"api_url,omitempty"`
	AppURL  string `json:"app_url,omitempty"`
	BaseURL string `json:"base_url,omitempty"` // overrides a lesson's ${baseURL}
	Token   string `json:"token,omitempty"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`

	path string
}

func configDir() (string, error) {
	if d := os.Getenv("CODEQUEST_CONFIG_DIR"); d != "" {
		return d, nil
	}
	d, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("find config directory: %w", err)
	}
	return filepath.Join(d, "codequest"), nil
}

func loadConfig() (*config, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	c := &config{path: filepath.Join(dir, "config.json")}
	b, err := os.ReadFile(c.path)
	if errors.Is(err, fs.ErrNotExist) {
		return c, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", c.path, err)
	}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, fmt.Errorf("read %s: %w (delete it and run `codequest login`)", c.path, err)
	}
	return c, nil
}

// save writes the config atomically with owner-only permissions.
func (c *config) save() error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create %s: %w", dir, err)
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".config-*.json")
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write config: %w", err)
	}
	if _, err := tmp.Write(append(b, '\n')); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write config: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	if err := os.Rename(tmp.Name(), c.path); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func pick(env, cfg, def string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	if cfg != "" {
		return cfg
	}
	return def
}

func (c *config) apiURL() string {
	return strings.TrimSuffix(pick("CODEQUEST_API_URL", c.APIURL, defaultAPIURL), "/")
}

func (c *config) appURL() string {
	return strings.TrimSuffix(pick("CODEQUEST_APP_URL", c.AppURL, defaultAppURL), "/")
}

// checkServiceURL refuses plain http except to this machine, so the token
// never crosses a network unencrypted.
func checkServiceURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return fmt.Errorf("%q is not a URL", raw)
	}
	switch u.Scheme {
	case "https":
		return nil
	case "http":
		h := u.Hostname()
		if h == "localhost" || h == "127.0.0.1" || h == "::1" {
			return nil
		}
		return fmt.Errorf("%s: use https (plain http is allowed only for localhost)", raw)
	default:
		return fmt.Errorf("%s: use https", raw)
	}
}
