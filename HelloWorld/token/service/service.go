package service

import (
	"encoding/base64"
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
	tokenFromDB := ts.dbService.GetRefreshToken(tokenData)
	var token string
	if token, err = ts.decodeAfterTransfer(tokenFromDB.RefreshUID); err != nil {
		log.Fatal("Декодирование после получения прошло неудачно!")
	}
	log.Println(token)
	test, _ := ts.decodeAfterTransfer(tokenData)
	log.Println(token == test)
	// cachedTokens := new(CachedTokens)
	// err = json.Unmarshal([]byte(cacheJSON), cachedTokens)

	// if err != nil || tokenUID != claims.UID {
	// 	return errors.New("token not found")
	// }

	return
}

//TODO стоит поменять способ создания. Каждый вид токена должен создаваться по-разному
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

func (ts *TokenService) createRefreshToken(userID uuid.UUID, expireMinutes int, secret string) (cypheredToken string,
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
	//Шифрование токена в base64 для передачи
	cypheredToken = ts.encodeToTransfer(string(cryptedToken))
	return
}

func (ts *TokenService) encodeToTransfer(signedToken string) (cypheredToken string) {
	cypheredToken = base64.StdEncoding.EncodeToString([]byte(signedToken))
	return
}

func (ts *TokenService) decodeAfterTransfer(cypheredToken string) (signedToken string, err error) {
	res, err := base64.StdEncoding.DecodeString(cypheredToken)
	signedToken = string(res)
	return signedToken, err
}
