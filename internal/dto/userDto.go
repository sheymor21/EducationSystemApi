package dto

type UserLoginRequest struct {
	Carnet   string `json:"carnet" validate:"required"`
	Password string `json:"password" validate:"required"`
} // @name UserLoginRequest
