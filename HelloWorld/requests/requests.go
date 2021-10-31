package requests

import "github.com/google/uuid"

type LoginRequest struct {
	UserID uuid.UUID `json:"userID"`
}

//TODO Подумать над изменением
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required" example:"refresh_token"`
}
