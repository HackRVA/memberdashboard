package in_memory

import (
	"time"

	"github.com/HackRVA/memberserver/internal/models"
)

func (i *In_memory) GetCommunications() []models.Communication {
	return []models.Communication{}
}
func (i *In_memory) GetCommunication(name string) (models.Communication, error) {
	return models.Communication{}, nil
}
func (i *In_memory) GetMostRecentCommunicationToMember(memberId string, commId int) (time.Time, error) {
	return time.Time{}, nil
}
func (i *In_memory) LogCommunication(communicationId int, memberId string) error {
	return nil
}
