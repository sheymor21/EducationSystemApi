package dto

type UserAddRequest struct {
	Carnet   string
	Username string
	Password string
}

type UserLoginRequest struct {
	Carnet   string `json:"carnet"`
	Password string `json:"password"`
} // @name UserLogin
