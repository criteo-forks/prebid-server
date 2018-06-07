package criteo

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGetDeviceType(t *testing.T) {

	// Setup:
	deviceTypeCases := []struct {
		deviceType string
		expected   string
	}{
		{"ios", "idfa"},
		{"Ios", "idfa"},
		{"IOS", "idfa"},
		{"android", "gaid"},
		{"unknown", "unknown"},
		{"", "unknown"},
		{"qwerty", "unknown"},
		{"qWerty", "unknown"},
		{"abc", "unknown"},
	}

	for _, uc := range deviceTypeCases {

		// Execute:
		result := getDeviceType(&(uc.deviceType))

		// Verify:
		if uc.expected != *result {
			t.Errorf("Bad getDeviceType for '%s'. Expected: %s, got %s", uc.deviceType, uc.expected, *result)
		}
	}
}

func TestCriteoRequestBuilderFillAllFields(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoRequestBuilderPartialFill(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoPublisherBuilderFillAllFields(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoPublisherBuilderPartialFill(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoUserBuilderFillAllFields(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoUserBuilderPartialFill(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoGdprConsentBuilderFillAllFields(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoGdprConsentBuilderPartialFill(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoRequestSlotBuilderFillAllFields(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoRequestSlotBuilderPartialFill(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoResponseSlotGetID(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoResponseSlotGetCreativeID(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestNewCriteoResponseFromBytes(t *testing.T) {
	// TODO: To implement
	// Setup:
	// Execute:
	// Verify:
}

func TestCriteoResponseBuilderFillAllFields(t *testing.T) {
	// Setup:
	var (
		dummyID    = "dummyId"
		dummyImpID = "dummyImpId"
	)

	criteoResponseCases := []struct {
		ID    *string
		ImpID *string
	}{
		{&dummyID, &dummyImpID},
		{&dummyID, nil},
		{nil, &dummyImpID},
		{nil, nil},
	}

	for _, uc := range criteoResponseCases {

		// Execute:
		criteoResponseBuilder := NewCriteoResponseBuilder()

		criteoResponseBuilder.Id(uc.ID)

		criteoResponseSlotBuilder := NewCriteoResponseSlotBuilder().ImpID(uc.ImpID)
		criteoResponseBuilder.Slots([]*CriteoResponseSlotBuilder{criteoResponseSlotBuilder})

		criteoResponseResult := criteoResponseBuilder.Build()

		// Verify:
		expectedResult := &criteoResponse{}
		expectedResult.ID = uc.ID
		expectedResult.Slots = []*criteoResponseSlot{
			&criteoResponseSlot{
				ImpID: uc.ImpID,
			},
		}

		if !reflect.DeepEqual(*expectedResult, *criteoResponseResult) {
			expectedJSON, _ := json.Marshal(expectedResult)
			resultJSON, _ := json.Marshal(criteoResponseResult)
			t.Errorf("Bad response from builder. Expected: %s, got %s", expectedJSON, resultJSON)
		}
	}
}

func TestCriteoResponseBuilderPartialFill(t *testing.T) {
	// Setup:
	var (
		dummyID    = "dummyId"
		dummyImpID = "dummyImpId"
	)

	criteoResponseCases := []struct {
		ID    *string
		ImpID *string
	}{
		{&dummyID, &dummyImpID},
		{&dummyID, nil},
		{nil, &dummyImpID},
		{nil, nil},
	}

	for _, uc := range criteoResponseCases {

		// Execute:
		expectedResult := &criteoResponse{}
		criteoResponseBuilder := NewCriteoResponseBuilder()
		if uc.ID != nil {
			criteoResponseBuilder.Id(uc.ID)
			expectedResult.ID = uc.ID
		}
		if uc.ImpID != nil {
			criteoResponseSlotBuilder := NewCriteoResponseSlotBuilder().ImpID(uc.ImpID)
			criteoResponseBuilder.Slots([]*CriteoResponseSlotBuilder{criteoResponseSlotBuilder})
			expectedResult.Slots = []*criteoResponseSlot{
				&criteoResponseSlot{
					ImpID: uc.ImpID,
				},
			}
		}

		criteoResponseResult := criteoResponseBuilder.Build()

		// Verify:
		if !reflect.DeepEqual(*expectedResult, *criteoResponseResult) {
			expectedJSON, _ := json.Marshal(expectedResult)
			resultJSON, _ := json.Marshal(criteoResponseResult)
			t.Errorf("Bad response from builder. Expected: %s, got %s", expectedJSON, resultJSON)
		}
	}
}

func TestCriteoResponseSlotBuilderFillAllFields(t *testing.T) {

	// Setup:
	var (
		dummyImpID      = "dummyImpId"
		dummyZoneID     = uint(123)
		dummyCpm        = 0.01
		dummyCurrency   = "USD"
		dummyWidth      = uint(720)
		dummyHeight     = uint(90)
		dummyCreativeID = "dummyCreativeId"
	)

	criteoResponseSlotCases := []struct {
		impID    *string
		zoneID   *uint
		cpm      *float64
		currency *string
		width    *uint
		height   *uint
		creative *string
	}{
		{&dummyImpID, &dummyZoneID, &dummyCpm, &dummyCurrency, &dummyWidth, &dummyHeight, &dummyCreativeID},
		{nil, &dummyZoneID, &dummyCpm, &dummyCurrency, &dummyWidth, &dummyHeight, &dummyCreativeID},
		{&dummyImpID, nil, &dummyCpm, &dummyCurrency, &dummyWidth, &dummyHeight, &dummyCreativeID},
		{&dummyImpID, &dummyZoneID, nil, &dummyCurrency, &dummyWidth, &dummyHeight, &dummyCreativeID},
		{&dummyImpID, &dummyZoneID, &dummyCpm, nil, &dummyWidth, &dummyHeight, &dummyCreativeID},
		{&dummyImpID, &dummyZoneID, &dummyCpm, &dummyCurrency, nil, &dummyHeight, &dummyCreativeID},
		{&dummyImpID, &dummyZoneID, &dummyCpm, &dummyCurrency, &dummyWidth, nil, &dummyCreativeID},
		{&dummyImpID, &dummyZoneID, &dummyCpm, &dummyCurrency, &dummyWidth, &dummyHeight, nil},
		{nil, nil, nil, nil, nil, nil, nil},
	}

	for _, uc := range criteoResponseSlotCases {

		// Execute:
		criteoResponseSlotBuilder := NewCriteoResponseSlotBuilder()
		criteoResponseSlotBuilder.ImpID(uc.impID)
		criteoResponseSlotBuilder.ZoneID(uc.zoneID)
		criteoResponseSlotBuilder.Cpm(uc.cpm)
		criteoResponseSlotBuilder.Currency(uc.currency)
		criteoResponseSlotBuilder.Width(uc.width)
		criteoResponseSlotBuilder.Height(uc.height)
		criteoResponseSlotBuilder.Creative(uc.creative)

		criteoResponseSlotResult := criteoResponseSlotBuilder.Build()

		// Verify:
		expectedResult := &criteoResponseSlot{
			ImpID:    uc.impID,
			ZoneID:   uc.zoneID,
			Cpm:      uc.cpm,
			Currency: uc.currency,
			Width:    uc.width,
			Height:   uc.height,
			Creative: uc.creative,
		}

		if !reflect.DeepEqual(*expectedResult, *criteoResponseSlotResult) {
			expectedJSON, _ := json.Marshal(expectedResult)
			resultJSON, _ := json.Marshal(criteoResponseSlotResult)
			t.Errorf("Bad response slot from builder. Expected: %s, got %s", expectedJSON, resultJSON)
		}
	}
}

func TestCriteoResponseSlotBuilderPartialFill(t *testing.T) {

	// Setup:
	var (
		dummyImpID      = "dummyImpId"
		dummyZoneID     = uint(123)
		dummyCpm        = 0.01
		dummyCurrency   = "USD"
		dummyWidth      = uint(720)
		dummyHeight     = uint(90)
		dummyCreativeID = "dummyCreativeId"
	)

	criteoResponseSlotCases := []struct {
		impID    *string
		zoneID   *uint
		cpm      *float64
		currency *string
		width    *uint
		height   *uint
		creative *string
	}{
		{&dummyImpID, &dummyZoneID, &dummyCpm, &dummyCurrency, &dummyWidth, &dummyHeight, &dummyCreativeID},
		{&dummyImpID, nil, nil, nil, nil, nil, nil},
		{nil, &dummyZoneID, nil, nil, nil, nil, nil},
		{nil, nil, &dummyCpm, nil, nil, nil, nil},
		{nil, nil, nil, &dummyCurrency, nil, nil, nil},
		{nil, nil, nil, nil, &dummyWidth, nil, nil},
		{nil, nil, nil, nil, nil, &dummyHeight, nil},
		{nil, nil, nil, nil, nil, nil, &dummyCreativeID},
		{nil, nil, nil, nil, nil, nil, nil},
	}

	for _, uc := range criteoResponseSlotCases {

		// Execute:
		expectedResult := &criteoResponseSlot{}
		criteoResponseSlotBuilder := NewCriteoResponseSlotBuilder()
		if uc.impID != nil {
			criteoResponseSlotBuilder.ImpID(uc.impID)
			expectedResult.ImpID = uc.impID
		}
		if uc.zoneID != nil {
			criteoResponseSlotBuilder.ZoneID(uc.zoneID)
			expectedResult.ZoneID = uc.zoneID
		}
		if uc.cpm != nil {
			criteoResponseSlotBuilder.Cpm(uc.cpm)
			expectedResult.Cpm = uc.cpm
		}
		if uc.currency != nil {
			criteoResponseSlotBuilder.Currency(uc.currency)
			expectedResult.Currency = uc.currency
		}
		if uc.width != nil {
			criteoResponseSlotBuilder.Width(uc.width)
			expectedResult.Width = uc.width
		}
		if uc.height != nil {
			criteoResponseSlotBuilder.Height(uc.height)
			expectedResult.Height = uc.height
		}
		if uc.creative != nil {
			criteoResponseSlotBuilder.Creative(uc.creative)
			expectedResult.Creative = uc.creative
		}

		criteoResponseSlotResult := criteoResponseSlotBuilder.Build()

		// Verify:

		if !reflect.DeepEqual(*expectedResult, *criteoResponseSlotResult) {
			expectedJSON, _ := json.Marshal(expectedResult)
			resultJSON, _ := json.Marshal(criteoResponseSlotResult)
			t.Errorf("Bad response slot from builder. Expected: %s, got %s", expectedJSON, resultJSON)
		}
	}
}
