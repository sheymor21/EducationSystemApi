package services

import (
	"calificationApi/internal/Utilities"
	"calificationApi/internal/database"
	"calificationApi/internal/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func HttpStudentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getStudent(w, r)
	case http.MethodPost:
		addStudent(w, r)
	case http.MethodPut:
		putStudent(w, r)
	case http.MethodDelete:
		deleteStudent(w, r)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func addStudent(w http.ResponseWriter, r *http.Request) {

	student := models.Student{}
	Utilities.ReadJson(r.Body, &student)
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())
	_, err := dbContext.Users.InsertOne(context.TODO(), student)
	if err != nil {
		log.Fatal(err)
	}

}

func getStudent(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	_, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	student, err := getStudentByCarnet(carnet)
	if errors.Is(err, mongo.ErrNoDocuments) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else if err != nil {
		log.Fatal(err)
	} else {

		Utilities.WriteJson(w, http.StatusOK, student)
	}
}

func putStudent(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	anyStudent := anyStudent(carnet)
	if !anyStudent {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {

		var student models.Student
		Utilities.ReadJson(r.Body, &student)
		dbContext, client := database.GetDatabaseConnection()
		defer database.CloseConnection(client, context.TODO())
		filter := bson.D{{"carnet", carnet}}
		update := bson.D{{"$set", student}}
		_, err := dbContext.Users.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			log.Fatal(err)
		}
	}

}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	carnet := r.URL.Query().Get("Carnet")
	anyStudent := anyStudent(carnet)
	if !anyStudent {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {

		filter := bson.D{{"carnet", carnet}}
		_, err := dbContext.Users.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func getStudentByCarnet(carnet string) (models.Student, error) {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	var student models.Student
	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Users.FindOne(context.TODO(), filter).Decode(&student)
	return student, err
}

func anyStudent(carnet string) bool {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Users.FindOne(context.TODO(), filter).Decode(&models.Student{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false
	} else if err != nil {
		log.Fatal(err)
		return false
	} else {

		return true
	}
}
