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

func TestShouldRejectRequest(t *testing.T) {
	// Setup:
	bidRequestCases := []struct {
		description string
		input       *openrtb.BidRequest
		expected    bool
	}{
		{
			"Default request, everything is empty. Should reject.",
			&openrtb.BidRequest{},
			true,
		},
		{
			"App set but no Device. Should reject.",
			&openrtb.BidRequest{
				App: &openrtb.App{},
			},
			true,
		},
		{
			"Device set without IFA but no App. Should reject.",
			&openrtb.BidRequest{
				Device: &openrtb.Device{},
			},
			true,
		},
		{
			"Device set without IFA and App. Should reject.",
			&openrtb.BidRequest{
				App:    &openrtb.App{},
				Device: &openrtb.Device{},
			},
			true,
		},
		{
			"Device set without IFA, App and Site (shouldn't happen but just in case). Should reject.",
			&openrtb.BidRequest{
				Site:   &openrtb.Site{},
				App:    &openrtb.App{},
				Device: &openrtb.Device{},
			},
			true,
		},
		{
			"Device IFA set but no App. Should reject.",
			&openrtb.BidRequest{
				Device: &openrtb.Device{
					IFA: "non-empty-ifa",
				},
			},
			true,
		},
		{
			"Device and Site set. Should reject.",
			&openrtb.BidRequest{
				Site: &openrtb.Site{},
				Device: &openrtb.Device{
					IFA: "non-empty-ifa",
				},
			},
			true,
		},
		{
			"Device IFA, App and Site set (shouldn't happen but just in case). Shouldn't reject.",
			&openrtb.BidRequest{
				Site: &openrtb.Site{},
				App:  &openrtb.App{},
				Device: &openrtb.Device{
					IFA: "non-empty-ifa",
				},
			},
			false,
		},
		{
			"In App-only request. Shouldn't reject",
			&openrtb.BidRequest{
				App: &openrtb.App{},
				Device: &openrtb.Device{
					IFA: "non-empty-ifa",
				},
			},
			false,
		},
	}

	for _, tc := range bidRequestCases {

		// Execute:
		result := shouldRejectRequest(tc.input)

		// Verify:
		if tc.expected != result {
			t.Errorf("shouldRejectRequest scenario: '%s' was incorrect, got %t, want %t.", tc.description, result, tc.expected)
		}
	}
}

func TestGetCriteoRequestHeaders(t *testing.T) {

	// Setup:
	dummyCookieID := "random-uid"
	dummyIP := "1.1.1.1"
	dummyUa := "A very random UA"
	HeadersTestCases := []struct {
		input    *criteoRequest
		expected http.Header
	}{
		{
			NewCriteoRequestBuilder().Build(),
			http.Header{},
		},
		{
			NewCriteoRequestBuilder().
				User(
					NewCriteoUserBuilder().
						CookieID(&dummyCookieID),
				).
				Build(),
			http.Header{
				"Cookie": []string{"uid=" + dummyCookieID},
			},
		},
		{
			NewCriteoRequestBuilder().
				User(
					NewCriteoUserBuilder().
						IP(&dummyIP),
				).
				Build(),
			http.Header{
				"X-Client-Ip": []string{dummyIP},
			},
		},
		{
			NewCriteoRequestBuilder().
				User(
					NewCriteoUserBuilder().
						Ua(&dummyUa),
				).
				Build(),
			http.Header{
				"User-Agent": []string{dummyUa},
			},
		},
		{
			NewCriteoRequestBuilder().
				User(
					NewCriteoUserBuilder().
						CookieID(&dummyCookieID).
						IP(&dummyIP).
						Ua(&dummyUa),
				).
				Build(),
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
			t.Errorf("getCriteoRequestHeaders was incorrect, got '%s', want '%s'.", result, uc.expected)
		}
	}
}

