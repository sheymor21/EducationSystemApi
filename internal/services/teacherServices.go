package services

import (
	"calificationApi/internal/Utilities"
	"calificationApi/internal/database"
	"calificationApi/internal/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func HttpTeacherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeacher(w, r)
	case http.MethodPost:
		addTeacher(w, r)
	case http.MethodPut:
		updateTeacher(w, r)
	case http.MethodDelete:
		deleteTeacher(w, r)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func addTeacher(w http.ResponseWriter, r *http.Request) {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	var teacher models.Teacher
	err := Utilities.ReadJson(w, r, &teacher)
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

	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	var teacher models.Teacher
	err := Utilities.ReadJson(w, r, &teacher)
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
}

func deleteTeacher(w http.ResponseWriter, r *http.Request) {

	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	carnet := r.URL.Query().Get("Carnet")
	filter := bson.M{"carnet": carnet}

	_, err := dbContext.Teachers.DeleteOne(context.TODO(), filter)
	if err != nil {
		httpInternalError(w, err.Error())
		return
	}
}

func getTeacher(w http.ResponseWriter, r *http.Request) {

	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	carnet := r.URL.Query().Get("Carnet")
	anyTeacher := anyTeacher(carnet)
	if !anyTeacher {
		httpNotFoundError(w, NewNotFoundMongoError("carnet").Error())
	} else {
		var teacher models.Teacher
		filter := bson.M{"carnet": carnet}
		err := dbContext.Teachers.FindOne(context.TODO(), filter).Decode(&teacher)
		if err != nil {
			httpNotFoundError(w, NewNotFoundMongoError("carnet").Error())
		}
		Utilities.WriteJson(w, http.StatusOK, teacher)
	}
}

func anyTeacher(carnet string) bool {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

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

func getTeacherIdByCarnet(dbContext *database.MongoClient, carnet string) (string, error) {
	var result struct {
		Id string `bson:"_id"`
	}

	filter := bson.D{{"carnet", carnet}}
	projection := bson.D{{"_id", 1}}
	op := options.FindOne().SetProjection(projection)
	err := dbContext.Teachers.FindOne(context.TODO(), filter, op).Decode(&result)
	if err != nil {

		return "", NewNotFoundMongoError("Carnet")
	}
	return result.Id, nil

}

func getTeacherCarnetById(id string) (string, error) {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	var result struct {
		Carnet string `bson:"carnet"`
	}
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}
	filter := bson.D{{"_id", hex}}
	projection := bson.D{{"carnet", 1}}
	op := options.FindOne().SetProjection(projection)
	err = dbContext.Teachers.FindOne(context.TODO(), filter, op).Decode(&result)
	if err != nil {
		return "", NewNotFoundMongoError("Carnet")
	}
	return result.Carnet, nil
}
