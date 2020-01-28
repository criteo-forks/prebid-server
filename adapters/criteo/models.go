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
	if publisher := newCriteoPublisher(request.App, request.Site); publisher != nil {
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

	// Extracting request slots
	if request.Imp != nil && len(request.Imp) > 0 {
		criteoSlots, slotsErr := newCriteoRequestSlots(request.Imp)
		if slotsErr != nil && len(slotsErr) > 0 {
			return nil, append(errs, slotsErr...)
		}
		criteoRequest.Slots = criteoSlots
	}

	return criteoRequest, errs
}

type criteoPublisher struct {
	SiteID   *string `json:"siteid,omitempty"` // TODO: make sure it's siteid and not publisherid
	BundleID *string `json:"bundleid,omitempty"`
	URL      *string `json:"url,omitempty"`
}

func newCriteoPublisher(app *openrtb.App, site *openrtb.Site) *criteoPublisher {
	if app == nil && site == nil {
		return nil
	}

	criteoPublisher := &criteoPublisher{}

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
	ImpID  *string `json:"impid,omitempty"`
	ZoneID *uint   `json:"zoneid,omitempty"`
}

func newCriteoRequestSlots(impressions []openrtb.Imp) ([]*criteoRequestSlot, []error) {
	var errs []error

	if impressions == nil || len(impressions) == 0 {
		return nil, []error{&errortypes.BadInput{
			Message: "Bid request impressions is nil or empty",
		}}
	}

	var criteoSlots = make([]*criteoRequestSlot, uint(len(impressions)))

	for i := 0; i < len(impressions); i++ {
		criteoSlots[i] = &criteoRequestSlot{}

		if impressions[i].ID != "" {
			criteoSlots[i].ImpID = &impressions[i].ID
		}

		var bidderExt adapters.ExtImpBidder
		json.Unmarshal(impressions[i].Ext, &bidderExt)
		if bidderExt.Bidder != nil {
			var criteoExt openrtb_ext.ExtImpCriteo
			json.Unmarshal(bidderExt.Bidder, &criteoExt)
			if criteoExt.ZoneID > 0 {
				criteoSlots[i].ZoneID = &criteoExt.ZoneID
			}
		}
	}

	return criteoSlots, errs
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
	CPM        *float64 `json:"cpm,omitempty"`
	Currency   *string  `json:"currency,omitempty"`
	Width      *uint    `json:"width,omitempty"`
	Height     *uint    `json:"height,omitempty"`
	Creative   *string  `json:"creative,omitempty"`
	CreativeID *string  `json:"creativeid,omitempty"`
}
