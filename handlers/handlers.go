package handlers

import (
	"auth/requests"
	"auth/responses"
	"auth/server"
	"auth/token"
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

	//Передаём данные запроса в loginRequest
	if err := c.Bind(loginRequest); err != nil {
		return err
	}

	//Создаём новую пару токенов
	accessToken, refreshToken, err := authHandler.tokenService.GenerateTokenPair(loginRequest)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken)

	return responses.Response(c, http.StatusOK, res)
}

func (authHandler *AuthHandler) Refresh(c echo.Context) (err error) {
	refreshRequest := new(requests.RefreshRequest)

	//Передаём данные запроса в refreshRequest
	if err := c.Bind(refreshRequest); err != nil {
		return err
	}

	//Проверяем, есть ли токен в базе данных
	if err := authHandler.tokenService.ValidateToken(refreshRequest.RefreshToken); err != nil {
		return responses.MessageResponse(c, http.StatusUnauthorized, "Пользователь не авторизован!")
	}

	//Получаем данные из токена
	var claims *token.JwtCustomClaims
	if claims, err = authHandler.tokenService.ParseToken(refreshRequest.AccessToken); err != nil {
		return responses.MessageResponse(c, http.StatusUnauthorized, "Пользователь не авторизован!")
	}

	//Удаляем токен из базы данных
	authHandler.tokenService.RemoveRefreshToken(refreshRequest.RefreshToken)

	loginRequest := new(requests.LoginRequest)
	loginRequest.UserID = claims.UserID

	//Создаём новую пару токенов
	accessToken, refreshToken, err := authHandler.tokenService.GenerateTokenPair(loginRequest)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken)

	return responses.Response(c, http.StatusOK, res)
}
