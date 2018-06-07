package criteo

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

type criteoRequest struct {
	ID          *string              `json:"id,omitempty"`
	Publisher   *criteoPublisher     `json:"publisher,omitempty"`
	User        *criteoUser          `json:"user,omitempty"`
	GdprConsent *criteoGdprConsent   `json:"gdprConsent,omitempty"`
	Slots       []*criteoRequestSlot `json:"slots,omitempty"`
}

// CriteoRequestBuilder allows to create `criteoRequest` objects
type CriteoRequestBuilder struct {
	id                        *string
	publisherBuilder          *CriteoPublisherBuilder
	userBuilder               *CriteoUserBuilder
	gdprConsentBuilder        *CriteoGdprConsentBuilder
	criteoRequestSlotBuilders []*CriteoRequestSlotBuilder
}

// NewCriteoRequestBuilder creates a new builder instance
func NewCriteoRequestBuilder() *CriteoRequestBuilder {
	return &CriteoRequestBuilder{}
}

// ID set `ID` for the `criteoRequest` to be created
func (b *CriteoRequestBuilder) ID(requestId *string) *CriteoRequestBuilder {
	b.id = requestId
	return b
}

// Publisher set `Publisher` for the `criteoRequest` to be created
func (b *CriteoRequestBuilder) Publisher(publisherBuilder *CriteoPublisherBuilder) *CriteoRequestBuilder {
	b.publisherBuilder = publisherBuilder
	return b
}

// User set `User` for the `criteoRequest` to be created
func (b *CriteoRequestBuilder) User(userBuilder *CriteoUserBuilder) *CriteoRequestBuilder {
	b.userBuilder = userBuilder
	return b
}

// GdprConsent set `GdprConsent` for the `criteoRequest` to be created
func (b *CriteoRequestBuilder) GdprConsent(gdprConsentBuilder *CriteoGdprConsentBuilder) *CriteoRequestBuilder {
	b.gdprConsentBuilder = gdprConsentBuilder
	return b
}

// Slots set `Slots` for the `criteoRequest` to be created
func (b *CriteoRequestBuilder) Slots(slotsBuilders []*CriteoRequestSlotBuilder) *CriteoRequestBuilder {
	b.criteoRequestSlotBuilders = slotsBuilders
	return b
}

// Build returns a `criteoRequest` instance based on the builder fields
func (b *CriteoRequestBuilder) Build() *criteoRequest {
	criteoReq := &criteoRequest{}

	// Fill request id if needed
	if b.id != nil {
		criteoReq.ID = b.id
	}

	// Fill publisher if needed
	if b.publisherBuilder != nil && *b.publisherBuilder != (CriteoPublisherBuilder{}) {
		criteoReq.Publisher = b.publisherBuilder.Build()
	}

	// Fill user if needed
	if b.userBuilder != nil && *b.userBuilder != (CriteoUserBuilder{}) {
		criteoReq.User = b.userBuilder.Build()
	}

	// Fill gdpr consent if needed
	if b.gdprConsentBuilder != nil && *b.gdprConsentBuilder != (CriteoGdprConsentBuilder{}) {
		criteoReq.GdprConsent = b.gdprConsentBuilder.Build()
	}

	// Fill slots if needed
	if b.criteoRequestSlotBuilders != nil && len(b.criteoRequestSlotBuilders) > 0 {
		criteoReq.Slots = make([]*criteoRequestSlot, len(b.criteoRequestSlotBuilders))
		for i := 0; i < len(b.criteoRequestSlotBuilders); i++ {
			criteoReq.Slots[i] = b.criteoRequestSlotBuilders[i].Build()
		}
	}

	return criteoReq
}

type criteoPublisher struct {
	BundleID *string `json:"bundleid,omitempty"`
	URL      *string `json:"url,omitempty"`
}

// CriteoPublisherBuilder allows to create `criteoPublisher` objects
type CriteoPublisherBuilder struct {
	bundleID *string
	url      *string
}

// NewCriteoPublisherBuilder creates a new builder instance
func NewCriteoPublisherBuilder() *CriteoPublisherBuilder {
	return &CriteoPublisherBuilder{}
}

