package db

import (
	"context"
	"fmt"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

// Service for saving and reading data of user image resize results.
// Store image as object with info about original image (path) and resized image (path, requested size params)

const (
	Retry     = 3
	SleepTime = 100 * time.Millisecond
)

type MongoDbService struct {
	logger          *logrus.Logger
	client          *mongo.Client
	ImageStore      string
	UsersCollection string
}

func NewMongoDbService(cfg *dto.MongoDbConfig, logger *logrus.Logger) *MongoDbService {
	service := &MongoDbService{
		logger:          logger,
		ImageStore:      cfg.Store,
		UsersCollection: cfg.Collection,
	}
	service.client = service.connect(cfg.Username, cfg.Password, cfg.Address)
	return service
}

func (m *MongoDbService) connect(username, password, address string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(
		fmt.Sprintf("mongodb+srv://%s:%s@%s?retryWrites=true&w=majority", username, password, address),
	))
	if err != nil {
		m.logger.Fatal(err)
	}
	return client
}

// Inserting total info of processed image to DB (original url, resized url, resize params)
func (m *MongoDbService) Insert(storeDto *dto.DbImageStoreDAO) error {
	if storeDto == nil {
		return fmt.Errorf("Nil data for inserting ")
	}
	col := m.client.Database(m.ImageStore).Collection(m.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	leftRetry := Retry
	currentSleepTime := SleepTime

	var err error
	for leftRetry > 0 {
		_, err = col.InsertOne(ctx, storeDto)
		if err != nil {
			m.logger.Warnf("Cannot save data to db. Retrying... Error: %v", err)
			leftRetry--
			time.Sleep(currentSleepTime)
			currentSleepTime += SleepTime
			continue
		}
		return nil
	}

	return err
}

// Searching image by imageId and size params
func (m *MongoDbService) GetImage(picId uint32, width, height int) *dto.DbImageStoreDAO {
	col := m.client.Database(m.ImageStore).Collection(m.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := dto.DbImageStoreDAO{}
	leftRetry := Retry
	currentSleepTime := SleepTime

	var err error
	for leftRetry > 0 {
		err = col.FindOne(ctx, bson.D{{"picid", picId}, {"resizedwidth", width}, {"resizedheight", height}}).Decode(&res)

		if err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				return nil
			}
			leftRetry--
			m.logger.Warnf("Cannot get data from db by request (picid: %d, width: %d, height: %d). Retrying... Err: %v", picId, width, height, err)
			time.Sleep(currentSleepTime)
			currentSleepTime += SleepTime
			continue
		}
		return &res

	}

	return nil
}

func (m *MongoDbService) GetImageByImageId(picId uint32) *dto.DbImageStoreDAO {
	col := m.client.Database(m.ImageStore).Collection(m.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := dto.DbImageStoreDAO{}
	leftRetry := Retry
	currentSleepTime := SleepTime

	var err error
	for leftRetry > 0 {
		err = col.FindOne(ctx, bson.D{{"picid", picId}}).Decode(&res)

		if err != nil {
			leftRetry--
			m.logger.Warnf("Cannot get data from db by request (picid: %d). Retrying...  Err: %v", picId, err)
			time.Sleep(currentSleepTime)
			currentSleepTime += SleepTime
			continue
		}
		return &res
	}

	return nil
}

// Collect all user images by userId
func (m *MongoDbService) FindAllPictureByUserId(userId string) []*dto.DbImageStoreDAO {
	col := m.client.Database(m.ImageStore).Collection(m.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logEntity := m.logger.WithFields(logrus.Fields{
		"userId": userId,
	})

	result := make([]*dto.DbImageStoreDAO, 0)

	leftRetry := Retry
	currentSleepTime := SleepTime

	for leftRetry > 0 {
		cursor, err := col.Find(ctx, bson.D{{"userid", userId}})
		if err != nil {
			leftRetry--
			logEntity.Warnf("Cannot get records from db. Retrying... Err: %v", err)
			time.Sleep(currentSleepTime)
			currentSleepTime += SleepTime
			continue
		}

		err = cursor.All(context.TODO(), &result)

		if err != nil {
			logEntity.Errorf("Cannot map cursor results to response array. Err: %v", err)
			return nil
		}

		return result
	}

	return nil
}
