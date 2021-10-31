package service

import (
	"auth/token"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseService struct {
	Client *mongo.Client
}

func NewService(client *mongo.Client) *DatabaseService {
	return &DatabaseService{
		Client: client,
	}
}

func (service *DatabaseService) InsertToken(RefreshTokenData *token.RefreshToken) {
	DBClient := service.Client
	// DBContext := DBSettings.Context
	collection := DBClient.Database("Cluster").Collection("RefreshToken")
	res, err := collection.InsertOne(context.TODO(), bson.M{"refreshuid": RefreshTokenData.RefreshUID})

	if err != nil {
		log.Fatal("Вставка токена в базу не удалась!")
		return
	}

	fmt.Println("Результат вставки токена: ", res)
}

func (service *DatabaseService) GetRefreshToken(RefreshTokenData string) (result token.RefreshToken, err error) {
	// DBSettings := service.server.DBSettings
	DBClient := service.Client
	// DBContext := DBSettings.Context
	collection := DBClient.Database("Cluster").Collection("RefreshToken")
	filter := bson.M{"refreshuid": RefreshTokenData}
	res := collection.FindOne(context.TODO(), filter)
	err = res.Decode(&result)

	if err != nil {
		log.Fatal("Получение токена не удалось: ", err)
	}

	fmt.Println("Результат получения токена: ", result)
	return
}

func (service *DatabaseService) RemoveRefreshToken(RefreshTokenData string) (err error) {
	// DBSettings := service.server.DBSettings
	DBClient := service.Client
	// DBContext := DBSettings.Context
	collection := DBClient.Database("Cluster").Collection("RefreshToken")
	filter := bson.M{"refreshuid": RefreshTokenData}
	var result *mongo.DeleteResult
	result, err = collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal("Удаление токена не удалось: ", err)
	}

	fmt.Println("Результат удаления токена: ", result)
	return
}