// Slots set `BundleID` for the `criteoPublisher` to be created
func (b *CriteoPublisherBuilder) BundleID(bundleID *string) *CriteoPublisherBuilder {
	b.bundleID = bundleID
	return b
}

// Slots set `URL` for the `criteoPublisher` to be created
func (b *CriteoPublisherBuilder) URL(url *string) *CriteoPublisherBuilder {
	b.url = url
	return b
}

// Build returns a `criteoPublisher` instance based on the builder fields
func (b *CriteoPublisherBuilder) Build() *criteoPublisher {
	return &criteoPublisher{
		BundleID: b.bundleID,
		URL:      b.url,
	}
}

type criteoUser struct {
	DeviceID     *string `json:"deviceid,omitempty"`
	DeviceOs     *string `json:"deviceos,omitempty"`
	DeviceIDType *string `json:"deviceidtype,omitempty"`
	CookieID     *string `json:"cookieuid,omitempty"`
	IP           *string `json:"ip,omitempty"`
	Ua           *string `json:"ua,omitempty"`
}

// CriteoUserBuilder allows to create `criteoUser` objects
type CriteoUserBuilder struct {
	deviceID     *string
	deviceOs     *string
	deviceIDType *string
	cookieID     *string
	ip           *string
	ua           *string
}

// NewCriteoUserBuilder creates a new builder instance
func NewCriteoUserBuilder() *CriteoUserBuilder {
	return &CriteoUserBuilder{}
}

// DeviceID set `DeviceID` for the `criteoUser` to be created
func (b *CriteoUserBuilder) DeviceID(deviceID *string) *CriteoUserBuilder {
	b.deviceID = deviceID
	return b
}

// DeviceOs set `DeviceOs` for the `criteoUser` to be created
func (b *CriteoUserBuilder) DeviceOs(deviceOs *string) *CriteoUserBuilder {
	b.deviceOs = deviceOs
	b.deviceIDType = getDeviceType(deviceOs)
	return b
}

// CookieID set `CookieID` for the `criteoUser` to be created
func (b *CriteoUserBuilder) CookieID(cookieID *string) *CriteoUserBuilder {
	b.cookieID = cookieID
	return b
}

// IP set `IP` for the `criteoUser` to be created
func (b *CriteoUserBuilder) IP(ip *string) *CriteoUserBuilder {
	b.ip = ip
	return b
}

// Ua set `Ua` for the `criteoUser` to be created
func (b *CriteoUserBuilder) Ua(ua *string) *CriteoUserBuilder {
	b.ua = ua
	return b
}

// Build returns a `criteoUser` instance based on the builder fields
func (b *CriteoUserBuilder) Build() *criteoUser {
	return &criteoUser{
		CookieID:     b.cookieID,
		DeviceID:     b.deviceID,
		DeviceIDType: b.deviceIDType,
		DeviceOs:     b.deviceOs,
		IP:           b.ip,
		Ua:           b.ua,
	}
}

type criteoGdprConsent struct {
	GdprApplies *bool   `json:"gdprApplies,omitempty"`
	ConsentData *string `json:"consentData,omitempty"`
}

// CriteoGdprConsentBuilder allows to create `criteoGdprConsent` objects
type CriteoGdprConsentBuilder struct {
	gdprApplies *bool
	consentData *string
}

// NewCriteoGdprConsentBuilder creates a new builder instance
func NewCriteoGdprConsentBuilder() *CriteoGdprConsentBuilder {
	return &CriteoGdprConsentBuilder{}
}

// GdprApplies set `GdprApplies` for the `criteoGdprConsent` to be created
func (b *CriteoGdprConsentBuilder) GdprApplies(gdprApplies *bool) *CriteoGdprConsentBuilder {
	b.gdprApplies = gdprApplies
	return b
}

// Consent set `Consent` for the `criteoGdprConsent` to be created
func (b *CriteoGdprConsentBuilder) Consent(consentData *string) *CriteoGdprConsentBuilder {
	b.consentData = consentData
	return b
}

// Build returns a `criteoGdprConsent` instance based on the builder fields
func (b *CriteoGdprConsentBuilder) Build() *criteoGdprConsent {
	return &criteoGdprConsent{
		GdprApplies: b.gdprApplies,
		ConsentData: b.consentData,
	}
}

