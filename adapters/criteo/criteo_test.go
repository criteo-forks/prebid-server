package criteo

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/adapters/adapterstest"
	"github.com/prebid/prebid-server/openrtb_ext"
)

func TestJsonSamples(t *testing.T) {
	adapterstest.RunJSONBidderTest(t, "criteotest", NewCriteoBidder(nil, "https://bidder.criteo.com/cdb?profileId=230"))
}

func TestNewCriteoRequestHeaders(t *testing.T) {
	// Setup:
	dummyCookieID := "random-uid"
	dummyIP := "1.1.1.1"
	dummyUa := "A very random UA"
	HeadersTestCases := []struct {
		input    *criteoRequest
		expected http.Header
	}{
		{
			&criteoRequest{},
			http.Header{},
		},
		{
			&criteoRequest{
				User: &criteoUser{
					CookieID: &dummyCookieID,
				},
			},
			http.Header{
				"Cookie": []string{"uid=" + dummyCookieID},
			},
		},
		{
			&criteoRequest{
				User: &criteoUser{
					IP: &dummyIP,
				},
			},
			http.Header{
				"X-Client-Ip": []string{dummyIP},
			},
		},
		{
			&criteoRequest{
				User: &criteoUser{
					UA: &dummyUa,
				},
			},
			http.Header{
				"User-Agent": []string{dummyUa},
			},
		},
		{
			&criteoRequest{
				User: &criteoUser{
					CookieID: &dummyCookieID,
					IP:       &dummyIP,
					UA:       &dummyUa,
				},
			},
			http.Header{
				"Cookie":      []string{"uid=" + dummyCookieID},
				"X-Client-Ip": []string{dummyIP},
				"User-Agent":  []string{dummyUa},
			},
		},
	}

	for _, uc := range HeadersTestCases {
		// Execute:
		result := getCriteoRequestHeaders(uc.input)

		// Verify:
		if !reflect.DeepEqual(uc.expected, result) {
			t.Errorf("newCriteoRequestHeaders was incorrect, got '%s', want '%s'.", result, uc.expected)
		}
	}
}

