package telnyx

import (
	"context"
	"fmt"
)

type WebhookDelivery struct {
	ID         string            `json:"id"`
	UserID     string            `json:"user_id"`
	RecordType string            `json:"record_type"`
	Status     string            `json:"status"`
	StartedAt  string            `json:"started_at"`
	FinishedAt string            `json:"finished_at"`
	Webhook    WebhookPayload    `json:"webhook"`
	Attempts   []DeliveryAttempt `json:"attempts"`
}

type WebhookPayload struct {
	RecordType string `json:"record_type"`
	EventType  string `json:"event_type"`
	ID         string `json:"id"`
	OccurredAt string `json:"occurred_at"`
}

type DeliveryAttempt struct {
	Status     string               `json:"status"`
	StartedAt  string               `json:"started_at"`
	FinishedAt string               `json:"finished_at"`
	Request    DeliveryHTTPRequest  `json:"http.request"`
	Response   DeliveryHTTPResponse `json:"http.response"`
	Errors     []int                `json:"errors,omitempty"`
}

type DeliveryHTTPRequest struct {
	URL     string      `json:"url"`
	Headers [][2]string `json:"headers"`
}

type DeliveryHTTPResponse struct {
	Status  int         `json:"status"`
	Headers [][2]string `json:"headers"`
	Body    string      `json:"body"`
}

type ListWebhookDeliveriesParams struct {
	FilterStatus        string
	FilterEventType     string
	FilterWebhook       string
	FilterStartedAtGTE  string
	FilterStartedAtLTE  string
	FilterFinishedAtGTE string
	FilterFinishedAtLTE string
	PageNumber          int
	PageSize            int
}

func (c *Client) GetWebhookDelivery(ctx context.Context, id string) (*WebhookDelivery, error) {
	var out struct {
		Data WebhookDelivery `json:"data"`
	}
	if err := c.get(ctx, fmt.Sprintf("/webhook_deliveries/%s", id), &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

func (c *Client) ListWebhookDeliveries(ctx context.Context, params ListWebhookDeliveriesParams) ([]WebhookDelivery, *PaginationMeta, error) {
	path := "/webhook_deliveries"
	sep := "?"
	addStr := func(k, v string) {
		if v != "" {
			path += fmt.Sprintf("%s%s=%s", sep, k, v)
			sep = "&"
		}
	}
	addInt := func(k string, v int) {
		if v > 0 {
			path += fmt.Sprintf("%s%s=%d", sep, k, v)
			sep = "&"
		}
	}
	addStr("filter[status.eq]", params.FilterStatus)
	addStr("filter[event_type]", params.FilterEventType)
	addStr("filter[webhook.contains]", params.FilterWebhook)
	addStr("filter[started_at.gte]", params.FilterStartedAtGTE)
	addStr("filter[started_at.lte]", params.FilterStartedAtLTE)
	addStr("filter[finished_at.gte]", params.FilterFinishedAtGTE)
	addStr("filter[finished_at.lte]", params.FilterFinishedAtLTE)
	addInt("page[number]", params.PageNumber)
	addInt("page[size]", params.PageSize)

	var out struct {
		Data []WebhookDelivery `json:"data"`
		Meta PaginationMeta    `json:"meta"`
	}
	if err := c.get(ctx, path, &out); err != nil {
		return nil, nil, err
	}
	return out.Data, &out.Meta, nil
}
