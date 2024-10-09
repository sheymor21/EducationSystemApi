package services

import (
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

func addTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher models.Teacher
	err := utilities.ReadJson(w, r, &teacher)
	if err != nil {
		httpInternalError(w, err.Error())
		log.Println(err)
		return
	}
	_, err = dbContext.Teachers.InsertOne(context.TODO(), teacher)
	if err != nil {
		httpInternalError(w, err.Error())
	}
}

func updateTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher models.Teacher
	err := utilities.ReadJson(w, r, &teacher)
	if err != nil {
		httpInternalError(w, err.Error())
		log.Println(err)
		return
	}
	filter := bson.M{"carnet": teacher.Carnet}
	update := bson.M{"$set": teacher}
	_, err = dbContext.Teachers.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		httpInternalError(w, err.Error())
		return
	}
}

func deleteTeacher(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	filter := bson.M{"carnet": carnet}

	_, err := dbContext.Teachers.DeleteOne(context.TODO(), filter)
	if err != nil {
		httpInternalError(w, err.Error())
		return
	}

}

func getTeacher(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	ch := make(chan bool)
	carnet := r.URL.Query().Get("Carnet")
	go anyTeacher(carnet, &wg, ch)
	if <-ch {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Error())
		return
	}
	wg.Wait()
	var teacher models.Teacher
	filter := bson.M{"carnet": carnet}
	err := dbContext.Teachers.FindOne(context.TODO(), filter).Decode(&teacher)
	if err != nil {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Error())
		return
	}
	utilities.WriteJson(w, http.StatusOK, teacher)

}

func getTeachers(w http.ResponseWriter, r *http.Request) {
	var teachers []models.Teacher
	find, findErr := dbContext.Teachers.Find(context.TODO(), bson.M{})
	if findErr != nil {
		log.Println(findErr)
		return
	}
	decodeErr := find.All(context.TODO(), &teachers)
	if decodeErr != nil {
		log.Println(decodeErr)
		return
	}
	utilities.WriteJson(w, http.StatusOK, teachers)
}
func anyTeacher(carnet string, wg *sync.WaitGroup, ch chan bool) {
	defer wg.Done()
	defer close(ch)
	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Teachers.FindOne(context.TODO(), filter).Decode(&models.Teacher{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		ch <- false
	} else if err != nil {
		log.Println(err)
		ch <- false
	} else {
		ch <- true
	}
}
