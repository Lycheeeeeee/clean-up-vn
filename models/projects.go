package models

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	u "github.com/Lycheeeeeee/clean-up-vn/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

var e = godotenv.Load()
var (
	S3_REGION = "ap-southeast-1"
	S3_BUCKET = "elasticbeanstalk-ap-southeast-1-429048041589"
	AwsID     = os.Getenv("id")
	AwsKey    = os.Getenv("key")
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
	TopicArn    string    `json:"topic_arn"`
}

func TimeDecoder(timer string) time.Time {
	var year, month, date, hour, min int
	stage1 := strings.Split(timer, "T")
	for i := 0; i < len(stage1); i++ {
		if i == 0 {
			stage2 := strings.Split(stage1[0], "-")
			year, _ = strconv.Atoi(stage2[0])
			// if err != nil {
			// 	log.Fatal(err)
			// }
			month, _ = strconv.Atoi(stage2[1])
			// if err != nil {
			// 	log.Fatal(err)
			// }
			date, _ = strconv.Atoi(stage2[2])
			// if err != nil {
			// 	log.Fatal(err)
			// }
		}
		if i == 1 {
			stage3 := strings.Split(stage1[1], ":")
			hour, _ = strconv.Atoi(stage3[0])
			// if err != nil {
			// 	log.Fatal(err)
			// }
			min, _ = strconv.Atoi(stage3[1])
			// if err != nil {
			// 	log.Fatal(err)
			// }
		}
	}
	then := time.Date(
		year, time.Month(month), date, hour, min, 00, 000000000, time.UTC)

	return then
}

func (project *Project) Create(timer string) map[string]interface{} {
	project.Time = TimeDecoder(timer)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			AwsID, AwsKey, ""), // token can be left blank for now
	})
	if err != nil {
		log.Fatal(err)
	}

	svc := sns.New(sess)
	dater := strings.Split(project.Time.String(), " ")
	topicName := strconv.Itoa(int(project.ID)) + "_" + dater[0]

	result, err := svc.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	project.TopicArn = *result.TopicArn
	GetDB().Create(project)
	response := u.Message(true, "Project has been created successfully")
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
	_, er := os.Create(filePath)
	if er != nil {
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

func (project *Project) GetVolunteerList(owner uint) map[string]interface{} {
	if project.Owner == owner {
		var projectnum = strconv.FormatUint(uint64(project.ID), 10)
		// return ReadFileFromS3("s3File/project_num_" + s + ".csv")
		s, err := session.NewSession(&aws.Config{
			Region: aws.String(S3_REGION),
			Credentials: credentials.NewStaticCredentials(
				AwsID, AwsKey, ""), // token can be left blank for now
		})
		if err != nil {
			log.Fatal(err)
		}
		svc := s3.New(s)
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(S3_BUCKET),
			Key:    aws.String("s3File/project_num_" + projectnum + ".csv"),
		})
		urlStr, err := req.Presign(15 * time.Minute)
		if err != nil {
			log.Println("Failed to sign request", err)
		}
		response := u.Message(true, "Successful generate presigned url for s3 file downloading")
		response["presignedurl"] = urlStr
		return response
	}

	return nil
}

func TestNotification() map[string]interface{} {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			AwsID, AwsKey, ""), // token can be left blank for now
	})
	if err != nil {
		log.Fatal(err)
	}

	svc := sns.New(sess)
	result, err := svc.ListTopics(nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	response := u.Message(true, "Successful")
	for i, t := range result.Topics {
		splitstring := strings.Split(*t.TopicArn, "_")
		mes := "Remember to help us green our earth at " + splitstring[1] + " GMT007"
		resu, err := svc.Publish(&sns.PublishInput{
			Message:  aws.String(mes),
			TopicArn: aws.String(*t.TopicArn),
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		responsename := "ID-" + strconv.Itoa(i)
		response[responsename] = *resu.MessageId
	}
	return response
}
