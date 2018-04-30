package usersync

// Criteo doesn't need user synchronization yet. The first version of
// the adapter is meant to work with Apps only.
func NewCriteoSyncer() Usersyncer {
	return &syncer{
		familyName: "criteo",
		syncInfo: &UsersyncInfo{
			URL:         "",
			Type:        "",
			SupportCORS: false,
		},
	}
}
