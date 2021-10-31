package routes

import (
	"auth/handlers"
	"auth/server"

	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
)

func ConfigureRoutes(server *server.Server) {
	//Создать обработчик событий по Auth
	authHandler := handlers.NewAuthHandler(server)
	//Запустить прослойку, выводящую логи по действиям
	server.Echo.Use(echoMW.Logger())
	//Присвоить путям login, refresh соответствующие обработчики
	server.Echo.POST("/login", authHandler.Login)
	server.Echo.POST("/refresh", authHandler.Refresh)
	//Тестовый путь, который логинится, а потом запускает обновление
	server.Echo.GET("/test", func(c echo.Context) error {
		loginResult := server.CheckLogin()
		server.CheckRefresh(loginResult)
		return nil
	})
}
