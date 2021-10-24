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

func Init() DBSettings {
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://dbUser:dbUserPassword@cluster.6qr8d.mongodb.net/Cluster?retryWrites=true&w=majority")
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		print("Инициализация базы данных не удалась!")
		log.Fatal(err)
	}
	return DBSettings{client, context.TODO()}
}
