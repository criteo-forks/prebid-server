package criteo

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/config"
	"github.com/prebid/prebid-server/openrtb_ext"
)

type CriteoAdapter struct {
	URI string
}

func (a *CriteoAdapter) MakeRequests(request *openrtb.BidRequest, extraRequestInfo *adapters.ExtraRequestInfo) ([]*adapters.RequestData, []error) {
	criteoRequest, errs := newCriteoRequest(request)
	if errs != nil || len(errs) > 0 {
		return nil, errs
	}

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
				headers.Add("X-Client-Ip", *criteoRequest.User.IP)
			}

			if criteoRequest.User.UA != nil {
				headers.Add("User-Agent", *criteoRequest.User.UA)
			}
		}
	}

	return headers
}

func (a *CriteoAdapter) MakeBids(internalRequest *openrtb.BidRequest, externalRequest *adapters.RequestData, response *adapters.ResponseData) (*adapters.BidderResponse, []error) {
	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	// TODO = Handle more bad response.StatusCode
	bidResponse, err := newCriteoResponseFromBytes(response.Body)
	if err != nil {
		return nil, []error{err}
	}

	bidderResponse := adapters.NewBidderResponse()
	bidderResponse.Bids = make([]*adapters.TypedBid, len(bidResponse.Slots))

	// TODO - support native bids (openrtb_ext.BidTypeNative)
	for i := 0; i < len(bidResponse.Slots); i++ {
		bidderResponse.Bids[i] = &adapters.TypedBid{
			Bid: &openrtb.Bid{
				ID:    *bidResponse.Slots[i].ID,
				ImpID: *bidResponse.Slots[i].ImpID,
				Price: *bidResponse.Slots[i].CPM,
				AdM:   *bidResponse.Slots[i].Creative,
				W:     uint64(*bidResponse.Slots[i].Width),
				H:     uint64(*bidResponse.Slots[i].Height),
				CrID:  *bidResponse.Slots[i].CreativeID,
			},
			BidType: openrtb_ext.BidTypeBanner,
		}
	}

	return bidderResponse, nil
}

func Builder(bidderName openrtb_ext.BidderName, config config.Adapter) (adapters.Bidder, error) {
	return &CriteoAdapter{
		URI: config.Endpoint,
	}, nil
}
