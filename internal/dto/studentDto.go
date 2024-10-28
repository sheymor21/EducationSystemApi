package dto

type StudentAddRequest struct {
	Carnet    string `json:"carnet" validate:"required,min=10,max=10"`
	FirstName string `json:"firstName" validate:"required,min=3,max=20"`
	LastName  string `json:"lastName" validate:"required,min=3,max=20"`
	Age       uint8  `json:"age" validate:"required,min=1,max=100"`
	Classroom string `json:"classroom" validate:"required,min=2,max=4"`
} // @name StudentAddRequest
