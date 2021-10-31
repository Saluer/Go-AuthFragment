package token

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtCustomClaims struct {
	UserID uuid.UUID `json:"userid"`
	jwt.StandardClaims
}

type RefreshToken struct {
	RefreshUID string `json:"refresh" bson:"refreshuid"`
}
