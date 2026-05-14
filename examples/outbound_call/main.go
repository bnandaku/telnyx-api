// Outbound call with answering machine detection.
// Dials a number, waits for the call.answered or call.machine.* events,
// and speaks a different message depending on whether a human or machine answered.
//
// Usage:
//
//	TELNYX_API_KEY=... TELNYX_PUBLIC_KEY=... \
//	  TELNYX_CONNECTION_ID=... FROM=+1555... TO=+1555... go run .
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	telnyx "github.com/bnandaku/telnyx-api"
)

func main() {
	apiKey := mustEnv("TELNYX_API_KEY")
	publicKey := mustEnv("TELNYX_PUBLIC_KEY")
	connectionID := mustEnv("TELNYX_CONNECTION_ID")
	from := mustEnv("FROM")
	to := mustEnv("TO")
	port := envOr("PORT", "8080")
	webhookURL := mustEnv("WEBHOOK_URL")

	client := telnyx.NewClient(apiKey)

	// Dial immediately on startup.
	call, err := client.Dial(context.Background(), telnyx.DialRequest{
		ConnectionID:              connectionID,
		From:                      from,
		To:                        to,
		AnsweringMachineDetection: "premium",
		WebhookURL:                webhookURL + "/webhooks",
	})
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	log.Printf("dialed: call_control_id=%s", call.CallControlID)

	wh, err := telnyx.NewWebhookHandler(publicKey)
	if err != nil {
		log.Fatalf("webhook handler: %v", err)
	}

	// Human answered.
	wh.On(telnyx.EventCallAnswered, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallAnsweredPayload)
		log.Printf("human answered: %s", p.CallControlID)
		_, err := client.Speak(ctx, p.CallControlID, telnyx.SpeakRequest{
			Payload: "Hello! This is a reminder about your appointment tomorrow at 10am. Press 1 to confirm or 2 to reschedule.",
			Voice:   "Telnyx.KokoroTTS.af",
		})
		return err
	})

	// Machine detected — leave a voicemail after the beep.
	wh.On(telnyx.EventCallMachinePremiumGreetingEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallMachinePremiumGreetingEndedPayload)
		log.Printf("machine beep detected on: %s", p.CallControlID)
		_, err := client.Speak(ctx, p.CallControlID, telnyx.SpeakRequest{
			Payload: "Hi, this is a reminder about your appointment tomorrow at 10am. Please call us back to confirm. Thank you!",
			Voice:   "Telnyx.KokoroTTS.af",
		})
		return err
	})

	wh.On(telnyx.EventCallSpeakEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallSpeakEndedPayload)
		_, err := client.Hangup(ctx, p.CallControlID, telnyx.HangupRequest{})
		return err
	})

	wh.On(telnyx.EventCallHangup, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallHangupPayload)
		log.Printf("call ended: cause=%s", p.HangupCause)
		return nil
	})

	http.Handle("/webhooks", wh)
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}
	return v
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
