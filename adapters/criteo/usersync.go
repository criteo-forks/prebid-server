package criteo

import "text/template"
import "github.com/prebid/prebid-server/adapters"
import "github.com/prebid/prebid-server/usersync"

// Criteo doesn't need user synchronization yet. The first version of
// the adapter is meant to work with Apps only.
func NewCriteoSyncer(temp *template.Template) usersync.Usersyncer {
	return adapters.NewSyncer("criteo", 91, temp, adapters.SyncTypeIframe)
}
