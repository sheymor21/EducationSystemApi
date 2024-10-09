package models

import (
	"calificationApi/internal/dto"
)

type Mark struct {
	ID        string `bson:"_id,omitempty"`
	StudentId string `bson:"student_id"`
	TeacherId string `bson:"teacher_id"`
	Grade     string `bson:"grade"`
	Mark      string `bson:"mark"`
	Semester  string `bson:"semester"`
}

func (m *Mark) ToGetRequest(studentCarnet string, teacherCarnet string) dto.MarksGetRequest {
	var mapper dto.MarksGetRequest
	{
		mapper.Mark = m.Mark
		mapper.ID = m.ID
		mapper.Grade = m.Grade
		mapper.Semester = m.Semester
		mapper.StudentCarnet = studentCarnet
		mapper.TeacherCarnet = teacherCarnet
	}
	return mapper
}
