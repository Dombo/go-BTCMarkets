package report

import (
	"BTCMarkets/lib"
	"errors"
	"fmt"
	"net/http"
)

const (
	pathBase = "reports"
)

type Report struct {
	ID string
	ContentURL string
	CreationTime BTCMarkets.SpecialDatetime
	Type string
	Status string
	Format string
}

// The API design returns a different response for creating a report vs retrieving it
// This behaviour is also undocumented
type CreatedReport struct {
	ID string `json:"reportID"`
}

type createReportRequest struct {
	Type                string	`json:"type"`
	Format              string	`json:"format"`
}

func Create(c *BTCMarkets.Client, format string) (*CreatedReport, error) {
	requestData, err := requestDataForReport("", format)
	if err != nil {
		return nil, err
	}

	report := &CreatedReport{}
	if err := c.Do(report, http.MethodPost, pathBase, requestData); err != nil {
		return nil, err
	}

	return report, nil
}

func Get(c *BTCMarkets.Client, id string) (*Report, error) {
	if id == "" {
		return nil, errors.New("id is required to retrieve a report")
	}

	path := fmt.Sprintf("%s/%s", pathBase, id)

	report := &Report{}
	if err := c.Do(report, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	return report, nil
}

func requestDataForReport(reportType string, reportFormat string) (*createReportRequest, error) {
	if reportType == "" {
		reportType = "TransactionReport"
	}
	if reportFormat == "" {
		return nil, errors.New("reportFormat is required")
	}

	request := &createReportRequest{
		Type:       reportType,
		Format:		reportFormat,
	}

	return request, nil
}