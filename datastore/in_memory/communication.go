package in_memory

import (
	"context"
	"time"

	"github.com/HackRVA/memberserver/models"
)

func (i *In_memory) GetCommunications(ctx context.Context) []models.Communication {
	return []models.Communication{}
}

func (i *In_memory) GetCommunication(ctx context.Context, name string) (models.Communication, error) {
	return models.Communication{}, nil
}

func (i *In_memory) GetMostRecentCommunicationToMember(ctx context.Context, memberId string, commId int) (time.Time, error) {
	return time.Time{}, nil
}

func (i *In_memory) LogCommunication(ctx context.Context, communicationId int, memberId string) error {
	return nil
}
