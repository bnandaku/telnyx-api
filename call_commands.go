package telnyx

import (
	"context"
	"fmt"
)

func (c *Client) action(ctx context.Context, callControlID, action string, body any) (*CallResponse, error) {
	var resp CallResponse
	path := fmt.Sprintf("/calls/%s/actions/%s", callControlID, action)
	if err := c.post(ctx, path, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AnswerRequest is the payload for answering an inbound call.
type AnswerRequest struct {
	BillingGroupID string                   `json:"billing_group_id,omitempty"`
	ClientState    string                   `json:"client_state,omitempty"`
	CommandID      string                   `json:"command_id,omitempty"`
	CustomHeaders  []CustomSipHeader        `json:"custom_headers,omitempty"`
	SIPHeaders     []CustomSipHeader        `json:"sip_headers,omitempty"`
	PreferredCodecs string                  `json:"preferred_codecs,omitempty"`
	SoundModifications *SoundModifications  `json:"sound_modifications,omitempty"`
	Assistant      *CallAssistantRequest    `json:"assistant,omitempty"`

	// Streaming
	StreamURL                     string `json:"stream_url,omitempty"`
	StreamTrack                   string `json:"stream_track,omitempty"`
	StreamCodec                   string `json:"stream_codec,omitempty"`
	StreamBidirectionalMode       string `json:"stream_bidirectional_mode,omitempty"`
	StreamBidirectionalCodec      string `json:"stream_bidirectional_codec,omitempty"`
	StreamBidirectionalTargetLegs string `json:"stream_bidirectional_target_legs,omitempty"`
	SendSilenceWhenIdle           bool   `json:"send_silence_when_idle,omitempty"`

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

	// Deepfake detection
	DeepfakeDetection *DeepfakeDetectionConfig `json:"deepfake_detection,omitempty"`

	WebhookURL        string            `json:"webhook_url,omitempty"`
	WebhookURLMethod  string            `json:"webhook_url_method,omitempty"`
	WebhookURLs       map[string]string `json:"webhook_urls,omitempty"`
	WebhookURLsMethod string            `json:"webhook_urls_method,omitempty"`
}

// Answer answers an inbound call.
func (c *Client) Answer(ctx context.Context, callControlID string, req AnswerRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "answer", req)
}

// HangupRequest is the payload for hanging up a call.
type HangupRequest struct {
	ClientState   string            `json:"client_state,omitempty"`
	CommandID     string            `json:"command_id,omitempty"`
	CustomHeaders []CustomSipHeader `json:"custom_headers,omitempty"`
}

// Hangup ends a call leg.
func (c *Client) Hangup(ctx context.Context, callControlID string, req HangupRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "hangup", req)
}

// RejectRequest rejects an inbound call.
type RejectRequest struct {
	Cause     string `json:"cause,omitempty"` // USER_BUSY or CALL_REJECTED
	ClientState string `json:"client_state,omitempty"`
	CommandID   string `json:"command_id,omitempty"`
}

// Reject rejects an inbound call without answering.
func (c *Client) Reject(ctx context.Context, callControlID string, req RejectRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "reject", req)
}

// BridgeRequest connects two call legs.
type BridgeRequest struct {
	CallControlID      string `json:"call_control_id,omitempty"`
	ClientState        string `json:"client_state,omitempty"`
	CommandID          string `json:"command_id,omitempty"`
	Queue              string `json:"queue,omitempty"`
	VideoRoomID        string `json:"video_room_id,omitempty"`
	VideoRoomContext   string `json:"video_room_context,omitempty"`
	PreventDoubleBridge bool  `json:"prevent_double_bridge,omitempty"`
	ParkAfterUnbridge  string `json:"park_after_unbridge,omitempty"`
	HoldAfterUnbridge  bool   `json:"hold_after_unbridge,omitempty"`
	PlayRingtone       bool   `json:"play_ringtone,omitempty"`
	Ringtone           string `json:"ringtone,omitempty"`
	MuteDTMF           string `json:"mute_dtmf,omitempty"`

	Record               string `json:"record,omitempty"`
	RecordChannels       string `json:"record_channels,omitempty"`
	RecordFormat         string `json:"record_format,omitempty"`
	RecordMaxLength      int    `json:"record_max_length,omitempty"`
	RecordTimeoutSecs    int    `json:"record_timeout_secs,omitempty"`
	RecordTrack          string `json:"record_track,omitempty"`
	RecordTrim           string `json:"record_trim,omitempty"`
	RecordCustomFileName string `json:"record_custom_file_name,omitempty"`
}

