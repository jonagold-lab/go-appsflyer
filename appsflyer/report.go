package appsflyer

import (
	"context"
)

// ReportService struct to hold individual service
type ReportService service

// Options parameter you can add to the url
type Options struct {
	From             string `url:"from,omitempty"`              // "2019-05-09"
	To               string `url:"to,omitempty"`                // "2019-05-09"
	Category         string `url:"category,omitempty"`          // standard || facebook || twitter
	MediaSource      bool   `url:"media_source,omitempty"`      // googleadwords_int || facebook
	EventName        string `url:"event_name,omitempty"`        // af_purchase,ftd
	ReTargeting      bool   `url:"reattr,omitempty"`            //
	Timezone         string `url:"timezone,omitempty"`          // +10:00
	AdditionalFields string `url:"additional_fields,omitempty"` // comma seperarted list (rejected_reason,rejected_reason_value)
	Currency         string `url:"currency,omitempty"`          // preferred
}

type parameter struct {
	*Options
	apiToken string `url:"api_token"` // gets set internally
}

// InstallReport to report installs
type InstallReport struct {
	AdditionalFields
}

// AdditionalFields will get created based on additional_fields input
type AdditionalFields struct {
	BlockedReason       string `json:"blocked_reason,omitempty"`
	BlockedReasonValue  string `json:"blocked_reason_value,omitempty"`
	BlockedSubReason    string `json:"blocked_sub_reason,omitempty"`
	BlockedReasonRule   string `json:"blocked_reason_rule,omitempty"`
	InstallAppStore     string `json:"install_app_store,omitempty"`
	CustomData          string `json:"custom_data,omitempty"`
	GPReferrer          string `json:"gp_referrer,omitempty"`
	GPClickTime         string `json:"gp_click_time,omitempty"`
	GPInstallBegin      string `json:"gp_install_begin,omitempty"`
	GPBroadcastReferrer string `json:"gp_broadcast_referrer,omitempty"`
	AmazonAID           string `json:"amazon_aid,omitempty"`
	KeywordMatchType    string `json:"keyword_match_type,omitempty"`
	NetworkAccountID    string `json:"network_account_id,omitempty"`
	RejectedReason      string `json:"rejected_reason,omitempty"`
	RejectedReasonValue string `json:"rejected_reason_value,omitempty"`
}

// TODO: Reasearch dynamic fields in sturcts for AdditionalFields
// https://stackoverflow.com/questions/40559250/golang-dynamic-creating-member-of-struct

// InstallReports returns reports of installs for options
func (s *ReportService) InstallReports(ctx context.Context, opt *Options) ([]*InstallReport, *Response, error) {
	u, err := addOptions("/installs_report/v5", parameter{
		Options:  opt,
		apiToken: s.client.APIToken,
	})
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	reports := []*InstallReport{}

	resp, err := s.client.Do(ctx, req, &reports)
	if err != nil {
		return nil, resp, err
	}

	return reports, resp, nil
}
