package telnyx

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func newTestKeyPair(t *testing.T) (ed25519.PublicKey, ed25519.PrivateKey) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return pub, priv
}

func signWebhook(t *testing.T, priv ed25519.PrivateKey, body []byte, ts time.Time) (sig, timestamp string) {
	t.Helper()
	timestamp = strconv.FormatInt(ts.Unix(), 10)
	msg := append([]byte(timestamp+"|"), body...)
	raw := ed25519.Sign(priv, msg)
	sig = base64.StdEncoding.EncodeToString(raw)
	return
}

func buildRequest(t *testing.T, body []byte, sig, timestamp string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/webhooks", bytes.NewReader(body))
	req.Header.Set(webhookSignatureHeader, sig)
	req.Header.Set(webhookTimestampHeader, timestamp)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestWebhookSignatureValid(t *testing.T) {
	pub, priv := newTestKeyPair(t)

	event := map[string]any{
		"data": map[string]any{
			"record_type": "event",
			"event_type":  "call.initiated",
			"id":          "abc-123",
			"occurred_at": time.Now().Format(time.RFC3339),
			"payload": map[string]any{
				"call_control_id":  "ctrl-1",
				"call_leg_id":      "leg-1",
				"call_session_id":  "sess-1",
				"connection_id":    "conn-1",
				"client_state":     "",
				"from":             "+15550001111",
				"to":               "+15550002222",
				"direction":        "inbound",
				"state":            "parked",
			},
		},
		"meta": map[string]any{
			"attempt":      1,
			"delivered_to": "https://example.com/webhooks",
		},
	}

	body, _ := json.Marshal(event)
	sig, ts := signWebhook(t, priv, body, time.Now())

	h, err := NewWebhookHandler(base64.StdEncoding.EncodeToString(pub))
	if err != nil {
		t.Fatalf("NewWebhookHandler: %v", err)
	}

	var received *CallInitiatedPayload
	h.On(EventCallInitiated, func(_ context.Context, _ Event, payload any) error {
		p, ok := payload.(*CallInitiatedPayload)
		if !ok {
			return fmt.Errorf("wrong payload type: %T", payload)
		}
		received = p
		return nil
	})

	w := httptest.NewRecorder()
	h.ServeHTTP(w, buildRequest(t, body, sig, ts))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if received == nil {
		t.Fatal("handler was not called")
	}
	if received.CallControlID != "ctrl-1" {
		t.Errorf("call_control_id = %q, want ctrl-1", received.CallControlID)
	}
	if received.From != "+15550001111" {
		t.Errorf("from = %q, want +15550001111", received.From)
	}
}

func TestWebhookSignatureInvalid(t *testing.T) {
	pub, _ := newTestKeyPair(t)
	_, otherPriv := newTestKeyPair(t)

	body := []byte(`{"data":{"event_type":"call.hangup","record_type":"event","id":"x","occurred_at":"2024-01-01T00:00:00Z","payload":{}},"meta":{}}`)
	sig, ts := signWebhook(t, otherPriv, body, time.Now())

	h, _ := NewWebhookHandler(base64.StdEncoding.EncodeToString(pub))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, buildRequest(t, body, sig, ts))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestWebhookReplayRejected(t *testing.T) {
	pub, priv := newTestKeyPair(t)

	body := []byte(`{"data":{"event_type":"call.hangup","record_type":"event","id":"x","occurred_at":"2024-01-01T00:00:00Z","payload":{}},"meta":{}}`)
	oldTime := time.Now().Add(-10 * time.Minute)
	sig, ts := signWebhook(t, priv, body, oldTime)

	h, _ := NewWebhookHandler(base64.StdEncoding.EncodeToString(pub))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, buildRequest(t, body, sig, ts))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for stale timestamp, got %d", w.Code)
	}
}

func TestWebhookFallback(t *testing.T) {
	pub, priv := newTestKeyPair(t)

	body, _ := json.Marshal(map[string]any{
		"data": map[string]any{
			"record_type": "event",
			"event_type":  "some.future.event",
			"id":          "y",
			"occurred_at": time.Now().Format(time.RFC3339),
			"payload":     map[string]any{},
		},
		"meta": map[string]any{"attempt": 1},
	})
	sig, ts := signWebhook(t, priv, body, time.Now())

	h, _ := NewWebhookHandler(base64.StdEncoding.EncodeToString(pub))
	called := false
	h.OnFallback(func(_ context.Context, _ Event, _ any) error {
		called = true
		return nil
	})

	w := httptest.NewRecorder()
	h.ServeHTTP(w, buildRequest(t, body, sig, ts))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !called {
		t.Fatal("fallback was not called")
	}
}

func TestNewWebhookHandlerInvalidKey(t *testing.T) {
	_, err := NewWebhookHandler("not-valid-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64 key")
	}

	// Wrong size key
	_, err = NewWebhookHandler(base64.StdEncoding.EncodeToString([]byte("tooshort")))
	if err == nil {
		t.Fatal("expected error for wrong key size")
	}
}
