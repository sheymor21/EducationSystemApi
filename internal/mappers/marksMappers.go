package mappers

import (
	"SchoolManagerApi/internal/dto"
	"SchoolManagerApi/internal/models"
	"SchoolManagerApi/internal/services/search"
	"SchoolManagerApi/internal/utilities"
	"errors"
)

func MarkToGetDto(m models.Mark) (dto.MarksGetRequest, error) {

	studentCarnet, studentErr := search.GetStudentCarnetById(m.StudentId)
	if studentErr != nil {
		return dto.MarksGetRequest{}, errors.New("student not found")
	}
	teacherCarnet, teacherErr := search.GetTeacherCarnetById(m.TeacherId)
	if teacherErr != nil {
		return dto.MarksGetRequest{}, errors.New("teacher not found")
	}

	var mapper dto.MarksGetRequest
	{
		mapper.Mark = m.Mark
		mapper.ID = m.ID
		mapper.Grade = m.Grade
		mapper.Semester = m.Semester
		mapper.StudentCarnet = studentCarnet
		mapper.TeacherCarnet = teacherCarnet
	}
	return mapper, nil
}

func MarkListToGetDto(marks []models.Mark) []dto.MarksGetRequest {
	var mapper []dto.MarksGetRequest
	for _, m := range marks {
		mark, err := MarkToGetDto(m)
		if err != nil {
			utilities.Log.Errorln(err)
			return nil
		}
		mapper = append(mapper, mark)
	}
	return mapper
}

func UpdateDtoToMark(dto dto.MarksUpdateRequest, id string) (models.Mark, error) {

	studentId, studentErr := search.GetStudentIdByCarnet(dto.StudentCarnet)
	if studentErr != nil {
		return models.Mark{}, errors.New("student not found")
	}
	teacherId, teacherErr := search.GetTeacherIdByCarnet(dto.TeacherCarnet)
	if teacherErr != nil {
		return models.Mark{}, errors.New("teacher not found")
	}

	var mapper models.Mark
	{
		mapper.ID = id
		mapper.Mark = dto.Mark
		mapper.Grade = dto.Grade
		mapper.Semester = dto.Semester
		mapper.StudentId = studentId
		mapper.TeacherId = teacherId
	}
	return mapper, nil
}
