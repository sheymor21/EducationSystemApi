package services

import (
	"SchoolManagerApi/internal/database"
	"SchoolManagerApi/internal/dto"
	"SchoolManagerApi/internal/mappers"
	"SchoolManagerApi/internal/models"
	"SchoolManagerApi/internal/server/customErrors"
	"SchoolManagerApi/internal/utilities"
	"SchoolManagerApi/internal/validations"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// @Summary Add a student
// @Description Add a new student to the database
// @Tags student
// @Accept  json
// @Produce  json
// @Param student body StudentAddRequest true "Add new Student"
// @Success 200
// @Failure 400 string error
// @Failure 500 string error
// @Router /student [post]
func addStudent(w http.ResponseWriter, r *http.Request) {
	studentDto := dto.StudentAddRequest{}
	jsonErr := utilities.ReadJson(w, r, &studentDto)
	if jsonErr != nil {
		httpInternalError(w, jsonErr.Error())
		utilities.Log.Errorln(jsonErr)
		return
	}
	student := models.Student{
		Carnet:    studentDto.Carnet,
		FirstName: studentDto.FirstName,
		LastName:  studentDto.LastName,
		Age:       studentDto.Age,
		Classroom: studentDto.Classroom,
	}
	_, insertStudentErr := dbContext.Student.InsertOne(context.TODO(), student)
	if insertStudentErr != nil {
		utilities.Log.Errorln(insertStudentErr)
		httpInternalError(w, insertStudentErr.Error())
		return
	}
	userErr := addUser(student.FirstName, student.LastName, student.Carnet, validations.TeacherRol)
	if userErr != nil {
		utilities.Log.Errorln(userErr)
		httpInternalError(w, userErr.Error())
		return
	}

}

// @Summary Retrieves students
// @Description Fetches a list of all students from the database and returns the data as JSON.
// @Tags students
// @Produce json
// @Success 200 {array} StudentGetRequest
// @Failure 500 string error
// @Router /students [get]
func getStudents(w http.ResponseWriter) {
	var students []models.Student
	find, findErr := dbContext.Student.Find(context.TODO(), bson.M{})
	if findErr != nil {
		utilities.Log.Errorln(findErr)
		httpInternalError(w, findErr.Error())
		return
	}
	decodeErr := find.All(context.TODO(), &students)
	if decodeErr != nil {
		utilities.Log.Errorln(decodeErr)
		httpInternalError(w, decodeErr.Error())
		return
	}
	studentsDto := mappers.StudentListToGetRequest(students)
	utilities.WriteJson(w, http.StatusOK, studentsDto)
}

// @Summary Get student by carnet
// @Description Retrieve a student's information from the database using their carnet
// @Tags student
// @Accept json
// @Produce json
// @Param Carnet query string true "StudentRol Carnet"
// @Success 200 {array} StudentGetRequest
// @Failure 404 {object} string "Student not found"
// @Router /student [get]
func getStudent(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	student, err := getStudentByCarnet(dbContext, carnet)
	studentDto := mappers.StudentToGetRequest(student)
	if err != nil {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
		return
	} else {

		utilities.WriteJson(w, http.StatusOK, studentDto)
	}
}

// @Summary Update a student
// @Description Update an existing student's information in the database
// @Tags student
// @Accept json
// @Produce json
// @Param Carnet query string true "Student Carnet"
// @Param student body StudentGetRequest true "Update Student"
// @Success 200
// @Failure 400 string error
// @Failure 404 {string} string "Student Not Found"
// @Failure 500 string error
// @Router /student [put]
func putStudent(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	studentExist := anyStudent(carnet)
	if studentExist {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
		return
	}

	var studentDto dto.StudentUpdateRequest
	err := utilities.ReadJson(w, r, &studentDto)
	if err != nil {
		httpInternalError(w, err.Error())
		utilities.Log.Errorln(err)
		return
	}
	student := mappers.StudentUpdateToModel(studentDto, carnet)
	filter := bson.D{{"carnet", carnet}}
	update := bson.D{{"$set", student}}
	_, err = dbContext.Student.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		utilities.Log.Errorln(err)
		httpInternalError(w, err.Error())
	}

}

// @Summary Delete a student
// @Description Delete a student from the database using their carnet
// @Tags student
// @Accept  json
// @Produce  json
// @Param Carnet query string true "Student Carnet"
// @Success 204
// @Failure 404 {object} string "Student Not Found"
// @Failure 500 string error
// @Router /student [delete]
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	studentExist := anyStudent(carnet)
	if studentExist {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
		return
	}

	filter := bson.D{{"carnet", carnet}}
	_, err := dbContext.Student.DeleteOne(context.TODO(), filter)
	if err != nil {
		utilities.Log.Errorln(err)
		httpInternalError(w, err.Error())
		return
	}
	utilities.WriteJson(w, http.StatusNoContent, nil)
}

func getStudentByCarnet(dbContext *database.MongoContext, carnet string) (models.Student, *customErrors.NotFoundMongoError) {

	var student models.Student
	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Student.FindOne(context.TODO(), filter).Decode(&student)
	if err != nil {
		return models.Student{}, customErrors.NewNotFoundMongoError("Carnet")
	}
	return student, nil
}

func anyStudent(carnet string) bool {
	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Student.FindOne(context.TODO(), filter).Decode(&models.Student{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false
	} else if err != nil {
		utilities.Log.Errorln(err)
		return false
	} else {
		return true
	}
}