func TestGetCriteoRequest(t *testing.T) {

	// Setup:
	var (
		dummyRequestID         = "random request ID"
		dummyPublisherBundleID = "bundleid"
		dummyPublisherURL      = "test.com"
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
			BundleID: &dummyPublisherBundleID,
			URL:      &dummyPublisherURL,
		},
		User: &criteoUser{
			DeviceID:     &dummyUserDeviceID,
			DeviceOs:     &dummyUserDeviceOs,
			DeviceIDType: &dummyUserDeviceIDType,
			CookieID:     &dummyUserCookieID,
			IP:           &dummyUserIP,
			Ua:           &dummyUserUA,
		},
		GdprConsent: &criteoGdprConsent{
			GdprApplies: &dummyGdprApplies,
			ConsentData: &dummyGdprConsentData,
		},
		Slots: []*criteoRequestSlot{
			&criteoRequestSlot{
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
			openrtb.Imp{
				ID:  dummySlotImpID,
				Ext: impExtJSON,
			},
		},
	}

	// Execute:
	result, err := getCriteoRequest(incomingRequest)

	// Verify:
	if err != nil {
		t.Errorf("getCriteoRequest has errors: %s", err)
	}

	requestResult := result.Build()
	if *expectedCriteoRequest.ID != *requestResult.ID ||
		!reflect.DeepEqual(*expectedCriteoRequest, *requestResult) ||
		!reflect.DeepEqual(*expectedCriteoRequest.Publisher, *requestResult.Publisher) ||
		!reflect.DeepEqual(*expectedCriteoRequest.User, *requestResult.User) ||
		!reflect.DeepEqual(*expectedCriteoRequest.GdprConsent, *requestResult.GdprConsent) ||
		len(expectedCriteoRequest.Slots) != len(requestResult.Slots) ||
		!reflect.DeepEqual(*expectedCriteoRequest.Slots[0], *requestResult.Slots[0]) {
		actualResultJSON, _ := json.Marshal(requestResult)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("getCriteoRequest was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
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

	// Execute:
	result := NewCriteoRequestBuilder().
		GdprConsent(getGdprConsent(incomingRequest.User, incomingRequest.Regs))

	// Verify:
	requestResult := result.Build()
	if !reflect.DeepEqual(*expectedCriteoRequest, *requestResult) {
		actualResultJSON, _ := json.Marshal(requestResult)
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
	)
	expectedCriteoRequest := &criteoRequest{
		User: &criteoUser{
			DeviceID:     &dummyUserDeviceID,
			DeviceOs:     &dummyUserDeviceOs,
			DeviceIDType: &dummyUserDeviceIDType,
			CookieID:     &dummyUserCookieID,
			IP:           &dummyUserIP,
			Ua:           &dummyUserUA,
		},
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
	}

	// Execute:
	result := NewCriteoRequestBuilder().
		User(getUser(incomingRequest.User, incomingRequest.Device))

	// Verify:
	requestResult := result.Build()
	if !reflect.DeepEqual(*expectedCriteoRequest.User, *requestResult.User) {
		actualResultJSON, _ := json.Marshal(requestResult)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("getUser was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
	}
}

func TestPublisher(t *testing.T) {

	// Setup:
	var (
		dummyPublisherBundleID = "bundleid"
		dummyPublisherURL      = "test.com"
	)
	expectedCriteoRequest := &criteoRequest{
		Publisher: &criteoPublisher{
			BundleID: &dummyPublisherBundleID,
			URL:      &dummyPublisherURL,
		},
	}

	incomingRequest := &openrtb.BidRequest{
		App: &openrtb.App{
			Bundle: dummyPublisherBundleID,
		},
		Site: &openrtb.Site{
			Page: dummyPublisherURL,
		},
	}

	// Execute:
	result := NewCriteoRequestBuilder().
		Publisher(getPublisher(incomingRequest.App, incomingRequest.Site))

	// Verify:
	requestResult := result.Build()
	if !reflect.DeepEqual(*expectedCriteoRequest.Publisher, *requestResult.Publisher) {
		actualResultJSON, _ := json.Marshal(requestResult)
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
			&criteoRequestSlot{
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
			openrtb.Imp{
				ID:  dummySlotImpID,
				Ext: impExtJSON,
			},
		},
	}

	// Execute:
	slots, err := getRequestSlots(incomingRequest.Imp)
	result := NewCriteoRequestBuilder().
		Slots(slots)

	// Verify:
	if err != nil {
		t.Errorf("getCriteoRequestSlots has errors: %s", err)
	}

	requestResult := result.Build()
	if len(expectedCriteoRequest.Slots) != len(requestResult.Slots) ||
		!reflect.DeepEqual(*expectedCriteoRequest.Slots[0], *requestResult.Slots[0]) {
		actualResultJSON, _ := json.Marshal(requestResult)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("getCriteoRequest was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
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
	slotsBuilders := NewCriteoRequestSlotsBuilders(uint(len(dummySlots)))

	for i := range dummySlots {

		// Build expected slots
		slotsBuilders[i].ImpID(&dummySlots[i].ID)
		slotsBuilders[i].ZoneID(&dummySlots[i].ZoneID)

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

	expectedCriteoRequest := NewCriteoRequestBuilder().Slots(slotsBuilders).Build()

	// Execute:
	slots, err := getRequestSlots(incomingRequest.Imp)
	result := NewCriteoRequestBuilder().Slots(slots)

	// Verify:
	if err != nil {
		t.Errorf("getCriteoRequestSlots has errors: %s", err)
	}

	requestResult := result.Build()
	if len(expectedCriteoRequest.Slots) != len(requestResult.Slots) ||
		!reflect.DeepEqual(expectedCriteoRequest.Slots, requestResult.Slots) {
		actualResultJSON, _ := json.Marshal(requestResult)
		expectedResultJSON, _ := json.Marshal(expectedCriteoRequest)
		t.Errorf("getCriteoRequest was incorrect, got '%s', want '%s'.", actualResultJSON, expectedResultJSON)
	}
}
