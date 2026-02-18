package program

import (
	"context"
	"fmt"
	"time"

	"github.com/deweysasser/zoom-recordings/auth"
	"github.com/deweysasser/zoom-recordings/zoom"
	"github.com/rs/zerolog"
)

// ListCmd lists available Zoom recordings.
type ListCmd struct {
	From string `help:"Start date (YYYY-MM-DD). Defaults to 24 hours ago." short:"f"`
	To   string `help:"End date (YYYY-MM-DD). Defaults to today." short:"t"`
}

// Run executes the list command.
func (cmd *ListCmd) Run(ctx context.Context, logger zerolog.Logger, opts *Options) error {
	token, err := auth.LoadToken()
	if err != nil {
		return err
	}
	if token == nil {
		return fmt.Errorf("not authenticated, run 'zoom-recordings login' first")
	}

	cfg := &auth.OAuthConfig{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		CallbackPort: opts.CallbackPort,
	}

	from, to := resolveDateRange(cmd.From, cmd.To)

	client := zoom.NewClient(auth.NewHTTPClient(ctx, cfg, token), logger)
	meetings, err := client.ListRecordings(ctx, from, to)
	if err != nil {
		return err
	}

	// Re-save token in case it was refreshed
	if newToken, err := auth.LoadToken(); err == nil && newToken != nil {
		_ = auth.SaveToken(newToken)
	}

	if len(meetings) == 0 {
		fmt.Println("No recordings found.")
		return nil
	}

	fmt.Printf("Recordings from %s to %s:\n\n", from, to)
	for _, m := range meetings {
		fmt.Printf("  %s  %s (%d min)\n", m.StartTime.Format("2006-01-02 15:04"), m.Topic, m.Duration)
		for _, f := range m.Files {
			fmt.Printf("    - %s (%s, %.1f MB)\n", f.RecordingType, f.FileExtension, float64(f.FileSize)/1024/1024)
		}
	}

	fmt.Printf("\nTotal: %d meetings\n", len(meetings))
	return nil
}

func resolveDateRange(from, to string) (string, string) {
	now := time.Now()
	if to == "" {
		to = now.Format("2006-01-02")
	}
	if from == "" {
		from = now.AddDate(0, 0, -1).Format("2006-01-02")
	}
	return from, to
}
