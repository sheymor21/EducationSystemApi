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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
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
			log.Fatal(err)
		}
	}()
	wg.Wait()
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		carnet := r.URL.Query().Get("Carnet")
		student, err := getStudentByCarnet(dbContext, carnet)
		if err != nil {
			httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
		} else {

			utilities.WriteJson(w, http.StatusOK, student)
		}
	}()
	wg.Wait()
}

func putStudent(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		carnet := r.URL.Query().Get("Carnet")
		anyStudent := anyStudent(carnet)
		if !anyStudent {
			httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
		} else {

			var student models.Student
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
				httpInternalError(w, err.Error())
				log.Fatal(err)
			}
		}
	}()
	wg.Wait()

}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		carnet := r.URL.Query().Get("Carnet")
		anyStudent := anyStudent(carnet)
		if !anyStudent {
			httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
		} else {

			filter := bson.D{{"carnet", carnet}}
			_, err := dbContext.Student.DeleteOne(context.TODO(), filter)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}()
	wg.Wait()
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
		log.Fatal(err)
		return false
	} else {

		return true
	}
}
