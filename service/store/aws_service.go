package store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

// Service for manage user files in Amazon S3 bucket
type AwsService struct {
	region string
	bucket string

	logger  *logrus.Logger
	session *session.Session
}

func NewAwsService(config *dto.AwsConfig, logger *logrus.Logger) *AwsService {

	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(config.AwsRegion),
			Credentials: credentials.NewStaticCredentials(config.AwsAccessKeyId, config.AwsSecretAccessKey, ""),
		},
	)
	if err != nil {
		panic(err)
	}

	return &AwsService{
		region:  config.AwsRegion,
		session: sess,
		logger:  logger,
		bucket:  config.AwsBucket,
	}
}

// Upload user files to bucket
func (m *AwsService) Upload(id uint32, userId string, data []*dto.FileInfoDto) (*dto.CloudResponseDto, error) {
	uploader := s3manager.NewUploader(m.session)

	target := fmt.Sprintf(fmt.Sprintf("%s/%d/", userId, id))

	respArr := make([]*dto.FileCloudStoreDto, 0)

	for _, v := range data {
		output, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(m.bucket),
			Key:    aws.String(fmt.Sprintf("%s/%s", target, v.Name)),
			Body:   v.Buffer,
		})
		if err != nil {
			m.logger.Errorf("Cannot upload original file to aws store: %v", err)
			return nil, err
		}

		respArr = append(respArr, &dto.FileCloudStoreDto{
			Id:   id,
			Name: v.Name,
			Type: v.Type,
			Url:  output.Location,
		})
	}

	return &dto.CloudResponseDto{Data: respArr}, nil
}

// Downloading file from amazon s3 bucket using userId and imageId, and original url path
func (m *AwsService) Download(url string, userId string, imageId uint32) (*os.File, error) {
	urls := strings.Split(url, "/")

	if len(urls) == 0 {
		return nil, fmt.Errorf("Incorrect url for file downloading ")
	}

	filepath := urls[len(urls)-1] // separate url to get file name

	file, _ := os.Create(filepath)

	downloader := s3manager.NewDownloader(m.session)
	numBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(fmt.Sprintf("%v/%v/%v", userId, imageId, filepath)),
	})
	if err != nil {
		m.logger.Errorf("Unable to download item %q, %v", filepath, err)
		return nil, err
	}
	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	return file, nil

}
