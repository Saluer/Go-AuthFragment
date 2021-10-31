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
	Client  *mongo.Client
	Context context.Context
}

func NewService(client *mongo.Client, context context.Context) *DatabaseService {
	return &DatabaseService{
		Client:  client,
		Context: context,
	}
}

func (service *DatabaseService) InsertToken(RefreshTokenData *token.RefreshToken) {
	DBClient := service.Client
	DBContext := service.Context
	//Получение коллекции refresh-токенов из заданной базы данных
	collection := DBClient.Database("Cluster").Collection("RefreshToken")
	//Вставка токена в базу
	res, err := collection.InsertOne(DBContext, bson.M{"refreshuid": RefreshTokenData.RefreshUID})

	if err != nil {
		log.Fatal("Вставка токена в базу не удалась!")
		return
	}

	fmt.Println("Результат вставки токена: ", res)
}

func (service *DatabaseService) GetRefreshToken(RefreshTokenData string) (result token.RefreshToken, err error) {
	DBClient := service.Client
	DBContext := service.Context
	//Получение коллекции refresh-токенов из заданной базы данных
	collection := DBClient.Database("Cluster").Collection("RefreshToken")
	//Создание критерия поиска токена
	filter := bson.M{"refreshuid": RefreshTokenData}
	//Найти токен и скопировать его обработанные данные в result
	res := collection.FindOne(DBContext, filter)
	err = res.Decode(&result)

	if err != nil {
		log.Fatal("Получение токена не удалось: ", err)
	}

	fmt.Println("Результат получения токена: ", result)
	return
}

func (service *DatabaseService) RemoveRefreshToken(RefreshTokenData string) (err error) {
	DBClient := service.Client
	DBContext := service.Context
	//Получение коллекции refresh-токенов из заданной базы данных
	collection := DBClient.Database("Cluster").Collection("RefreshToken")
	//Создание критерия поиска токена
	filter := bson.M{"refreshuid": RefreshTokenData}
	//Удалить токен с указанными данными
	var result *mongo.DeleteResult
	result, err = collection.DeleteOne(DBContext, filter)

	if err != nil {
		log.Fatal("Удаление токена не удалось: ", err)
	}

	fmt.Println("Результат удаления токена: ", result)
	return
}