// Bridge connects this call leg to another.
func (c *Client) Bridge(ctx context.Context, callControlID string, req BridgeRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "bridge", req)
}

// GatherRequest collects DTMF digits from the caller.
type GatherRequest struct {
	MinimumDigits           int    `json:"minimum_digits,omitempty"`
	MaximumDigits           int    `json:"maximum_digits,omitempty"`
	TimeoutMillis           int    `json:"timeout_millis,omitempty"`
	InterDigitTimeoutMillis int    `json:"inter_digit_timeout_millis,omitempty"`
	InitialTimeoutMillis    int    `json:"initial_timeout_millis,omitempty"`
	TerminatingDigit        string `json:"terminating_digit,omitempty"`
	ValidDigits             string `json:"valid_digits,omitempty"`
	GatherID                string `json:"gather_id,omitempty"`
	ClientState             string `json:"client_state,omitempty"`
	CommandID               string `json:"command_id,omitempty"`
}

// Gather starts a DTMF gather operation.
func (c *Client) Gather(ctx context.Context, callControlID string, req GatherRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "gather", req)
}

// GatherStop cancels an in-progress gather.
func (c *Client) GatherStop(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "gather_stop", body)
}

// SpeakRequest synthesises text-to-speech on a call.
type SpeakRequest struct {
	Payload      string `json:"payload"`
	Voice        string `json:"voice"`
	PayloadType  string `json:"payload_type,omitempty"`  // text | ssml
	ServiceLevel string `json:"service_level,omitempty"` // basic | premium
	Language     string `json:"language,omitempty"`
	Stop         string `json:"stop,omitempty"`        // current | all
	Loop         any    `json:"loop,omitempty"`        // int or "infinity"
	TargetLegs   string `json:"target_legs,omitempty"` // self | opposite | both
	VoiceSettings *VoiceSettings `json:"voice_settings,omitempty"`
	ClientState  string `json:"client_state,omitempty"`
	CommandID    string `json:"command_id,omitempty"`
}

// Speak plays synthesised text-to-speech on the call.
func (c *Client) Speak(ctx context.Context, callControlID string, req SpeakRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "speak", req)
}

// PlayAudioRequest plays an audio file on the call.
type PlayAudioRequest struct {
	AudioURL        string `json:"audio_url,omitempty"`
	MediaName       string `json:"media_name,omitempty"`
	PlaybackContent string `json:"playback_content,omitempty"`
	Loop            any    `json:"loop,omitempty"` // int or "infinity"
	Overlay         bool   `json:"overlay,omitempty"`
	Stop            string `json:"stop,omitempty"`       // current | all
	TargetLegs      string `json:"target_legs,omitempty"`
	CacheAudio      *bool  `json:"cache_audio,omitempty"`
	AudioType       string `json:"audio_type,omitempty"`
	ClientState     string `json:"client_state,omitempty"`
	CommandID       string `json:"command_id,omitempty"`
}

// PlayAudio plays an audio file on the call.
func (c *Client) PlayAudio(ctx context.Context, callControlID string, req PlayAudioRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "playback_start", req)
}

// PlaybackStop stops audio playback on the call.
func (c *Client) PlaybackStop(ctx context.Context, callControlID string, stop, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"stop":         stop,
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "playback_stop", body)
}

