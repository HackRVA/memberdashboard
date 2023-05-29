package in_memory

import (
	"errors"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"
)

func (i *In_memory) GetUser(email string) (models.UserResponse, error) {

	for _, k := range i.Members {
		if k.Email == email {
			resources := []models.Resource{}
			for _, r := range k.Resources {
				resources = append(resources, models.Resource{
					Name: r.Name,
					ID:   r.ResourceID,
				})
			}
			return models.UserResponse{
				Email:     k.Email,
				Resources: resources,
			}, nil
		}
	}
	return models.UserResponse{}, errors.New("error getting user: not found")
}

func (i *In_memory) UserSignin(email string, password string) error {
	return nil
}
func (i *In_memory) RegisterUser(creds models.Credentials) error {
	if _, ok := i.Members[creds.Email]; ok {
		return errors.New("error registering user")
	}
	return nil
}
