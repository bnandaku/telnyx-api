package telnyx

import (
	"context"
	"fmt"
)

// AuditLogChange holds the before/after values of a single changed field.
type AuditLogChange struct {
	Field string `json:"field"`
	From  any    `json:"from"`
	To    any    `json:"to"`
}

// AuditLog is a single audit event recorded by Telnyx.
type AuditLog struct {
	ID                  string           `json:"id"`
	UserID              string           `json:"user_id"`
	RecordType          string           `json:"record_type"`
	ResourceID          string           `json:"resource_id"`
	AlternateResourceID string           `json:"alternate_resource_id"`
	ChangeMadeBy        string           `json:"change_made_by"` // telnyx | account_manager | account_owner | organization_member
	ChangeType          string           `json:"change_type"`
	Changes             []AuditLogChange `json:"changes"`
	OrganizationID      string           `json:"organization_id"`
	CreatedAt           string           `json:"created_at"`
}

type auditLogListResponse struct {
	Data []AuditLog     `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// ListAuditLogsParams filters the audit log list.
type ListAuditLogsParams struct {
	PageNumber    int
	PageSize      int
	CreatedBefore string // ISO 8601
	CreatedAfter  string // ISO 8601
	Sort          string // asc | desc
}

// ListAuditLogs returns audit log entries for the account.
func (c *Client) ListAuditLogs(ctx context.Context, params ListAuditLogsParams) ([]AuditLog, *PaginationMeta, error) {
	path := "/audit_logs"
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
	add("filter[created_before]", params.CreatedBefore)
	add("filter[created_after]", params.CreatedAfter)
	add("sort", params.Sort)

	var resp auditLogListResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Meta, nil
}
