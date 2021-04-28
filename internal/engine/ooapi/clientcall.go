// Code generated by go generate; DO NOT EDIT.
// 2021-04-26 14:35:19.81671012 +0200 CEST m=+0.000117403

package ooapi

//go:generate go run ./internal/generator -file clientcall.go

import (
	"context"

	"github.com/ooni/probe-cli/v3/internal/engine/ooapi/apimodel"
)

func (c *Client) newCheckReportIDCaller() callerForCheckReportIDAPI {
	return &simpleCheckReportIDAPI{
		BaseURL:      c.BaseURL,
		HTTPClient:   c.HTTPClient,
		JSONCodec:    c.JSONCodec,
		RequestMaker: c.RequestMaker,
		UserAgent:    c.UserAgent,
	}
}

// CheckReportID calls the CheckReportID API.
func (c *Client) CheckReportID(
	ctx context.Context, req *apimodel.CheckReportIDRequest,
) (*apimodel.CheckReportIDResponse, error) {
	api := c.newCheckReportIDCaller()
	return api.Call(ctx, req)
}

func (c *Client) newCheckInCaller() callerForCheckInAPI {
	return &simpleCheckInAPI{
		BaseURL:      c.BaseURL,
		HTTPClient:   c.HTTPClient,
		JSONCodec:    c.JSONCodec,
		RequestMaker: c.RequestMaker,
		UserAgent:    c.UserAgent,
	}
}

// CheckIn calls the CheckIn API.
func (c *Client) CheckIn(
	ctx context.Context, req *apimodel.CheckInRequest,
) (*apimodel.CheckInResponse, error) {
	api := c.newCheckInCaller()
	return api.Call(ctx, req)
}

func (c *Client) newMeasurementMetaCaller() callerForMeasurementMetaAPI {
	return &withCacheMeasurementMetaAPI{
		API: &simpleMeasurementMetaAPI{
			BaseURL:      c.BaseURL,
			HTTPClient:   c.HTTPClient,
			JSONCodec:    c.JSONCodec,
			RequestMaker: c.RequestMaker,
			UserAgent:    c.UserAgent,
		},
		GobCodec: c.GobCodec,
		KVStore:  c.KVStore,
	}
}

// MeasurementMeta calls the MeasurementMeta API.
func (c *Client) MeasurementMeta(
	ctx context.Context, req *apimodel.MeasurementMetaRequest,
) (*apimodel.MeasurementMetaResponse, error) {
	api := c.newMeasurementMetaCaller()
	return api.Call(ctx, req)
}

func (c *Client) newTestHelpersCaller() callerForTestHelpersAPI {
	return &simpleTestHelpersAPI{
		BaseURL:      c.BaseURL,
		HTTPClient:   c.HTTPClient,
		JSONCodec:    c.JSONCodec,
		RequestMaker: c.RequestMaker,
		UserAgent:    c.UserAgent,
	}
}

// TestHelpers calls the TestHelpers API.
func (c *Client) TestHelpers(
	ctx context.Context, req *apimodel.TestHelpersRequest,
) (apimodel.TestHelpersResponse, error) {
	api := c.newTestHelpersCaller()
	return api.Call(ctx, req)
}

func (c *Client) newPsiphonConfigCaller() callerForPsiphonConfigAPI {
	return &withLoginPsiphonConfigAPI{
		API: &simplePsiphonConfigAPI{
			BaseURL:      c.BaseURL,
			HTTPClient:   c.HTTPClient,
			JSONCodec:    c.JSONCodec,
			RequestMaker: c.RequestMaker,
			UserAgent:    c.UserAgent,
		},
		JSONCodec: c.JSONCodec,
		KVStore:   c.KVStore,
		RegisterAPI: &simpleRegisterAPI{
			BaseURL:      c.BaseURL,
			HTTPClient:   c.HTTPClient,
			JSONCodec:    c.JSONCodec,
			RequestMaker: c.RequestMaker,
			UserAgent:    c.UserAgent,
		},
		LoginAPI: &simpleLoginAPI{
			BaseURL:      c.BaseURL,
			HTTPClient:   c.HTTPClient,
			JSONCodec:    c.JSONCodec,
			RequestMaker: c.RequestMaker,
			UserAgent:    c.UserAgent,
		},
	}
}

