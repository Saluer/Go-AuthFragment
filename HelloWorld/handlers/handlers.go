package handlers

import (
	"auth/responses"
	"auth/server"
	ts "auth/token/service"
	"net/http"
)

type AuthHandler struct {
	LoginHandler   http.HandlerFunc
	RefreshHandler http.HandlerFunc
	AccessHandler  http.HandlerFunc
}

func NewAuthHandler(server *server.Server) *AuthHandler {
	var LoginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, refreshToken, exp, err := ts.GenerateTokenPair(server)

		if err != nil {
			w.Write([]byte("Случилась ошибка"))
		}

		res := responses.NewLoginResponse(accessToken, refreshToken, exp)

		w.Write([]byte(res.AccessToken))
		w.Write([]byte(res.RefreshToken))
	})
	var RefreshHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	var AccessHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	return &AuthHandler{LoginHandler: LoginHandler, RefreshHandler: RefreshHandler, AccessHandler: AccessHandler}
}
