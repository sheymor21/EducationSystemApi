package services

import (
	"calificationApi/internal/Dto"
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

func HttpMarkHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMark(w, r)
	case http.MethodPost:
		addMark(w, r)
	case http.MethodPut:
		updateMark(w, r)
	case http.MethodDelete:
		deleteMark(w, r)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func HttpMarksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMarksByStudentCarnet(w, r)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func addMark(w http.ResponseWriter, r *http.Request) {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	var input Dto.MarkAddRequest
	err := Utilities.ReadJson(w, r, &input)
	if err != nil {
		httpInternalError(w, err.Error())
		log.Println(err)
		return
	}
	teacherId, teacherErr := getTeacherIdByCarnet(dbContext, input.TeacherCarnet)
	studentId, studentErr := getStudentIdByCarnet(dbContext, input.StudentCarnet)
	if teacherErr != nil || studentErr != nil {
		if teacherErr == nil {
			httpInternalError(w, studentErr.Error())
		} else {
			httpInternalError(w, teacherErr.Error())
		}
	} else {

		var mark models.Mark
		{
			//mark.ID = "8"
			mark.TeacherId = teacherId
			mark.StudentId = studentId
			mark.Grade = input.Grade
			mark.Semester = input.Semester
			mark.Mark = input.Mark

		}
		_, err := dbContext.Marks.InsertOne(context.TODO(), mark)
		if err != nil {
			httpInternalError(w, err.Error())
		}
	}

}

func getMarksByStudentCarnet(w http.ResponseWriter, r *http.Request) {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	carnet := r.URL.Query().Get("Carnet")
	studentId, err := getStudentIdByCarnet(dbContext, carnet)
	if err != nil {
		httpNotFoundError(w, err.Error())
	} else {

		var marks []Dto.MarksGetRequest
		filter := bson.D{{"student_id", studentId}}
		cursor, err := dbContext.Marks.Find(context.TODO(), filter)
		if err != nil {
			httpNotFoundError(w, NewNotFoundMongoError("carnet").msg)
		} else {

			for cursor.Next(context.TODO()) {
				var dbMark models.Mark
				err := cursor.Decode(&dbMark)
				if err != nil {
					httpInternalError(w, err.Error())
				} else {

					teacherCarnet, err := getTeacherCarnetById(dbMark.TeacherId)
					if err != nil {
						httpNotFoundError(w, NewNotFoundMongoError("teacherCarnet").msg)
					} else {

						var mark Dto.MarksGetRequest
						{
							mark.ID = dbMark.ID
							mark.Semester = dbMark.Semester
							mark.Grade = dbMark.Grade
							mark.Mark = dbMark.Mark
							mark.TeacherCarnet = teacherCarnet
							mark.StudentCarnet = carnet
						}
						marks = append(marks, mark)
					}
				}
			}
			Utilities.WriteJson(w, http.StatusOK, marks)
		}

	}

}

func getMark(w http.ResponseWriter, r *http.Request) {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	id := r.URL.Query().Get("id")
	anyMark := anyMark(id)
	if !anyMark {
		httpNotFoundError(w, NewNotFoundMongoError("id").msg)
	} else {
		filter := bson.M{"id": id}
		var mark models.Mark
		err := dbContext.Marks.FindOne(context.TODO(), filter).Decode(&mark)
		if err != nil {
			httpNotFoundError(w, NewNotFoundMongoError("id").msg)
		} else {
			Utilities.WriteJson(w, http.StatusOK, mark)
		}
	}
}

func deleteMark(w http.ResponseWriter, r *http.Request) {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	id := r.URL.Query().Get("id")
	anyMark := anyMark(id)
	if !anyMark {
		httpNotFoundError(w, NewNotFoundMongoError("id").msg)
	} else {
		filter := bson.M{"id": id}
		_, err := dbContext.Marks.DeleteOne(context.TODO(), filter)
		if err != nil {
			httpInternalError(w, err.Error())
		}
	}

}

func updateMark(w http.ResponseWriter, r *http.Request) {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	carnet := r.URL.Query().Get("Carnet")
	anyMark := anyMark(carnet)
	if !anyMark {
		httpNotFoundError(w, NewNotFoundMongoError("carnet").msg)
	} else {
		var mark models.Mark
		err := Utilities.ReadJson(w, r, &mark)
		if err != nil {
			httpInternalError(w, err.Error())
			log.Println(err)
			return
		}
		filter := bson.M{"carnet": carnet}
		update := bson.M{"$set": mark}
		_, err = dbContext.Marks.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			httpInternalError(w, err.Error())
		}
	}
}

func anyMark(carnet string) bool {
	dbContext, client := database.GetDatabaseConnection()
	defer database.CloseConnection(client, context.TODO())

	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Student.FindOne(context.TODO(), filter).Decode(&models.Mark{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false
	} else if err != nil {
		log.Fatal(err)
		return false
	} else {

		return true
	}
}
