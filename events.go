package telnyx

import (
	"encoding/json"
	"time"
)

// EventType identifies the kind of webhook event.
type EventType string

const (
	// Call state events
	EventCallInitiated EventType = "call.initiated"
	EventCallAnswered  EventType = "call.answered"
	EventCallHangup    EventType = "call.hangup"
	EventCallHold      EventType = "call.hold"
	EventCallUnhold    EventType = "call.unhold"
	EventCallBridged   EventType = "call.bridged"

	// Audio playback events
	EventCallPlaybackStarted EventType = "call.playback.started"
	EventCallPlaybackEnded   EventType = "call.playback.ended"
	EventCallSpeakStarted    EventType = "call.speak.started"
	EventCallSpeakEnded      EventType = "call.speak.ended"

	// DTMF & gather events
	EventCallDTMFReceived EventType = "call.dtmf.received"
	EventCallGatherEnded  EventType = "call.gather.ended"

	// Recording events
	EventCallRecordingSaved              EventType = "call.recording.saved"
	EventCallRecordingError              EventType = "call.recording.error"
	EventCallRecordingTranscriptionSaved EventType = "call.recording.transcription.saved"

	// Answering machine detection events
	EventCallMachineDetectionEnded        EventType = "call.machine.detection.ended"
	EventCallMachineGreetingEnded         EventType = "call.machine.greeting.ended"
	EventCallMachinePremiumDetectionEnded EventType = "call.machine.premium.detection.ended"
	EventCallMachinePremiumGreetingEnded  EventType = "call.machine.premium.greeting.ended"

	// Media forking events
	EventCallForkStarted EventType = "call.fork.started"
	EventCallForkStopped EventType = "call.fork.stopped"

	// Queue events
	EventCallEnqueued   EventType = "call.enqueued"
	EventCallDequeued   EventType = "call.dequeued"
	EventCallLeftQueue  EventType = "call.left_queue"

	// Transcription events
	EventCallTranscription EventType = "call.transcription"

	// Streaming events
	EventStreamingStarted EventType = "streaming.started"
	EventStreamingStopped EventType = "streaming.stopped"
	EventStreamingFailed  EventType = "streaming.failed"

	// SIP REFER events
	EventCallReferStarted   EventType = "call.refer.started"
	EventCallReferCompleted EventType = "call.refer.completed"
	EventCallReferFailed    EventType = "call.refer.failed"

	// SIPREC events
	EventCallSIPRECStarted EventType = "call.siprec.started"
	EventCallSIPRECStopped EventType = "call.siprec.stopped"
	EventCallSIPRECFailed  EventType = "call.siprec.failed"

	// Conversational AI events
	EventCallAIGatherEnded                  EventType = "call.ai_gather.ended"
	EventCallAIGatherMessageHistoryUpdated  EventType = "call.ai_gather.message_history.updated"
	EventCallAIGatherPartialResults         EventType = "call.ai_gather.partial_results"
	EventCallConversationEnded              EventType = "call.conversation.ended"
	EventCallConversationInsightsGenerated  EventType = "call.conversation.insights_generated"

	// Deepfake detection events
	EventCallDeepfakeDetectionResult EventType = "call.deepfake_detection.result"
	EventCallDeepfakeDetectionError  EventType = "call.deepfake_detection.error"

	// Conference events
	EventConferenceCreated              EventType = "conference.created"
	EventConferenceEnded                EventType = "conference.ended"
	EventConferenceFloorChanged         EventType = "conference.floor_changed"
	EventConferenceParticipantJoined    EventType = "conference.participant.joined"
	EventConferenceParticipantLeft      EventType = "conference.participant.left"
	EventConferencePlaybackStarted      EventType = "conference.playback.started"
	EventConferencePlaybackEnded        EventType = "conference.playback.ended"
	EventConferenceSpeakStarted         EventType = "conference.speak.started"
	EventConferenceSpeakEnded           EventType = "conference.speak.ended"
	EventConferenceRecordingSaved       EventType = "conference.recording.saved"
)

// Event is the top-level webhook envelope.
type Event struct {
	Data EventData `json:"data"`
	Meta EventMeta `json:"meta"`
}

// EventData is the data object within the webhook envelope.
type EventData struct {
	RecordType string          `json:"record_type"`
	EventType  EventType       `json:"event_type"`
	ID         string          `json:"id"`
	OccurredAt time.Time       `json:"occurred_at"`
	Payload    json.RawMessage `json:"payload"`
}

// EventMeta holds delivery metadata for the webhook.
type EventMeta struct {
	Attempt     int    `json:"attempt"`
	DeliveredTo string `json:"delivered_to"`
}

