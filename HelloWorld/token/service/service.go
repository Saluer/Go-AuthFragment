package service

import (
	"fmt"
	"time"

	dbs "auth/database/service"
	"auth/server"
	t "auth/token"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
	fmt.Println("В функции генерации токенов")
	var refreshUID string
	userID := uuid.New().String()
	if accessToken, exp, err = createToken(userID, ExpireAccessMinutes,
		secret); err != nil {
		return
	}

	if refreshToken, _, err = createToken(userID, ExpireRefreshMinutes,
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
func createToken(userID string, expireMinutes int, secret string) (token string,
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
	token, err = jwtToken.SignedString([]byte(secret))

	return
}
