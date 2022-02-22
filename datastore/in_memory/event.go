package in_memory

import "memberserver/api/models"

func (store *In_memory) LogAccessEvent(event models.LogMessage) error {
	return nil
}
