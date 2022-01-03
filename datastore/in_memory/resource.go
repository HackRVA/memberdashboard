package in_memory

import "memberserver/api/models"

var Resources = map[string]models.Resource{}

func (store *In_memory) GetResources() []models.Resource {
	resources := []models.Resource{}
	for _, v := range Resources {
		resources = append(resources, v)
	}
	return resources
}

func (store *In_memory) GetResourceACL(models.Resource) ([]string, error) {
	return []string{}, nil
}

func (store *In_memory) GetResourceACLWithMemberInfo(models.Resource) ([]models.Member, error) {
	return []models.Member{{
		ID:   "123",
		Name: "test",
	}}, nil
}

func (store *In_memory) GetMembersAccess(models.Member) ([]models.MemberAccess, error) {
	return []models.MemberAccess{}, nil
}
func (store *In_memory) GetInactiveMembersByResource() ([]models.MemberAccess, error) {
	return []models.MemberAccess{}, nil
}

func (store *In_memory) RegisterResource(name string, address string, isDefault bool) (models.Resource, error) {
	Resources[name] = models.Resource{
		Name:      name,
		Address:   address,
		IsDefault: isDefault,
	}

	return Resources[name], nil
}

func (store *In_memory) GetResourceByID(ID string) (models.Resource, error) {
	return models.Resource{}, nil
}
func (store *In_memory) GetResourceByName(resourceName string) (models.Resource, error) {
	return models.Resource{}, nil
}
func (store *In_memory) UpdateResource(id string, name string, address string, isDefault bool) (*models.Resource, error) {
	return &models.Resource{}, nil
}
func (store *In_memory) DeleteResource(id string) error {
	return nil
}
func (store *In_memory) AddMultipleMembersToResource(emails []string, resourceID string) ([]models.MemberResourceRelation, error) {
	return []models.MemberResourceRelation{}, nil
}
func (store *In_memory) AddUserToDefaultResources(email string) ([]models.MemberResourceRelation, error) {
	return []models.MemberResourceRelation{}, nil
}
func (store *In_memory) GetMemberResourceRelation(m models.Member, r models.Resource) (models.MemberResourceRelation, error) {
	return models.MemberResourceRelation{}, nil
}
func (store *In_memory) RemoveUserFromResource(email string, resourceID string) error {
	return nil
}
