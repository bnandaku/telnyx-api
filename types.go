package telnyx

// CustomSipHeader is a name/value pair added to SIP messages.
type CustomSipHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// SoundModifications adjusts audio characteristics.
type SoundModifications struct {
	Pitch    float64 `json:"pitch,omitempty"`
	Semitone float64 `json:"semitone,omitempty"`
	Octaves  float64 `json:"octaves,omitempty"`
	Track    string  `json:"track,omitempty"`
}

// TranscriptionConfig controls speech-to-text settings.
type TranscriptionConfig struct {
	TranscriptionEngine             string `json:"transcription_engine,omitempty"`
	TranscriptionLanguage           string `json:"transcription_language,omitempty"`
	TranscriptionProfanityFilter    bool   `json:"transcription_profanity_filter,omitempty"`
	TranscriptionSpeakerDiarization bool   `json:"transcription_speaker_diarization,omitempty"`
	TranscriptionMinSpeakerCount    int    `json:"transcription_min_speaker_count,omitempty"`
	TranscriptionMaxSpeakerCount    int    `json:"transcription_max_speaker_count,omitempty"`
}

// AMDConfig tunes answering machine detection behaviour.
type AMDConfig struct {
	AfterGreetingSilenceMillis int    `json:"after_greeting_silence_millis,omitempty"`
	BetweenWordsSilenceMillis  int    `json:"between_words_silence_millis,omitempty"`
	GreetingDurationMillis     int    `json:"greeting_duration_millis,omitempty"`
	InitialSilenceMillis       int    `json:"initial_silence_millis,omitempty"`
	MaximumNumberOfWords       int    `json:"maximum_number_of_words,omitempty"`
	MaximumWordLengthMillis    int    `json:"maximum_word_length_millis,omitempty"`
	SilenceThreshold           int    `json:"silence_threshold,omitempty"`
	TotalAnalysisTimeMillis    int    `json:"total_analysis_time_millis,omitempty"`
	TotalWaitingTimeMillis     int    `json:"total_waiting_time_millis,omitempty"`
	Greeting                   string `json:"greeting,omitempty"`
}

// DeepfakeDetectionConfig controls deepfake/voice-auth detection.
type DeepfakeDetectionConfig struct {
	Enabled    bool `json:"enabled,omitempty"`
	Timeout    int  `json:"timeout,omitempty"`
	RTPTimeout int  `json:"rtp_timeout,omitempty"`
}

// StreamConfig holds media streaming parameters.
type StreamConfig struct {
	StreamURL                      string `json:"stream_url,omitempty"`
	StreamTrack                    string `json:"stream_track,omitempty"`
	StreamCodec                    string `json:"stream_codec,omitempty"`
	StreamBidirectionalMode        string `json:"stream_bidirectional_mode,omitempty"`
	StreamBidirectionalCodec       string `json:"stream_bidirectional_codec,omitempty"`
	StreamBidirectionalTargetLegs  string `json:"stream_bidirectional_target_legs,omitempty"`
	SendSilenceWhenIdle            bool   `json:"send_silence_when_idle,omitempty"`
}

// RecordConfig holds common recording parameters.
type RecordConfig struct {
	Record               string `json:"record,omitempty"`
	RecordChannels       string `json:"record_channels,omitempty"`
	RecordFormat         string `json:"record_format,omitempty"`
	RecordMaxLength      int    `json:"record_max_length,omitempty"`
	RecordTimeoutSecs    int    `json:"record_timeout_secs,omitempty"`
	RecordTrack          string `json:"record_track,omitempty"`
	RecordTrim           string `json:"record_trim,omitempty"`
	RecordCustomFileName string `json:"record_custom_file_name,omitempty"`
}

// WebhookConfig overrides webhook delivery for a single command.
type WebhookConfig struct {
	WebhookURL             string            `json:"webhook_url,omitempty"`
	WebhookURLMethod       string            `json:"webhook_url_method,omitempty"`
	WebhookURLs            map[string]string `json:"webhook_urls,omitempty"`
	WebhookURLsMethod      string            `json:"webhook_urls_method,omitempty"`
	WebhookRetriesPolicies map[string]any    `json:"webhook_retries_policies,omitempty"`
}

// CallAssistantRequest identifies the AI assistant to attach to a call.
type CallAssistantRequest struct {
	ID string `json:"id"`
}

// CallResponse is returned by most call control command endpoints.
type CallResponse struct {
	Data struct {
		Result         string `json:"result"`
		ConversationID string `json:"conversation_id,omitempty"`
		RecordingID    string `json:"recording_id,omitempty"`
	} `json:"data"`
}

// Call is the full call object returned by Dial or GetCall.
type Call struct {
	CallControlID  string `json:"call_control_id"`
	CallLegID      string `json:"call_leg_id"`
	CallSessionID  string `json:"call_session_id"`
	IsAlive        bool   `json:"is_alive"`
	RecordType     string `json:"record_type"`
	ClientState    string `json:"client_state"`
	CallDuration   int    `json:"call_duration"`
	RecordingID    string `json:"recording_id"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	ConnectionID   string `json:"connection_id"`
}

// InterruptionSettings controls user interruption behaviour during AI speech.
type InterruptionSettings struct {
	InterruptionSensitivity string `json:"interruption_sensitivity,omitempty"`
	HoldAfterInterrupt      bool   `json:"hold_after_interrupt,omitempty"`
}

// MessageHistoryEntry is a pre-seeded conversation message for AI assistants.
type MessageHistoryEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Participant is added to an AI conversation at start.
type Participant struct {
	CallControlID string `json:"call_control_id"`
	Role          string `json:"role,omitempty"`
}

// VoiceSettings controls TTS voice provider and parameters.
type VoiceSettings struct {
	Provider   string         `json:"provider,omitempty"`
	VoiceID    string         `json:"voice_id,omitempty"`
	ExtraParam map[string]any `json:"extra_param,omitempty"`
}
