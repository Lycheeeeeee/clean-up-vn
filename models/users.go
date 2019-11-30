package models

import (
	"os"
	"strings"

	u "github.com/Lycheeeeeee/clean-up-vn/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Social       string `json:"social"`
	Displayname  string `json:"displayname"`
	Phonenumber  string `json:"phonenumber"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Issubscribed bool   `json:"issubscribed"`
	Token        string `json:"token";sql:"-"`
}

func (user *User) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &User{}

	//check for errors and duplicate emails
	err := GetDB().Table("users").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (user *User) CreateAccount() map[string]interface{} {
	if user.Social != "" {
		temp := &User{}

		//check for errors and duplicate emails
		err := GetDB().Table("users").Where("social = ?", user.Social).First(temp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return u.Message(false, "Connection error. Please retry")
		}
		if temp.Social != "" {
			return u.Message(false, "Social id has been registed")
		}

		GetDB().Create(user)
	} else {
		if resp, ok := user.Validate(); !ok {
			return resp
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)

		GetDB().Create(user)

		if user.ID <= 0 {
			return u.Message(false, "Failed to create account, connection error.")
		}
	}
	//Create new JWT token for the newly registered account
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["user"] = user
	return response
}

func Login(email, password string) map[string]interface{} {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	user.Password = ""

	//Create JWT token
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}
func GetTokenFromSocial(socialid string) map[string]interface{} {
	user := &User{}
	err := GetDB().Table("users").Where("social = ?", socialid).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Social address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
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
