package seed

import (
	"github.com/isd-sgcu/johnjud-backend/src/app/model/user"

	"github.com/bxcodec/faker/v3"
)

func (s Seed) UserSeed1705075734828() error {
	for i := 0; i < 3; i++ {
		usr := user.User{
			Email:     faker.Email(),
			Password:  faker.Password(),
			Firstname: faker.FirstName(),
			Lastname:  faker.LastName(),
			Role:      "admin",
		}
		err := s.db.Create(&usr).Error

		if err != nil {
			return err
		}
	}
	return nil
}
