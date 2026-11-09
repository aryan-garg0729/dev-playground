package auth

import "context"

// Provider is an OAuth 2.0 web-flow sign-in (GitHub, Google). Every provider
// must return a verified email: accounts are matched and linked by it.
type Provider interface {
	// Name is the route and identity key: "github", "google".
	Name() string
	// AuthCodeURL is where to send the browser; state is checked on return.
	AuthCodeURL(state string) string
	// Exchange trades the callback code for an access token.
	Exchange(ctx context.Context, code string) (string, error)
	// Profile fetches the user, with a verified email.
	Profile(ctx context.Context, token string) (Profile, error)
}

var (
	_ Provider = (*GitHub)(nil)
	_ Provider = (*Google)(nil)
)
