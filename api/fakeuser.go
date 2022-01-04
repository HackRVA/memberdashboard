package api

import (
	"memberserver/api/models"
	"os"
)

func isFakeUser() bool {
	return len(os.Getenv("FAKE_USER")) > 0
}

func getFakeUserProfile() models.Member {
	return models.Member{
		Name:  "Fake User",
		Email: "fakeuser",
		Resources: []models.MemberResource{
			{Name: "admin"},
		},
	}
}
