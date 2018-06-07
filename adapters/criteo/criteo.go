package criteo

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type CriteoAdapter struct {
	URI  string
	http *adapters.HTTPAdapter
}

func (a *CriteoAdapter) MakeRequests(request *openrtb.BidRequest) ([]*adapters.RequestData, []error) {

	criteoRequestBuilder, errs := getCriteoRequest(request)
	if errs != nil || len(errs) > 0 {
		return nil, errs
	}

	criteoRequest := criteoRequestBuilder.Build()

	jsonRequest, err := json.Marshal(criteoRequest)
	if err != nil {
		return nil, []error{err}
	}

	rqData := adapters.RequestData{
		Method:  "POST",
		Uri:     a.URI,
		Body:    jsonRequest,
		Headers: getCriteoRequestHeaders(criteoRequest),
	}

	return []*adapters.RequestData{&rqData}, nil
}

func getCriteoRequestHeaders(criteoRequest *criteoRequest) http.Header {
	headers := http.Header{}

	if criteoRequest != nil {
		if criteoRequest.User != nil {
			if criteoRequest.User.CookieID != nil {
				buf := bytes.NewBuffer([]byte("uid="))
				buf.Write([]byte(*criteoRequest.User.CookieID))
				headers.Add("Cookie", buf.String())
			}

			if criteoRequest.User.IP != nil {
				headers.Add("X-Client-IP", *criteoRequest.User.IP)
			}

			if criteoRequest.User.Ua != nil {
				headers.Add("User-Agent", *criteoRequest.User.Ua)
			}
		}
	}

	return headers
}

func getCriteoRequest(request *openrtb.BidRequest) (*CriteoRequestBuilder, []error) {
	var errs []error

	if request == nil {
		return nil, []error{&adapters.BadInputError{
			Message: "Bid request is nil",
		}}
	}

	var criteoReqBuilder = NewCriteoRequestBuilder()

	criteoReqBuilder.ID(&request.ID)

	criteoReqBuilder.Publisher(getPublisher(request.App, request.Site))

	criteoReqBuilder.User(getUser(request.User, request.Device))

	criteoReqBuilder.GdprConsent(getGdprConsent(request.User, request.Regs))

	// Extracting request slots
	if request.Imp != nil && len(request.Imp) > 0 {
		slotsBuilders, slotsErr := getRequestSlots(request.Imp)
		if slotsErr != nil && len(slotsErr) > 0 {
			return nil, append(errs, slotsErr...)
		}
		criteoReqBuilder.Slots(slotsBuilders)
	}

	return criteoReqBuilder, errs
}

func getGdprConsent(user *openrtb.User, regs *openrtb.Regs) *CriteoGdprConsentBuilder {

	if user == nil && regs == nil {
		return nil
	}

	gdprConsentBuilder := &CriteoGdprConsentBuilder{}

	if user != nil {
		if user.Ext != nil {
			var userExt *openrtb_ext.ExtUser
			json.Unmarshal(user.Ext, &userExt)
			if userExt.Consent != "" {
				gdprConsentBuilder.Consent(&userExt.Consent)
			}
		}
	}

	if regs != nil {
		if regs.Ext != nil {
			var regsExt *openrtb_ext.ExtRegs
			json.Unmarshal(regs.Ext, &regsExt)
			if regsExt.GDPR != nil {
				gdprApplies := bool((*regsExt.GDPR & 1) == 1)
				gdprConsentBuilder.GdprApplies(&gdprApplies)
			}
		}
	}

	return gdprConsentBuilder
}

func getUser(user *openrtb.User, device *openrtb.Device) *CriteoUserBuilder {

	if user == nil && device == nil {
		return nil
	}

	userBuilder := &CriteoUserBuilder{}

	if user != nil {
		if user.BuyerUID != "" {
			userBuilder.CookieID(&user.BuyerUID)
		}
	}

	if device != nil {
		if device.IFA != "" {
			userBuilder.DeviceID(&device.IFA)
		}
		if device.OS != "" {
			userBuilder.DeviceOs(&device.OS)
		}
		if device.IP != "" {
			userBuilder.IP(&device.IP)
		}
		if device.UA != "" {
			userBuilder.Ua(&device.UA)
		}
	}

	return userBuilder
}

func getPublisher(app *openrtb.App, site *openrtb.Site) *CriteoPublisherBuilder {

	if app == nil && site == nil {
		return nil
	}

	publisherBuilder := &CriteoPublisherBuilder{}

	if app != nil {
		if app.Bundle != "" {
			publisherBuilder.BundleID(&app.Bundle)
		}
	}

	if site != nil {
		if site.Page != "" {
			publisherBuilder.URL(&site.Page)
		}
	}

	return publisherBuilder
}

func getRequestSlots(impressions []openrtb.Imp) ([]*CriteoRequestSlotBuilder, []error) {

	var errs []error

	if impressions == nil || len(impressions) == 0 {
		return nil, []error{&adapters.BadInputError{
			Message: "Bid request impressions is nil or empty",
		}}
	}

	var slotsBuilders = NewCriteoRequestSlotsBuilders(uint(len(impressions)))

	for i := 0; i < len(impressions); i++ {

		if impressions[i].ID != "" {
			slotsBuilders[i].ImpID(&impressions[i].ID)
		}

		var bidderExt adapters.ExtImpBidder
		json.Unmarshal(impressions[i].Ext, &bidderExt)
		if bidderExt.Bidder != nil {
			var criteoExt openrtb_ext.ExtImpCriteo
			json.Unmarshal(bidderExt.Bidder, &criteoExt)
			if criteoExt.ZoneID > 0 {
				slotsBuilders[i].ZoneID(&criteoExt.ZoneID)
			}
		}
	}

	return slotsBuilders, errs
}

func (a *CriteoAdapter) MakeBids(
	internalRequest *openrtb.BidRequest,
	externalRequest *adapters.RequestData,
	response *adapters.ResponseData,
) (
	*adapters.BidderResponse,
	[]error,
) {
	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	// TODO = Handle more bad response.StatusCode

	bidResponse, err := NewCriteoResponseFromBytes(response.Body)
	if err != nil {
		return nil, []error{err}
	}

	bidderResponse := adapters.NewBidderResponse()
	bidderResponse.Bids = make([]*adapters.TypedBid, len(bidResponse.Slots))

	// TODO - support native bids (openrtb_ext.BidTypeNative)
	for i := 0; i < len(bidResponse.Slots); i++ {
		bidderResponse.Bids[i] = &adapters.TypedBid{
			Bid: &openrtb.Bid{
				ID:    *bidResponse.Slots[i].GetID(),
				ImpID: *bidResponse.Slots[i].ImpID,
				Price: *bidResponse.Slots[i].Cpm,
				AdM:   *bidResponse.Slots[i].Creative,
				W:     uint64(*bidResponse.Slots[i].Width),
				H:     uint64(*bidResponse.Slots[i].Height),
				CrID:  *bidResponse.Slots[i].GetCreativeID(),
			},
			BidType: openrtb_ext.BidTypeBanner,
		}
	}
	return bidderResponse, nil
}

func NewCriteoBidder(client *http.Client, endpoint string) *CriteoAdapter {
	a := &adapters.HTTPAdapter{Client: client}
	return &CriteoAdapter{
		http: a,
		URI:  endpoint,
	}
}
