// Number lookup — looks up a phone number and prints carrier, CNAM,
// and LERG/portability data.
//
// Usage:
//
//	TELNYX_API_KEY=... go run . +15550001111
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	telnyx "github.com/bnandaku/telnyx-api"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: number_lookup <phone_number>")
	}
	phoneNumber := os.Args[1]
	apiKey := mustEnv("TELNYX_API_KEY")

	client := telnyx.NewClient(apiKey)

	result, err := client.LookupNumber(context.Background(), phoneNumber)
	if err != nil {
		log.Fatalf("lookup: %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	fmt.Printf("Phone number : %s\n", result.PhoneNumber)
	fmt.Printf("Country      : %s\n", result.CountryCode)
	fmt.Printf("National fmt : %s\n", result.NationalFormat)

	if result.Carrier != nil {
		fmt.Printf("Carrier      : %s (%s)\n", result.Carrier.Name, result.Carrier.Type)
		fmt.Printf("Network      : %s\n", result.Carrier.NormalizedCarrier)
	}

	if result.CallerName != nil && result.CallerName.CallerName != "" {
		fmt.Printf("Caller name  : %s\n", result.CallerName.CallerName)
	}

	if result.Portability != nil {
		p := result.Portability
		fmt.Printf("LRN          : %s\n", p.LRN)
		fmt.Printf("Ported       : %s\n", p.PortedStatus)
		fmt.Printf("OCN          : %s\n", p.OCN)
		fmt.Printf("Location     : %s, %s\n", p.City, p.State)
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}
	return v
}
