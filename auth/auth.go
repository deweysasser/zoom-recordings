package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
)

var zoomEndpoint = oauth2.Endpoint{
	AuthURL:  "https://zoom.us/oauth/authorize",
	TokenURL: "https://zoom.us/oauth/token",
}

// OAuthConfig holds the configuration for the OAuth flow.
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	CallbackPort int
}

func (c *OAuthConfig) oauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint:     zoomEndpoint,
		RedirectURL:  fmt.Sprintf("http://localhost:%d/oauth/callback", c.CallbackPort),
		Scopes: []string{
			"cloud_recording:read:list_user_recordings",
			"cloud_recording:read:list_recording_files",
		},
	}
}

// Authenticate performs the full OAuth2 authorization code flow.
// It starts a local HTTP server, opens the browser for user consent,
// and returns the resulting token.
func Authenticate(ctx context.Context, cfg *OAuthConfig) (*oauth2.Token, error) {
	conf := cfg.oauth2Config()

	state, err := randomState()
	if err != nil {
		return nil, fmt.Errorf("generating state: %w", err)
	}

	tokenCh := make(chan *oauth2.Token, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "invalid state", http.StatusBadRequest)
			errCh <- fmt.Errorf("state mismatch in OAuth callback")
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			errCh <- fmt.Errorf("missing code in OAuth callback")
			return
		}

		token, err := conf.Exchange(ctx, code)
		if err != nil {
			http.Error(w, "token exchange failed", http.StatusInternalServerError)
			errCh <- fmt.Errorf("exchanging code for token: %w", err)
			return
		}

		fmt.Fprintln(w, "Authentication successful! You can close this window.")
		tokenCh <- token
	})

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.CallbackPort))
	if err != nil {
		return nil, fmt.Errorf("starting callback server: %w", err)
	}

	server := &http.Server{Handler: mux}
	go func() { _ = server.Serve(listener) }()
	defer server.Shutdown(ctx)

	authURL := conf.AuthCodeURL(state)
	fmt.Printf("Opening browser for Zoom authentication...\n")
	fmt.Printf("If the browser doesn't open, visit:\n%s\n\n", authURL)

	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Could not open browser automatically: %v\n", err)
	}

	select {
	case token := <-tokenCh:
		return token, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// NewHTTPClient creates an HTTP client that automatically handles OAuth2 tokens.
func NewHTTPClient(ctx context.Context, cfg *OAuthConfig, token *oauth2.Token) *http.Client {
	conf := cfg.oauth2Config()
	return conf.Client(ctx, token)
}

func randomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	default:
		return fmt.Errorf("unsupported platform %s", runtime.GOOS)
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
