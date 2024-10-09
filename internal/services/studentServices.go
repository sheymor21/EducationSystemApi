package services

import (
	"calificationApi/internal/database"
	"calificationApi/internal/dto"
	"calificationApi/internal/models"
	"calificationApi/internal/server/customErrors"
	"calificationApi/internal/utilities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"sync"
)

func addStudent(w http.ResponseWriter, r *http.Request) {
	studentDto := dto.StudentAddDto{}
	err := utilities.ReadJson(w, r, &studentDto)
	if err != nil {
		httpInternalError(w, err.Error())
		log.Println(err)
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
		log.Println(err)
		httpInternalError(w, err.Error())
	}
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	var student []models.Student
	find, findErr := dbContext.Student.Find(context.TODO(), bson.M{})
	if findErr != nil {
		log.Println(findErr)
		httpInternalError(w, findErr.Error())
		return
	}
	decodeErr := find.All(context.TODO(), &student)
	if decodeErr != nil {
		log.Println(decodeErr)
		httpInternalError(w, decodeErr.Error())
		return
	}
	utilities.WriteJson(w, http.StatusOK, student)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	student, err := getStudentByCarnet(dbContext, carnet)
	if err != nil {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
	} else {

		utilities.WriteJson(w, http.StatusOK, student)
	}
}

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
		log.Println(err)
		return
	}
	filter := bson.D{{"carnet", carnet}}
	update := bson.D{{"$set", student}}
	_, err = dbContext.Student.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
		httpInternalError(w, err.Error())
	}

}

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
		log.Println(err)
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
		log.Println(err)
	} else {
		ch <- true
	}
}
