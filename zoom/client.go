package zoom

import (
	"net/http"

	"github.com/rs/zerolog"
)

const defaultBaseURL = "https://api.zoom.us/v2"

// Client is a Zoom API client.
type Client struct {
	HTTP    *http.Client
	BaseURL string
	Logger  zerolog.Logger
}

// NewClient creates a new Zoom API client.
func NewClient(httpClient *http.Client, logger zerolog.Logger) *Client {
	return &Client{
		HTTP:    httpClient,
		BaseURL: defaultBaseURL,
		Logger:  logger,
	}
}
