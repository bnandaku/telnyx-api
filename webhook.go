package telnyx

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	webhookSignatureHeader = "Telnyx-Signature-Ed25519"
	webhookTimestampHeader = "Telnyx-Timestamp"

	// maxTimestampAge rejects webhooks older than 5 minutes to prevent replay attacks.
	maxTimestampAge = 5 * time.Minute
)

// HandlerFunc is a callback invoked when a matching webhook event arrives.
// The raw Event envelope and the decoded payload are both provided.
type HandlerFunc func(ctx context.Context, event Event, payload any) error

// WebhookHandler verifies and dispatches Telnyx webhook events.
type WebhookHandler struct {
	publicKey ed25519.PublicKey
	handlers  map[EventType][]HandlerFunc
	fallback  HandlerFunc // called for unregistered event types
}

// NewWebhookHandler creates a handler that verifies signatures with publicKey.
// publicKey must be the base64-encoded ED25519 public key from your Telnyx portal.
func NewWebhookHandler(base64PublicKey string) (*WebhookHandler, error) {
	raw, err := base64.StdEncoding.DecodeString(base64PublicKey)
	if err != nil {
		return nil, fmt.Errorf("telnyx: decode public key: %w", err)
	}
	if len(raw) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("telnyx: public key must be %d bytes, got %d", ed25519.PublicKeySize, len(raw))
	}
	return &WebhookHandler{
		publicKey: ed25519.PublicKey(raw),
		handlers:  make(map[EventType][]HandlerFunc),
	}, nil
}

// On registers a handler for a specific event type.
// Multiple handlers for the same type are called in registration order.
func (h *WebhookHandler) On(eventType EventType, fn HandlerFunc) {
	h.handlers[eventType] = append(h.handlers[eventType], fn)
}

// OnFallback registers a handler called for any unregistered event type.
func (h *WebhookHandler) OnFallback(fn HandlerFunc) {
	h.fallback = fn
}

