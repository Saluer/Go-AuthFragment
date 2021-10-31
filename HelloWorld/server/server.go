package server

import (
	db "auth/database"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo *echo.Echo
	*db.DBSettings
}

func NewServer() *Server {
	return &Server{
		Echo:       echo.New(),
		DBSettings: db.Init(),
	}
}

func (server *Server) Start(addr string) error {
	//Запустить сервер по данному порту
	return server.Echo.Start(":" + addr)

}
