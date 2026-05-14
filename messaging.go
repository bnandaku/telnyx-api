package telnyx

import (
	"context"
	"fmt"
)

// MessageAddress represents a phone number endpoint on a message.
type MessageAddress struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name,omitempty"`
	Carrier     string `json:"carrier,omitempty"`
	LineType    string `json:"line_type,omitempty"`
}

// MessageMedia is a single media attachment on an MMS message.
type MessageMedia struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	SHA256      string `json:"sha256"`
	Size        int    `json:"size"`
}

// MessageCost holds billing info for a sent message.
type MessageCost struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// MessageRecipient holds per-recipient delivery info on an outbound message.
type MessageRecipient struct {
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	Carrier     string `json:"carrier,omitempty"`
	LineType    string `json:"line_type,omitempty"`
}

// Message is the full message object returned by the API.
type Message struct {
	RecordType         string             `json:"record_type"`
	ID                 string             `json:"id"`
	Direction          string             `json:"direction"`
	MessagingProfileID string             `json:"messaging_profile_id"`
	From               MessageAddress     `json:"from"`
	To                 []MessageRecipient `json:"to"`
	Text               string             `json:"text"`
	Media              []MessageMedia     `json:"media"`
	WebhookURL         string             `json:"webhook_url"`
	WebhookFailoverURL string             `json:"webhook_failover_url"`
	Encoding           string             `json:"encoding"`
	Parts              int                `json:"parts"`
	Tags               []string           `json:"tags"`
	Type               string             `json:"type"` // SMS | MMS | RCS
	Cost               *MessageCost       `json:"cost"`
	ValidUntil         string             `json:"valid_until"`
	ReceivedAt         string             `json:"received_at"`
	SentAt             string             `json:"sent_at"`
	CompletedAt        string             `json:"completed_at"`
	CreatedAt          string             `json:"created_at"`
	UpdatedAt          string             `json:"updated_at"`
	Errors             []ErrorDetail      `json:"errors"`
}

type messageResponse struct {
	Data Message `json:"data"`
}

// SendMessageRequest is the payload for sending an SMS or MMS.
type SendMessageRequest struct {
	// From is the sending number (E.164) or messaging profile ID.
	From string `json:"from"`
	// To is one or more recipient numbers (E.164). Use a string for SMS,
	// a []string for group MMS.
	To any `json:"to"`

	Text      string   `json:"text,omitempty"`
	MediaURLs []string `json:"media_urls,omitempty"`
	Subject   string   `json:"subject,omitempty"` // MMS subject

	MessagingProfileID  string   `json:"messaging_profile_id,omitempty"`
	WebhookURL          string   `json:"webhook_url,omitempty"`
	WebhookFailoverURL  string   `json:"webhook_failover_url,omitempty"`
	UseProfileWebhooks  bool     `json:"use_profile_webhooks,omitempty"`
	Type                string   `json:"type,omitempty"` // SMS | MMS | RCS
	AutoDetectType      bool     `json:"auto_detect,omitempty"`
	Tags                []string `json:"tags,omitempty"`
	ValidUntil          string   `json:"valid_until,omitempty"`
}

