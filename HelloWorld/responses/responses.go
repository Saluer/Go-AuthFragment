package responses

import "github.com/labstack/echo/v4"

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Exp          int64  `json:"exp"`
}

type Error struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type Data struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewLoginResponse(token, refreshToken string, exp int64) *LoginResponse {
	return &LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		Exp:          exp,
	}
}

func MessageResponse(c echo.Context, statusCode int, message string) error {
	return Response(c, statusCode, Data{
		Code:    statusCode,
		Message: message,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return Response(c, statusCode, Error{
		Code:  statusCode,
		Error: message,
	})
}

func Response(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, data)
}
