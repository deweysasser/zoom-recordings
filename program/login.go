package program

import (
	"context"

	"github.com/deweysasser/zoom-recordings/auth"
	"github.com/rs/zerolog"
)

// LoginCmd authenticates with Zoom via OAuth.
type LoginCmd struct{}

// Run executes the login command.
func (cmd *LoginCmd) Run(ctx context.Context, logger zerolog.Logger, opts *Options) error {
	cfg := &auth.OAuthConfig{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		CallbackPort: opts.CallbackPort,
	}

	logger.Info().Msg("Starting Zoom OAuth authentication...")

	token, err := auth.Authenticate(ctx, cfg)
	if err != nil {
		return err
	}

	if err := auth.SaveToken(token); err != nil {
		return err
	}

	logger.Info().Msg("Authentication successful, token saved")
	return nil
}
