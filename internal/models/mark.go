package models

import (
	"calificationApi/internal/dto"
	"calificationApi/internal/services/search"
)

type Mark struct {
	ID        string `bson:"_id,omitempty"`
	StudentId string `bson:"student_id"`
	TeacherId string `bson:"teacher_id"`
	Grade     string `bson:"grade"`
	Mark      string `bson:"mark"`
	Semester  string `bson:"semester"`
}

type Marks []Mark

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

func (m *Marks) ToGetRequest() []dto.MarksGetRequest {
	var mapper []dto.MarksGetRequest
	for _, value := range *m {
		studentCarnet, studentErr := search.GetStudentCarnetById(value.StudentId)
		if studentErr != nil {
			return nil
		}
		teacherCarnet, teacherErr := search.GetTeacherCarnetById(value.TeacherId)
		if teacherErr != nil {
			return nil
		}
		markDto := value.ToGetRequest(studentCarnet, teacherCarnet)
		mapper = append(mapper, markDto)
	}
	return mapper
}
