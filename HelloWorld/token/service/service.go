package service

import (
	"encoding/base64"
	"log"
	"time"

	dbs "auth/database/service"
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

func GenerateTokenPair(server *server.Server) (accessToken string,
	refreshToken string,
	exp int64,
	err error,
) {
	var refreshUID string
	userID := uuid.New().String()
	if accessToken, exp, err = createAccessToken(userID, ExpireAccessMinutes,
		secret); err != nil {
		return
	}

	if refreshToken, _, err = createRefreshToken(userID, ExpireRefreshMinutes,
		secret); err != nil {
		return
	}

	//Использовать кэшированные токены для добавления в базу. Вернее, один refresh
	dbService := dbs.NewService(server)
	refreshTokenData := t.RefreshToken{RefreshUID: refreshUID, TokenText: refreshToken}
	dbService.InsertToken(refreshTokenData)
	return
}

//TODO стоит поменять способ создания. Каждый вид токена должен создаваться по-разному
func createAccessToken(userID string, expireMinutes int, secret string) (signedToken string,
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

func createRefreshToken(userID string, expireMinutes int, secret string) (cypheredToken string,
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

	//? Так ли нужно?
	//Шифрование в bcrypt пароля
	var cryptedPassword []byte
	if cryptedPassword, err = bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost); err != nil {
		log.Print("Пароль не был удачно зашифрован!")
		return
	}

	//Создание нового токена
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	var signedToken string
	//Подпись токена зашифрованным паролем
	if signedToken, err = jwtToken.SignedString([]byte(cryptedPassword)); err != nil {
		log.Print("Токен не был удачно подписан!")
		return
	}
	//Шифрование в base64 для передачи
	cypheredToken = base64.StdEncoding.EncodeToString([]byte(signedToken))
	return
}
