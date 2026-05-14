package telnyx

import (
	"context"
	"fmt"
)

// Recording is a call recording stored by Telnyx.
type Recording struct {
	RecordType     string       `json:"record_type"`
	ID             string       `json:"id"`
	CallLegID      string       `json:"call_leg_id"`
	CallSessionID  string       `json:"call_session_id"`
	ConnectionID   string       `json:"connection_id"`
	Channels       string       `json:"channels"`   // single | dual
	Format         string       `json:"format"`     // wav | mp3
	Status         string       `json:"status"`     // processing | completed | failed
	Source         string       `json:"source"`     // call_leg | conference
	DurationMillis int          `json:"duration_millis"`
	Urls           RecordingURLs `json:"download_urls"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
}

type recordingResponse struct {
	Data Recording `json:"data"`
}

type recordingListResponse struct {
	Data []Recording    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// ListRecordingsParams filters the recordings list.
type ListRecordingsParams struct {
	PageNumber    int
	PageSize      int
	FilterCallLegID    string
	FilterCallSessionID string
	FilterStatus   string
	FilterSource   string
}

// ListRecordings returns all recordings for the account.
func (c *Client) ListRecordings(ctx context.Context, params ListRecordingsParams) ([]Recording, *PaginationMeta, error) {
	path := "/recordings"
	sep := "?"
	add := func(k, v string) {
		if v != "" {
			path += fmt.Sprintf("%s%s=%s", sep, k, v)
			sep = "&"
		}
	}
	if params.PageNumber > 0 {
		path += fmt.Sprintf("%spage[number]=%d", sep, params.PageNumber)
		sep = "&"
	}
	if params.PageSize > 0 {
		path += fmt.Sprintf("%spage[size]=%d", sep, params.PageSize)
		sep = "&"
	}
	add("filter[call_leg_id]", params.FilterCallLegID)
	add("filter[call_session_id]", params.FilterCallSessionID)
	add("filter[status]", params.FilterStatus)
	add("filter[source]", params.FilterSource)

	var resp recordingListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetRecording retrieves a recording by ID.
func (c *Client) GetRecording(ctx context.Context, id string) (*Recording, error) {
	var resp recordingResponse
	if err := c.get(ctx, fmt.Sprintf("/recordings/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeleteRecording permanently removes a recording.
func (c *Client) DeleteRecording(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/recordings/%s", id))
}

// Transcription holds a recording transcription result.
type Transcription struct {
	RecordType        string              `json:"record_type"`
	ID                string              `json:"id"`
	RecordingID       string              `json:"recording_id"`
	Status            string              `json:"status"` // completed | failed | processing
	TranscriptionText string              `json:"transcription_text"`
	Speakers          []TranscriptSpeaker `json:"speakers"`
	CreatedAt         string              `json:"created_at"`
	UpdatedAt         string              `json:"updated_at"`
}

// TranscriptSpeaker holds per-speaker segments in a transcription.
type TranscriptSpeaker struct {
	SpeakerID string              `json:"speaker_id"`
	Segments  []TranscriptSegment `json:"segments"`
}

// TranscriptSegment is a timed speech segment within a transcription.
type TranscriptSegment struct {
	Text       string  `json:"text"`
	StartTime  float64 `json:"start_time"`
	EndTime    float64 `json:"end_time"`
	Confidence float64 `json:"confidence"`
}

type transcriptionResponse struct {
	Data Transcription `json:"data"`
}

type transcriptionListResponse struct {
	Data []Transcription `json:"data"`
	Meta PaginationMeta  `json:"meta"`
}

// GetTranscription retrieves a transcription by ID.
func (c *Client) GetTranscription(ctx context.Context, id string) (*Transcription, error) {
	var resp transcriptionResponse
	if err := c.get(ctx, fmt.Sprintf("/transcriptions/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListTranscriptions returns transcriptions, optionally filtered by recording.
func (c *Client) ListTranscriptions(ctx context.Context, recordingID string, pageNumber, pageSize int) ([]Transcription, *PaginationMeta, error) {
	path := fmt.Sprintf("/transcriptions?page[number]=%d&page[size]=%d", pageNumber, pageSize)
	if recordingID != "" {
		path += fmt.Sprintf("&filter[recording_id]=%s", recordingID)
	}
	var resp transcriptionListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}
