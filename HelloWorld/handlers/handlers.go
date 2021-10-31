package handlers

import (
	"auth/requests"
	"auth/responses"
	"auth/server"
	ts "auth/token/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	server       *server.Server
	tokenService *ts.TokenService
}

func NewAuthHandler(server *server.Server) *AuthHandler {
	return &AuthHandler{server: server, tokenService: ts.NewTokenService(server)}
}

func (authHandler *AuthHandler) Login(c echo.Context) error {
	loginRequest := new(requests.LoginRequest)
	//Прикрепляем к контексту loginRequest
	if err := c.Bind(loginRequest); err != nil {
		return err
	}
	accessToken, refreshToken, exp, err := authHandler.tokenService.GenerateTokenPair(loginRequest)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return responses.Response(c, http.StatusOK, res)
}

func (authHandler *AuthHandler) Refresh(c echo.Context) error {
	refreshRequest := new(requests.RefreshRequest)

	//Прикрепляем к контексту refreshRequest
	if err := c.Bind(refreshRequest); err != nil {
		return err
	}

	if err := authHandler.tokenService.ValidateToken(refreshRequest.RefreshToken); err != nil {
		return responses.MessageResponse(c, http.StatusUnauthorized, "Пользователь не авторизован!")
	}

	loginRequest := new(requests.LoginRequest)
	//Прикрепляем к контексту loginRequest
	if err := c.Bind(loginRequest); err != nil {
		return err
	}
	accessToken, refreshToken, exp, err := authHandler.tokenService.GenerateTokenPair(loginRequest)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return responses.Response(c, http.StatusOK, res)
}

func (authHandler *AuthHandler) Access(c echo.Context) error {
	return nil
}
