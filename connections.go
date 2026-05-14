package telnyx

import (
	"context"
	"fmt"
)

// RTCPSettings configures RTCP behaviour on a connection.
type RTCPSettings struct {
	Port            string `json:"port,omitempty"`
	CaptureEnabled  bool   `json:"capture_enabled,omitempty"`
	ReportFrequency int    `json:"report_frequency_secs,omitempty"`
}

// JitterBuffer configures RTP jitter buffering on a connection.
type JitterBuffer struct {
	Enabled bool `json:"enabled,omitempty"`
	MinMS   int  `json:"min_ms,omitempty"`
	MaxMS   int  `json:"max_ms,omitempty"`
}

// ConnectionInbound holds inbound call settings for a connection.
type ConnectionInbound struct {
	ANINumberFormat          string   `json:"ani_number_format,omitempty"`
	DNISNumberFormat         string   `json:"dnis_number_format,omitempty"`
	Codecs                   []string `json:"codecs,omitempty"`
	ChannelLimit             int      `json:"channel_limit,omitempty"`
	GenerateRingbackTone     bool     `json:"generate_ringback_tone,omitempty"`
	ISUPHeadersEnabled       bool     `json:"isup_headers_enabled,omitempty"`
	PRACKEnabled             bool     `json:"prack_enabled,omitempty"`
	PrivacyZoneEnabled       bool     `json:"privacy_zone_enabled,omitempty"`
	SIPCompactHeadersEnabled bool     `json:"sip_compact_headers_enabled,omitempty"`
	Timeout1xxSecs           int      `json:"timeout_1xx_secs,omitempty"`
	Timeout2xxSecs           int      `json:"timeout_2xx_secs,omitempty"`
}

// ConnectionOutbound holds outbound call settings for a connection.
type ConnectionOutbound struct {
	ANIOverride     string `json:"ani_override,omitempty"`
	ANIOverrideType string `json:"ani_override_type,omitempty"` // always | normal | emergency
	CallParkingEnabled bool `json:"call_parking_enabled,omitempty"`
	ChannelLimit    int    `json:"channel_limit,omitempty"`
	InstantRingbackEnabled bool `json:"instant_ringback_enabled,omitempty"`
	LocalCustomSDAEnabled  bool `json:"localization_enabled,omitempty"`
	OutboundVoiceProfileID string `json:"outbound_voice_profile_id,omitempty"`
}

// CredentialConnection is a SIP credential-based connection.
type CredentialConnection struct {
	RecordType                    string             `json:"record_type"`
	ID                            string             `json:"id"`
	Active                        bool               `json:"active"`
	UserName                      string             `json:"user_name"`
	ConnectionName                string             `json:"connection_name"`
	SIPURICallingPreference       string             `json:"sip_uri_calling_preference"`
	AnchorSiteOverride            string             `json:"anchorsite_override"`
	DefaultOnHoldComfortNoiseEnabled bool            `json:"default_on_hold_comfort_noise_enabled"`
	DTMFType                      string             `json:"dtmf_type"`
	EncodeContactHeaderEnabled    bool               `json:"encode_contact_header_enabled"`
	EncryptedMedia                string             `json:"encrypted_media"`
	OnnetT38PassthroughEnabled    bool               `json:"onnet_t38_passthrough_enabled"`
	WebhookEventURL               string             `json:"webhook_event_url"`
	WebhookEventFailoverURL       string             `json:"webhook_event_failover_url"`
	WebhookAPIVersion             string             `json:"webhook_api_version"`
	WebhookTimeoutSecs            int                `json:"webhook_timeout_secs"`
	NoiseSuppression              string             `json:"noise_suppression"`
	RTCPSettings                  *RTCPSettings      `json:"rtcp_settings,omitempty"`
	Inbound                       *ConnectionInbound `json:"inbound,omitempty"`
	Outbound                      *ConnectionOutbound `json:"outbound,omitempty"`
	JitterBuffer                  *JitterBuffer      `json:"jitter_buffer,omitempty"`
	Tags                          []string           `json:"tags"`
	CreatedAt                     string             `json:"created_at"`
	UpdatedAt                     string             `json:"updated_at"`
}