func TestNewCriteoRequest(t *testing.T) {
	// Setup:
	var (
		dummyRequestID         = "random request ID"
		dummyPublisherBundleID = "bundleid"
		dummyPublisherURL      = "test.com"
		dummyPublisherSiteID   = "siteid"
		dummyUserDeviceID      = "random-device-id"
		dummyUserDeviceOs      = "android"
		dummyUserDeviceIDType  = "gaid"
		dummyUserCookieID      = "random-cookie-id"
		dummyUserIP            = "1.1.1.1"
		dummyUserUA            = "random UA"
		dummyGdprApplies       = true
		dummyGdprAppliesUint   = int8(1)
		dummyGdprConsentData   = "randomconsentdata"
		dummySlotImpID         = "fake-imp-id-1"
		dummySlotZoneID        = uint(1)
	)
	// The request doesn't make any sense but aims to fill every single criteo request fields
	expectedCriteoRequest := &criteoRequest{
		ID: &dummyRequestID,
		Publisher: &criteoPublisher{
			SiteID:   &dummyPublisherSiteID,
			BundleID: &dummyPublisherBundleID,
			URL:      &dummyPublisherURL,
		},
		User: &criteoUser{
			DeviceID:     &dummyUserDeviceID,
			DeviceOs:     &dummyUserDeviceOs,
			DeviceIDType: &dummyUserDeviceIDType,
			CookieID:     &dummyUserCookieID,
			IP:           &dummyUserIP,
			UA:           &dummyUserUA,
		},
		GdprConsent: &criteoGdprConsent{
			GdprApplies: &dummyGdprApplies,
			ConsentData: &dummyGdprConsentData,
		},
		Slots: []*criteoRequestSlot{
			{
				ImpID:  &dummySlotImpID,
				ZoneID: &dummySlotZoneID,
			},
		},
	}

	userExtJSON, _ := json.Marshal(&openrtb_ext.ExtUser{
		Consent: dummyGdprConsentData,
	})
	regsExtJSON, _ := json.Marshal(&openrtb_ext.ExtRegs{
		GDPR: &dummyGdprAppliesUint,
	})
	bidderExtJSON, _ := json.Marshal(&openrtb_ext.ExtImpCriteo{
		ZoneID: dummySlotZoneID,
	})
	impExtJSON, _ := json.Marshal(&adapters.ExtImpBidder{
		Bidder: bidderExtJSON,
	})
	incomingRequest := &openrtb.BidRequest{
		ID: dummyRequestID,
		App: &openrtb.App{
			Bundle: dummyPublisherBundleID,
		},
		Site: &openrtb.Site{
			ID:   dummyPublisherSiteID,
			Page: dummyPublisherURL,
		},
		User: &openrtb.User{
			BuyerUID: dummyUserCookieID,
			Ext:      userExtJSON,
		},
		Regs: &openrtb.Regs{
			Ext: regsExtJSON,
		},
		Device: &openrtb.Device{
			IFA: dummyUserDeviceID,
			OS:  dummyUserDeviceOs,
			IP:  dummyUserIP,
			UA:  dummyUserUA,
		},
		Imp: []openrtb.Imp{
			{
				ID:  dummySlotImpID,
				Ext: impExtJSON,
			},
		},
	}

	// Execute:
	result, err := newCriteoRequest(incomingRequest)

	// Verify:
	if err != nil {
		t.Errorf("newCriteoRequest has errors: %s", err)
	}

	if *expectedCriteoRequest.ID != *result.ID ||
		!reflect.DeepEqual(*expectedCriteoRequest, *result) ||
		!reflect.DeepEqual(*expectedCriteoRequest.Publisher, *result.Publisher) ||
		!reflect.DeepEqual(*expectedCriteoRequest.User, *result.User) ||
		!reflect.DeepEqual(*expectedCriteoRequest.GdprConsent, *result.GdprConsent) ||
		len(expectedCriteoRequest.Slots) != len(result.Slots) ||
		!reflect.DeepEqual(*expectedCriteoRequest.Slots[0], *result.Slots[0]) {
		actualResultJSON, _ := json.Marshal(result)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("newCriteoRequest was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
	}
}

func TestGetGdprConsent(t *testing.T) {
	// Setup:
	var (
		dummyGdprApplies     = true
		dummyGdprConsentData = "randomconsentdata"
		dummyGdprAppliesUint = int8(1)
	)

	expectedCriteoRequest := &criteoRequest{
		GdprConsent: &criteoGdprConsent{
			GdprApplies: &dummyGdprApplies,
			ConsentData: &dummyGdprConsentData,
		},
	}

	userExtJSON, _ := json.Marshal(&openrtb_ext.ExtUser{
		Consent: dummyGdprConsentData,
	})
	regsExtJSON, _ := json.Marshal(&openrtb_ext.ExtRegs{
		GDPR: &dummyGdprAppliesUint,
	})
	incomingRequest := &openrtb.BidRequest{
		User: &openrtb.User{
			Ext: userExtJSON,
		},
		Regs: &openrtb.Regs{
			Ext: regsExtJSON,
		},
	}

	var regsExt *openrtb_ext.ExtRegs
	if incomingRequest.Regs != nil {
		json.Unmarshal(incomingRequest.Regs.Ext, &regsExt)
	}

	// Execute:
	result := &criteoRequest{
		GdprConsent: newCriteoGdprConsent(incomingRequest.User, regsExt),
	}

	// Verify:
	if !reflect.DeepEqual(*expectedCriteoRequest, *result) {
		actualResultJSON, _ := json.Marshal(result)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("getGdprConsent was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
	}
}

func TestGetUser(t *testing.T) {
	// Setup:
	var (
		dummyUserDeviceID     = "random-device-id"
		dummyUserDeviceOs     = "android"
		dummyUserDeviceIDType = "gaid"
		dummyUserCookieID     = "random-cookie-id"
		dummyUserIP           = "1.1.1.1"
		dummyUserUA           = "random UA"
		dummyCcpaString       = "1YYY"
	)
	expectedCriteoRequest := &criteoRequest{
		User: &criteoUser{
			DeviceID:     &dummyUserDeviceID,
			DeviceOs:     &dummyUserDeviceOs,
			DeviceIDType: &dummyUserDeviceIDType,
			CookieID:     &dummyUserCookieID,
			IP:           &dummyUserIP,
			UA:           &dummyUserUA,
			UspIab:       &dummyCcpaString,
		},
	}

	regsExt := &openrtb_ext.ExtRegs{
		USPrivacy: dummyCcpaString,
	}
	regsExtData, err := json.Marshal(regsExt)
	if err != nil {
		t.Errorf("cannot marshal regsExt data")
	}

	incomingRequest := &openrtb.BidRequest{
		User: &openrtb.User{
			BuyerUID: dummyUserCookieID,
		},
		Device: &openrtb.Device{
			IFA: dummyUserDeviceID,
			OS:  dummyUserDeviceOs,
			IP:  dummyUserIP,
			UA:  dummyUserUA,
		},
		Regs: &openrtb.Regs{
			Ext: regsExtData,
		},
	}

	// Execute:
	result := &criteoRequest{
		User: newCriteoUser(incomingRequest.User, incomingRequest.Device, regsExt),
	}

	// Verify:
	if !reflect.DeepEqual(*expectedCriteoRequest.User, *result.User) {
		actualResultJSON, _ := json.Marshal(result)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("getUser was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
	}
}

func TestPublisher(t *testing.T) {
	// Setup:
	var (
		dummyPublisherSiteID   = "siteid"
		dummyPublisherBundleID = "bundleid"
		dummyPublisherURL      = "test.com"
	)
	expectedCriteoRequest := &criteoRequest{
		Publisher: &criteoPublisher{
			SiteID:   &dummyPublisherSiteID,
			BundleID: &dummyPublisherBundleID,
			URL:      &dummyPublisherURL,
		},
	}

	incomingRequest := &openrtb.BidRequest{
		App: &openrtb.App{
			Bundle: dummyPublisherBundleID,
		},
		Site: &openrtb.Site{
			ID:   dummyPublisherSiteID,
			Page: dummyPublisherURL,
		},
	}

	// Execute:
	result := &criteoRequest{
		Publisher: newCriteoPublisher(incomingRequest.App, incomingRequest.Site),
	}

	// Verify:
	if !reflect.DeepEqual(*expectedCriteoRequest.Publisher, *result.Publisher) {
		actualResultJSON, _ := json.Marshal(result)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("getPublisher was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
	}
}

func TestGetRequestSlots(t *testing.T) {
	// Setup:
	var (
		dummySlotImpID  = "fake-imp-id-1"
		dummySlotZoneID = uint(1)
	)
	expectedCriteoRequest := &criteoRequest{
		Slots: []*criteoRequestSlot{
			{
				ImpID:  &dummySlotImpID,
				ZoneID: &dummySlotZoneID,
			},
		},
	}

	bidderExtJSON, _ := json.Marshal(&openrtb_ext.ExtImpCriteo{
		ZoneID: dummySlotZoneID,
	})
	impExtJSON, _ := json.Marshal(&adapters.ExtImpBidder{
		Bidder: bidderExtJSON,
	})
	incomingRequest := &openrtb.BidRequest{
		Imp: []openrtb.Imp{
			{
				ID:  dummySlotImpID,
				Ext: impExtJSON,
			},
		},
	}

	// Execute:
	slots, err := newCriteoRequestSlots(incomingRequest.Imp)
	result := &criteoRequest{
		Slots: slots,
	}

	// Verify:
	if err != nil {
		t.Errorf("newCriteoRequestSlots has errors: %s", err)
	}

	if len(expectedCriteoRequest.Slots) != len(result.Slots) ||
		!reflect.DeepEqual(*expectedCriteoRequest.Slots[0], *result.Slots[0]) {
		actualResultJSON, _ := json.Marshal(result)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("newCriteoRequest was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
	}
}

func TestGetRequestMultipleSlots(t *testing.T) {
	// Setup:
	dummySlots := []struct {
		ID     string
		ZoneID uint
	}{
		{"fake-imp-id-1", uint(1)},
		{"fake-imp-id-2", uint(2)},
		{"fake-imp-id-3", uint(3)},
		{"fake-imp-id-4", uint(4)},
		{"fake-imp-id-5", uint(5)},
	}

	incomingRequest := &openrtb.BidRequest{
		Imp: make([]openrtb.Imp, len(dummySlots)),
	}
	slots := make([]*criteoRequestSlot, uint(len(dummySlots)))

	for i := range dummySlots {
		// Build expected slots
		slots[i] = &criteoRequestSlot{
			ImpID:  &dummySlots[i].ID,
			ZoneID: &dummySlots[i].ZoneID,
		}

		// Build incoming request imps
		bidderExtJSON, _ := json.Marshal(&openrtb_ext.ExtImpCriteo{
			ZoneID: dummySlots[i].ZoneID,
		})
		impExtJSON, _ := json.Marshal(&adapters.ExtImpBidder{
			Bidder: bidderExtJSON,
		})
		incomingRequest.Imp[i] = openrtb.Imp{
			ID:  dummySlots[i].ID,
			Ext: impExtJSON,
		}
	}

	expectedCriteoRequestSlots, err := newCriteoRequestSlots(incomingRequest.Imp)
	expectedCriteoRequest := &criteoRequest{
		Slots: expectedCriteoRequestSlots,
	}

	// Execute:
	slotsResult, err := newCriteoRequestSlots(incomingRequest.Imp)
	result := &criteoRequest{
		Slots: slotsResult,
	}

	// Verify:
	if err != nil {
		t.Errorf("newCriteoRequestSlots has errors: %s", err)
	}

	if len(expectedCriteoRequest.Slots) != len(result.Slots) ||
		!reflect.DeepEqual(expectedCriteoRequest.Slots, result.Slots) {
		actualResultJSON, _ := json.Marshal(result)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("newCriteoRequest was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
	}
}
