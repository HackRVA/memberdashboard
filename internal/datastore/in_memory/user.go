package in_memory

import (
	"errors"
	"memberserver/internal/models"
)

var testUsers = map[string]models.Member{
	"test": {
		Email: "test",
	},
}

func (i In_memory) GetUser(email string) (models.UserResponse, error) {

	for _, k := range testUsers {
		if k.Email == email {
			return models.UserResponse{
				Email: k.Email,
			}, nil
		}
	}
	return models.UserResponse{}, errors.New("error getting user: not found")
}

func (i In_memory) UserSignin(email string, password string) error {
	return nil
}
func (i In_memory) RegisterUser(creds models.Credentials) error {
	if _, ok := testUsers[creds.Email]; ok {
		return errors.New("error registering user")
	}
	return nil
}
