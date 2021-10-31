package main

import (
	"auth/routes"
	"auth/server"
	"fmt"
)

func main() {
	server := server.NewServer()
	fmt.Println("Начало")
	//Запустить конфигурацию путей
	routes.ConfigureRoutes(server)
	fmt.Println("Отконфигурированы пути")
	server.Start("4000")
}
