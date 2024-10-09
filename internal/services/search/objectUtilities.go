package search

import (
	"calificationApi/internal/database"
	"calificationApi/internal/server/customErrors"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var dbContext = database.GetMongoContext()

func GetTeacherIdByCarnet(carnet string) (string, error) {
	var result struct {
		Id string `bson:"_id"`
	}
	filter := bson.D{{"carnet", carnet}}
	projection := bson.D{{"_id", 1}}
	op := options.FindOne().SetProjection(projection)
	err := dbContext.Teachers.FindOne(context.TODO(), filter, op).Decode(&result)
	if err != nil {

		return "", customErrors.NewNotFoundMongoError("Carnet")
	}
	return result.Id, nil

}

func GetTeacherCarnetById(id string) (string, error) {

	var result struct {
		Carnet string `bson:"carnet"`
	}
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return "", err
	}
	filter := bson.D{{"_id", hex}}
	projection := bson.D{{"carnet", 1}}
	op := options.FindOne().SetProjection(projection)
	err = dbContext.Teachers.FindOne(context.TODO(), filter, op).Decode(&result)
	if err != nil {
		return "", customErrors.NewNotFoundMongoError("Carnet")
	}
	return result.Carnet, nil
}

func GetStudentIdByCarnet(carnet string) (string, error) {
	var result struct {
		Id string `bson:"_id"`
	}

	filter := bson.D{{"carnet", carnet}}
	projection := bson.D{{"_id", 1}}
	op := options.FindOne().SetProjection(projection)
	err := dbContext.Student.FindOne(context.TODO(), filter, op).Decode(&result)
	if err != nil {
		return "", customErrors.NewNotFoundMongoError("Carnet")

	}
	return result.Id, nil
}

func GetStudentCarnetById(id string) (string, *customErrors.NotFoundMongoError) {
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
		return "", customErrors.NewNotFoundMongoError("Carnet")
	}
	return result.Carnet, nil
}