// PsiphonConfig calls the PsiphonConfig API.
func (c *Client) PsiphonConfig(
	ctx context.Context, req *apimodel.PsiphonConfigRequest,
) (apimodel.PsiphonConfigResponse, error) {
	api := c.newPsiphonConfigCaller()
	return api.Call(ctx, req)
}

func (c *Client) newTorTargetsCaller() callerForTorTargetsAPI {
	return &withLoginTorTargetsAPI{
		API: &simpleTorTargetsAPI{
			BaseURL:      c.BaseURL,
			HTTPClient:   c.HTTPClient,
			JSONCodec:    c.JSONCodec,
			RequestMaker: c.RequestMaker,
			UserAgent:    c.UserAgent,
		},
		JSONCodec: c.JSONCodec,
		KVStore:   c.KVStore,
		RegisterAPI: &simpleRegisterAPI{
			BaseURL:      c.BaseURL,
			HTTPClient:   c.HTTPClient,
			JSONCodec:    c.JSONCodec,
			RequestMaker: c.RequestMaker,
			UserAgent:    c.UserAgent,
		},
		LoginAPI: &simpleLoginAPI{
			BaseURL:      c.BaseURL,
			HTTPClient:   c.HTTPClient,
			JSONCodec:    c.JSONCodec,
			RequestMaker: c.RequestMaker,
			UserAgent:    c.UserAgent,
		},
	}
}

// TorTargets calls the TorTargets API.
func (c *Client) TorTargets(
	ctx context.Context, req *apimodel.TorTargetsRequest,
) (apimodel.TorTargetsResponse, error) {
	api := c.newTorTargetsCaller()
	return api.Call(ctx, req)
}

func (c *Client) newURLsCaller() callerForURLsAPI {
	return &simpleURLsAPI{
		BaseURL:      c.BaseURL,
		HTTPClient:   c.HTTPClient,
		JSONCodec:    c.JSONCodec,
		RequestMaker: c.RequestMaker,
		UserAgent:    c.UserAgent,
	}
}

// URLs calls the URLs API.
func (c *Client) URLs(
	ctx context.Context, req *apimodel.URLsRequest,
) (*apimodel.URLsResponse, error) {
	api := c.newURLsCaller()
	return api.Call(ctx, req)
}

func (c *Client) newOpenReportCaller() callerForOpenReportAPI {
	return &simpleOpenReportAPI{
		BaseURL:      c.BaseURL,
		HTTPClient:   c.HTTPClient,
		JSONCodec:    c.JSONCodec,
		RequestMaker: c.RequestMaker,
		UserAgent:    c.UserAgent,
	}
}

// OpenReport calls the OpenReport API.
func (c *Client) OpenReport(
	ctx context.Context, req *apimodel.OpenReportRequest,
) (*apimodel.OpenReportResponse, error) {
	api := c.newOpenReportCaller()
	return api.Call(ctx, req)
}

func (c *Client) newSubmitMeasurementCaller() callerForSubmitMeasurementAPI {
	return &simpleSubmitMeasurementAPI{
		BaseURL:      c.BaseURL,
		HTTPClient:   c.HTTPClient,
		JSONCodec:    c.JSONCodec,
		RequestMaker: c.RequestMaker,
		UserAgent:    c.UserAgent,
	}
}

// SubmitMeasurement calls the SubmitMeasurement API.
func (c *Client) SubmitMeasurement(
	ctx context.Context, req *apimodel.SubmitMeasurementRequest,
) (*apimodel.SubmitMeasurementResponse, error) {
	api := c.newSubmitMeasurementCaller()
	return api.Call(ctx, req)
}
