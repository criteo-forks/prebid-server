package criteo

import (
	"encoding/json"
	"strings"

	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type criteoRequest struct {
	ID          *string              `json:"id,omitempty"`
	Publisher   *criteoPublisher     `json:"publisher,omitempty"`
	User        *criteoUser          `json:"user,omitempty"`
	GdprConsent *criteoGdprConsent   `json:"gdprconsent,omitempty"`
	Slots       []*criteoRequestSlot `json:"slots,omitempty"`
}

func newCriteoRequest(request *openrtb.BidRequest) (*criteoRequest, []error) {
	var errs []error

	if request == nil {
		return nil, []error{&errortypes.BadInput{
			Message: "Bid request is nil",
		}}
	}

	criteoRequest := &criteoRequest{}

	if request.ID != "" {
		criteoRequest.ID = &request.ID
	}

	// Extracting request slots
	if request.Imp != nil && len(request.Imp) > 0 {
		criteoSlots, slotsErr := newCriteoRequestSlots(request.Imp)
		if slotsErr != nil && len(slotsErr) > 0 {
			return nil, append(errs, slotsErr...)
		}
		criteoRequest.Slots = criteoSlots
	}

	var networkId *uint
	for _, criteoSlot := range criteoRequest.Slots {
		if networkId == nil && criteoSlot.NetworkID != nil && *criteoSlot.NetworkID > 0 {
			networkId = criteoSlot.NetworkID
		} else if networkId != nil && criteoSlot.NetworkID != nil && *criteoSlot.NetworkID != *networkId {
			return nil, []error{&errortypes.BadInput{
				Message: "Bid request has slots coming with several network IDs which is not allowed",
			}}
		}
	}

	if publisher := newCriteoPublisher(networkId, request.App, request.Site); publisher != nil {
		criteoRequest.Publisher = publisher
	}

	var regsExt *openrtb_ext.ExtRegs
	if request.Regs != nil {
		json.Unmarshal(request.Regs.Ext, &regsExt)
	}

	if user := newCriteoUser(request.User, request.Device, regsExt); user != nil {
		criteoRequest.User = user
	}
	if gdprConsent := newCriteoGdprConsent(request.User, regsExt); gdprConsent != nil {
		criteoRequest.GdprConsent = gdprConsent
	}

	return criteoRequest, errs
}

type criteoPublisher struct {
	SiteID    *string `json:"siteid,omitempty"` // TODO: make sure it's siteid and not publisherid
	BundleID  *string `json:"bundleid,omitempty"`
	URL       *string `json:"url,omitempty"`
	NetworkID *uint   `json:"networkid,omitempty"`
}

func newCriteoPublisher(networkId *uint, app *openrtb.App, site *openrtb.Site) *criteoPublisher {
	if app == nil && site == nil {
		return nil
	}

	criteoPublisher := &criteoPublisher{}

	if networkId != nil && *networkId > 0 {
		criteoPublisher.NetworkID = networkId
	}

	if app != nil {
		if app.Bundle != "" {
			criteoPublisher.BundleID = &app.Bundle
		}
	}

	if site != nil {
		if site.ID != "" {
			criteoPublisher.SiteID = &site.ID
		}
		if site.Page != "" {
			criteoPublisher.URL = &site.Page
		}
	}

	return criteoPublisher
}

type criteoUserExt struct {
	Eids []openrtb_ext.ExtUserEid `json:"eids,omitempty"`
}

type criteoUser struct {
	DeviceID     *string `json:"deviceid,omitempty"`
	DeviceOs     *string `json:"deviceos,omitempty"`
	DeviceIDType *string `json:"deviceidtype,omitempty"`
	CookieID     *string `json:"cookieuid,omitempty"`
	UID          *string `json:"uid,omitempty"`
	IP           *string `json:"ip,omitempty"`
	UA           *string `json:"ua,omitempty"`
	UspIab       *string `json:"uspIab,omitempty"`
}

func newCriteoUser(user *openrtb.User, device *openrtb.Device, regsExt *openrtb_ext.ExtRegs) *criteoUser {
	if user == nil && device == nil {
		return nil
	}

	criteoUser := &criteoUser{}

	if user != nil {
		if user.BuyerUID != "" {
			criteoUser.CookieID = &user.BuyerUID
		}
		if user.Ext != nil {
			var criteoUserExt *criteoUserExt
			json.Unmarshal(user.Ext, &criteoUserExt)
			if criteoUserExt.Eids != nil && len(criteoUserExt.Eids) > 0 {
				for _, eid := range criteoUserExt.Eids {
					if eid.Source == "criteo.com" && len(eid.Uids) > 0 {
						criteoUser.UID = &eid.Uids[0].ID
						break
					}
				}
			}
		}
	}

	if device != nil {
		deviceType := getDeviceType(device.OS)
		criteoUser.DeviceIDType = &deviceType

		if device.OS != "" {
			criteoUser.DeviceOs = &device.OS
		}
		if device.IFA != "" {
			criteoUser.DeviceID = &device.IFA
		}
		if device.IP != "" {
			criteoUser.IP = &device.IP
		}
		if device.UA != "" {
			criteoUser.UA = &device.UA
		}
	}

	if regsExt != nil {
		if regsExt.USPrivacy != "" {
			criteoUser.UspIab = &regsExt.USPrivacy // CCPA
		}
	}

	return criteoUser
}

type criteoGdprConsent struct {
	GdprApplies *bool   `json:"gdprapplies,omitempty"`
	ConsentData *string `json:"consentdata,omitempty"`
}

func newCriteoGdprConsent(user *openrtb.User, regsExt *openrtb_ext.ExtRegs) *criteoGdprConsent {
	if user == nil && regsExt == nil {
		return nil
	}

	var consent *criteoGdprConsent
	if user != nil {
		if user.Ext != nil {
			var userExt *openrtb_ext.ExtUser
			json.Unmarshal(user.Ext, &userExt)
			if userExt.Consent != "" {
				if consent == nil {
					consent = &criteoGdprConsent{}
				}
				consent.ConsentData = &userExt.Consent
			}
		}
	}

	if regsExt != nil {
		if regsExt.GDPR != nil {
			if consent == nil {
				consent = &criteoGdprConsent{}
			}
			gdprApplies := bool((*regsExt.GDPR & 1) == 1)
			consent.GdprApplies = &gdprApplies
		}
	}

	return consent
}

type criteoRequestSlot struct {
	ImpID       *string              `json:"impid,omitempty"`
	ZoneID      *uint                `json:"zoneid,omitempty"`
	NetworkID   *uint                `json:"networkid,omitempty"`
	PlacementID *string              `json:"placement,omitempty"`
	Sizes       []*criteoRequestSize `json:"sizes,omitempty"`
}

func newCriteoRequestSlots(impressions []openrtb.Imp) ([]*criteoRequestSlot, []error) {
	var errs []error

	if impressions == nil || len(impressions) == 0 {
		return nil, []error{&errortypes.BadInput{
			Message: "Bid request impressions is nil or empty",
		}}
	}

	// TODO: Criteo slot should comes either with:
	//   - `zoneid`
	// OR
	//   - `networkid`, `placementid`, `sizes`
	//
	// if not, criteo will reject the slot.
	// Would be nice preventing PBS to send such slots when conditions aren't met

	var criteoSlots = make([]*criteoRequestSlot, uint(len(impressions)))

	for i := 0; i < len(impressions); i++ {
		criteoSlots[i] = &criteoRequestSlot{}

		if impressions[i].ID != "" {
			criteoSlots[i].ImpID = &impressions[i].ID
		}

		if impressions[i].Banner != nil {
			if impressions[i].Banner.Format != nil {
				criteoSlots[i].Sizes = make([]*criteoRequestSize, len(impressions[i].Banner.Format))
				for idx, format := range impressions[i].Banner.Format {
					// TODO: handle properly uint conversion
					criteoSlots[i].Sizes[idx] = newCriteoRequestSize(uint(format.W), uint(format.H))
				}
			}
		}

		var bidderExt adapters.ExtImpBidder
		json.Unmarshal(impressions[i].Ext, &bidderExt)
		if bidderExt.Bidder != nil {
			var criteoExt openrtb_ext.ExtImpCriteo
			json.Unmarshal(bidderExt.Bidder, &criteoExt)
			if criteoExt.ZoneID > 0 {
				criteoSlots[i].ZoneID = &criteoExt.ZoneID
			}
			if criteoExt.NetworkID > 0 {
				criteoSlots[i].NetworkID = &criteoExt.NetworkID
			}
			if criteoExt.PlacementID != "" {
				criteoSlots[i].PlacementID = &criteoExt.PlacementID
			}
		}
	}

	return criteoSlots, errs
}

type criteoRequestSize struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}