// TransferRequest transfers a call to a new destination.
type TransferRequest struct {
	To          string `json:"to"`
	From        string `json:"from,omitempty"`
	FromDisplayName string `json:"from_display_name,omitempty"`
	Privacy     string `json:"privacy,omitempty"`
	AudioURL    string `json:"audio_url,omitempty"`
	MediaName   string `json:"media_name,omitempty"`
	TimeoutSecs int    `json:"timeout_secs,omitempty"`
	TimeLimitSecs int  `json:"time_limit_secs,omitempty"`
	ParkAfterUnbridge string `json:"park_after_unbridge,omitempty"`
	EarlyMedia  bool   `json:"early_media,omitempty"`
	MuteDTMF    string `json:"mute_dtmf,omitempty"`

	AnsweringMachineDetection       string     `json:"answering_machine_detection,omitempty"`
	AnsweringMachineDetectionConfig *AMDConfig `json:"answering_machine_detection_config,omitempty"`

	CustomHeaders        []CustomSipHeader `json:"custom_headers,omitempty"`
	SIPHeaders           []CustomSipHeader `json:"sip_headers,omitempty"`
	SIPAuthUsername      string            `json:"sip_auth_username,omitempty"`
	SIPAuthPassword      string            `json:"sip_auth_password,omitempty"`
	SIPTransportProtocol string            `json:"sip_transport_protocol,omitempty"`
	SIPRegion            string            `json:"sip_region,omitempty"`
	PreferredCodecs      string            `json:"preferred_codecs,omitempty"`
	MediaEncryption      string            `json:"media_encryption,omitempty"`
	SoundModifications   *SoundModifications `json:"sound_modifications,omitempty"`

	Record               string `json:"record,omitempty"`
	RecordChannels       string `json:"record_channels,omitempty"`
	RecordFormat         string `json:"record_format,omitempty"`
	RecordMaxLength      int    `json:"record_max_length,omitempty"`
	RecordTimeoutSecs    int    `json:"record_timeout_secs,omitempty"`
	RecordTrack          string `json:"record_track,omitempty"`
	RecordTrim           string `json:"record_trim,omitempty"`
	RecordCustomFileName string `json:"record_custom_file_name,omitempty"`

	ClientState          string            `json:"client_state,omitempty"`
	TargetLegClientState string            `json:"target_leg_client_state,omitempty"`
	CommandID            string            `json:"command_id,omitempty"`
	WebhookURL           string            `json:"webhook_url,omitempty"`
	WebhookURLMethod     string            `json:"webhook_url_method,omitempty"`
	WebhookURLs          map[string]string `json:"webhook_urls,omitempty"`
	WebhookURLsMethod    string            `json:"webhook_urls_method,omitempty"`
}

// Transfer redirects a call to a new destination.
func (c *Client) Transfer(ctx context.Context, callControlID string, req TransferRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "transfer", req)
}

// RecordStartRequest starts recording a call.
type RecordStartRequest struct {
	Format         string `json:"format"`   // wav | mp3
	Channels       string `json:"channels"` // single | dual
	ClientState    string `json:"client_state,omitempty"`
	CommandID      string `json:"command_id,omitempty"`
	PlayBeep       bool   `json:"play_beep,omitempty"`
	MaxLength      int    `json:"max_length,omitempty"`
	TimeoutSecs    int    `json:"timeout_secs,omitempty"`
	RecordingTrack string `json:"recording_track,omitempty"` // both | inbound | outbound
	Trim           string `json:"trim,omitempty"`
	CustomFileName string `json:"custom_file_name,omitempty"`

	Transcription                   bool   `json:"transcription,omitempty"`
	TranscriptionEngine             string `json:"transcription_engine,omitempty"`
	TranscriptionLanguage           string `json:"transcription_language,omitempty"`
	TranscriptionProfanityFilter    bool   `json:"transcription_profanity_filter,omitempty"`
	TranscriptionSpeakerDiarization bool   `json:"transcription_speaker_diarization,omitempty"`
	TranscriptionMinSpeakerCount    int    `json:"transcription_min_speaker_count,omitempty"`
	TranscriptionMaxSpeakerCount    int    `json:"transcription_max_speaker_count,omitempty"`
}

// RecordStart begins recording a call.
func (c *Client) RecordStart(ctx context.Context, callControlID string, req RecordStartRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "record_start", req)
}

// RecordStop stops an active recording.
func (c *Client) RecordStop(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "record_stop", body)
}

// RecordPause pauses an active recording.
func (c *Client) RecordPause(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "record_pause", body)
}

// RecordResume resumes a paused recording.
func (c *Client) RecordResume(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "record_resume", body)
}

