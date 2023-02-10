package criteo

import (
	"github.com/prebid/prebid-server/adapters/adapterstest"
	"testing"
)

func TestJsonSamples(t *testing.T) {

	// Setup:
	bidder := &adapter{
		endpoint: "https://bidder.criteo.com/openrtb/prebid/auction/request?profile=230",
	}

	// Execute & Verify:
	adapterstest.RunJSONBidderTest(t, "criteotest", bidder)
}
