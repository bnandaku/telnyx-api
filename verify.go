package telnyx

import (
	"context"
	"fmt"
	"net/url"
)

type Verification struct {
	ID              string `json:"id"`
	RecordType      string `json:"record_type"`
	Type            string `json:"type"`
	PhoneNumber     string `json:"phone_number"`
	VerifyProfileID string `json:"verify_profile_id"`
	CustomCode      string `json:"custom_code,omitempty"`
	TimeoutSecs     int    `json:"timeout_secs,omitempty"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type CreateVerificationSMSRequest struct {
	PhoneNumber     string `json:"phone_number"`
	VerifyProfileID string `json:"verify_profile_id"`
	CustomCode      string `json:"custom_code,omitempty"`
	TimeoutSecs     int    `json:"timeout_secs,omitempty"`
}

type CreateVerificationCallRequest struct {
	PhoneNumber     string `json:"phone_number"`
	VerifyProfileID string `json:"verify_profile_id"`
	CustomCode      string `json:"custom_code,omitempty"`
	TimeoutSecs     int    `json:"timeout_secs,omitempty"`
	Extension       string `json:"extension,omitempty"`
}

type CreateVerificationFlashcallRequest struct {
	PhoneNumber     string `json:"phone_number"`
	VerifyProfileID string `json:"verify_profile_id"`
	TimeoutSecs     int    `json:"timeout_secs,omitempty"`
}

type VerifyCodeByIDRequest struct {
	Code   string `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
}

type VerifyCodeByPhoneRequest struct {
	Code            string `json:"code"`
	VerifyProfileID string `json:"verify_profile_id"`
}

type VerifyCodeResponse struct {
	PhoneNumber  string `json:"phone_number"`
	ResponseCode string `json:"response_code"`
}

type VerifyProfile struct {
	ID                     string      `json:"id"`
	RecordType             string      `json:"record_type"`
	Name                   string      `json:"name"`
	Language               string      `json:"language,omitempty"`
	WebhookURL             string      `json:"webhook_url,omitempty"`
	WebhookFailoverURL     string      `json:"webhook_failover_url,omitempty"`
	DailySpendLimitEnabled bool        `json:"daily_spend_limit_enabled"`
	DailySpendLimit        float64     `json:"daily_spend_limit,omitempty"`
	SMS                    interface{} `json:"sms,omitempty"`
	Call                   interface{} `json:"call,omitempty"`
	CreatedAt              string      `json:"created_at"`
	UpdatedAt              string      `json:"updated_at"`
}

type CreateVerifyProfileRequest struct {
	Name                   string  `json:"name"`
	Language               string  `json:"language,omitempty"`
	WebhookURL             string  `json:"webhook_url,omitempty"`
	WebhookFailoverURL     string  `json:"webhook_failover_url,omitempty"`
	DailySpendLimitEnabled bool    `json:"daily_spend_limit_enabled,omitempty"`
	DailySpendLimit        float64 `json:"daily_spend_limit,omitempty"`
}

func (c *Client) SendVerificationSMS(ctx context.Context, req CreateVerificationSMSRequest) (*Verification, error) {
	var out struct {
		Data Verification `json:"data"`
	}
	if err := c.post(ctx, "/verifications/sms", req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) SendVerificationCall(ctx context.Context, req CreateVerificationCallRequest) (*Verification, error) {
	var out struct {
		Data Verification `json:"data"`
	}
	if err := c.post(ctx, "/verifications/call", req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) SendVerificationFlashcall(ctx context.Context, req CreateVerificationFlashcallRequest) (*Verification, error) {
	var out struct {
		Data Verification `json:"data"`
	}
	if err := c.post(ctx, "/verifications/flashcall", req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetVerification(ctx context.Context, id string) (*Verification, error) {
	var out struct {
		Data Verification `json:"data"`
	}
	if err := c.get(ctx, fmt.Sprintf("/verifications/%s", id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetVerificationsByPhone(ctx context.Context, phoneNumber string, pageNumber, pageSize int) ([]Verification, *PaginationMeta, error) {
	path := fmt.Sprintf("/verifications/by_phone_number/%s", url.PathEscape(phoneNumber))
	sep := "?"
	if pageNumber > 0 {
		path += fmt.Sprintf("%spage[number]=%d", sep, pageNumber)
		sep = "&"
	}
	if pageSize > 0 {
		path += fmt.Sprintf("%spage[size]=%d", sep, pageSize)
	}
	var out struct {
		Data []Verification `json:"data"`
		Meta PaginationMeta `json:"meta"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Data, &out.Meta, nil
}

func (c *Client) VerifyCodeByID(ctx context.Context, id string, req VerifyCodeByIDRequest) (*VerifyCodeResponse, error) {
	var out struct {
		Data VerifyCodeResponse `json:"data"`
	}
	if err := c.post(ctx, fmt.Sprintf("/verifications/%s/actions/verify", id), req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) VerifyCodeByPhone(ctx context.Context, phoneNumber string, req VerifyCodeByPhoneRequest) (*VerifyCodeResponse, error) {
	var out struct {
		Data VerifyCodeResponse `json:"data"`
	}
	if err := c.post(ctx, fmt.Sprintf("/verifications/by_phone_number/%s/actions/verify", url.PathEscape(phoneNumber)), req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) CreateVerifyProfile(ctx context.Context, req CreateVerifyProfileRequest) (*VerifyProfile, error) {
	var out struct {
		Data VerifyProfile `json:"data"`
	}
	if err := c.post(ctx, "/verify_profiles", req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) GetVerifyProfile(ctx context.Context, id string) (*VerifyProfile, error) {
	var out struct {
		Data VerifyProfile `json:"data"`
	}
	if err := c.get(ctx, fmt.Sprintf("/verify_profiles/%s", id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) UpdateVerifyProfile(ctx context.Context, id string, req CreateVerifyProfileRequest) (*VerifyProfile, error) {
	var out struct {
		Data VerifyProfile `json:"data"`
	}
	if err := c.patch(ctx, fmt.Sprintf("/verify_profiles/%s", id), req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) DeleteVerifyProfile(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/verify_profiles/%s", id))
}

func (c *Client) ListVerifyProfiles(ctx context.Context, pageNumber, pageSize int) ([]VerifyProfile, *PaginationMeta, error) {
	path := "/verify_profiles"
	sep := "?"
	if pageNumber > 0 {
		path += fmt.Sprintf("%spage[number]=%d", sep, pageNumber)
		sep = "&"
	}
	if pageSize > 0 {
		path += fmt.Sprintf("%spage[size]=%d", sep, pageSize)
	}
	var out struct {
		Data []VerifyProfile `json:"data"`
		Meta PaginationMeta  `json:"meta"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Data, &out.Meta, nil
}
