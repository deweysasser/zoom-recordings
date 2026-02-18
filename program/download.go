package program

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/deweysasser/zoom-recordings/auth"
	"github.com/deweysasser/zoom-recordings/zoom"
	"github.com/rs/zerolog"
)

// DownloadCmd downloads Zoom recordings.
type DownloadCmd struct {
	From      string `help:"Start date (YYYY-MM-DD). Defaults to 24 hours ago." short:"f"`
	To        string `help:"End date (YYYY-MM-DD). Defaults to today." short:"t"`
	OutputDir string `help:"Directory to save recordings to." default:"." short:"o"`
}

// Run executes the download command.
func (cmd *DownloadCmd) Run(ctx context.Context, logger zerolog.Logger, opts *Options) error {
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

	httpClient := auth.NewHTTPClient(ctx, cfg, token)
	client := zoom.NewClient(httpClient, logger)

	meetings, err := client.ListRecordings(ctx, from, to)
	if err != nil {
		return err
	}

	if len(meetings) == 0 {
		logger.Info().Msg("No recordings found")
		return nil
	}

	if err := os.MkdirAll(cmd.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	var downloaded, skipped int

	for _, m := range meetings {
		date := m.StartTime.Format("2006-01-02")
		topic := sanitizeFilename(m.Topic)

		for _, f := range m.Files {
			if f.Status != "completed" {
				logger.Debug().Str("status", f.Status).Str("file", f.ID).Msg("Skipping incomplete recording")
				continue
			}

			ext := strings.ToLower(f.FileExtension)
			if ext == "" {
				ext = strings.ToLower(f.FileType)
			}

			filename := fmt.Sprintf("%s_%s_%s.%s", date, topic, strings.ToLower(f.RecordingType), ext)
			destPath := filepath.Join(cmd.OutputDir, filename)

			// Skip if file exists with matching size
			if info, err := os.Stat(destPath); err == nil && info.Size() == f.FileSize {
				logger.Debug().Str("file", filename).Msg("Skipping existing file")
				skipped++
				continue
			}

			logger.Info().Str("file", filename).Int64("size_mb", f.FileSize/1024/1024).Msg("Downloading")

			if err := client.DownloadFile(ctx, f.DownloadURL, token.AccessToken, destPath); err != nil {
				logger.Error().Err(err).Str("file", filename).Msg("Failed to download")
				continue
			}

			downloaded++
		}
	}

	// Re-save token in case it was refreshed
	if err := auth.SaveToken(token); err != nil {
		logger.Warn().Err(err).Msg("Failed to save refreshed token")
	}

	logger.Info().
		Int("downloaded", downloaded).
		Int("skipped", skipped).
		Str("output_dir", cmd.OutputDir).
		Msg("Download complete")

	return nil
}

var nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9-]+`)

func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, " ", "-")
	name = nonAlphanumeric.ReplaceAllString(name, "")
	name = strings.ToLower(name)
	if name == "" {
		name = "untitled"
	}
	return name
}
