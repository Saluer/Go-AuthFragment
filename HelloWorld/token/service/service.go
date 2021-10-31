package service

import (
	"fmt"
	"log"

	dbs "auth/database/service"
	"auth/requests"
	"auth/server"
	t "auth/token"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//Пароль, с помощью которого мы шифруем токены
const secret = "secretPassword"

type TokenService struct {
	dbService *dbs.DatabaseService
}

func NewTokenService(server *server.Server) *TokenService {
	return &TokenService{dbService: dbs.NewService(server.Client, server.Context)}
}

func (ts *TokenService) GenerateTokenPair(lr *requests.LoginRequest) (accessToken string,
	refreshToken string,
	err error,
) {

	if accessToken, err = ts.createAccessToken(lr.UserID, secret); err != nil {
		return
	}

	if refreshToken, err = ts.createRefreshToken(lr.UserID, secret); err != nil {
		return
	}

	//Вставить токен в базу данных
	refreshTokenData := &t.RefreshToken{RefreshUID: refreshToken}
	ts.dbService.InsertToken(refreshTokenData)
	return
}

//Проверить, есть ли токен в базе данных
func (ts *TokenService) ValidateToken(tokenData string) (
	err error,
) {
	//Если не выдаст ошибку, то токен есть в базе и пользоватерь авторизован
	if _, err = ts.dbService.GetRefreshToken(tokenData); err != nil {
		fmt.Println("Получение refresh-токена из базы не удалось")
	}
	return
}

//Получить данные из токена
func (ts *TokenService) ParseToken(tokenString string) (
	claims *t.JwtCustomClaims,
	err error,
) {
	token, err := jwt.ParseWithClaims(tokenString, &t.JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("ошибка получения данных по алгоритму: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*t.JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

//Удалить refresh-токен из базы данных
func (ts *TokenService) RemoveRefreshToken(tokenString string) (err error) {
	err = ts.dbService.RemoveRefreshToken(tokenString)
	return err
}

func (ts *TokenService) createAccessToken(userID uuid.UUID, secret string) (signedToken string,
	err error,
) {

	claims := t.JwtCustomClaims{
		UserID:         userID,
		StandardClaims: jwt.StandardClaims{},
	}

	//Создание нового токена
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err = jwtToken.SignedString([]byte(secret))

	return
}

func (ts *TokenService) createRefreshToken(userID uuid.UUID, secret string) (cryptedTokenString string,
	err error,
) {

	claims := &t.JwtCustomClaims{
		UserID:         userID,
		StandardClaims: jwt.StandardClaims{},
	}

	//Создание нового токена
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	//Подпись токена зашифрованным паролем
	var signedToken string
	if signedToken, err = jwtToken.SignedString([]byte(secret)); err != nil {
		log.Print("Токен не был удачно подписан!")
		return
	}

	// Хэширование с помощью bcrypt подписанного токена
	var cryptedToken []byte
	if cryptedToken, err = bcrypt.GenerateFromPassword([]byte(signedToken), bcrypt.DefaultCost); err != nil {
		log.Print("Пароль не был удачно зашифрован!")
		return
	}
	cryptedTokenString = string(cryptedToken)
	return
}
