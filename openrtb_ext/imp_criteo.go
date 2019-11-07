package openrtb_ext

// ExtImpCriteo defines the contract for bidrequest.imp[i].ext.criteo
type ExtImpCriteo struct {
	PlacementID string `json:"placement"`
	ZoneID      uint   `json:"zoneId"`
	NetworkID   uint   `json:"networkId"`
}
