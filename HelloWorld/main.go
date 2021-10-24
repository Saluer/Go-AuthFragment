package main

import (
	"auth/routes"
	"auth/server"
	"fmt"
)

func main() {
	server := server.NewServer(":4000")
	//Запустить работу с путями
	fmt.Println("Начало")
	router := routes.ConfigureRoutes(server)
	fmt.Println("Отконфигурированы пути")
	server.Start(router)
}
