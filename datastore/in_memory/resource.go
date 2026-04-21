package in_memory

import (
	"context"

	"github.com/HackRVA/memberserver/models"
)

// ensureResources initializes the map if it's nil. Caller must hold the write lock.
func (store *In_memory) ensureResources() {
	if store.resources == nil {
		store.resources = map[string]models.Resource{}
	}
}

func (store *In_memory) GetResources(ctx context.Context) []models.Resource {
	store.mu.RLock()
	defer store.mu.RUnlock()
	resources := []models.Resource{}
	for _, v := range store.resources {
		resources = append(resources, v)
	}
	return resources
}

func (store *In_memory) GetResourceACL(ctx context.Context, r models.Resource) ([]string, error) {
	return []string{}, nil
}

func (store *In_memory) GetResourceACLWithMemberInfo(ctx context.Context, r models.Resource) ([]models.Member, error) {
	return []models.Member{{
		ID:   "123",
		Name: "test",
	}}, nil
}

func (store *In_memory) GetMembersAccess(ctx context.Context, m models.Member) ([]models.MemberAccess, error) {
	return []models.MemberAccess{}, nil
}

func (store *In_memory) GetInactiveMembersByResource(ctx context.Context) ([]models.MemberAccess, error) {
	return []models.MemberAccess{}, nil
}

func (store *In_memory) GetActiveMembersByResource(ctx context.Context) ([]models.MemberAccess, error) {
	return []models.MemberAccess{}, nil
}

func (store *In_memory) RegisterResource(ctx context.Context, name string, address string, isDefault bool) (models.Resource, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.ensureResources()
	store.resources[name] = models.Resource{
		Name:      name,
		Address:   address,
		IsDefault: isDefault,
	}

	return store.resources[name], nil
}

func (store *In_memory) GetResourceByID(ctx context.Context, ID string) (models.Resource, error) {
	return models.Resource{}, nil
}

func (store *In_memory) GetResourceByName(ctx context.Context, resourceName string) (models.Resource, error) {
	return models.Resource{}, nil
}

func (store *In_memory) UpdateResource(ctx context.Context, res models.Resource) (*models.Resource, error) {
	return &models.Resource{}, nil
}

func (store *In_memory) DeleteResource(ctx context.Context, id string) error {
	return nil
}

func (store *In_memory) AddMultipleMembersToResource(ctx context.Context, emails []string, resourceID string) ([]models.MemberResourceRelation, error) {
	return []models.MemberResourceRelation{}, nil
}

func (store *In_memory) AddUserToDefaultResources(ctx context.Context, email string) ([]models.MemberResourceRelation, error) {
	return []models.MemberResourceRelation{}, nil
}

func (store *In_memory) RemoveUserFromResource(ctx context.Context, email string, resourceID string) error {
	return nil
}
