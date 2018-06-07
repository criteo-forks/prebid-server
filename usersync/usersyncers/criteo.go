package usersyncers

// Criteo doesn't need user synchronization yet. The first version of
// the adapter is meant to work with Apps only.
func NewCriteoSyncer(usersyncURL string, externalURL string) *syncer {

	return &syncer{
		familyName:   "criteo",
		gdprVendorID: 91,
		syncType:     "",
		syncEndpointBuilder: func(gdpr string, consent string) string {
			return ""
		},
	}
}
