package in_memory

import (
	"context"
	"errors"

	"github.com/HackRVA/memberserver/models"
)

func (i *In_memory) GetUser(ctx context.Context, email string) (models.UserResponse, error) {
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

func (i *In_memory) UserSignin(ctx context.Context, email string, password string) error {
	return nil
}

func (i *In_memory) RegisterUser(ctx context.Context, creds models.Credentials) error {
	if _, ok := i.Members[creds.Email]; ok {
		return errors.New("error registering user")
	}
	return nil
}
