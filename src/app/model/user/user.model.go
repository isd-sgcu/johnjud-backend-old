package user

import "github.com/isd-sgcu/johnjud-backend/src/app/model"

type User struct {
	model.Base
	Email     string `json:"email" gorm:"tinytext"`
	Password  string `json:"password" gorm:"tinytext"`
	Firstname string `json:"firstname" gorm:"tinytext"`
	Lastname  string `json:"lastname" gorm:"tinytext"`
	Role      string `json:"role" gorm:"tinytext"`
}