// ServeHTTP implements http.Handler. It verifies the signature, decodes the
// event, dispatches to registered handlers, and writes 200 on success.
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	if err := h.verifySignature(r.Header, body); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "cannot decode event", http.StatusBadRequest)
		return
	}

	payload, err := decodePayload(event.Data.EventType, event.Data.Payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot decode payload: %s", err), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	fns := h.handlers[event.Data.EventType]
	if len(fns) == 0 && h.fallback != nil {
		fns = []HandlerFunc{h.fallback}
	}
	for _, fn := range fns {
		if err := fn(ctx, event, payload); err != nil {
			http.Error(w, "handler error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// verifySignature validates the ED25519 signature and timestamp freshness.
func (h *WebhookHandler) verifySignature(headers http.Header, body []byte) error {
	sigB64 := headers.Get(webhookSignatureHeader)
	if sigB64 == "" {
		return &ErrSignatureVerification{Reason: "missing " + webhookSignatureHeader + " header"}
	}
	timestamp := headers.Get(webhookTimestampHeader)
	if timestamp == "" {
		return &ErrSignatureVerification{Reason: "missing " + webhookTimestampHeader + " header"}
	}

	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return &ErrSignatureVerification{Reason: "signature is not valid base64"}
	}

	// Telnyx signs: timestamp + "|" + body
	msg := append([]byte(timestamp+"|"), body...)
	if !ed25519.Verify(h.publicKey, msg, sig) {
		return &ErrSignatureVerification{Reason: "signature mismatch"}
	}

	// Replay protection: reject events older than maxTimestampAge.
	var ts int64
	if _, err := fmt.Sscanf(timestamp, "%d", &ts); err != nil {
		return &ErrSignatureVerification{Reason: "timestamp is not a unix integer"}
	}
	eventTime := time.Unix(ts, 0)
	if age := time.Since(eventTime); age < 0 || age > maxTimestampAge {
		return &ErrSignatureVerification{Reason: fmt.Sprintf("timestamp too old or in the future: %s", age)}
	}

	return nil
}

// decodePayload unmarshals the raw payload JSON into the concrete event struct.
func decodePayload(et EventType, raw json.RawMessage) (any, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var dest any
	switch et {
	case EventCallInitiated:
		dest = new(CallInitiatedPayload)
	case EventCallAnswered:
		dest = new(CallAnsweredPayload)
	case EventCallHangup:
		dest = new(CallHangupPayload)
	case EventCallHold:
		dest = new(CallHoldPayload)
	case EventCallUnhold:
		dest = new(CallUnholdPayload)
	case EventCallBridged:
		dest = new(CallBridgedPayload)
	case EventCallPlaybackStarted:
		dest = new(CallPlaybackStartedPayload)
	case EventCallPlaybackEnded:
		dest = new(CallPlaybackEndedPayload)
	case EventCallSpeakStarted:
		dest = new(CallSpeakStartedPayload)
	case EventCallSpeakEnded:
		dest = new(CallSpeakEndedPayload)
	case EventCallDTMFReceived:
		dest = new(CallDTMFReceivedPayload)
	case EventCallGatherEnded:
		dest = new(CallGatherEndedPayload)
	case EventCallRecordingSaved:
		dest = new(CallRecordingSavedPayload)
	case EventCallRecordingError:
		dest = new(CallRecordingErrorPayload)
	case EventCallRecordingTranscriptionSaved:
		dest = new(CallRecordingTranscriptionSavedPayload)
	case EventCallMachineDetectionEnded:
		dest = new(CallMachineDetectionEndedPayload)
	case EventCallMachineGreetingEnded:
		dest = new(CallMachineGreetingEndedPayload)
	case EventCallMachinePremiumDetectionEnded:
		dest = new(CallMachinePremiumDetectionEndedPayload)
	case EventCallMachinePremiumGreetingEnded:
		dest = new(CallMachinePremiumGreetingEndedPayload)
	case EventCallForkStarted:
		dest = new(CallForkStartedPayload)
	case EventCallForkStopped:
		dest = new(CallForkStoppedPayload)
	case EventCallEnqueued:
		dest = new(CallEnqueuedPayload)
	case EventCallDequeued:
		dest = new(CallDequeuedPayload)
	case EventCallLeftQueue:
		dest = new(CallLeftQueuePayload)
	case EventCallTranscription:
		dest = new(CallTranscriptionPayload)
	case EventStreamingStarted:
		dest = new(StreamingStartedPayload)
	case EventStreamingStopped:
		dest = new(StreamingStoppedPayload)
	case EventStreamingFailed:
		dest = new(StreamingFailedPayload)
	case EventCallReferStarted:
		dest = new(CallReferStartedPayload)
	case EventCallReferCompleted:
		dest = new(CallReferCompletedPayload)
	case EventCallReferFailed:
		dest = new(CallReferFailedPayload)
	case EventCallSIPRECStarted:
		dest = new(CallSIPRECStartedPayload)
	case EventCallSIPRECStopped:
		dest = new(CallSIPRECStoppedPayload)
	case EventCallSIPRECFailed:
		dest = new(CallSIPRECFailedPayload)
	case EventCallAIGatherEnded:
		dest = new(CallAIGatherEndedPayload)
	case EventCallAIGatherMessageHistoryUpdated:
		dest = new(CallAIGatherMessageHistoryUpdatedPayload)
	case EventCallAIGatherPartialResults:
		dest = new(CallAIGatherPartialResultsPayload)
	case EventCallConversationEnded:
		dest = new(CallConversationEndedPayload)
	case EventCallConversationInsightsGenerated:
		dest = new(CallConversationInsightsGeneratedPayload)
	case EventCallDeepfakeDetectionResult:
		dest = new(DeepfakeDetectionResultPayload)
	case EventCallDeepfakeDetectionError:
		dest = new(DeepfakeDetectionErrorPayload)
	case EventConferenceCreated:
		dest = new(ConferenceCreatedPayload)
	case EventConferenceEnded:
		dest = new(ConferenceEndedPayload)
	case EventConferenceFloorChanged:
		dest = new(ConferenceFloorChangedPayload)
	case EventConferenceParticipantJoined:
		dest = new(ConferenceParticipantJoinedPayload)
	case EventConferenceParticipantLeft:
		dest = new(ConferenceParticipantLeftPayload)
	case EventConferencePlaybackStarted:
		dest = new(ConferencePlaybackStartedPayload)
	case EventConferencePlaybackEnded:
		dest = new(ConferencePlaybackEndedPayload)
	case EventConferenceSpeakStarted:
		dest = new(ConferenceSpeakStartedPayload)
	case EventConferenceSpeakEnded:
		dest = new(ConferenceSpeakEndedPayload)
	case EventConferenceRecordingSaved:
		dest = new(ConferenceRecordingSavedPayload)
	default:
		// Return raw bytes for unknown types so callers can handle them.
		return raw, nil
	}

	if err := json.Unmarshal(raw, dest); err != nil {
		return nil, fmt.Errorf("event %s: %w", et, err)
	}
	return dest, nil
}
