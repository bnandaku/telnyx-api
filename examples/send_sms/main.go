// Sends an SMS and prints the delivery status from the webhook.
//
// Usage:
//
//	TELNYX_API_KEY=... TELNYX_PUBLIC_KEY=... \
//	  FROM=+1555... TO=+1555... go run .
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
	from := mustEnv("FROM")
	to := mustEnv("TO")
	port := envOr("PORT", "8080")

	client := telnyx.NewClient(apiKey)

	// Send SMS.
	msg, err := client.SendMessage(context.Background(), telnyx.SendMessageRequest{
		From: from,
		To:   to,
		Text: "Hello from Telnyx! Reply STOP to unsubscribe.",
	})
	if err != nil {
		log.Fatalf("send message: %v", err)
	}
	log.Printf("sent: id=%s status=%s", msg.ID, msg.To[0].Status)

	// Receive delivery webhooks.
	wh, err := telnyx.NewWebhookHandler(publicKey)
	if err != nil {
		log.Fatalf("webhook handler: %v", err)
	}

	wh.OnFallback(func(ctx context.Context, event telnyx.Event, payload any) error {
		log.Printf("event: %s", event.Data.EventType)
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
