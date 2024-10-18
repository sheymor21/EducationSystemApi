package dto

type MarkAddRequest struct {
	StudentCarnet string `json:"student_carnet" validate:"required,min=10,max=10"`
	TeacherCarnet string `json:"teacher_carnet" validate:"required"`
	Grade         string `json:"grade" validate:"required"`
	Mark          string `json:"mark" validate:"required"`
	Semester      string `json:"semester" validate:"required"`
}

type MarksGetRequest struct {
	ID            string `bson:"_id,omitempty"`
	StudentCarnet string `json:"student_carnet,omitempty"`
	TeacherCarnet string `json:"teacher_carnet,omitempty"`
	Grade         string `json:"grade,omitempty"`
	Mark          string `json:"mark,omitempty"`
	Semester      string `json:"semester,omitempty"`
}

type MarksUpdateRequest struct {
	StudentCarnet string
	TeacherCarnet string
	Grade         string
	Mark          string
	Semester      string
}
