package db

import (
	"context"
	"fmt"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDbService struct {
	logger          *logrus.Logger
	client          *mongo.Client
	ImageStore      string
	UsersCollection string
}

// TODO add reties
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

func (m *MongoDbService) Insert(storeDto *dto.DbImageStoreDAO) error {
	if storeDto == nil {
		return fmt.Errorf("Nil data for inserting ")
	}
	col := m.client.Database(m.ImageStore).Collection(m.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, storeDto)
	if err != nil {
		m.logger.Errorf("Cannot save data to db: %v", err)
		return err
	}

	return nil
}

func (m *MongoDbService) GetPicture(picId uint32, width, height int) *dto.DbImageStoreDAO {
	col := m.client.Database(m.ImageStore).Collection(m.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := dto.DbImageStoreDAO{}
	err := col.FindOne(ctx, bson.D{{"picid", picId}, {"resizedwidth", width}, {"resizedheight", height}}).Decode(&res)

	if err != nil {
		m.logger.Warnf("Cannot get data from db by request (picid: %d, width: %d, height: %d), err: %v", picId, width, height, err)
		return nil
	}

	return &res
}

func (m *MongoDbService) FindAllPictureByUserId(userId string) []*dto.DbImageStoreDAO {
	col := m.client.Database(m.ImageStore).Collection(m.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	logEntity := m.logger.WithFields(logrus.Fields{
		"userId": userId,
	})

	result := make([]*dto.DbImageStoreDAO, 0)

	cursor, err := col.Find(ctx, bson.D{{"userid", userId}})
	if err != nil {
		logEntity.Errorf("Cannot get records from db. Err: %v", err)
		return nil
	}

	err = cursor.All(context.TODO(), &result)

	if err != nil {
		logEntity.Errorf("Cannot map cursor results to response array. Err: %v", err)
		return nil
	}

	return result

}
