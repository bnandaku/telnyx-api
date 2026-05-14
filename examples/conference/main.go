// Conference bridge — creates a conference when the first participant dials in
// and joins each subsequent caller.
//
// Usage:
//
//	TELNYX_API_KEY=... TELNYX_PUBLIC_KEY=... go run .
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	telnyx "github.com/bnandaku/telnyx-api"
)

var (
	mu           sync.Mutex
	conferenceID string
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
		_, err := client.Answer(ctx, p.CallControlID, telnyx.AnswerRequest{})
		return err
	})

	wh.On(telnyx.EventCallAnswered, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallAnsweredPayload)
		log.Printf("caller answered: %s", p.From)

		mu.Lock()
		confID := conferenceID
		mu.Unlock()

		if confID == "" {
			// First caller — create the conference.
			conf, err := client.CreateConference(ctx, telnyx.CreateConferenceRequest{
				CallControlID:   p.CallControlID,
				Name:            "Bridge",
				BeepEnabled:     "on_enter",
				MaxParticipants: 10,
			})
			if err != nil {
				return err
			}
			mu.Lock()
			conferenceID = conf.ID
			mu.Unlock()
			log.Printf("conference created: %s", conf.ID)
		} else {
			// Subsequent callers — join existing conference.
			_, err := client.JoinConference(ctx, confID, telnyx.JoinConferenceRequest{
				CallControlID: p.CallControlID,
			})
			if err != nil {
				return err
			}
			log.Printf("joined conference: %s", confID)
		}
		return nil
	})

	wh.On(telnyx.EventConferenceParticipantLeft, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.ConferenceParticipantLeftPayload)
		log.Printf("participant left conference: %s", p.CallControlID)
		return nil
	})

	wh.On(telnyx.EventConferenceEnded, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.ConferenceEndedPayload)
		log.Printf("conference ended: %s reason=%s", p.ConferenceID, p.EndReason)
		mu.Lock()
		conferenceID = ""
		mu.Unlock()
		return nil
	})

	wh.On(telnyx.EventCallHangup, func(ctx context.Context, event telnyx.Event, payload any) error {
		p := payload.(*telnyx.CallHangupPayload)
		log.Printf("call ended: %s", p.HangupCause)
		return nil
	})

	http.Handle("/webhooks", wh)
	log.Printf("conference server listening on :%s", port)
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