type criteoRequestSlot struct {
	ImpID  *string `json:"impid,omitempty"`
	ZoneID *uint   `json:"zoneid,omitempty"`
}

// CriteoRequestSlotBuilder allows to create `criteoRequestSlot` objects
type CriteoRequestSlotBuilder struct {
	impID  *string
	zoneID *uint
}

// NewCriteoRequestSlotBuilder creates a new builder instance
func NewCriteoRequestSlotBuilder() *CriteoRequestSlotBuilder {
	return &CriteoRequestSlotBuilder{}
}

// NewCriteoRequestSlotsBuilders creates a list of n `CriteoRequestSlotBuilders``
func NewCriteoRequestSlotsBuilders(size uint) []*CriteoRequestSlotBuilder {
	requestSlotsBuilders := make([]*CriteoRequestSlotBuilder, size)
	for i := uint(0); i < size; i++ {
		requestSlotsBuilders[i] = NewCriteoRequestSlotBuilder()
	}
	return requestSlotsBuilders
}

// ImpID set `ImpID` for the `criteoRequestSlot` to be created
func (b *CriteoRequestSlotBuilder) ImpID(impID *string) *CriteoRequestSlotBuilder {
	b.impID = impID
	return b
}

// ZoneID set `ZoneID` for the `criteoRequestSlot` to be created
func (b *CriteoRequestSlotBuilder) ZoneID(zoneID *uint) *CriteoRequestSlotBuilder {
	b.zoneID = zoneID
	return b
}

// Build returns a `criteoRequestSlot` instance based on the builder fields
func (b *CriteoRequestSlotBuilder) Build() *criteoRequestSlot {
	return &criteoRequestSlot{
		ImpID:  b.impID,
		ZoneID: b.zoneID,
	}
}

func getDeviceType(os *string) *string {
	deviceType := map[string]string{
		"ios":     "idfa",
		"android": "gaid",
		"unknown": "unknown",
	}

	if os != nil {
		dtype, ok := deviceType[strings.ToLower(*os)]
		if ok {
			return &dtype
		}
	}

	ret := deviceType["unknown"]
	return &ret
}

type criteoResponse struct {
	ID    *string               `json:"id,omitempty"`
	Slots []*criteoResponseSlot `json:"slots,omitempty"`
}

// NewCriteoResponseFromBytes creates a `criteoResponse` from JSON bytes
func NewCriteoResponseFromBytes(bytes []byte) (*criteoResponse, error) {
	var err error
	var bidResponse *criteoResponse

	err = json.Unmarshal(bytes, &bidResponse)
	if err != nil {
		return nil, err
	}

	return bidResponse, nil
}

// CriteoResponseBuilder allows to create `criteoResponse` objects
type CriteoResponseBuilder struct {
	id                          *string
	criteoResponseSlotsBuilders []*CriteoResponseSlotBuilder
}

// NewCriteoResponseBuilder creates a new builder instance
func NewCriteoResponseBuilder() *CriteoResponseBuilder {
	return &CriteoResponseBuilder{}
}

// ID set `ID` for the `criteoResponse` to be created
func (b *CriteoResponseBuilder) Id(id *string) *CriteoResponseBuilder {
	b.id = id
	return b
}

// Slots set `Slots` for the `criteoResponse` to be created
func (b *CriteoResponseBuilder) Slots(criteoResponseSlotsBuilders []*CriteoResponseSlotBuilder) *CriteoResponseBuilder {
	b.criteoResponseSlotsBuilders = criteoResponseSlotsBuilders
	return b
}

// Build returns a `criteoResponse` instance based on the builder fields
func (b *CriteoResponseBuilder) Build() *criteoResponse {
	criteoResponse := &criteoResponse{
		ID: b.id,
	}

	if b.criteoResponseSlotsBuilders != nil && len(b.criteoResponseSlotsBuilders) > 0 {
		criteoResponse.Slots = make([]*criteoResponseSlot, len(b.criteoResponseSlotsBuilders))
		for i := 0; i < len(b.criteoResponseSlotsBuilders); i++ {
			criteoResponse.Slots[i] = b.criteoResponseSlotsBuilders[i].Build()
		}
	}

	return criteoResponse
}