type credentialConnectionResponse struct {
	Data CredentialConnection `json:"data"`
}

type credentialConnectionListResponse struct {
	Data []CredentialConnection `json:"data"`
	Meta PaginationMeta         `json:"meta"`
}

// CreateCredentialConnectionRequest creates a SIP credential connection.
type CreateCredentialConnectionRequest struct {
	UserName                      string             `json:"user_name"`
	Password                      string             `json:"password"`
	ConnectionName                string             `json:"connection_name,omitempty"`
	Active                        bool               `json:"active,omitempty"`
	AnchorSiteOverride            string             `json:"anchorsite_override,omitempty"`
	SIPURICallingPreference       string             `json:"sip_uri_calling_preference,omitempty"` // disabled | unrestricted | internal
	DefaultOnHoldComfortNoiseEnabled bool            `json:"default_on_hold_comfort_noise_enabled,omitempty"`
	DTMFType                      string             `json:"dtmf_type,omitempty"` // RFC 2833 | Inband | SIP INFO
	EncodeContactHeaderEnabled    bool               `json:"encode_contact_header_enabled,omitempty"`
	EncryptedMedia                string             `json:"encrypted_media,omitempty"` // SRTP or omit
	OnnetT38PassthroughEnabled    bool               `json:"onnet_t38_passthrough_enabled,omitempty"`
	IOSPushCredentialID           string             `json:"ios_push_credential_id,omitempty"`
	AndroidPushCredentialID       string             `json:"android_push_credential_id,omitempty"`
	WebhookEventURL               string             `json:"webhook_event_url,omitempty"`
	WebhookEventFailoverURL       string             `json:"webhook_event_failover_url,omitempty"`
	WebhookAPIVersion             string             `json:"webhook_api_version,omitempty"` // 1 | 2 | texml
	WebhookTimeoutSecs            int                `json:"webhook_timeout_secs,omitempty"`
	CallCostInWebhooks            bool               `json:"call_cost_in_webhooks,omitempty"`
	NoiseSuppression              string             `json:"noise_suppression,omitempty"` // inbound | outbound | both | disabled
	RTCPSettings                  *RTCPSettings      `json:"rtcp_settings,omitempty"`
	Inbound                       *ConnectionInbound `json:"inbound,omitempty"`
	Outbound                      *ConnectionOutbound `json:"outbound,omitempty"`
	JitterBuffer                  *JitterBuffer      `json:"jitter_buffer,omitempty"`
	Tags                          []string           `json:"tags,omitempty"`
}

