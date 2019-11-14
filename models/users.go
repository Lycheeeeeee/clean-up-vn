package models

import (
	u "github.com/Lycheeeeeee/clean-up-vn/utils"
)

type User struct {
	ID           uint   `json:"id"`
	Displayname  string `json:"displayname"`
	Email        string `json:"email"`
	Issubscribed bool   `json:"issubscribed"`
}

func (user *User) Create() map[string]interface{} {
	user.Issubscribed = false
	GetDB().Create(user)
	response := u.Message(true, "User has been registered")
	response["user"] = user
	return response
}

func GetUser(u string) *User {
	usr := &User{}
	err := GetDB().Table("users").Where("id = ?", u).First(usr).Error
	if err != nil {
		return nil
	}
	return usr
}

func GetAllUsers() []*User {
	usrs := make([]*User, 0)
	err := GetDB().Table("users").Find(&usrs).Error
	if err != nil {
		return nil
	}
	return usrs
}

func (user *User) UpdateSub() map[string]interface{} {
	GetDB().Table("users").Where("id = ?", user.ID).Update("issubscribed", user.Issubscribed)
	response := u.Message(true, "User has been updated")
	response["user"] = user
	return response
}
