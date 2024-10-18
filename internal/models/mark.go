package models

type Mark struct {
	ID        string `bson:"_id,omitempty"`
	StudentId string `bson:"student_id"`
	TeacherId string `bson:"teacher_id"`
	Grade     string `bson:"grade"`
	Mark      string `bson:"mark"`
	Semester  string `bson:"semester"`
}
