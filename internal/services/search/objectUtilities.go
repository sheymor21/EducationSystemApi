package search

import (
	"SchoolManagerApi/internal/database"
	"SchoolManagerApi/internal/server/customErrors"
	"SchoolManagerApi/internal/utilities"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		utilities.Log.Errorln(err)
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
		utilities.Log.Errorln(err)
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
