package model

import (
	_ "gorm.io/gorm"
)

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CreateUser(user *User) {
	db.Create(user)
}

func FindUser(u *User) User {
	var user User
	db.Where(u).First(&user)
	return user
}
