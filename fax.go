package telnyx

import (
	"context"
	"fmt"
)

type FaxApplication struct {
	ID                      string   `json:"id"`
	RecordType              string   `json:"record_type"`
	ApplicationName         string   `json:"application_name"`
	Active                  bool     `json:"active"`
	AnchorSiteOverride      string   `json:"anchorsite_override,omitempty"`
	WebhookEventURL         string   `json:"webhook_event_url"`
	WebhookEventFailoverURL string   `json:"webhook_event_failover_url,omitempty"`
	WebhookTimeoutSecs      int      `json:"webhook_timeout_secs,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at"`
}

type CreateFaxApplicationRequest struct {
	ApplicationName         string   `json:"application_name"`
	WebhookEventURL         string   `json:"webhook_event_url"`
	Active                  bool     `json:"active,omitempty"`
	AnchorSiteOverride      string   `json:"anchorsite_override,omitempty"`
	WebhookEventFailoverURL string   `json:"webhook_event_failover_url,omitempty"`
	WebhookTimeoutSecs      int      `json:"webhook_timeout_secs,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
}

type UpdateFaxApplicationRequest struct {
	ApplicationName         string   `json:"application_name,omitempty"`
	WebhookEventURL         string   `json:"webhook_event_url,omitempty"`
	Active                  *bool    `json:"active,omitempty"`
	AnchorSiteOverride      string   `json:"anchorsite_override,omitempty"`
	WebhookEventFailoverURL string   `json:"webhook_event_failover_url,omitempty"`
	WebhookTimeoutSecs      int      `json:"webhook_timeout_secs,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
	FaxEmailRecipient       string   `json:"fax_email_recipient,omitempty"`
}

type Fax struct {
	ID              string `json:"id"`
	RecordType      string `json:"record_type"`
	ConnectionID    string `json:"connection_id"`
	Direction       string `json:"direction"`
	MediaURL        string `json:"media_url,omitempty"`
	MediaName       string `json:"media_name,omitempty"`
	To              string `json:"to"`
	From            string `json:"from"`
	FromDisplayName string `json:"from_display_name,omitempty"`
	Quality         string `json:"quality,omitempty"`
	Status          string `json:"status"`
	WebhookURL      string `json:"webhook_url,omitempty"`
	StoreMedia      bool   `json:"store_media,omitempty"`
	StoredMediaURL  string `json:"stored_media_url,omitempty"`
	PreviewURL      string `json:"preview_url,omitempty"`
	ClientState     string `json:"client_state,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type SendFaxRequest struct {
	ConnectionID    string `json:"connection_id"`
	To              string `json:"to"`
	From            string `json:"from"`
	MediaURL        string `json:"media_url,omitempty"`
	MediaName       string `json:"media_name,omitempty"`
	FromDisplayName string `json:"from_display_name,omitempty"`
	Quality         string `json:"quality,omitempty"`
	T38Enabled      *bool  `json:"t38_enabled,omitempty"`
	Monochrome      bool   `json:"monochrome,omitempty"`
	BlackThreshold  int    `json:"black_threshold,omitempty"`
	StoreMedia      bool   `json:"store_media,omitempty"`
	StorePreview    bool   `json:"store_preview,omitempty"`
	PreviewFormat   string `json:"preview_format,omitempty"`
	WebhookURL      string `json:"webhook_url,omitempty"`
	ClientState     string `json:"client_state,omitempty"`
}

type ListFaxesParams struct {
	PageNumber int
	PageSize   int
	FilterTo   string
	FilterFrom string
}

func (c *Client) CreateFaxApplication(ctx context.Context, req CreateFaxApplicationRequest) (*FaxApplication, error) {
	var out struct {
		Data FaxApplication `json:"data"`
	}
	if err := c.post(ctx, "/fax_applications", req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetFaxApplication(ctx context.Context, id string) (*FaxApplication, error) {
	var out struct {
		Data FaxApplication `json:"data"`
	}
	if err := c.get(ctx, fmt.Sprintf("/fax_applications/%s", id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateFaxApplication(ctx context.Context, id string, req UpdateFaxApplicationRequest) (*FaxApplication, error) {
	var out struct {
		Data FaxApplication `json:"data"`
	}
	if err := c.patch(ctx, fmt.Sprintf("/fax_applications/%s", id), req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) DeleteFaxApplication(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/fax_applications/%s", id))
}

func (c *Client) ListFaxApplications(ctx context.Context, pageNumber, pageSize int) ([]FaxApplication, *PaginationMeta, error) {
	path := "/fax_applications"
	sep := "?"
	add := func(k string, v int) {
		if v > 0 {
			path += fmt.Sprintf("%s%s=%d", sep, k, v)
			sep = "&"
		}
	}
	add("page[number]", pageNumber)
	add("page[size]", pageSize)

	var out struct {
		Data []FaxApplication `json:"data"`
		Meta PaginationMeta   `json:"meta"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Data, &out.Meta, nil
}

func (c *Client) SendFax(ctx context.Context, req SendFaxRequest) (*Fax, error) {
	var out struct {
		Data Fax `json:"data"`
	}
	if err := c.post(ctx, "/faxes", req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetFax(ctx context.Context, id string) (*Fax, error) {
	var out struct {
		Data Fax `json:"data"`
	}
	if err := c.get(ctx, fmt.Sprintf("/faxes/%s", id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) DeleteFax(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/faxes/%s", id))
}

func (c *Client) CancelFax(ctx context.Context, id string) (*Fax, error) {
	var out struct {
		Data Fax `json:"data"`
	}
	if err := c.post(ctx, fmt.Sprintf("/faxes/%s/actions/cancel", id), nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) RefreshFax(ctx context.Context, id string) (*Fax, error) {
	var out struct {
		Data Fax `json:"data"`
	}
	if err := c.post(ctx, fmt.Sprintf("/faxes/%s/actions/refresh", id), nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) ListFaxes(ctx context.Context, params ListFaxesParams) ([]Fax, *PaginationMeta, error) {
	path := "/faxes"
	sep := "?"
	addStr := func(k, v string) {
		if v != "" {
			path += fmt.Sprintf("%s%s=%s", sep, k, v)
			sep = "&"
		}
	}
	addInt := func(k string, v int) {
		if v > 0 {
			path += fmt.Sprintf("%s%s=%d", sep, k, v)
			sep = "&"
		}
	}
	addInt("page[number]", params.PageNumber)
	addInt("page[size]", params.PageSize)
	addStr("filter[to]", params.FilterTo)
	addStr("filter[from]", params.FilterFrom)

	var out struct {
		Data []Fax          `json:"data"`
		Meta PaginationMeta `json:"meta"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Data, &out.Meta, nil
}