// BasePayload contains fields present on every call webhook payload.
type BasePayload struct {
	CallControlID string `json:"call_control_id"`
	CallLegID     string `json:"call_leg_id"`
	CallSessionID string `json:"call_session_id"`
	ConnectionID  string `json:"connection_id"`
	ClientState   string `json:"client_state"`
	From          string `json:"from"`
	To            string `json:"to"`
}

// CallInitiatedPayload is emitted when a new call leg is created.
type CallInitiatedPayload struct {
	BasePayload
	Direction       string `json:"direction"`   // inbound | outbound
	State           string `json:"state"`
	StartTime       string `json:"start_time"`
	FromDisplayName string `json:"from_display_name"`
	ToDisplayName   string `json:"to_display_name"`
	SIPHeaders      []CustomSipHeader `json:"sip_headers"`
}

// CallAnsweredPayload is emitted when the call is answered.
type CallAnsweredPayload struct {
	BasePayload
	State           string `json:"state"`
	StartTime       string `json:"start_time"`
	FromDisplayName string `json:"from_display_name"`
	ToDisplayName   string `json:"to_display_name"`
}

// CallHangupPayload is emitted when a call leg ends.
type CallHangupPayload struct {
	BasePayload
	HangupCause     string `json:"hangup_cause"`
	HangupSource    string `json:"hangup_source"`
	SIPHangupCause  string `json:"sip_hangup_cause"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
}

// CallHoldPayload is emitted when a call is placed on hold.
type CallHoldPayload struct {
	BasePayload
}

// CallUnholdPayload is emitted when a call is taken off hold.
type CallUnholdPayload struct {
	BasePayload
}

// CallBridgedPayload is emitted when two call legs are connected.
type CallBridgedPayload struct {
	BasePayload
	BridgedTo string `json:"bridged_to"`
	State     string `json:"state"`
}

// CallPlaybackStartedPayload is emitted when audio playback begins.
type CallPlaybackStartedPayload struct {
	BasePayload
	MediaURL    string `json:"media_url"`
	Overlay     bool   `json:"overlay"`
	TargetLegs  string `json:"target_legs"`
}

// CallPlaybackEndedPayload is emitted when audio playback finishes.
type CallPlaybackEndedPayload struct {
	BasePayload
	MediaURL   string `json:"media_url"`
	Status     string `json:"status"`
	TargetLegs string `json:"target_legs"`
}

// CallSpeakStartedPayload is emitted when TTS starts.
type CallSpeakStartedPayload struct {
	BasePayload
	TargetLegs string `json:"target_legs"`
}

// CallSpeakEndedPayload is emitted when TTS finishes.
type CallSpeakEndedPayload struct {
	BasePayload
	Status     string `json:"status"`
	TargetLegs string `json:"target_legs"`
}

// CallDTMFReceivedPayload is emitted when a DTMF digit is pressed.
type CallDTMFReceivedPayload struct {
	BasePayload
	Digit string `json:"digit"`
}

// CallGatherEndedPayload is emitted when a gather operation completes.
type CallGatherEndedPayload struct {
	BasePayload
	Digits   string `json:"digits"`
	Status   string `json:"status"` // valid | invalid | call_hangup | cancelled | timeout | term_digit
	GatherID string `json:"gather_id"`
}

// CallRecordingSavedPayload is emitted when a recording is saved.
type CallRecordingSavedPayload struct {
	BasePayload
	RecordingID         string              `json:"recording_id"`
	RecordingStartedAt  string              `json:"recording_started_at"`
	RecordingEndedAt    string              `json:"recording_ended_at"`
	RecordingURL        RecordingURLs       `json:"recording_urls"`
	Channels            string              `json:"channels"`
	DurationMillis      int                 `json:"duration_millis"`
	Transcription       *RecordingTranscript `json:"transcription,omitempty"`
}

// RecordingURLs holds the download URLs for a recording.
type RecordingURLs struct {
	MP3 string `json:"mp3"`
	WAV string `json:"wav"`
}

// RecordingTranscript holds transcription data for a recording.
type RecordingTranscript struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

// CallRecordingErrorPayload is emitted on a recording error.
type CallRecordingErrorPayload struct {
	BasePayload
	RecordingID string `json:"recording_id"`
	Reason      string `json:"reason"`
}

// CallRecordingTranscriptionSavedPayload is emitted when a transcription is saved.
type CallRecordingTranscriptionSavedPayload struct {
	BasePayload
	RecordingID        string `json:"recording_id"`
	TranscriptionURL   string `json:"transcription_url"`
	TranscriptionText  string `json:"transcription_text"`
}

// AMDResult holds answering machine detection outcome.
type AMDResult struct {
	Result   string `json:"result"`  // human | machine | not_sure | silence
	Beep     bool   `json:"beep"`
	Silence  bool   `json:"silence"`
}

// CallMachineDetectionEndedPayload is emitted when AMD completes.
type CallMachineDetectionEndedPayload struct {
	BasePayload
	AMDResult string `json:"result"` // human | machine | not_sure | silence
}

// CallMachineGreetingEndedPayload is emitted when a machine beep is detected.
type CallMachineGreetingEndedPayload struct {
	BasePayload
	AMDResult string `json:"result"`
}

// CallMachinePremiumDetectionEndedPayload is emitted for premium AMD.
type CallMachinePremiumDetectionEndedPayload struct {
	BasePayload
	Result        string `json:"result"`
	ReasonCode    int    `json:"reason_code"`
}

// CallMachinePremiumGreetingEndedPayload is emitted for premium AMD beep detection.
type CallMachinePremiumGreetingEndedPayload struct {
	BasePayload
	Result string `json:"result"`
}

// CallForkStartedPayload is emitted when media forking begins.
type CallForkStartedPayload struct {
	BasePayload
}

// CallForkStoppedPayload is emitted when media forking ends.
type CallForkStoppedPayload struct {
	BasePayload
}

// CallEnqueuedPayload is emitted when a call enters a queue.
type CallEnqueuedPayload struct {
	BasePayload
	QueueName    string `json:"queue_name"`
	QueueSize    int    `json:"queue_size"`
	EnqueuedAt   string `json:"enqueued_at"`
}

// CallDequeuedPayload is emitted when a call leaves a queue.
type CallDequeuedPayload struct {
	BasePayload
	QueueName       string `json:"queue_name"`
	DequeuedAt      string `json:"dequeued_at"`
	LeaveQueueCause string `json:"leave_queue_cause"`
}

// CallLeftQueuePayload is an alias event for call.left_queue.
type CallLeftQueuePayload struct {
	BasePayload
	QueueName       string `json:"queue_name"`
	LeaveQueueCause string `json:"leave_queue_cause"`
}

// TranscriptionData holds a single transcription segment.
type TranscriptionData struct {
	Transcript     string  `json:"transcript"`
	Confidence     float64 `json:"confidence"`
	IsFinal        bool    `json:"is_final"`
	Language       string  `json:"language"`
	TranscriptFrom string  `json:"transcript_from"` // customer | agent
}

// CallTranscriptionPayload is emitted for real-time transcription segments.
type CallTranscriptionPayload struct {
	BasePayload
	TranscriptionData TranscriptionData `json:"transcription_data"`
}

// StreamingStartedPayload is emitted when WebSocket streaming begins.
type StreamingStartedPayload struct {
	BasePayload
	StreamURL  string `json:"stream_url"`
	StreamID   string `json:"stream_id"`
	StreamTrack string `json:"stream_track"`
}

// StreamingStoppedPayload is emitted when WebSocket streaming ends.
type StreamingStoppedPayload struct {
	BasePayload
	StreamID string `json:"stream_id"`
}

// StreamingFailedPayload is emitted when WebSocket streaming fails.
type StreamingFailedPayload struct {
	BasePayload
	StreamID string `json:"stream_id"`
	Reason   string `json:"reason"`
}

// CallReferStartedPayload is emitted when a SIP REFER begins.
type CallReferStartedPayload struct {
	BasePayload
	SIPReferEventType string `json:"sip_refer_event_type"`
}

// CallReferCompletedPayload is emitted when a SIP REFER succeeds.
type CallReferCompletedPayload struct {
	BasePayload
	SIPReferEventType string `json:"sip_refer_event_type"`
}

// CallReferFailedPayload is emitted when a SIP REFER fails.
type CallReferFailedPayload struct {
	BasePayload
	SIPReferEventType string `json:"sip_refer_event_type"`
	FailureReason     string `json:"failure_reason"`
}

// CallSIPRECStartedPayload is emitted when SIPREC recording begins.
type CallSIPRECStartedPayload struct {
	BasePayload
	SIPRECID string `json:"siprec_id"`
}

// CallSIPRECStoppedPayload is emitted when SIPREC recording ends.
type CallSIPRECStoppedPayload struct {
	BasePayload
	SIPRECID string `json:"siprec_id"`
}

// CallSIPRECFailedPayload is emitted when a SIPREC session fails.
type CallSIPRECFailedPayload struct {
	BasePayload
	SIPRECID string `json:"siprec_id"`
	Reason   string `json:"reason"`
}

// AIGatherMessage is a single message in an AI conversation.
type AIGatherMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CallAIGatherEndedPayload is emitted when an AI gather session ends.
type CallAIGatherEndedPayload struct {
	BasePayload
	ConversationID string            `json:"conversation_id"`
	Status         string            `json:"status"`
	MessageHistory []AIGatherMessage `json:"message_history"`
}

// CallAIGatherMessageHistoryUpdatedPayload is emitted on conversation history changes.
type CallAIGatherMessageHistoryUpdatedPayload struct {
	BasePayload
	ConversationID string            `json:"conversation_id"`
	MessageHistory []AIGatherMessage `json:"message_history"`
}

// CallAIGatherPartialResultsPayload is emitted with partial AI transcription.
type CallAIGatherPartialResultsPayload struct {
	BasePayload
	ConversationID string `json:"conversation_id"`
	Transcript     string `json:"transcript"`
	Role           string `json:"role"`
}

// CallConversationEndedPayload is emitted when an AI conversation ends.
type CallConversationEndedPayload struct {
	BasePayload
	ConversationID string            `json:"conversation_id"`
	Status         string            `json:"status"`
	EndReason      string            `json:"end_reason"`
	MessageHistory []AIGatherMessage `json:"message_history"`
}

// CallConversationInsightsGeneratedPayload is emitted when AI insights are ready.
type CallConversationInsightsGeneratedPayload struct {
	BasePayload
	ConversationID string         `json:"conversation_id"`
	Insights       map[string]any `json:"insights"`
}

// DeepfakeDetectionResultPayload is emitted with a deepfake detection verdict.
type DeepfakeDetectionResultPayload struct {
	BasePayload
	Result     string  `json:"result"`     // genuine | deepfake | unknown
	Confidence float64 `json:"confidence"`
}

// DeepfakeDetectionErrorPayload is emitted when deepfake detection fails.
type DeepfakeDetectionErrorPayload struct {
	BasePayload
	Reason string `json:"reason"`
}

// ConferenceBasePayload contains fields common to all conference webhooks.
type ConferenceBasePayload struct {
	ConferenceID   string `json:"conference_id"`
	ConnectionID   string `json:"connection_id"`
	OccurredAt     string `json:"occurred_at"`
	ClientState    string `json:"client_state"`
}

// ConferenceCreatedPayload is emitted when a conference is created.
type ConferenceCreatedPayload struct {
	ConferenceBasePayload
	Name string `json:"name"`
}

// ConferenceEndedPayload is emitted when a conference ends.
type ConferenceEndedPayload struct {
	ConferenceBasePayload
	Name       string `json:"name"`
	EndReason  string `json:"end_reason"`
}

// ConferenceFloorChangedPayload is emitted when the conference floor holder changes.
type ConferenceFloorChangedPayload struct {
	ConferenceBasePayload
	Floor    *ConferenceParticipant `json:"floor"`
	PrevFloor *ConferenceParticipant `json:"prev_floor"`
}

// ConferenceParticipant holds identity info for a conference participant.
type ConferenceParticipant struct {
	CallControlID string `json:"call_control_id"`
	CallLegID     string `json:"call_leg_id"`
}

// ConferenceParticipantJoinedPayload is emitted when a participant joins.
type ConferenceParticipantJoinedPayload struct {
	ConferenceBasePayload
	CallControlID  string `json:"call_control_id"`
	CallLegID      string `json:"call_leg_id"`
	CallSessionID  string `json:"call_session_id"`
	From           string `json:"from"`
	To             string `json:"to"`
}

// ConferenceParticipantLeftPayload is emitted when a participant leaves.
type ConferenceParticipantLeftPayload struct {
	ConferenceBasePayload
	CallControlID string `json:"call_control_id"`
	CallLegID     string `json:"call_leg_id"`
}

// ConferencePlaybackStartedPayload is emitted when conference audio playback begins.
type ConferencePlaybackStartedPayload struct {
	ConferenceBasePayload
	MediaURL string `json:"media_url"`
}

// ConferencePlaybackEndedPayload is emitted when conference audio playback ends.
type ConferencePlaybackEndedPayload struct {
	ConferenceBasePayload
	MediaURL string `json:"media_url"`
	Status   string `json:"status"`
}

// ConferenceSpeakStartedPayload is emitted when conference TTS starts.
type ConferenceSpeakStartedPayload struct {
	ConferenceBasePayload
}

// ConferenceSpeakEndedPayload is emitted when conference TTS ends.
type ConferenceSpeakEndedPayload struct {
	ConferenceBasePayload
	Status string `json:"status"`
}

// ConferenceRecordingSavedPayload is emitted when a conference recording is saved.
type ConferenceRecordingSavedPayload struct {
	ConferenceBasePayload
	RecordingID        string       `json:"recording_id"`
	RecordingStartedAt string       `json:"recording_started_at"`
	RecordingEndedAt   string       `json:"recording_ended_at"`
	RecordingURL       RecordingURLs `json:"recording_urls"`
	Channels           string       `json:"channels"`
	DurationMillis     int          `json:"duration_millis"`
}
