package telnyx

import (
	"context"
	"fmt"
)

// DialRequest is the payload for initiating an outbound call.
type DialRequest struct {
	ConnectionID string `json:"connection_id"`
	To           any    `json:"to"` // string or []string
	From         string `json:"from"`

	FromDisplayName string `json:"from_display_name,omitempty"`
	Privacy         string `json:"privacy,omitempty"`
	TimeoutSecs     int    `json:"timeout_secs,omitempty"`
	TimeLimitSecs   int    `json:"time_limit_secs,omitempty"`

	// Answering machine detection
	AnsweringMachineDetection       string     `json:"answering_machine_detection,omitempty"`
	AnsweringMachineDetectionConfig *AMDConfig `json:"answering_machine_detection_config,omitempty"`

	// Deepfake detection
	DeepfakeDetection *DeepfakeDetectionConfig `json:"deepfake_detection,omitempty"`

	// Media streaming
	StreamURL                     string `json:"stream_url,omitempty"`
	StreamTrack                   string `json:"stream_track,omitempty"`
	StreamCodec                   string `json:"stream_codec,omitempty"`
	StreamBidirectionalMode       string `json:"stream_bidirectional_mode,omitempty"`
	StreamBidirectionalCodec      string `json:"stream_bidirectional_codec,omitempty"`
	StreamBidirectionalTargetLegs string `json:"stream_bidirectional_target_legs,omitempty"`
	SendSilenceWhenIdle           bool   `json:"send_silence_when_idle,omitempty"`

	// Audio
	AudioURL        string `json:"audio_url,omitempty"`
	MediaName       string `json:"media_name,omitempty"`
	PreferredCodecs string `json:"preferred_codecs,omitempty"`
	MediaEncryption string `json:"media_encryption,omitempty"`

	// Recording
	Record               string `json:"record,omitempty"`
	RecordChannels       string `json:"record_channels,omitempty"`
	RecordFormat         string `json:"record_format,omitempty"`
	RecordMaxLength      int    `json:"record_max_length,omitempty"`
	RecordTimeoutSecs    int    `json:"record_timeout_secs,omitempty"`
	RecordTrack          string `json:"record_track,omitempty"`
	RecordTrim           string `json:"record_trim,omitempty"`
	RecordCustomFileName string `json:"record_custom_file_name,omitempty"`

	// Transcription
	Transcription       bool                 `json:"transcription,omitempty"`
	TranscriptionConfig *TranscriptionConfig `json:"transcription_config,omitempty"`

	// SIP
	SIPAuthUsername      string            `json:"sip_auth_username,omitempty"`
	SIPAuthPassword      string            `json:"sip_auth_password,omitempty"`
	SIPHeaders           []CustomSipHeader `json:"sip_headers,omitempty"`
	CustomHeaders        []CustomSipHeader `json:"custom_headers,omitempty"`
	SIPTransportProtocol string            `json:"sip_transport_protocol,omitempty"`
	SIPRegion            string            `json:"sip_region,omitempty"`

	// Supervision
	SuperviseCallControlID string `json:"supervise_call_control_id,omitempty"`
	SupervisorRole         string `json:"supervisor_role,omitempty"`

	// Bridge
	BridgeIntent    bool   `json:"bridge_intent,omitempty"`
	BridgeOnAnswer  bool   `json:"bridge_on_answer,omitempty"`
	LinkTo          string `json:"link_to,omitempty"`

	// AI
	Assistant *CallAssistantRequest `json:"assistant,omitempty"`

	// Misc
	SoundModifications *SoundModifications `json:"sound_modifications,omitempty"`
	BillingGroupID     string              `json:"billing_group_id,omitempty"`
	ClientState        string              `json:"client_state,omitempty"`
	CommandID          string              `json:"command_id,omitempty"`

	WebhookURL             string            `json:"webhook_url,omitempty"`
	WebhookURLMethod       string            `json:"webhook_url_method,omitempty"`
	WebhookURLs            map[string]string `json:"webhook_urls,omitempty"`
	WebhookURLsMethod      string            `json:"webhook_urls_method,omitempty"`
	WebhookRetriesPolicies map[string]any    `json:"webhook_retries_policies,omitempty"`
}

// DialResponse wraps the call object returned from a Dial.
type DialResponse struct {
	Data Call `json:"data"`
}

// Dial initiates an outbound call.
func (c *Client) Dial(ctx context.Context, req DialRequest) (*Call, error) {
	var resp DialResponse
	if err := c.post(ctx, "/calls", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetCall retrieves an active call by its call_control_id.
func (c *Client) GetCall(ctx context.Context, callControlID string) (*Call, error) {
	var resp DialResponse
	if err := c.get(ctx, fmt.Sprintf("/calls/%s", callControlID), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
