package criteo

import (
	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/openrtb_ext"
	"testing"
)

// This file actually intends to test static/bidder-params/criteo.json
//
// These also validate the format of the external API: request.imp[i].ext.criteo

// TestValidParams makes sure that the criteo schema accepts all imp.ext fields which we intend to support.
func TestValidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator("../../static/bidder-params")
	if err != nil {
		t.Fatalf("Failed to fetch the json-schemas. %v", err)
	}

	for _, validParam := range validParams {
		if err := validator.Validate(openrtb_ext.BidderCriteo, openrtb.RawJSON(validParam)); err != nil {
			t.Errorf("Schema rejected criteo params: %s", validParam)
		}
	}
}

// TestInvalidParams makes sure that the criteo schema rejects all the imp.ext fields we don't support.
func TestInvalidParams(t *testing.T) {
	validator, err := openrtb_ext.NewBidderParamsValidator("../../static/bidder-params")
	if err != nil {
		t.Fatalf("Failed to fetch the json-schemas. %v", err)
	}

	for _, invalidParam := range invalidParams {
		if err := validator.Validate(openrtb_ext.BidderCriteo, openrtb.RawJSON(invalidParam)); err == nil {
			t.Errorf("Schema allowed unexpected params: %s", invalidParam)
		}
	}
}

var validParams = []string{
	`{"zoneid": 123456}`,
}

var invalidParams = []string{
	``,
	`null`,
	`true`,
	`5`,
	`4.2`,
	`[]`,
	`{}`,
	`{"zoneid": "abcdef"}`,
	`{"zoneid": "123456"}`,
}