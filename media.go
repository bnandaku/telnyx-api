package telnyx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// MediaItem represents a stored media file on the Telnyx network.
type MediaItem struct {
	MediaName   string `json:"media_name"`
	ExpiresAt   string `json:"expires_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ContentType string `json:"content_type"`
}

type mediaResponse struct {
	Data MediaItem `json:"data"`
}

type mediaListResponse struct {
	Data []MediaItem    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// UploadMediaFromURL stores a remote file (≤20 MB) on the Telnyx network.
// mediaName is the unique identifier used in call commands via media_name.
// ttlSecs sets the expiry (default 2 days; max ~20 years). Pass 0 for default.
func (c *Client) UploadMediaFromURL(ctx context.Context, mediaURL, mediaName string, ttlSecs int) (*MediaItem, error) {
	body := map[string]any{"media_url": mediaURL}
	if mediaName != "" {
		body["media_name"] = mediaName
	}
	if ttlSecs > 0 {
		body["ttl_secs"] = ttlSecs
	}
	var resp mediaResponse
	if err := c.post(ctx, "/media", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UploadMediaFile uploads a local file (≤20 MB) to the Telnyx network via multipart.
// mediaName is the unique identifier; filename is used for the form field.
// ttlSecs sets the expiry (default 2 days). Pass 0 for default.
func (c *Client) UploadMediaFile(ctx context.Context, mediaName string, ttlSecs int, filename string, content io.Reader) (*MediaItem, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("media", filepath.Base(filename))
	if err != nil {
		return nil, fmt.Errorf("telnyx: create form file: %w", err)
	}
	if _, err = io.Copy(fw, content); err != nil {
		return nil, fmt.Errorf("telnyx: write media: %w", err)
	}
	if mediaName != "" {
		_ = w.WriteField("media_name", mediaName)
	}
	if ttlSecs > 0 {
		_ = w.WriteField("ttl_secs", fmt.Sprintf("%d", ttlSecs))
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/media", &buf)
	if err != nil {
		return nil, fmt.Errorf("telnyx: build upload request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("telnyx: upload media: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp)
	}

	var result mediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("telnyx: decode response: %w", err)
	}
	return &result.Data, nil
}

// GetMedia retrieves metadata for a stored media file by media_name.
func (c *Client) GetMedia(ctx context.Context, mediaName string) (*MediaItem, error) {
	var resp mediaResponse
	if err := c.get(ctx, fmt.Sprintf("/media/%s", mediaName), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListMedia returns all stored media files, optionally filtered by content type.
// contentTypes examples: "audio/mpeg", "audio/wav". Pass nil for no filter.
func (c *Client) ListMedia(ctx context.Context, pageNumber, pageSize int, contentTypes []string) ([]MediaItem, *PaginationMeta, error) {
	path := fmt.Sprintf("/media?page[number]=%d&page[size]=%d", pageNumber, pageSize)
	for _, ct := range contentTypes {
		path += fmt.Sprintf("&filter[content_type][]=%s", ct)
	}
	var resp mediaListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// DeleteMedia removes a stored media file.
func (c *Client) DeleteMedia(ctx context.Context, mediaName string) error {
	return c.delete(ctx, fmt.Sprintf("/media/%s", mediaName))
}
