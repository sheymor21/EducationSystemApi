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
	err := Utilities.ReadJson(w, r, &student)
	if err != nil {
		httpInternalError(w, err.Error())
		log.Println(err)
		return
	}
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	_, err = dbContext.Student.InsertOne(context.TODO(), student)
	if err != nil {
		log.Fatal(err)
	}

}

func getStudent(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	student, err := getStudentByCarnet(dbContext, carnet)
	if err != nil {
		httpNotFoundError(w, NewNotFoundMongoError("carnet").msg)
	} else {

		Utilities.WriteJson(w, http.StatusOK, student)
	}
}

func putStudent(w http.ResponseWriter, r *http.Request) {
	carnet := r.URL.Query().Get("Carnet")
	anyStudent := anyStudent(carnet)
	if !anyStudent {
		httpNotFoundError(w, NewNotFoundMongoError("carnet").msg)
	} else {

		var student models.Student
		err := Utilities.ReadJson(w, r, &student)
		if err != nil {
			httpInternalError(w, err.Error())
			log.Println(err)
			return
		}
		dbContext, client := database.GetDatabaseConnection()
		defer database.CloseConnection(client, context.TODO())
		filter := bson.D{{"carnet", carnet}}
		update := bson.D{{"$set", student}}
		_, err = dbContext.Student.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			httpInternalError(w, err.Error())
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
		httpNotFoundError(w, NewNotFoundMongoError("carnet").msg)
	} else {

		filter := bson.D{{"carnet", carnet}}
		_, err := dbContext.Student.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func getStudentByCarnet(dbContext *database.MongoClient, carnet string) (models.Student, *NotFoundMongoError) {

	var student models.Student
	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Student.FindOne(context.TODO(), filter).Decode(&student)
	if err != nil {
		return models.Student{}, NewNotFoundMongoError("Carnet")
	}
	return student, nil
}

func anyStudent(carnet string) bool {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

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

func getStudentIdByCarnet(dbContext *database.MongoClient, carnet string) (string, error) {
	var result struct {
		Id string `bson:"_id"`
	}

	filter := bson.D{{"carnet", carnet}}
	projection := bson.D{{"_id", 1}}
	op := options.FindOne().SetProjection(projection)
	err := dbContext.Student.FindOne(context.TODO(), filter, op).Decode(&result)
	if err != nil {
		return "", NewNotFoundMongoError("Carnet")

	}
	return result.Id, nil
}

func getStudentCarnetById(dbContext *database.MongoClient, id string) (string, *NotFoundMongoError) {
	var result struct {
		Carnet string `bson:"carnet"`
	}

	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}

	filter := bson.D{{"_id", hex}}
	projection := bson.D{{"carnet", 1}}
	op := options.FindOne().SetProjection(projection)
	err = dbContext.Student.FindOne(context.TODO(), filter, op).Decode(&result)
	if err != nil {
		return "", NewNotFoundMongoError("Id")
	}
	return result.Carnet, nil
}
