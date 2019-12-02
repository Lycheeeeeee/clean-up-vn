package models

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	u "github.com/Lycheeeeeee/clean-up-vn/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type UserInFile struct {
	Displayname string
	Email       string
}

type UserProject struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	UserID      uint   `gorm:"primary_key" json:"user_id"`
}

func ReadFileFromS3(fileName string) (b []byte, err error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			AwsID, AwsKey, ""), // token can be left blank for now
	})
	if err != nil {
		// log.Fatal(err)
	}
	file, err := os.Create(fileName)
	downloader := s3manager.NewDownloader(s)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(S3_BUCKET),
			Key:    aws.String(fileName),
		})
	if err != nil {
		// log.Fatalf("Unable to download item %q, %v", fileName, err)
		return
	}
	byteFormat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	return byteFormat, nil
}

func (userproject *UserProject) CreateUserProject() map[string]interface{} {
	pro := &Project{}
	err := GetDB().Table("projects").Where("id = ?", userproject.ID).First(pro).Error
	if err != nil {
		return nil
	}

	usr := &User{}
	er := GetDB().Table("users").Where("id = ?", userproject.UserID).First(usr).Error
	if er != nil {
		return nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			AwsID, AwsKey, ""), // token can be left blank for now
	})
	if err != nil {
		log.Fatal(err)
	}

	svc := sns.New(sess)
	_, e := svc.Subscribe(&sns.SubscribeInput{
		Endpoint:              aws.String(usr.Email),
		Protocol:              aws.String("email"),
		ReturnSubscriptionArn: aws.Bool(true),
		TopicArn:              aws.String(pro.TopicArn),
	})
	if e != nil {
		fmt.Println(err.Error())
	}
	GetDB().Create(userproject)
	response := u.Message(true, "User has been added to the project")
	fileName := "project_num_" + strconv.FormatUint(uint64(userproject.ID), 10) + ".csv"
	dir := "s3File"
	filePath := filepath.Join(dir, fileName)

	ReadFileFromS3(filePath)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0664)

	if err != nil {
		log.Println(err)
	}
	//defer f.Close()

	if _, err := f.WriteString(usr.Displayname + "," + usr.Email + "\n"); err != nil {
		log.Println(err)
	}
	f.Close()
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			AwsID, AwsKey, ""), // token can be left blank for now
	})
	if err != nil {
		log.Fatal(err)
	}

	// Upload
	err = AddFileToS3(s, filePath)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(filePath)
	if err != nil {
		log.Fatal(err)
	}

	response["userproject"] = userproject
	return response
}

func (userproject *UserProject) LeaveProject() map[string]interface{} {
	GetDB().Delete(userproject)
	targetUser := &User{}
	usrinfiles := []UserInFile{}

	usrpros := make([]*UserProject, 0)
	err := GetDB().Table("user_projects").Find(&usrpros).Error
	if err != nil {
		return nil
	}
	//db.Table("users").Select("users.displayname, users.email").Joins("left join user_projects on user_projects.user_id = users.id ").Where("user_projects.id = ?", userproject.ID).Scan(&usrinfiles)

	GetDB().Table("users").Where("id = ?", userproject.UserID).First(targetUser)

	response := u.Message(true, "User has been removed from the project")
	fileName := "project_num_" + strconv.FormatUint(uint64(userproject.ID), 10) + ".csv"
	dir := "s3File"
	filePath := filepath.Join(dir, fileName)
	// dat, err := ioutil.ReadFile(fileName)
	// check(err)
	singleValueByte, err := ReadFileFromS3(filePath)
	splitString := bytes.Split(singleValueByte, []byte("\n"))
	for i := 0; i < len(splitString)-1; i++ {
		usrAttrs := bytes.Split(splitString[i], []byte(","))
		usrinfile := UserInFile{}
		fmt.Println(string(usrAttrs[1]))
		fmt.Println(targetUser.Email)
		if string(usrAttrs[1]) != targetUser.Email {
			fmt.Println("through")
			usrinfile.Displayname = string(usrAttrs[0])
			usrinfile.Email = string(usrAttrs[1])
			usrinfiles = append(usrinfiles, usrinfile)
		}
	}

	// f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0664)

	// if _, err := f.WriteString("aasdasdasdsada"); err != nil {
	// 	log.Println(err)
	// }

	if len(usrinfiles) == 0 {
		err = ioutil.WriteFile(filePath, []byte(""), 0644)
		if err != nil {
			log.Println(err)
		}
	} else {
		for _, uInFile := range usrinfiles {
			err = ioutil.WriteFile(filePath, []byte(uInFile.Displayname+","+uInFile.Email), 0644)
			if err != nil {
				log.Println(err)
			}
		}
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
	err = AddFileToS3(s, filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(filePath)
	if err != nil {
		log.Fatal(err)
	}

	response["userproject"] = userproject
	return response
}

// func (user *User)Runreport() map[string]interface{} {
// 	if user.ID ==1{

// 	}

// 	GetDB().Create(user)

// 	response := u.Message(true, "User has been registered")
// 	response["user"] = user
// 	return response
// }
