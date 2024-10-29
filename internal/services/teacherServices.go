package services

import (
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

// @Summary Add a new teacher
// @Description Inserts a new teacher record to the database
// @Accept json
// @Produce json
// @Param teacher body TeacherAddRequest true "New Teacher"
// @Success 200
// @Failure 500 string error
// @Router /teacher [post]
// @Tags teacher
func addTeacher(w http.ResponseWriter, r *http.Request) {
	var teacherDto dto.TeacherAddRequest
	jsonErr := utilities.ReadJson(w, r, &teacherDto)
	customErrors.ThrowHttpError(jsonErr, w, "", http.StatusInternalServerError)
	teacher := mappers.TeacherAddToModel(teacherDto)
	_, dbErr := dbContext.Teachers.InsertOne(context.TODO(), teacher)
	customErrors.ThrowHttpError(dbErr, w, "", http.StatusInternalServerError)

	userErr := addUser(teacher.FirstName, teacher.LastName, teacher.Carnet, validations.TeacherRol)
	customErrors.ThrowHttpError(userErr, w, "", http.StatusInternalServerError)
}

// @Summary Update an existing teacher
// @Description Updates the information of an existing teacher in the database
// @Accept json
// @Produce json
// @Param teacher body TeacherUpdateRequest true "Updated Teacher"
// @Success 200
// @Failure 500 string error
// @Router /teacher [put]
// @Tags teacher
func updateTeacher(w http.ResponseWriter, r *http.Request) {
	var teacherDto dto.TeacherUpdateRequest
	jsonErr := utilities.ReadJson(w, r, &teacherDto)
	customErrors.ThrowHttpError(jsonErr, w, "", http.StatusInternalServerError)

	teacher := mappers.TeacherUpdateToModel(teacherDto)
	filter := bson.M{"carnet": teacher.Carnet}
	update := bson.M{"$set": teacher}
	_, dbErr := dbContext.Teachers.UpdateOne(context.TODO(), filter, update)
	customErrors.ThrowHttpError(dbErr, w, "", http.StatusInternalServerError)
}

// @Summary Delete a teacher
// @Description Deletes an existing teacher record from the database using the "Carnet" query parameter.
// @Param Carnet query string true "Teacher Carnet"
// @Success 204 "No Content"
// @Failure 500 string error
// @Router /teacher [delete]
// @Tags teacher
func deleteTeacher(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	filter := bson.M{"carnet": carnet}

	_, dbErr := dbContext.Teachers.DeleteOne(context.TODO(), filter)
	customErrors.ThrowHttpError(dbErr, w, "", http.StatusInternalServerError)
	utilities.WriteJson(w, http.StatusNoContent, nil)

}

// @Summary Get a teacher's details
// @Description Fetches the information of a teacher from the database using the "Carnet" query parameter.
// @Param Carnet query string true "Teacher Carnet"
// @Success 200 {object} TeacherGetRequest
// @Failure 404 {string} string "Teacher Not Found"
// @Failure 500 string error
// @Router /teacher [get]
// @Tags teacher
func getTeacher(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	teacherExist := anyTeacher(carnet)
	if teacherExist {
		http.Error(w, "Not found this carnet", http.StatusNotFound)
		return
	}
	var teacher models.Teacher
	filter := bson.M{"carnet": carnet}
	dbErr := dbContext.Teachers.FindOne(context.TODO(), filter).Decode(&teacher)
	customErrors.ThrowHttpError(dbErr, w, "Not found this carnet", http.StatusNotFound)
	teacherDto := mappers.TeacherToGetRequest(teacher)
	utilities.WriteJson(w, http.StatusOK, teacherDto)

}

// @Summary Retrieve all teachers
// @Description Fetch all teacher records from the database and return them as a JSON payload
// @Success 200 {array} TeacherGetRequest "List of teachers"
// @Failure 500 string error
// @Router /teachers [get]
// @Tags teachers
func getTeachers(w http.ResponseWriter) {
	var teachers []models.Teacher
	find, findErr := dbContext.Teachers.Find(context.TODO(), bson.M{})
	if findErr != nil {
		utilities.Log.Errorln(findErr)
		return
	}
	decodeErr := find.All(context.TODO(), &teachers)
	if decodeErr != nil {
		utilities.Log.Errorln(decodeErr)
		return
	}
	teacherDto := mappers.TeacherListToGetRequest(teachers)
	utilities.WriteJson(w, http.StatusOK, teacherDto)
}

func anyTeacher(carnet string) bool {
	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Teachers.FindOne(context.TODO(), filter).Decode(&models.Teacher{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false
	} else if err != nil {
		utilities.Log.Errorln(err)
		return false
	} else {
		return true
	}
}
