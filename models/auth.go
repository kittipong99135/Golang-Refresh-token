package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string
	Password string
	Name     string
	Phone    string
	Age      int
	Rank     string
	Status   string
	Role     string
}

type RequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
