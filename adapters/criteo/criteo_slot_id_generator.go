package criteo

import "github.com/gofrs/uuid"

type slotIDGenerator interface {
	NewSlotID() (string, error)
}

type criteoSlotIDGenerator struct{}

func newCriteoSlotIDGenerator() criteoSlotIDGenerator {
	return criteoSlotIDGenerator{}
}

func (g criteoSlotIDGenerator) NewSlotID() (string, error) {
	guid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return guid.String(), nil
}