// SendMessage sends an SMS, MMS, or RCS message.
func (c *Client) SendMessage(ctx context.Context, req SendMessageRequest) (*Message, error) {
	var resp messageResponse
	if err := c.post(ctx, "/messages", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetMessage retrieves a message by ID.
func (c *Client) GetMessage(ctx context.Context, id string) (*Message, error) {
	var resp messageResponse
	if err := c.get(ctx, fmt.Sprintf("/messages/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListMessagesParams filters the message list.
type ListMessagesParams struct {
	PageNumber  int
	PageSize    int
	DateRangeStart string
	DateRangeEnd   string
	Direction   string // inbound | outbound
	Status      string
}

type messageListResponse struct {
	Data []Message `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta holds standard pagination fields.
type PaginationMeta struct {
	TotalResults int `json:"total_results"`
	TotalPages   int `json:"total_pages"`
	PageSize     int `json:"page_size"`
	PageNumber   int `json:"page_number"`
}

// ListMessages returns messages for the account with optional filters.
func (c *Client) ListMessages(ctx context.Context, params ListMessagesParams) ([]Message, *PaginationMeta, error) {
	path := "/messages"
	sep := "?"
	if params.PageNumber > 0 {
		path += fmt.Sprintf("%spage[number]=%d", sep, params.PageNumber)
		sep = "&"
	}
	if params.PageSize > 0 {
		path += fmt.Sprintf("%spage[size]=%d", sep, params.PageSize)
		sep = "&"
	}
	if params.Direction != "" {
		path += fmt.Sprintf("%sfilter[direction]=%s", sep, params.Direction)
		sep = "&"
	}
	if params.DateRangeStart != "" {
		path += fmt.Sprintf("%sfilter[date_range_start]=%s", sep, params.DateRangeStart)
		sep = "&"
	}
	if params.DateRangeEnd != "" {
		path += fmt.Sprintf("%sfilter[date_range_end]=%s", sep, params.DateRangeEnd)
	}

	var resp messageListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// ScheduleMessageRequest schedules a message for future delivery.
type ScheduleMessageRequest struct {
	SendMessageRequest
	SendAt string `json:"send_at"` // ISO 8601
}

// ScheduleMessage schedules an SMS/MMS for delivery at a future time.
func (c *Client) ScheduleMessage(ctx context.Context, req ScheduleMessageRequest) (*Message, error) {
	var resp messageResponse
	if err := c.post(ctx, "/messages", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// CancelScheduledMessage cancels a message that hasn't been sent yet.
func (c *Client) CancelScheduledMessage(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/messages/%s", id))
}

// --- Messaging Profiles ---

// MessagingProfile configures routing, webhooks, and features for messaging.
type MessagingProfile struct {
	RecordType             string `json:"record_type"`
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Enabled                bool   `json:"enabled"`
	WebhookURL             string `json:"webhook_url"`
	WebhookFailoverURL     string `json:"webhook_failover_url"`
	WebhookAPIVersion      string `json:"webhook_api_version"`
	WhitelistedDestinations []string `json:"whitelisted_destinations"`
	MMSFallBackToSMS       bool   `json:"mms_fall_back_to_sms"`
	MMSTranscoding         bool   `json:"mms_transcoding"`
	SmartEncoding          bool   `json:"smart_encoding"`
	MobileOnly             bool   `json:"mobile_only"`
	DailySpendLimit        string `json:"daily_spend_limit"`
	DailySpendLimitEnabled bool   `json:"daily_spend_limit_enabled"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
}

type messagingProfileResponse struct {
	Data MessagingProfile `json:"data"`
}

type messagingProfileListResponse struct {
	Data []MessagingProfile `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

// CreateMessagingProfileRequest is the payload for creating a messaging profile.
type CreateMessagingProfileRequest struct {
	Name                    string   `json:"name"`
	Enabled                 bool     `json:"enabled,omitempty"`
	WebhookURL              string   `json:"webhook_url,omitempty"`
	WebhookFailoverURL      string   `json:"webhook_failover_url,omitempty"`
	WebhookAPIVersion       string   `json:"webhook_api_version,omitempty"` // 1 | 2 | 2010-04-01
	WhitelistedDestinations []string `json:"whitelisted_destinations,omitempty"`
	MMSFallBackToSMS        bool     `json:"mms_fall_back_to_sms,omitempty"`
	MMSTranscoding          bool     `json:"mms_transcoding,omitempty"`
	SmartEncoding           bool     `json:"smart_encoding,omitempty"`
	MobileOnly              bool     `json:"mobile_only,omitempty"`
	DailySpendLimit         string   `json:"daily_spend_limit,omitempty"`
	DailySpendLimitEnabled  bool     `json:"daily_spend_limit_enabled,omitempty"`
}

// CreateMessagingProfile creates a new messaging profile.
func (c *Client) CreateMessagingProfile(ctx context.Context, req CreateMessagingProfileRequest) (*MessagingProfile, error) {
	var resp messagingProfileResponse
	if err := c.post(ctx, "/messaging_profiles", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetMessagingProfile retrieves a messaging profile by ID.
func (c *Client) GetMessagingProfile(ctx context.Context, id string) (*MessagingProfile, error) {
	var resp messagingProfileResponse
	if err := c.get(ctx, fmt.Sprintf("/messaging_profiles/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UpdateMessagingProfileRequest carries partial updates for a messaging profile.
type UpdateMessagingProfileRequest struct {
	Name                    string   `json:"name,omitempty"`
	Enabled                 *bool    `json:"enabled,omitempty"`
	WebhookURL              string   `json:"webhook_url,omitempty"`
	WebhookFailoverURL      string   `json:"webhook_failover_url,omitempty"`
	WebhookAPIVersion       string   `json:"webhook_api_version,omitempty"`
	WhitelistedDestinations []string `json:"whitelisted_destinations,omitempty"`
	MMSFallBackToSMS        *bool    `json:"mms_fall_back_to_sms,omitempty"`
	MMSTranscoding          *bool    `json:"mms_transcoding,omitempty"`
	SmartEncoding           *bool    `json:"smart_encoding,omitempty"`
	DailySpendLimit         string   `json:"daily_spend_limit,omitempty"`
	DailySpendLimitEnabled  *bool    `json:"daily_spend_limit_enabled,omitempty"`
}

// UpdateMessagingProfile applies partial updates to a messaging profile.
func (c *Client) UpdateMessagingProfile(ctx context.Context, id string, req UpdateMessagingProfileRequest) (*MessagingProfile, error) {
	var resp messagingProfileResponse
	if err := c.patch(ctx, fmt.Sprintf("/messaging_profiles/%s", id), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeleteMessagingProfile permanently removes a messaging profile.
func (c *Client) DeleteMessagingProfile(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/messaging_profiles/%s", id))
}

// ListMessagingProfiles returns all messaging profiles for the account.
func (c *Client) ListMessagingProfiles(ctx context.Context, pageNumber, pageSize int) ([]MessagingProfile, *PaginationMeta, error) {
	path := fmt.Sprintf("/messaging_profiles?page[number]=%d&page[size]=%d", pageNumber, pageSize)
	var resp messagingProfileListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}
