package token

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	UserID string `json:"userid"`
	ID     string `json:"id"`
	jwt.StandardClaims
}

type RefreshToken struct {
	RefreshUID string `json:"refresh"`
	TokenText  string `json:"tokenText"`
}
