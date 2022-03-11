package in_memory

import "memberserver/internal/models"

func (store *In_memory) LogAccessEvent(event models.LogMessage) error {
	return nil
}
