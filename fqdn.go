package telnyx

import (
	"context"
	"fmt"
)

type FQDNConnection struct {
	ID                               string        `json:"id"`
	RecordType                       string        `json:"record_type"`
	ConnectionName                   string        `json:"connection_name"`
	Active                           bool          `json:"active"`
	AnchorSiteOverride               string        `json:"anchorsite_override,omitempty"`
	TransportProtocol                string        `json:"transport_protocol,omitempty"`
	DefaultOnHoldComfortNoiseEnabled bool          `json:"default_on_hold_comfort_noise_enabled,omitempty"`
	DTMFType                         string        `json:"dtmf_type,omitempty"`
	EncodeContactHeaderEnabled       bool          `json:"encode_contact_header_enabled,omitempty"`
	EncryptedMedia                   string        `json:"encrypted_media,omitempty"`
	MicrosoftTeamsSBC                bool          `json:"microsoft_teams_sbc,omitempty"`
	OnnetT38PassthroughEnabled       bool          `json:"onnet_t38_passthrough_enabled,omitempty"`
	Tags                             []string      `json:"tags,omitempty"`
	WebhookEventURL                  string        `json:"webhook_event_url,omitempty"`
	WebhookEventFailoverURL          string        `json:"webhook_event_failover_url,omitempty"`
	WebhookAPIVersion                string        `json:"webhook_api_version,omitempty"`
	WebhookTimeoutSecs               int           `json:"webhook_timeout_secs,omitempty"`
	CallCostInWebhooks               bool          `json:"call_cost_in_webhooks,omitempty"`
	NoiseSuppression                 string        `json:"noise_suppression,omitempty"`
	RTCPSettings                     *RTCPSettings `json:"rtcp_settings,omitempty"`
	Inbound                          *FQDNInbound  `json:"inbound,omitempty"`
	Outbound                         *FQDNOutbound `json:"outbound,omitempty"`
	CreatedAt                        string        `json:"created_at"`
	UpdatedAt                        string        `json:"updated_at"`
}

type FQDNInbound struct {
	ANINumberFormat             string   `json:"ani_number_format,omitempty"`
	DNISNumberFormat            string   `json:"dnis_number_format,omitempty"`
	Codecs                      []string `json:"codecs,omitempty"`
	DefaultRoutingMethod        string   `json:"default_routing_method,omitempty"`
	DefaultPrimaryFQDNID        string   `json:"default_primary_fqdn_id,omitempty"`
	DefaultSecondaryFQDNID      string   `json:"default_secondary_fqdn_id,omitempty"`
	DefaultTertiaryFQDNID       string   `json:"default_tertiary_fqdn_id,omitempty"`
	ChannelLimit                *int     `json:"channel_limit,omitempty"`
	GenerateRingbackTone        bool     `json:"generate_ringback_tone,omitempty"`
	ISUPHeadersEnabled          bool     `json:"isup_headers_enabled,omitempty"`
	PRACKEnabled                bool     `json:"prack_enabled,omitempty"`
	SIPCompactHeadersEnabled    bool     `json:"sip_compact_headers_enabled,omitempty"`
	SIPRegion                   string   `json:"sip_region,omitempty"`
	SIPSubdomain                string   `json:"sip_subdomain,omitempty"`
	SIPSubdomainReceiveSettings string   `json:"sip_subdomain_receive_settings,omitempty"`
	Timeout1XXSecs              int      `json:"timeout_1xx_secs,omitempty"`
	Timeout2XXSecs              int      `json:"timeout_2xx_secs,omitempty"`
	SHAKENSTIREnabled           bool     `json:"shaken_stir_enabled,omitempty"`
}

type FQDNOutbound struct {
	ANIOverride            string `json:"ani_override,omitempty"`
	ANIOverrideType        string `json:"ani_override_type,omitempty"`
	CallParkingEnabled     *bool  `json:"call_parking_enabled,omitempty"`
	ChannelLimit           int    `json:"channel_limit,omitempty"`
	GenerateRingbackTone   bool   `json:"generate_ringback_tone,omitempty"`
	InstantRingbackEnabled bool   `json:"instant_ringback_enabled,omitempty"`
	IPAuthenticationMethod string `json:"ip_authentication_method,omitempty"`
	IPAuthenticationToken  string `json:"ip_authentication_token,omitempty"`
	Localization           string `json:"localization,omitempty"`
	OutboundVoiceProfileID string `json:"outbound_voice_profile_id,omitempty"`
	T38ReinviteSource      string `json:"t38_reinvite_source,omitempty"`
	TechPrefix             string `json:"tech_prefix,omitempty"`
	EncryptedMedia         string `json:"encrypted_media,omitempty"`
	Timeout1XXSecs         int    `json:"timeout_1xx_secs,omitempty"`
	Timeout2XXSecs         int    `json:"timeout_2xx_secs,omitempty"`
}

