package requests

import "github.com/google/uuid"

type LoginRequest struct {
	UserID uuid.UUID `json:"userID"`
}

type RefreshRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
