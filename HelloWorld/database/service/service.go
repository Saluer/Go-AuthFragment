package service

import (
	"auth/server"
	"auth/token"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type DatabaseService struct {
	server *server.Server
}

func NewService(server *server.Server) *DatabaseService {
	return &DatabaseService{
		server: server,
	}
}

func (service *DatabaseService) InsertToken(RefreshTokenData token.RefreshToken) {
	DBSettings := service.server.DBSettings
	DBClient := DBSettings.Client
	// DBContext := DBSettings.Context
	collection := DBClient.Database("Cluster").Collection("RefreshToken")
	res, err := collection.InsertOne(context.TODO(), bson.M{"text": RefreshTokenData})

	if err != nil {
		print("Вставка токена в базу не удалась!")
		return
	}

	fmt.Println("Результат вставки токена: ", res)
}
