package store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/senseyman/image-media-processor/dto/response_dto"
	"github.com/sirupsen/logrus"
)

type AwsStoreManager struct {
	region string
	bucket string

	logger  *logrus.Logger
	session *session.Session
}

func NewAwsStoreManager(config *dto.AwsConfig, logger *logrus.Logger) *AwsStoreManager {

	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(config.AwsRegion),
			Credentials: credentials.NewStaticCredentials(config.AwsAccessKeyId, config.AwsSecretAccessKey, ""),
		},
	)
	if err != nil {
		panic(err)
	}

	return &AwsStoreManager{
		region:  config.AwsRegion,
		session: sess,
		logger:  logger,
		bucket:  config.AwsBucket,
	}
}

func (m *AwsStoreManager) Upload(id uint32, userId string, data []*dto.FileInfoDto) (*response_dto.CloudResponseDto, error) {
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

	return &response_dto.CloudResponseDto{Data: respArr}, nil
}
