package telnyx

import (
	"context"
	"fmt"
)

// PhoneNumber is a Telnyx phone number object.
type PhoneNumber struct {
	RecordType              string   `json:"record_type"`
	ID                      string   `json:"id"`
	PhoneNumber             string   `json:"phone_number"`
	Status                  string   `json:"status"` // active | pending | inactive
	Tags                    []string `json:"tags"`
	ExternalPin             string   `json:"external_pin"`
	ConnectionID            string   `json:"connection_id"`
	ConnectionName          string   `json:"connection_name"`
	MessagingProfileID      string   `json:"messaging_profile_id"`
	BillingGroupID          string   `json:"billing_group_id"`
	EmergencyEnabled        bool     `json:"emergency_enabled"`
	EmergencyAddressID      string   `json:"emergency_address_id"`
	PurchasedAt             string   `json:"purchased_at"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at"`
	CountryCode             string   `json:"country_code"`
	NationalDestinationCode string   `json:"national_destination_code"`
	PhoneNumberType         string   `json:"phone_number_type"` // local | toll_free | mobile | national | shared_cost
	Features                []PhoneNumberFeature `json:"features"`
	InboundCallScreening    string   `json:"inbound_call_screening"`
	CustomerReference       string   `json:"customer_reference"`
	HDVoiceEnabled          bool     `json:"hd_voice_enabled"`
}

// PhoneNumberFeature describes a capability of a phone number.
type PhoneNumberFeature struct {
	Name string `json:"name"` // sms | mms | voice | fax
}

type phoneNumberResponse struct {
	Data PhoneNumber `json:"data"`
}