func newCriteoRequestSize(width uint, height uint) *criteoRequestSize {
	return &criteoRequestSize{
		Width:  width,
		Height: height,
	}
}

func getDeviceType(os string) string {
	deviceType := map[string]string{
		"ios":     "idfa",
		"android": "gaid",
		"unknown": "unknown",
	}

	if os != "" {
		dtype, ok := deviceType[strings.ToLower(os)]
		if ok {
			return dtype
		}
	}

	return deviceType["unknown"]
}

type criteoResponse struct {
	ID    *string               `json:"id,omitempty"`
	Slots []*criteoResponseSlot `json:"slots,omitempty"`
}

func newCriteoResponseFromBytes(bytes []byte) (*criteoResponse, error) {
	var err error
	var bidResponse *criteoResponse

	err = json.Unmarshal(bytes, &bidResponse)
	if err != nil {
		return nil, err
	}

	return bidResponse, nil
}

type criteoResponseSlot struct {
	ID         *string  `json:"id,omitempty"`
	ImpID      *string  `json:"impid,omitempty"`
	ZoneID     *uint    `json:"zoneid,omitempty"`
	NetworkID  *uint    `json:"networkid,omitempty"`
	CPM        *float64 `json:"cpm,omitempty"`
	Currency   *string  `json:"currency,omitempty"`
	Width      *uint    `json:"width,omitempty"`
	Height     *uint    `json:"height,omitempty"`
	Creative   *string  `json:"creative,omitempty"`
	CreativeID *string  `json:"creativeid,omitempty"`
}
