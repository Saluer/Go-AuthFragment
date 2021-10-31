package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBSettings struct {
	Client  *mongo.Client
	Context context.Context
}

func Init() *DBSettings {
	//Подключиться к базе данных mongoDB в облаке
	//Доступно для любого IP-адреса
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://dbUser:dbUserPassword@cluster.6qr8d.mongodb.net/Cluster?retryWrites=true&w=majority")
	//Создание контекста
	context := context.TODO()
	//Создание клиента подключения к базе данных
	client, err := mongo.Connect(context, clientOptions)
	if err != nil {
		print("Инициализация базы данных не удалась!")
		log.Fatal(err)
	}
	return &DBSettings{client, context}
}
