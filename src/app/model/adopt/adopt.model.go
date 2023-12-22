package adopt

import (
	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-backend/src/app/model"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/user"
)

type Adopt struct {
	model.Base
	PetID  *uuid.UUID `json:"pet_id" gorm:"index:idx_name,unique"`
	Pet    *pet.Pet   `json:"pet" gorm:"foreignKey:PetID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	UserID *uuid.UUID `json:"user_id" gorm:"index:idx_name,unique"`
	User   *user.User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
}