// CreateCredentialConnection creates a new SIP credential connection.
func (c *Client) CreateCredentialConnection(ctx context.Context, req CreateCredentialConnectionRequest) (*CredentialConnection, error) {
	var resp credentialConnectionResponse
	if err := c.post(ctx, "/credential_connections", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetCredentialConnection retrieves a credential connection by ID.
func (c *Client) GetCredentialConnection(ctx context.Context, id string) (*CredentialConnection, error) {
	var resp credentialConnectionResponse
	if err := c.get(ctx, fmt.Sprintf("/credential_connections/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UpdateCredentialConnectionRequest carries updatable fields.
type UpdateCredentialConnectionRequest struct {
	UserName                      string             `json:"user_name,omitempty"`
	Password                      string             `json:"password,omitempty"`
	ConnectionName                string             `json:"connection_name,omitempty"`
	Active                        *bool              `json:"active,omitempty"`
	WebhookEventURL               string             `json:"webhook_event_url,omitempty"`
	WebhookEventFailoverURL       string             `json:"webhook_event_failover_url,omitempty"`
	WebhookAPIVersion             string             `json:"webhook_api_version,omitempty"`
	WebhookTimeoutSecs            int                `json:"webhook_timeout_secs,omitempty"`
	NoiseSuppression              string             `json:"noise_suppression,omitempty"`
	Tags                          []string           `json:"tags,omitempty"`
	RTCPSettings                  *RTCPSettings      `json:"rtcp_settings,omitempty"`
	Inbound                       *ConnectionInbound `json:"inbound,omitempty"`
	Outbound                      *ConnectionOutbound `json:"outbound,omitempty"`
	JitterBuffer                  *JitterBuffer      `json:"jitter_buffer,omitempty"`
}

// UpdateCredentialConnection applies partial updates to a credential connection.
func (c *Client) UpdateCredentialConnection(ctx context.Context, id string, req UpdateCredentialConnectionRequest) (*CredentialConnection, error) {
	var resp credentialConnectionResponse
	if err := c.patch(ctx, fmt.Sprintf("/credential_connections/%s", id), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeleteCredentialConnection removes a credential connection.
func (c *Client) DeleteCredentialConnection(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/credential_connections/%s", id))
}

// ListCredentialConnections returns all credential connections.
func (c *Client) ListCredentialConnections(ctx context.Context, pageNumber, pageSize int) ([]CredentialConnection, *PaginationMeta, error) {
	path := fmt.Sprintf("/credential_connections?page[number]=%d&page[size]=%d", pageNumber, pageSize)
	var resp credentialConnectionListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// IPConnection is a SIP IP-authenticated connection.
type IPConnection struct {
	RecordType                    string             `json:"record_type"`
	ID                            string             `json:"id"`
	Active                        bool               `json:"active"`
	ConnectionName                string             `json:"connection_name"`
	TransportProtocol             string             `json:"transport_protocol"`
	AnchorSiteOverride            string             `json:"anchorsite_override"`
	DefaultOnHoldComfortNoiseEnabled bool            `json:"default_on_hold_comfort_noise_enabled"`
	DTMFType                      string             `json:"dtmf_type"`
	EncodeContactHeaderEnabled    bool               `json:"encode_contact_header_enabled"`
	EncryptedMedia                string             `json:"encrypted_media"`
	OnnetT38PassthroughEnabled    bool               `json:"onnet_t38_passthrough_enabled"`
	WebhookEventURL               string             `json:"webhook_event_url"`
	WebhookEventFailoverURL       string             `json:"webhook_event_failover_url"`
	WebhookAPIVersion             string             `json:"webhook_api_version"`
	WebhookTimeoutSecs            int                `json:"webhook_timeout_secs"`
	NoiseSuppression              string             `json:"noise_suppression"`
	RTCPSettings                  *RTCPSettings      `json:"rtcp_settings,omitempty"`
	Inbound                       *ConnectionInbound `json:"inbound,omitempty"`
	Outbound                      *ConnectionOutbound `json:"outbound,omitempty"`
	JitterBuffer                  *JitterBuffer      `json:"jitter_buffer,omitempty"`
	Tags                          []string           `json:"tags"`
	CreatedAt                     string             `json:"created_at"`
	UpdatedAt                     string             `json:"updated_at"`
}

type ipConnectionResponse struct {
	Data IPConnection `json:"data"`
}

type ipConnectionListResponse struct {
	Data []IPConnection `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// CreateIPConnectionRequest creates a SIP IP-authenticated connection.
// All fields are optional.
type CreateIPConnectionRequest struct {
	ConnectionName                string             `json:"connection_name,omitempty"`
	Active                        bool               `json:"active,omitempty"`
	TransportProtocol             string             `json:"transport_protocol,omitempty"` // UDP | TCP | TLS
	AnchorSiteOverride            string             `json:"anchorsite_override,omitempty"`
	DefaultOnHoldComfortNoiseEnabled bool            `json:"default_on_hold_comfort_noise_enabled,omitempty"`
	DTMFType                      string             `json:"dtmf_type,omitempty"` // RFC 2833 | Inband | SIP INFO
	EncodeContactHeaderEnabled    bool               `json:"encode_contact_header_enabled,omitempty"`
	EncryptedMedia                string             `json:"encrypted_media,omitempty"`
	OnnetT38PassthroughEnabled    bool               `json:"onnet_t38_passthrough_enabled,omitempty"`
	IOSPushCredentialID           string             `json:"ios_push_credential_id,omitempty"`
	AndroidPushCredentialID       string             `json:"android_push_credential_id,omitempty"`
	WebhookEventURL               string             `json:"webhook_event_url,omitempty"`
	WebhookEventFailoverURL       string             `json:"webhook_event_failover_url,omitempty"`
	WebhookAPIVersion             string             `json:"webhook_api_version,omitempty"` // 1 | 2
	WebhookTimeoutSecs            int                `json:"webhook_timeout_secs,omitempty"`
	CallCostInWebhooks            bool               `json:"call_cost_in_webhooks,omitempty"`
	NoiseSuppression              string             `json:"noise_suppression,omitempty"` // inbound | outbound | both | disabled
	RTCPSettings                  *RTCPSettings      `json:"rtcp_settings,omitempty"`
	Inbound                       *ConnectionInbound `json:"inbound,omitempty"`
	Outbound                      *ConnectionOutbound `json:"outbound,omitempty"`
	JitterBuffer                  *JitterBuffer      `json:"jitter_buffer,omitempty"`
	Tags                          []string           `json:"tags,omitempty"`
}

// CreateIPConnection creates a new IP-authenticated SIP connection.
func (c *Client) CreateIPConnection(ctx context.Context, req CreateIPConnectionRequest) (*IPConnection, error) {
	var resp ipConnectionResponse
	if err := c.post(ctx, "/ip_connections", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetIPConnection retrieves an IP connection by ID.
func (c *Client) GetIPConnection(ctx context.Context, id string) (*IPConnection, error) {
	var resp ipConnectionResponse
	if err := c.get(ctx, fmt.Sprintf("/ip_connections/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UpdateIPConnectionRequest carries updatable fields.
type UpdateIPConnectionRequest struct {
	ConnectionName          string             `json:"connection_name,omitempty"`
	Active                  *bool              `json:"active,omitempty"`
	TransportProtocol       string             `json:"transport_protocol,omitempty"`
	WebhookEventURL         string             `json:"webhook_event_url,omitempty"`
	WebhookEventFailoverURL string             `json:"webhook_event_failover_url,omitempty"`
	WebhookAPIVersion       string             `json:"webhook_api_version,omitempty"`
	WebhookTimeoutSecs      int                `json:"webhook_timeout_secs,omitempty"`
	NoiseSuppression        string             `json:"noise_suppression,omitempty"`
	Tags                    []string           `json:"tags,omitempty"`
	RTCPSettings            *RTCPSettings      `json:"rtcp_settings,omitempty"`
	Inbound                 *ConnectionInbound `json:"inbound,omitempty"`
	Outbound                *ConnectionOutbound `json:"outbound,omitempty"`
	JitterBuffer            *JitterBuffer      `json:"jitter_buffer,omitempty"`
}

// UpdateIPConnection applies partial updates to an IP connection.
func (c *Client) UpdateIPConnection(ctx context.Context, id string, req UpdateIPConnectionRequest) (*IPConnection, error) {
	var resp ipConnectionResponse
	if err := c.patch(ctx, fmt.Sprintf("/ip_connections/%s", id), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeleteIPConnection removes an IP connection.
func (c *Client) DeleteIPConnection(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/ip_connections/%s", id))
}

// ListIPConnections returns all IP connections.
func (c *Client) ListIPConnections(ctx context.Context, pageNumber, pageSize int) ([]IPConnection, *PaginationMeta, error) {
	path := fmt.Sprintf("/ip_connections?page[number]=%d&page[size]=%d", pageNumber, pageSize)
	var resp ipConnectionListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}
