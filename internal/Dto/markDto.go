package Dto

type MarkAddRequest struct {
	StudentCarnet string
	TeacherCarnet string
	Grade         string
	Mark          string
	Semester      string
}

type MarksGetRequest struct {
	ID            string `bson:"_id,omitempty"`
	StudentCarnet string `json:"student_carnet,omitempty"`
	TeacherCarnet string `json:"teacher_carnet,omitempty"`
	Grade         string `json:"grade,omitempty"`
	Mark          string `json:"mark,omitempty"`
	Semester      string `json:"semester,omitempty"`
}
