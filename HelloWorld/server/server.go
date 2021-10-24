package server

import (
	"net/http"
	"os"

	db "auth/database"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	port string
	db.DBSettings
}

func NewServer(port string) *Server {
	return &Server{
		port:       port,
		DBSettings: db.Init(),
	}
}

func (server *Server) Start(router *mux.Router) {
	//Запустить сервер по данному порту
	http.ListenAndServe(server.port, handlers.LoggingHandler(os.Stdout, router))

}
