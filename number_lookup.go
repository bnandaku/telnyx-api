package telnyx

import (
	"context"
	"fmt"
	"net/url"
)

// NumberLookupResult is the full response from a phone number lookup.
type NumberLookupResult struct {
	RecordType     string              `json:"record_type"`
	CountryCode    string              `json:"country_code"`
	NationalFormat string              `json:"national_format"`
	PhoneNumber    string              `json:"phone_number"`
	Carrier        *CarrierInfo        `json:"carrier,omitempty"`
	CallerName     *CallerNameInfo     `json:"caller_name,omitempty"`
	Portability    *PortabilityInfo    `json:"portability,omitempty"`
}

// CarrierInfo holds carrier/network data for a number.
type CarrierInfo struct {
	MobileCountryCode string `json:"mobile_country_code"`
	MobileNetworkCode string `json:"mobile_network_code"`
	Name              string `json:"name"`
	Type              string `json:"type"` // mobile | voip | landline | toll free | unknown
	NormalizedCarrier string `json:"normalized_carrier"`
	ErrorCode         string `json:"error_code"`
}

// CallerNameInfo holds CNAM data for a number.
type CallerNameInfo struct {
	CallerName string `json:"caller_name"`
	ErrorCode  string `json:"error_code"`
}

// PortabilityInfo holds LERG / LNP data for a number.
type PortabilityInfo struct {
	LRN          string `json:"lrn"`
	PortedStatus string `json:"ported_status"` // Y | N
	PortedDate   string `json:"ported_date"`
	OCN          string `json:"ocn"`
	LineType     string `json:"line_type"`
	SPID         string `json:"spid"`
	AltSPID      string `json:"altspid"`
	City         string `json:"city"`
	State        string `json:"state"`
}

type numberLookupResponse struct {
	Data NumberLookupResult `json:"data"`
}

// LookupNumber performs a full phone number lookup including carrier, caller
// name (CNAM), and portability (LERG) data.
//
// Carrier and caller-name lookups are made concurrently and merged into a
// single result. Portability data is returned by the carrier lookup.
func (c *Client) LookupNumber(ctx context.Context, phoneNumber string) (*NumberLookupResult, error) {
	type result struct {
		data *NumberLookupResult
		err  error
	}

	carrierCh := make(chan result, 1)
	cnamCh := make(chan result, 1)

	go func() {
		var resp numberLookupResponse
		path := fmt.Sprintf("/number_lookup/%s?type=carrier", url.PathEscape(phoneNumber))
		err := c.get(ctx, path, &resp)
		carrierCh <- result{&resp.Data, err}
	}()

	go func() {
		var resp numberLookupResponse
		path := fmt.Sprintf("/number_lookup/%s?type=caller-name", url.PathEscape(phoneNumber))
		err := c.get(ctx, path, &resp)
		cnamCh <- result{&resp.Data, err}
	}()

	carrierRes := <-carrierCh
	cnamRes := <-cnamCh

	if carrierRes.err != nil {
		return nil, fmt.Errorf("carrier lookup: %w", carrierRes.err)
	}

	merged := carrierRes.data
	if cnamRes.err == nil && cnamRes.data != nil {
		merged.CallerName = cnamRes.data.CallerName
	}

	return merged, nil
}
