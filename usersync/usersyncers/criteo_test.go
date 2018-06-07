package usersyncers

import (
	"testing"
)

func TestCriteoSyncer(t *testing.T) {
	syncer := NewCriteoSyncer("", "")
	info := syncer.GetUsersyncInfo("", "")
	if info.URL != "" {
		t.Fatalf("should be an empty string but was: %s", info.URL)
	}
	if info.Type != "" {
		t.Fatalf("should be an empty string but was: %s", info.Type)
	}
	if info.SupportCORS != false {
		t.Fatalf("should have been false")
	}
	if syncer.FamilyName() != "criteo" {
		t.Errorf("FamilyName '%s' != 'criteo'", syncer.FamilyName())
	}
}
