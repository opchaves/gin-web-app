package fixture

import (
	"github.com/go-faker/faker/v4"
	"github.com/opchaves/gin-web-app/app/model"
)

func GetMockUser() *model.User {
	return &model.User{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Password:  faker.Password(),
	}
}
