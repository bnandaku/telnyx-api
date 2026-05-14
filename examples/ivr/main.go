// IVR webhook server — answers inbound calls, plays a menu, and routes
// based on DTMF input.
//
// Usage:
//
//	TELNYX_API_KEY=... TELNYX_PUBLIC_KEY=... go run .
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	telnyx "github.com/bnandaku/telnyx-api"
)

func main() {
	apiKey := mustEnv("TELNYX_API_KEY")
	publicKey := mustEnv("TELNYX_PUBLIC_KEY")
	port := envOr("PORT", "8080")

	client := telnyx.NewClient(apiKey)

	wh, err := telnyx.NewWebhookHandler(publicKey)
	if err != nil {
		log.Fatalf("webhook handler: %v", err)
	}

	wh.On(telnyx.EventCallInitiated, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallInitiatedPayload)
		if p.Direction != "inbound" {
			return nil
		}
		log.Printf("inbound call from %s", p.From)
		_, err := client.Answer(ctx, p.CallControlID, telnyx.AnswerRequest{})
		return err
	})

	wh.On(telnyx.EventCallAnswered, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallAnsweredPayload)
		log.Printf("call answered: %s", p.CallControlID)
		_, err := client.Speak(ctx, p.CallControlID, telnyx.SpeakRequest{
			Payload:     "Welcome. Press 1 for sales, press 2 for support, or press 3 to leave a voicemail.",
			Voice:       "Telnyx.KokoroTTS.af",
			PayloadType: "text",
			ClientState: encode("menu"),
		})
		if err != nil {
			return err
		}
		return nil
	})

	wh.On(telnyx.EventCallSpeakEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallSpeakEndedPayload)
		if decode(p.ClientState) != "menu" {
			return nil
		}
		_, err := client.Gather(ctx, p.CallControlID, telnyx.GatherRequest{
			MinimumDigits:           1,
			MaximumDigits:           1,
			TimeoutMillis:           10000,
			InterDigitTimeoutMillis: 3000,
			ValidDigits:             "123",
			ClientState:             encode("gather"),
		})
		return err
	})

	wh.On(telnyx.EventCallGatherEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallGatherEndedPayload)
		if decode(p.ClientState) != "gather" {
			return nil
		}
		if p.Status != "valid" {
			_, err := client.Speak(ctx, p.CallControlID, telnyx.SpeakRequest{
				Payload:     "We did not receive your selection. Goodbye.",
				Voice:       "Telnyx.KokoroTTS.af",
				ClientState: encode("goodbye"),
			})
			return err
		}

		var msg string
		switch p.Digits {
		case "1":
			msg = "Connecting you to sales. Please hold."
		case "2":
			msg = "Connecting you to support. Please hold."
		case "3":
			msg = "Please leave your message after the tone."
		default:
			msg = "Invalid selection. Goodbye."
		}

		_, err := client.Speak(ctx, p.CallControlID, telnyx.SpeakRequest{
			Payload:     msg,
			Voice:       "Telnyx.KokoroTTS.af",
			ClientState: encode("routed:" + p.Digits),
		})
		return err
	})

	wh.On(telnyx.EventCallSpeakEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallSpeakEndedPayload)
		state := decode(p.ClientState)
		if state == "goodbye" || state == "routed:3" {
			_, err := client.Hangup(ctx, p.CallControlID, telnyx.HangupRequest{})
			return err
		}
		return nil
	})

	wh.On(telnyx.EventCallHangup, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallHangupPayload)
		log.Printf("call ended: %s cause=%s", p.CallControlID, p.HangupCause)
		return nil
	})

	http.Handle("/webhooks", wh)
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Printf("IVR server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func decode(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(b)
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
