package criteo

import (
	"encoding/json"
	"net/http"
	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/openrtb_ext"
)

// TODO: (26/04/2018) The profileId is temporary and should be updated soon
const CRITEO_URL = "https://bidder.criteo.com/cdb?profileId=154"

type CriteoAdapter struct {}

// Criteo bidder request and response structures
type CriteoRequest struct {
	Publisher CriteoPublisher     `json:"publisher"`
	User      CriteoUser          `json:"user"`
	Slots     []CriteoRequestSlot `json:"slots"`
}

type CriteoResponse struct {
	Slots []CriteoResponseSlot `json:"slots"`
	Exd   CriteoExtension      `json:"exd"`
}

type CriteoPublisher struct {
	BundleId string `json:"bundleid"`
}

type CriteoUser struct {
	DeviceId     string `json:"deviceid"`
	DeviceOs     string `json:"deviceos"`
	DeviceIdType string `json:"deviceidtype"`
}

type CriteoRequestSlot struct {
	ImpId  string `json:"impid"`
	ZoneId uint   `json:"zoneid"`
}

type CriteoResponseSlot struct {
	ImpId    string  `json:"impid"`
	ZoneId   uint    `json:"zoneid"`
	Cpm      float64 `json:"cpm"`
	Currency string  `json:"currency"`
	Width    uint    `json:"width"`
	Height   uint    `json:"height"`
	Creative string  `json:"creative"`
}

type CriteoExtension struct {
	GenImpIds []string              `json:"genimpids"`
	Slots     []CriteoExtensionSlot `json:"slots"`
}

type CriteoExtensionSlot struct {
	GenImpId string `json:"imp_id"`
	ImpId    string `json:"ad_unit_id"`
	ZoneId   uint   `json:"zone_id"`
}

func (a *CriteoAdapter) MakeRequests(request *openrtb.BidRequest) ([]*adapters.RequestData, []error) {
	rqjson, err := json.Marshal(CriteoRequest{
		Publisher: CriteoPublisher{
			BundleId: request.App.Bundle,
		},
		User: CriteoUser{
			DeviceId: request.Device.IFA,
			DeviceOs: request.Device.OS,
			DeviceIdType: getDeviceType(request.Device.OS),
		},
		Slots: getRequestSlots(request.Imp),
	})
	if err != nil {
		return nil, []error{err}
	}

	rqData := adapters.RequestData{
		Method: "POST",
		Uri: CRITEO_URL,
		Body: rqjson,
		Headers: http.Header{},
	}
	return []*adapters.RequestData{&rqData}, nil
}

func getDeviceType(os string) string {
	deviceType := map[string]string{
		"ios": "idfa",
		"android": "gaid",
	}

	dtype, ok := deviceType[os]
	if !ok {
		return "unknown_device_type";
	}
	return dtype
}

func getRequestSlots(impressions []openrtb.Imp) []CriteoRequestSlot {
	var slots []CriteoRequestSlot

	for _, imp := range(impressions) {
		var bidderExt adapters.ExtImpBidder
		var criteoExt openrtb_ext.ExtImpCriteo

		// TODO - Handle unmarshalling errors
		json.Unmarshal(imp.Ext, &bidderExt)
		json.Unmarshal(bidderExt.Bidder, &criteoExt)

		slots = append(slots, CriteoRequestSlot{
			ImpId: imp.ID,
			ZoneId: criteoExt.ZoneId,
		})
	}
	return slots
}

func (a *CriteoAdapter) MakeBids(
	internalRequest *openrtb.BidRequest,
	externalRequest *adapters.RequestData,
	response *adapters.ResponseData,
)(
	[]*adapters.TypedBid,
	[]error,
){
	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	// TODO = Handle more bad response.StatusCode

	var bidResp CriteoResponse
	// TODO - Handle unmarshalling errors
	if err := json.Unmarshal(response.Body, &bidResp); err != nil {
		return nil, []error{err}
	}

	// map[ImpId] => GenImpId
	var genImpId = make(map[string]string)
	for _, slot := range(bidResp.Exd.Slots) {
		genImpId[slot.ImpId] = slot.GenImpId
	}

	var bids []*adapters.TypedBid
	// TODO - support native bids (openrtb_ext.BidTypeNative)
	for _, slot := range(bidResp.Slots) {
		bid := adapters.TypedBid{
			Bid: &openrtb.Bid{
				ID: genImpId[slot.ImpId],
				ImpID: slot.ImpId,
				Price: slot.Cpm,
				AdM: slot.Creative,
				W: uint64(slot.Width),
				H: uint64(slot.Height),
			},
			BidType: openrtb_ext.BidTypeBanner,
		}
		bids = append(bids, &bid)
	}
	return bids, nil
}

func NewCriteoBidder() *CriteoAdapter {
	return &CriteoAdapter{}
}
