package telnyx

import (
	"context"
	"fmt"
)

// BillingGroup organises numbers and usage under a named cost centre.
type BillingGroup struct {
	RecordType     string `json:"record_type"`
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	DeletedAt      string `json:"deleted_at"`
}

type billingGroupResponse struct {
	Data BillingGroup `json:"data"`
}

type billingGroupListResponse struct {
	Data []BillingGroup `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// CreateBillingGroup creates a new billing group.
func (c *Client) CreateBillingGroup(ctx context.Context, name string) (*BillingGroup, error) {
	var resp billingGroupResponse
	body := map[string]string{"name": name}
	if err := c.post(ctx, "/billing_groups", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetBillingGroup retrieves a billing group by ID.
func (c *Client) GetBillingGroup(ctx context.Context, id string) (*BillingGroup, error) {
	var resp billingGroupResponse
	if err := c.get(ctx, fmt.Sprintf("/billing_groups/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UpdateBillingGroup renames a billing group.
func (c *Client) UpdateBillingGroup(ctx context.Context, id, name string) (*BillingGroup, error) {
	var resp billingGroupResponse
	body := map[string]string{"name": name}
	if err := c.patch(ctx, fmt.Sprintf("/billing_groups/%s", id), body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeleteBillingGroup removes a billing group.
func (c *Client) DeleteBillingGroup(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/billing_groups/%s", id))
}

// ListBillingGroups returns all billing groups.
func (c *Client) ListBillingGroups(ctx context.Context, pageNumber, pageSize int) ([]BillingGroup, *PaginationMeta, error) {
	path := fmt.Sprintf("/billing_groups?page[number]=%d&page[size]=%d", pageNumber, pageSize)
	var resp billingGroupListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}