// HoldRequest places a call on hold.
type HoldRequest struct {
	AudioURL    string `json:"audio_url,omitempty"`
	MediaName   string `json:"media_name,omitempty"`
	ClientState string `json:"client_state,omitempty"`
	CommandID   string `json:"command_id,omitempty"`
}

// Hold places a call leg on hold.
func (c *Client) Hold(ctx context.Context, callControlID string, req HoldRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "hold", req)
}

// Unhold takes a call off hold.
func (c *Client) Unhold(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "unhold", body)
}

// MuteRequest mutes audio on a call leg.
type MuteRequest struct {
	Mute        string `json:"mute,omitempty"` // self | opposite | both
	ClientState string `json:"client_state,omitempty"`
	CommandID   string `json:"command_id,omitempty"`
}

// Mute mutes audio on the call.
func (c *Client) Mute(ctx context.Context, callControlID string, req MuteRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "mute", req)
}

// Unmute unmutes a previously muted call.
func (c *Client) Unmute(ctx context.Context, callControlID string, req MuteRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "unmute", req)
}

// ForkStartRequest begins media forking to a remote target.
type ForkStartRequest struct {
	Target          string `json:"target"`
	RxMode          string `json:"rx,omitempty"`
	TxMode          string `json:"tx,omitempty"`
	StreamingOptions string `json:"streaming_options,omitempty"`
	ClientState     string `json:"client_state,omitempty"`
	CommandID       string `json:"command_id,omitempty"`
}

// ForkStart starts media forking.
func (c *Client) ForkStart(ctx context.Context, callControlID string, req ForkStartRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "fork_start", req)
}

// ForkStop stops media forking.
func (c *Client) ForkStop(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "fork_stop", body)
}

// TranscriptionStartRequest enables live transcription on a call.
type TranscriptionStartRequest struct {
	TranscriptionEngine   string `json:"transcription_engine,omitempty"`
	TranscriptionLanguage string `json:"transcription_language,omitempty"`
	InterimResults        bool   `json:"interim_results,omitempty"`
	ClientState           string `json:"client_state,omitempty"`
	CommandID             string `json:"command_id,omitempty"`
}

// TranscriptionStart enables real-time transcription.
func (c *Client) TranscriptionStart(ctx context.Context, callControlID string, req TranscriptionStartRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "transcription_start", req)
}

// TranscriptionStop disables real-time transcription.
func (c *Client) TranscriptionStop(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "transcription_stop", body)
}

// StreamingStartRequest starts media streaming over WebSocket.
type StreamingStartRequest struct {
	StreamURL                     string `json:"stream_url"`
	StreamTrack                   string `json:"stream_track,omitempty"`
	StreamBidirectionalMode       string `json:"stream_bidirectional_mode,omitempty"`
	StreamBidirectionalCodec      string `json:"stream_bidirectional_codec,omitempty"`
	StreamBidirectionalTargetLegs string `json:"stream_bidirectional_target_legs,omitempty"`
	ClientState                   string `json:"client_state,omitempty"`
	CommandID                     string `json:"command_id,omitempty"`
	SendSilenceWhenIdle           bool   `json:"send_silence_when_idle,omitempty"`
}

// StreamingStart starts WebSocket media streaming.
func (c *Client) StreamingStart(ctx context.Context, callControlID string, req StreamingStartRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "streaming_start", req)
}

// StreamingStop stops WebSocket media streaming.
func (c *Client) StreamingStop(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "streaming_stop", body)
}

// SendDTMFRequest sends DTMF digits on a call.
type SendDTMFRequest struct {
	Digits      string `json:"digits"`
	DurationMS  int    `json:"duration_millis,omitempty"`
	ClientState string `json:"client_state,omitempty"`
	CommandID   string `json:"command_id,omitempty"`
}

// SendDTMF sends DTMF tones on the call.
func (c *Client) SendDTMF(ctx context.Context, callControlID string, req SendDTMFRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "send_dtmf", req)
}

// ReferRequest performs a SIP REFER.
type ReferRequest struct {
	SIPAddress  string            `json:"sip_address"`
	SIPHeaders  []CustomSipHeader `json:"sip_headers,omitempty"`
	ClientState string            `json:"client_state,omitempty"`
	CommandID   string            `json:"command_id,omitempty"`
}

