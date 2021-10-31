package service

import (
	"fmt"
	"log"
	"time"

	dbs "auth/database/service"
	"auth/requests"
	"auth/server"
	t "auth/token"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const ExpireAccessMinutes = 30
const ExpireRefreshMinutes = 2 * 60
const AutoLogoffMinutes = 10

const secret = "secretPassword"

type TokenService struct {
	server    *server.Server
	dbService *dbs.DatabaseService
}

func NewTokenService(server *server.Server) *TokenService {
	return &TokenService{server: server, dbService: dbs.NewService(server.Client)}
}

func (ts *TokenService) GenerateTokenPair(lr *requests.LoginRequest) (accessToken string,
	refreshToken string,
	exp int64,
	err error,
) {

	if accessToken, exp, err = ts.createAccessToken(lr.UserID, ExpireAccessMinutes,
		secret); err != nil {
		return
	}

	if refreshToken, _, err = ts.createRefreshToken(lr.UserID, ExpireRefreshMinutes,
		secret); err != nil {
		return
	}

	//Использовать кэшированные токены для добавления в базу. Вернее, один refresh
	refreshTokenData := &t.RefreshToken{RefreshUID: refreshToken}
	ts.dbService.InsertToken(refreshTokenData)
	return
}

func (ts *TokenService) ValidateToken(tokenData string) (
	err error,
) {
	//Если не выдаст ошибку, то токен есть в базе и пользоватерь авторизован
	if _, err = ts.dbService.GetRefreshToken(tokenData); err != nil {
		fmt.Println("Получение refresh-токена из базы не удалось")
	}
	return
}

func (ts *TokenService) ParseToken(tokenString string) (
	claims *t.JwtCustomClaims,
	err error,
) {
	token, err := jwt.ParseWithClaims(tokenString, &t.JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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

func (ts *TokenService) RemoveRefreshToken(tokenString string) (
	err error,
) {
	err = ts.dbService.RemoveRefreshToken(tokenString)
	return err
}

func (ts *TokenService) createAccessToken(userID uuid.UUID, expireMinutes int, secret string) (signedToken string,
	exp int64,
	err error,
) {
	exp = time.Now().Add(time.Minute * time.Duration(expireMinutes)).Unix()

	claims := t.JwtCustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err = jwtToken.SignedString([]byte(secret))

	return
}

func (ts *TokenService) createRefreshToken(userID uuid.UUID, expireMinutes int, secret string) (cryptedTokenString string,
	exp int64,
	err error,
) {
	exp = time.Now().Add(time.Minute * time.Duration(expireMinutes)).Unix()

	claims := &t.JwtCustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	//Создание нового токена
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	var signedToken string
	//Подпись токена зашифрованным паролем
	if signedToken, err = jwtToken.SignedString([]byte(secret)); err != nil {
		log.Print("Токен не был удачно подписан!")
		return
	}
	fmt.Printf("signedToken: %v\n", signedToken)

	//? Так ли нужно?
	// Шифрование в bcrypt подписанного токена
	var cryptedToken []byte
	if cryptedToken, err = bcrypt.GenerateFromPassword([]byte(signedToken), bcrypt.DefaultCost); err != nil {
		log.Print("Пароль не был удачно зашифрован!")
		return
	}
	cryptedTokenString = string(cryptedToken)
	return
}
