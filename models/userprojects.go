package models

import (
	"log"
	"os"
	"strconv"

	u "github.com/Lycheeeeeee/clean-up-vn/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type UserProject struct {
	ID     uint `gorm:"primary_key" json:"id"`
	UserID uint `gorm:"primary_key" json:"user_id"`
}

func (userproject *UserProject) CreateUserProject() map[string]interface{} {
	GetDB().Create(userproject)
	response := u.Message(true, "User has been added to the project")
	fileName := "project_num_" + strconv.FormatUint(uint64(userproject.ID), 10) + ".csv"
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0664)

	if err != nil {
		log.Println(err)
	}

	defer f.Close()
	usr := &User{}
	GetDB().Table("users").Where("id = ?", userproject.UserID).First(usr)

	if _, err := f.WriteString(usr.Displayname + "," + usr.Email + "\n"); err != nil {
		log.Println(err)
	}

	s, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			AwsID, AwsKey, ""), // token can be left blank for now
	})
	if err != nil {
		log.Fatal(err)
	}

	// Upload
	err = AddFileToS3(s, fileName)
	if err != nil {
		log.Fatal(err)
	}

	response["userproject"] = userproject
	return response
}
