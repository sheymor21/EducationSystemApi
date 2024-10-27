package services

import (
	"SchoolManagerApi/internal/database"
	"SchoolManagerApi/internal/dto"
	"SchoolManagerApi/internal/models"
	"SchoolManagerApi/internal/server/customErrors"
	"SchoolManagerApi/internal/utilities"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sync"
)

// addStudent godoc
// @Summary Add a student
// @Description Add a new student to the database
// @Tags student
// @Accept  json
// @Produce  json
// @Param student body StudentAddRequest true "Add Student"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /student [post]
func addStudent(w http.ResponseWriter, r *http.Request) {
	studentDto := dto.StudentAddDto{}
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
	_, err = dbContext.Student.InsertOne(context.TODO(), student)
	if err != nil {
		utilities.Log.Println(err)
		httpInternalError(w, err.Error())
	_, insertStudentErr := dbContext.Student.InsertOne(context.TODO(), student)
	if insertStudentErr != nil {
		utilities.Log.Errorln(insertStudentErr)
		httpInternalError(w, insertStudentErr.Error())
		return
	}
	}
}

// getStudents retrieves a list of students from the database and writes it as a JSON response.
// @Summary Retrieves students
// @Description Fetches a list of all students from the database and returns the data as JSON.
// @Tags students
// @Produce json
// @Success 200 {array} models.Student
// @Failure 500 {string} string "Internal Server Error"
// @Router /students [get]
func getStudents(w http.ResponseWriter) {
	var student []models.Student
	find, findErr := dbContext.Student.Find(context.TODO(), bson.M{})
	if findErr != nil {
		utilities.Log.Errorln(findErr)
		httpInternalError(w, findErr.Error())
		return
	}
	decodeErr := find.All(context.TODO(), &student)
	if decodeErr != nil {
		utilities.Log.Errorln(decodeErr)
		httpInternalError(w, decodeErr.Error())
		return
	}
	utilities.WriteJson(w, http.StatusOK, student)
}

// getStudent godoc
// @Summary Get student by carnet
// @Description Retrieve a student's information from the database using their carnet
// @Tags student
// @Accept json
// @Produce json
// @Param Carnet query string true "Student Carnet"
// @Success 200 {object} models.Student
// @Failure 404 {object} string "Student not found"
// @Router /student [get]
func getStudent(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	student, err := getStudentByCarnet(dbContext, carnet)
	if err != nil {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
		return
	} else {

		utilities.WriteJson(w, http.StatusOK, student)
	}
}

// putStudent godoc
// @Summary Update a student
// @Description Update an existing student's information in the database
// @Tags student
// @Accept json
// @Produce json
// @Param Carnet query string true "Student Carnet"
// @Param student body models.Student true "Update Student"
// @Success 200 {object} models.Student
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /student [put]
func putStudent(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	ch := make(chan bool)
	carnet := r.URL.Query().Get("Carnet")
	wg.Add(1)
	go anyStudent(carnet, &wg, ch)
	if !<-ch {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
		return
	}
	wg.Wait()

	var student models.Student
	student.Carnet = carnet
	err := utilities.ReadJson(w, r, &student)
	if err != nil {
		httpInternalError(w, err.Error())
		utilities.Log.Errorln(err)
		return
	}
	filter := bson.D{{"carnet", carnet}}
	update := bson.D{{"$set", student}}
	_, err = dbContext.Student.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		utilities.Log.Errorln(err)
		httpInternalError(w, err.Error())
	}

}

// deleteStudent godoc
// @Summary Delete a student
// @Description Delete a student from the database using their carnet
// @Tags student
// @Accept  json
// @Produce  json
// @Param Carnet query string true "Student Carnet"
// @Success 200 "Student deleted successfully"
// @Failure 404 {object} string "Student not found"
// @Failure 500 {object} string "Internal server error"
// @Router /student [delete]
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	ch := make(chan bool)
	carnet := r.URL.Query().Get("Carnet")
	wg.Add(1)
	go anyStudent(carnet, &wg, ch)
	if <-ch {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
		return
	}
	wg.Wait()

	filter := bson.D{{"carnet", carnet}}
	_, err := dbContext.Student.DeleteOne(context.TODO(), filter)
	if err != nil {
		utilities.Log.Errorln(err)
		httpInternalError(w, err.Error())
		return
	}
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

func anyStudent(carnet string, wg *sync.WaitGroup, ch chan bool) {
	defer wg.Done()
	defer close(ch)
	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Student.FindOne(context.TODO(), filter).Decode(&models.Student{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		ch <- false
	} else if err != nil {
		ch <- false
		utilities.Log.Errorln(err)
	} else {
		ch <- true
	}
}
