package dto

type UserAddRequest struct {
	Carnet   string
	Username string
	Password string
}

type UserLoginRequest struct {
	Carnet   string `json:"carnet" validate:"required"`
	Password string `json:"password" validate:"required"`
} // @name UserLoginRequest
