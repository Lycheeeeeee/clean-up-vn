package models

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	u "github.com/Lycheeeeeee/clean-up-vn/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

const (
	S3_REGION = "ap-southeast-1"
	S3_BUCKET = "elasticbeanstalk-ap-southeast-1-429048041589"
	AwsID     = "AKIAJWFR247HPACDCOBA"
	AwsKey    = "7npSdtT85NGkGpufy/If9s3pumRE8qleAyKHVG3y"
)

func AddFileToS3(s *session.Session, fileDir string) error {
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3_BUCKET),
		Key:                  aws.String(fileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

type UUID [16]byte

func NewUUID() uuid.UUID {
	unique, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
	}
	return unique
}

type Project struct {
	gorm.Model
	Name        string    `json:"name"`
	Longtitude  float64   `json:"longtitude"`
	Latitude    float64   `json:"latitude"`
	Description string    `json:"description"`
	Owner       uint      `json:"owner"`
	Status      string    `json:"status"`
	Time        time.Time `json:"time"`
	Result      int       `json:"result"`
}

func (project *Project) Create() map[string]interface{} {
	GetDB().Create(project)
	response := u.Message(true, "Project has been created")
	response["project"] = project
	fileName := "project_num_" + strconv.FormatUint(uint64(project.ID), 10) + ".csv"
	// os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	dir := "s3File"
	filePath := filepath.Join(dir, fileName)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			panic("directory does not exist")
		}
	}
	fmt.Printf("creating file:%v", fileName)
	_, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
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
	return response
}

func GetAllProjects() []*Project {
	pros := make([]*Project, 0)
	err := GetDB().Table("projects").Find(&pros).Error
	if err != nil {
		return nil
	}
	return pros
}
func GetProject(u string) *Project {
	pro := &Project{}
	err := GetDB().Table("projects").Where("id = ?", u).First(pro).Error
	if err != nil {
		return nil
	}
	return pro
}

func (project *Project) InputResultNCloseProject() map[string]interface{} {
	GetDB().Table("projects").Where("id = ?", project.ID).Updates(map[string]interface{}{"status": "close", "result": project.Result})
	response := u.Message(true, "Project has been updated")
	response["project"] = project
	return response
}

func (project *Project) GetVolunteerList(owner uint) (result []byte, err error) {
	if project.Owner == owner {
		var s = strconv.FormatUint(uint64(project.ID), 10)
		return ReadFileFromS3("s3File/project_num" + s + ".csv")
	}
	return
}
