package mappers

import (
	"SchoolManagerApi/internal/dto"
	"SchoolManagerApi/internal/models"
)

func TeacherToGetRequest(model models.Teacher) dto.TeacherGetRequest {
	return dto.TeacherGetRequest{
		Carnet:    model.Carnet,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Age:       model.Age,
		Classroom: model.Classroom,
	}
}

func TeacherListToGetRequest(model []models.Teacher) []dto.TeacherGetRequest {

	var mapper []dto.TeacherGetRequest
	for _, v := range model {
		mapper = append(mapper, TeacherToGetRequest(v))
	}
	return mapper
}