type phoneNumberListResponse struct {
	Data []PhoneNumber  `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// ListPhoneNumbersParams filters the phone number list.
type ListPhoneNumbersParams struct {
	PageNumber         int
	PageSize           int
	FilterTag          string
	FilterPhoneNumber  string
	FilterStatus       string
	FilterConnectionID string
	FilterVoiceConnectionNameContains string
	Sort               string // purchased_at | phone_number | connection_name | usage_payment_method
}

// ListPhoneNumbers returns all phone numbers on the account.
func (c *Client) ListPhoneNumbers(ctx context.Context, params ListPhoneNumbersParams) ([]PhoneNumber, *PaginationMeta, error) {
	path := "/phone_numbers"
	sep := "?"
	add := func(k, v string) {
		if v != "" {
			path += fmt.Sprintf("%s%s=%s", sep, k, v)
			sep = "&"
		}
	}
	if params.PageNumber > 0 {
		path += fmt.Sprintf("%spage[number]=%d", sep, params.PageNumber)
		sep = "&"
	}
	if params.PageSize > 0 {
		path += fmt.Sprintf("%spage[size]=%d", sep, params.PageSize)
		sep = "&"
	}
	add("filter[tag]", params.FilterTag)
	add("filter[phone_number]", params.FilterPhoneNumber)
	add("filter[status]", params.FilterStatus)
	add("filter[connection_id]", params.FilterConnectionID)
	add("filter[voice.connection_name][contains]", params.FilterVoiceConnectionNameContains)
	add("sort", params.Sort)

	var resp phoneNumberListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// GetPhoneNumber retrieves a phone number by ID.
func (c *Client) GetPhoneNumber(ctx context.Context, id string) (*PhoneNumber, error) {
	var resp phoneNumberResponse
	if err := c.get(ctx, fmt.Sprintf("/phone_numbers/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UpdatePhoneNumberRequest carries updatable fields on a phone number.
type UpdatePhoneNumberRequest struct {
	ConnectionID            string   `json:"connection_id,omitempty"`
	MessagingProfileID      string   `json:"messaging_profile_id,omitempty"`
	BillingGroupID          string   `json:"billing_group_id,omitempty"`
	ExternalPin             string   `json:"external_pin,omitempty"`
	Tags                    []string `json:"tags,omitempty"`
	CustomerReference       string   `json:"customer_reference,omitempty"`
	HDVoiceEnabled          *bool    `json:"hd_voice_enabled,omitempty"`
	InboundCallScreening    string   `json:"inbound_call_screening,omitempty"`
	EmergencyEnabled        *bool    `json:"emergency_enabled,omitempty"`
	EmergencyAddressID      string   `json:"emergency_address_id,omitempty"`
}

// UpdatePhoneNumber applies partial updates to a phone number.
func (c *Client) UpdatePhoneNumber(ctx context.Context, id string, req UpdatePhoneNumberRequest) (*PhoneNumber, error) {
	var resp phoneNumberResponse
	if err := c.patch(ctx, fmt.Sprintf("/phone_numbers/%s", id), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeletePhoneNumber releases a phone number back to Telnyx.
func (c *Client) DeletePhoneNumber(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/phone_numbers/%s", id))
}

// AvailablePhoneNumber represents a number available for purchase.
type AvailablePhoneNumber struct {
	RecordType       string   `json:"record_type"`
	PhoneNumber      string   `json:"phone_number"`
	RegionInformation []RegionInfo `json:"region_information"`
	Costs            []NumberCost `json:"costs"`
	Features         []PhoneNumberFeature `json:"features"`
	Quickship        bool     `json:"quickship"`
	Reservable       bool     `json:"reservable"`
	VanityFormat     string   `json:"vanity_format"`
}

// RegionInfo holds geographic info for an available number.
type RegionInfo struct {
	RegionType string `json:"region_type"`
	RegionName string `json:"region_name"`
}

// NumberCost holds pricing info for an available number.
type NumberCost struct {
	RecordType string `json:"record_type"`
	Currency   string `json:"currency"`
	Amount     string `json:"amount"`
	Type       string `json:"type"` // monthly | upfront | per-minute
}

type availableNumberListResponse struct {
	Data []AvailablePhoneNumber `json:"data"`
	Meta PaginationMeta         `json:"meta"`
}

// SearchAvailableNumbersParams filters available number searches.
type SearchAvailableNumbersParams struct {
	FilterCountryCode         string
	FilterNationalDestinationCode string
	FilterPhoneNumberType     string // local | toll_free | mobile | national
	FilterAreaCode            string
	FilterLocality            string
	FilterAdministrativeArea  string
	FilterFeatures            string // sms,mms,voice
	FilterLimit               int
	FilterQuickship           bool
	FilterReservable          bool
	FilterExcludeHeldNumbers  bool
}

// SearchAvailableNumbers finds phone numbers available for purchase.
func (c *Client) SearchAvailableNumbers(ctx context.Context, params SearchAvailableNumbersParams) ([]AvailablePhoneNumber, *PaginationMeta, error) {
	path := "/available_phone_numbers"
	sep := "?"
	add := func(k, v string) {
		if v != "" {
			path += fmt.Sprintf("%s%s=%s", sep, k, v)
			sep = "&"
		}
	}
	add("filter[country_code]", params.FilterCountryCode)
	add("filter[national_destination_code]", params.FilterNationalDestinationCode)
	add("filter[phone_number_type]", params.FilterPhoneNumberType)
	add("filter[features]", params.FilterFeatures)
	add("filter[locality]", params.FilterLocality)
	add("filter[administrative_area]", params.FilterAdministrativeArea)
	add("filter[phone_number][starts_with]", params.FilterAreaCode)
	if params.FilterLimit > 0 {
		path += fmt.Sprintf("%sfilter[limit]=%d", sep, params.FilterLimit)
		sep = "&"
	}
	if params.FilterQuickship {
		path += fmt.Sprintf("%sfilter[quickship]=true", sep)
		sep = "&"
	}
	if params.FilterReservable {
		path += fmt.Sprintf("%sfilter[reservable]=true", sep)
	}

	var resp availableNumberListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}

// NumberOrderRequest places an order for one or more phone numbers.
type NumberOrderRequest struct {
	PhoneNumbers       []NumberOrderItem `json:"phone_numbers"`
	ConnectionID       string            `json:"connection_id,omitempty"`
	MessagingProfileID string            `json:"messaging_profile_id,omitempty"`
	BillingGroupID     string            `json:"billing_group_id,omitempty"`
	CustomerReference  string            `json:"customer_reference,omitempty"`
}

// NumberOrderItem is a single number in a purchase order.
type NumberOrderItem struct {
	PhoneNumber        string `json:"phone_number"`
	RegulatoryGroupID  string `json:"regulatory_group_id,omitempty"`
}

// NumberOrder is the order object returned after purchasing numbers.
type NumberOrder struct {
	RecordType         string            `json:"record_type"`
	ID                 string            `json:"id"`
	PhoneNumbers       []OrderedNumber   `json:"phone_numbers"`
	PhoneNumbersCount  int               `json:"phone_numbers_count"`
	ConnectionID       string            `json:"connection_id"`
	MessagingProfileID string            `json:"messaging_profile_id"`
	Status             string            `json:"status"` // pending | success | failure
	CustomerReference  string            `json:"customer_reference"`
	CreatedAt          string            `json:"created_at"`
	UpdatedAt          string            `json:"updated_at"`
}

// OrderedNumber holds per-number status in an order.
type OrderedNumber struct {
	PhoneNumber        string `json:"phone_number"`
	RegulatoryGroupID  string `json:"regulatory_group_id"`
	Status             string `json:"status"`
	Type               string `json:"type"`
}

type numberOrderResponse struct {
	Data NumberOrder `json:"data"`
}

// OrderPhoneNumbers purchases one or more available phone numbers.
func (c *Client) OrderPhoneNumbers(ctx context.Context, req NumberOrderRequest) (*NumberOrder, error) {
	var resp numberOrderResponse
	if err := c.post(ctx, "/number_orders", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetNumberOrder retrieves a number order by ID.
func (c *Client) GetNumberOrder(ctx context.Context, id string) (*NumberOrder, error) {
	var resp numberOrderResponse
	if err := c.get(ctx, fmt.Sprintf("/number_orders/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
