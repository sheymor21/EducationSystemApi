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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()
	wg.Wait()
}

func updateTeacher(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
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
		}
	}()
	wg.Wait()
}

func deleteTeacher(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		carnet := r.URL.Query().Get("Carnet")
		filter := bson.M{"carnet": carnet}

		_, err := dbContext.Teachers.DeleteOne(context.TODO(), filter)
		if err != nil {
			httpInternalError(w, err.Error())
			return
		}
	}()
	wg.Wait()

}

func getTeacher(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		carnet := r.URL.Query().Get("Carnet")
		anyTeacher := anyTeacher(carnet)
		if !anyTeacher {
			httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Error())
		} else {
			var teacher models.Teacher
			filter := bson.M{"carnet": carnet}
			err := dbContext.Teachers.FindOne(context.TODO(), filter).Decode(&teacher)
			if err != nil {
				httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Error())
			}
			utilities.WriteJson(w, http.StatusOK, teacher)
		}
	}()
	wg.Wait()

}
func anyTeacher(carnet string) bool {

	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Teachers.FindOne(context.TODO(), filter).Decode(&models.Teacher{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false
	} else if err != nil {
		log.Fatal(err)
		return false
	} else {

		return true
	}
}