// Refer sends a SIP REFER to transfer the call.
func (c *Client) Refer(ctx context.Context, callControlID string, req ReferRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "refer", req)
}

// SIPRECStartRequest starts a SIPREC session.
type SIPRECStartRequest struct {
	SIPRECDestinationAddress string            `json:"siprec_destination_address"`
	SIPRECTrack              string            `json:"siprec_track,omitempty"`
	SIPHeaders               []CustomSipHeader `json:"sip_headers,omitempty"`
	ClientState              string            `json:"client_state,omitempty"`
	CommandID                string            `json:"command_id,omitempty"`
}

// SIPRECStart starts a SIPREC recording session.
func (c *Client) SIPRECStart(ctx context.Context, callControlID string, req SIPRECStartRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "siprec_start", req)
}

// SIPRECStop stops a SIPREC session.
func (c *Client) SIPRECStop(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "siprec_stop", body)
}

// EnqueueRequest places a call in a queue.
type EnqueueRequest struct {
	QueueName      string `json:"queue_name"`
	MaxSize        int    `json:"max_size,omitempty"`
	MaxWaitTimeSecs int   `json:"max_wait_time_secs,omitempty"`
	ClientState    string `json:"client_state,omitempty"`
	CommandID      string `json:"command_id,omitempty"`
}

// Enqueue places a call into a named queue.
func (c *Client) Enqueue(ctx context.Context, callControlID string, req EnqueueRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "enqueue", req)
}

// LeaveQueue removes a call from its queue.
func (c *Client) LeaveQueue(ctx context.Context, callControlID string, clientState, commandID string) (*CallResponse, error) {
	body := map[string]string{
		"client_state": clientState,
		"command_id":   commandID,
	}
	return c.action(ctx, callControlID, "leave_queue", body)
}

// StartAIAssistantRequest attaches a conversational AI assistant to a call.
type StartAIAssistantRequest struct {
	Assistant            CallAssistantRequest   `json:"assistant"`
	Voice                string                 `json:"voice,omitempty"`
	VoiceSettings        *VoiceSettings         `json:"voice_settings,omitempty"`
	Greeting             string                 `json:"greeting,omitempty"`
	InterruptionSettings *InterruptionSettings  `json:"interruption_settings,omitempty"`
	Transcription        *TranscriptionConfig   `json:"transcription,omitempty"`
	MessageHistory       []MessageHistoryEntry  `json:"message_history,omitempty"`
	SendMessageHistoryUpdates bool              `json:"send_message_history_updates,omitempty"`
	Participants         []Participant          `json:"participants,omitempty"`
	ClientState          string                 `json:"client_state,omitempty"`
	CommandID            string                 `json:"command_id,omitempty"`
}

// StartAIAssistantResponse includes the created conversation ID.
type StartAIAssistantResponse struct {
	Data struct {
		Result         string `json:"result"`
		ConversationID string `json:"conversation_id"`
	} `json:"data"`
}

// StartAIAssistant attaches an AI assistant to an active call.
func (c *Client) StartAIAssistant(ctx context.Context, callControlID string, req StartAIAssistantRequest) (*StartAIAssistantResponse, error) {
	var resp StartAIAssistantResponse
	path := fmt.Sprintf("/calls/%s/actions/ai_assistant_start", callControlID)
	if err := c.post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StopAIAssistantRequest detaches the AI assistant from a call.
type StopAIAssistantRequest struct {
	ClientState string `json:"client_state,omitempty"`
	CommandID   string `json:"command_id,omitempty"`
}

// StopAIAssistant detaches the AI assistant from an active call.
func (c *Client) StopAIAssistant(ctx context.Context, callControlID string, req StopAIAssistantRequest) (*CallResponse, error) {
	return c.action(ctx, callControlID, "ai_assistant_stop", req)
}

// UpdateClientStateRequest updates the client_state on an active call.
type UpdateClientStateRequest struct {
	ClientState string `json:"client_state"`
}

// UpdateClientState updates the passthrough client state on an active call.
func (c *Client) UpdateClientState(ctx context.Context, callControlID string, clientState string) (*CallResponse, error) {
	return c.action(ctx, callControlID, "client_state_update", UpdateClientStateRequest{ClientState: clientState})
}
