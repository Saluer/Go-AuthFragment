package main

import (
	"auth/routes"
	"auth/server"
	"fmt"
)

func main() {
	server := server.NewServer()
	fmt.Println("Начало")
	//Запустить работу с путями
	routes.ConfigureRoutes(server)
	fmt.Println("Отконфигурированы пути")
	if err := server.Start("4000"); err != nil {
		fmt.Println(err)
	}

	// server.CheckRefresh(loginResult)
}