type CreateFQDNConnectionRequest struct {
	ConnectionName                   string        `json:"connection_name"`
	Active                           bool          `json:"active,omitempty"`
	AnchorSiteOverride               string        `json:"anchorsite_override,omitempty"`
	TransportProtocol                string        `json:"transport_protocol,omitempty"`
	DefaultOnHoldComfortNoiseEnabled bool          `json:"default_on_hold_comfort_noise_enabled,omitempty"`
	DTMFType                         string        `json:"dtmf_type,omitempty"`
	EncodeContactHeaderEnabled       bool          `json:"encode_contact_header_enabled,omitempty"`
	EncryptedMedia                   string        `json:"encrypted_media,omitempty"`
	MicrosoftTeamsSBC                bool          `json:"microsoft_teams_sbc,omitempty"`
	OnnetT38PassthroughEnabled       bool          `json:"onnet_t38_passthrough_enabled,omitempty"`
	Tags                             []string      `json:"tags,omitempty"`
	IOSPushCredentialID              string        `json:"ios_push_credential_id,omitempty"`
	AndroidPushCredentialID          string        `json:"android_push_credential_id,omitempty"`
	WebhookEventURL                  string        `json:"webhook_event_url,omitempty"`
	WebhookEventFailoverURL          string        `json:"webhook_event_failover_url,omitempty"`
	WebhookAPIVersion                string        `json:"webhook_api_version,omitempty"`
	WebhookTimeoutSecs               int           `json:"webhook_timeout_secs,omitempty"`
	CallCostInWebhooks               bool          `json:"call_cost_in_webhooks,omitempty"`
	NoiseSuppression                 string        `json:"noise_suppression,omitempty"`
	RTCPSettings                     *RTCPSettings `json:"rtcp_settings,omitempty"`
	Inbound                          *FQDNInbound  `json:"inbound,omitempty"`
	Outbound                         *FQDNOutbound `json:"outbound,omitempty"`
}

type FQDN struct {
	ID            string `json:"id"`
	RecordType    string `json:"record_type"`
	ConnectionID  string `json:"connection_id"`
	FQDN          string `json:"fqdn"`
	Port          int    `json:"port,omitempty"`
	DNSRecordType string `json:"dns_record_type"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type CreateFQDNRequest struct {
	ConnectionID  string `json:"connection_id"`
	FQDN          string `json:"fqdn"`
	DNSRecordType string `json:"dns_record_type"`
	Port          int    `json:"port,omitempty"`
}

type ListFQDNConnectionsParams struct {
	FilterNameContains      string
	FilterFQDN              string
	FilterOutboundProfileID string
	PageNumber              int
	PageSize                int
	Sort                    string
}

func (c *Client) CreateFQDNConnection(ctx context.Context, req CreateFQDNConnectionRequest) (*FQDNConnection, error) {
	var out struct {
		Data FQDNConnection `json:"data"`
	}
	if err := c.post(ctx, "/fqdn_connections", req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetFQDNConnection(ctx context.Context, id string) (*FQDNConnection, error) {
	var out struct {
		Data FQDNConnection `json:"data"`
	}
	if err := c.get(ctx, fmt.Sprintf("/fqdn_connections/%s", id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateFQDNConnection(ctx context.Context, id string, req CreateFQDNConnectionRequest) (*FQDNConnection, error) {
	var out struct {
		Data FQDNConnection `json:"data"`
	}
	if err := c.patch(ctx, fmt.Sprintf("/fqdn_connections/%s", id), req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) DeleteFQDNConnection(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/fqdn_connections/%s", id))
}

func (c *Client) ListFQDNConnections(ctx context.Context, params ListFQDNConnectionsParams) ([]FQDNConnection, *PaginationMeta, error) {
	path := "/fqdn_connections"
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
	addStr("filter[connection_name.contains]", params.FilterNameContains)
	addStr("filter[fqdn]", params.FilterFQDN)
	addStr("filter[outbound_voice_profile_id]", params.FilterOutboundProfileID)
	addInt("page[number]", params.PageNumber)
	addInt("page[size]", params.PageSize)
	addStr("sort", params.Sort)

	var out struct {
		Data []FQDNConnection `json:"data"`
		Meta PaginationMeta   `json:"meta"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Data, &out.Meta, nil
}

func (c *Client) CreateFQDN(ctx context.Context, req CreateFQDNRequest) (*FQDN, error) {
	var out struct {
		Data FQDN `json:"data"`
	}
	if err := c.post(ctx, "/fqdns", req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetFQDN(ctx context.Context, id string) (*FQDN, error) {
	var out struct {
		Data FQDN `json:"data"`
	}
	if err := c.get(ctx, fmt.Sprintf("/fqdns/%s", id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateFQDN(ctx context.Context, id string, req CreateFQDNRequest) (*FQDN, error) {
	var out struct {
		Data FQDN `json:"data"`
	}
	if err := c.patch(ctx, fmt.Sprintf("/fqdns/%s", id), req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) DeleteFQDN(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/fqdns/%s", id))
}

func (c *Client) ListFQDNs(ctx context.Context, connectionID string, pageNumber, pageSize int) ([]FQDN, *PaginationMeta, error) {
	path := "/fqdns"
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
	addStr("filter[connection_id]", connectionID)
	addInt("page[number]", pageNumber)
	addInt("page[size]", pageSize)

	var out struct {
		Data []FQDN         `json:"data"`
		Meta PaginationMeta `json:"meta"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Data, &out.Meta, nil
}
