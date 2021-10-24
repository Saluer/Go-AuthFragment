package routes

import (
	"auth/handlers"
	"auth/server"

	"github.com/gorilla/mux"
)

func ConfigureRoutes(server *server.Server) *mux.Router {
	//Создать роутер
	router := mux.NewRouter()
	//Создать обработчик событий по Auth
	authHandler := handlers.NewAuthHandler(server)
	//Присвоить путям login, refresh соответствующие обработчики
	router.Handle("/login", authHandler.LoginHandler).Methods("GET")
	router.Handle("/refresh", authHandler.RefreshHandler).Methods("GET")

	//Создать прослойку с jwt
	// Присвоить путю access соответствующий обработчик, обёрнутый в прослойку

	return router
}
