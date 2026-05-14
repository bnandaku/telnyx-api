package telnyx

import (
	"context"
	"fmt"
)

// CDRStatus represents the processing state of a CDR report.
type CDRStatus int

const (
	CDRStatusPending  CDRStatus = 1
	CDRStatusComplete CDRStatus = 2
	CDRStatusFailed   CDRStatus = 3
	CDRStatusExpired  CDRStatus = 4
)

// CDRCallType filters call direction in a CDR report.
const (
	CDRCallTypeInbound  = 1
	CDRCallTypeOutbound = 2
)

// CDRRecordType filters call status in a CDR report.
const (
	CDRRecordTypeComplete   = 1
	CDRRecordTypeIncomplete = 2
	CDRRecordTypeErrors     = 3
)

// CDRFilter is an advanced filter applied to a CDR report.
type CDRFilter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// CDRReport is the report object returned by the API.
type CDRReport struct {
	RecordType       string      `json:"record_type"`
	ID               string      `json:"id"`
	ReportName       string      `json:"report_name"`
	StartTime        string      `json:"start_time"`
	EndTime          string      `json:"end_time"`
	Timezone         string      `json:"timezone"`
	Status           CDRStatus   `json:"status"`
	ReportURL        string      `json:"report_url"`
	Source           string      `json:"source"`
	CallTypes        []int       `json:"call_types"`
	RecordTypes      []int       `json:"record_types"`
	Connections      []int64     `json:"connections"`
	Filters          []CDRFilter `json:"filters"`
	ManagedAccounts  []string    `json:"managed_accounts"`
	Retry            int         `json:"retry"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        string      `json:"updated_at"`
}

type cdrReportResponse struct {
	Data CDRReport `json:"data"`
}

type cdrReportListResponse struct {
	Data []CDRReport    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// CreateCDRReportRequest specifies the parameters for a new CDR report.
type CreateCDRReportRequest struct {
	StartTime                  string      `json:"start_time"`
	EndTime                    string      `json:"end_time"`
	Timezone                   string      `json:"timezone,omitempty"`
	ReportName                 string      `json:"report_name,omitempty"`
	Source                     string      `json:"source,omitempty"` // calls | call-control | fax-api | webrtc
	CallTypes                  []int       `json:"call_types,omitempty"`
	RecordTypes                []int       `json:"record_types,omitempty"`
	Connections                []int64     `json:"connections,omitempty"`
	Fields                     []string    `json:"fields,omitempty"`
	Filters                    []CDRFilter `json:"filters,omitempty"`
	IncludeAllMetadata         bool        `json:"include_all_metadata,omitempty"`
	ManagedAccounts            []string    `json:"managed_accounts,omitempty"`
	SelectAllManagedAccounts   bool        `json:"select_all_managed_accounts,omitempty"`
}

// CreateCDRReport requests a new CDR report. The report is generated
// asynchronously — poll GetCDRReport until Status == CDRStatusComplete.
func (c *Client) CreateCDRReport(ctx context.Context, req CreateCDRReportRequest) (*CDRReport, error) {
	var resp cdrReportResponse
	if err := c.post(ctx, "/reports/cdr_requests", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetCDRReport retrieves a CDR report by ID.
func (c *Client) GetCDRReport(ctx context.Context, id string) (*CDRReport, error) {
	var resp cdrReportResponse
	if err := c.get(ctx, fmt.Sprintf("/reports/cdr_requests/%s", id), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListCDRReports returns all CDR report requests for the account.
func (c *Client) ListCDRReports(ctx context.Context) ([]CDRReport, error) {
	var resp cdrReportListResponse
	if err := c.get(ctx, "/reports/cdr_requests", &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// DeleteCDRReport removes a CDR report request.
func (c *Client) DeleteCDRReport(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/reports/cdr_requests/%s", id))
}
