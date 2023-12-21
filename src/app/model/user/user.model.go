package user

import "github.com/isd-sgcu/johnjud-backend/src/app/model"

type User struct {
	model.Base
	Email     string `json:"email" gorm:"type:tinytext"`
	Password  string `json:"password" gorm:"type:tinytext"`
	Firstname string `json:"firstname" gorm:"type:tinytext"`
	Lastname  string `json:"lastname" gorm:"type:tinytext"`
	Role      string `json:"role" gorm:"type:tinytext"`
}
