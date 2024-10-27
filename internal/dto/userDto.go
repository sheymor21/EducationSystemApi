package dto

type UserAddRequest struct {
	Carnet   string
	Username string
	Password string
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
} // @name UserLogin
