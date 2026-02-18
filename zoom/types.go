package zoom

import "time"

// RecordingsResponse is the top-level response from the Zoom list recordings API.
type RecordingsResponse struct {
	From          string    `json:"from"`
	To            string    `json:"to"`
	TotalRecords  int       `json:"total_records"`
	NextPageToken string    `json:"next_page_token"`
	Meetings      []Meeting `json:"meetings"`
}

// Meeting represents a single Zoom meeting with its recordings.
type Meeting struct {
	UUID      string          `json:"uuid"`
	Topic     string          `json:"topic"`
	StartTime time.Time       `json:"start_time"`
	Duration  int             `json:"duration"`
	Files     []RecordingFile `json:"recording_files"`
}

// RecordingFile represents a single recording file within a meeting.
type RecordingFile struct {
	ID            string `json:"id"`
	FileType      string `json:"file_type"`
	FileExtension string `json:"file_extension"`
	DownloadURL   string `json:"download_url"`
	Status        string `json:"status"`
	RecordingType string `json:"recording_type"`
	FileSize      int64  `json:"file_size"`
}
