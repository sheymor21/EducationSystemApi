package mappers

import (
	"SchoolManagerApi/internal/dto"
	"SchoolManagerApi/internal/models"
)

func StudentToGetRequest(s models.Student) dto.StudentGetRequest {
	return dto.StudentGetRequest{
		Carnet:    s.Carnet,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Age:       s.Age,
		Classroom: s.Classroom,
	}
}

func StudentListToGetRequest(s []models.Student) []dto.StudentGetRequest {
	var mapper []dto.StudentGetRequest
	for _, v := range s {
		mapper = append(mapper, StudentToGetRequest(v))
	}
	return mapper
}

func StudentUpdateToModel(d dto.StudentUpdateRequest, carnet string) models.Student {
	return models.Student{
		Carnet:    carnet,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Age:       d.Age,
		Classroom: d.Classroom,
	}
}

func StudentAddToModel(s dto.StudentAddRequest) models.Student {
	return models.Student{

		Carnet:    s.Carnet,
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Age:       s.Age,
		Classroom: s.Classroom,
	}
}
