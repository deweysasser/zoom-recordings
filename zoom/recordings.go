package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ListRecordings fetches all recordings for the authenticated user within the given date range.
// It handles pagination automatically via next_page_token.
func (c *Client) ListRecordings(ctx context.Context, from, to string) ([]Meeting, error) {
	var allMeetings []Meeting
	pageToken := ""

	for {
		params := url.Values{}
		params.Set("from", from)
		params.Set("to", to)
		params.Set("page_size", "300")
		if pageToken != "" {
			params.Set("next_page_token", pageToken)
		}

		reqURL := fmt.Sprintf("%s/users/me/recordings?%s", c.BaseURL, params.Encode())

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		resp, err := c.HTTP.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetching recordings: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("reading response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		}

		var result RecordingsResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("decoding response: %w", err)
		}

		c.Logger.Debug().
			Int("meetings", len(result.Meetings)).
			Int("total_records", result.TotalRecords).
			Msg("Fetched recordings page")

		allMeetings = append(allMeetings, result.Meetings...)

		if result.NextPageToken == "" {
			break
		}
		pageToken = result.NextPageToken
	}

	return allMeetings, nil
}
