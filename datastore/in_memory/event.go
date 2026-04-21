package in_memory

import (
	"context"

	"github.com/HackRVA/memberserver/models"
)

func (store *In_memory) LogAccessEvent(ctx context.Context, event models.LogMessage) error {
	return nil
}