type criteoResponseSlot struct {
	ImpID    *string  `json:"impid,omitempty"`
	ZoneID   *uint    `json:"zoneid,omitempty"`
	Cpm      *float64 `json:"cpm,omitempty"`
	Currency *string  `json:"currency,omitempty"`
	Width    *uint    `json:"width,omitempty"`
	Height   *uint    `json:"height,omitempty"`
	Creative *string  `json:"creative,omitempty"`
}

// GetID returns the ID for the given `criteoResponseSlot`
func (s *criteoResponseSlot) GetID() *string {
	// TODO: This might be subject for change, maybe we should generate the ID on our end directly

	// Generate the ID from impression ID and zone ID
	w := bytes.NewBuffer([]byte(""))
	if s.ZoneID != nil {
		w.Write([]byte(strconv.Itoa(int(*s.ZoneID))))
		w.Write([]byte("-"))
	}
	if s.ImpID != nil {
		w.Write([]byte(*s.ImpID))
	}
	id := w.String()
	return &id
}

// GetCreativeID returns the creative ID for the given `criteoResponseSlot`
func (s *criteoResponseSlot) GetCreativeID() *string {
	// TODO: This might be subject for change, maybe we should generate the creative ID on our end directly

	// Generate the creative ID from width and height
	w := bytes.NewBuffer([]byte("CR-"))
	if s.Width != nil {
		w.Write([]byte(strconv.Itoa(int(*s.Width))))
		w.Write([]byte("x"))
	}
	if s.Height != nil {
		w.Write([]byte(strconv.Itoa(int(*s.Height))))
	}
	creativeID := w.String()
	return &creativeID
}

// CriteoResponseSlotBuilder allows to create `criteoResponseSlot` objects
type CriteoResponseSlotBuilder struct {
	impID    *string
	zoneID   *uint
	cpm      *float64
	currency *string
	width    *uint
	height   *uint
	creative *string
}

// NewCriteoResponseSlotBuilder creates a new builder instance
func NewCriteoResponseSlotBuilder() *CriteoResponseSlotBuilder {
	return &CriteoResponseSlotBuilder{}
}

// ImpID set `ImpID` for the `criteoResponseSlot` to be created
func (b *CriteoResponseSlotBuilder) ImpID(impID *string) *CriteoResponseSlotBuilder {
	b.impID = impID
	return b
}

// ZoneID set `ZoneID` for the `criteoResponseSlot` to be created
func (b *CriteoResponseSlotBuilder) ZoneID(zoneID *uint) *CriteoResponseSlotBuilder {
	b.zoneID = zoneID
	return b
}

// Cpm set `Cpm` for the `criteoResponseSlot` to be created
func (b *CriteoResponseSlotBuilder) Cpm(cpm *float64) *CriteoResponseSlotBuilder {
	b.cpm = cpm
	return b
}

// Currency set `Currency` for the `criteoResponseSlot` to be created
func (b *CriteoResponseSlotBuilder) Currency(currency *string) *CriteoResponseSlotBuilder {
	b.currency = currency
	return b
}

// Width set `Width` for the `criteoResponseSlot` to be created
func (b *CriteoResponseSlotBuilder) Width(width *uint) *CriteoResponseSlotBuilder {
	b.width = width
	return b
}

// Height set `Height` for the `criteoResponseSlot` to be created
func (b *CriteoResponseSlotBuilder) Height(height *uint) *CriteoResponseSlotBuilder {
	b.height = height
	return b
}

// Creative set `Creative` for the `criteoResponseSlot` to be created
func (b *CriteoResponseSlotBuilder) Creative(creative *string) *CriteoResponseSlotBuilder {
	b.creative = creative
	return b
}

// Build returns a `criteoResponseSlot` instance based on the builder fields
func (b *CriteoResponseSlotBuilder) Build() *criteoResponseSlot {
	return &criteoResponseSlot{
		ImpID:    b.impID,
		ZoneID:   b.zoneID,
		Cpm:      b.cpm,
		Currency: b.currency,
		Width:    b.width,
		Height:   b.height,
		Creative: b.creative,
	}
}
