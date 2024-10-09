package services

import (
	"calificationApi/internal/dto"
	"calificationApi/internal/models"
	"calificationApi/internal/server/customErrors"
	"calificationApi/internal/services/search"
	"calificationApi/internal/utilities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"sync"
)

func addMark(w http.ResponseWriter, r *http.Request) {

	var input dto.MarkAddRequest
	err := utilities.ReadJson(w, r, &input)
	if err != nil {
		httpInternalError(w, err.Error())
		log.Println(err)
		return
	}
	teacherId, teacherErr := search.GetTeacherIdByCarnet(input.TeacherCarnet)
	studentId, studentErr := search.GetStudentIdByCarnet(dbContext, input.StudentCarnet)
	if teacherErr != nil || studentErr != nil {
		if teacherErr == nil {
			httpInternalError(w, studentErr.Error())
		} else {
			httpInternalError(w, teacherErr.Error())
		}
	} else {

		var mark models.Mark
		{
			mark.TeacherId = teacherId
			mark.StudentId = studentId
			mark.Grade = input.Grade
			mark.Semester = input.Semester
			mark.Mark = input.Mark

		}
		_, err = dbContext.Marks.InsertOne(context.TODO(), mark)
		if err != nil {
			httpInternalError(w, err.Error())
		}
	}

}

func getMarksByStudentCarnet(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	studentIdCh := make(chan string)
	marksCh := make(chan []dto.MarksGetRequest)
	studentCarnet := r.URL.Query().Get("Carnet")
	wg.Add(2)
	go func(carnet string) {
		defer wg.Done()
		defer close(studentIdCh)
		studentId, err := search.GetStudentIdByCarnet(dbContext, studentCarnet)
		if err != nil {
			httpNotFoundError(w, customErrors.NewNotFoundMongoError("studentCarnet").Msg)
			return
		}
		studentIdCh <- studentId

	}(studentCarnet)

	go func(carnet string) {
		defer wg.Done()
		defer close(marksCh)
		select {
		case studentId := <-studentIdCh:
			var marks []dto.MarksGetRequest
			filter := bson.D{{"student_id", studentId}}
			cursor, markFindErr := dbContext.Marks.Find(context.TODO(), filter)
			if markFindErr != nil {
				httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
				return
			}
			for cursor.Next(context.TODO()) {
				var dbMark models.Mark
				decodeErr := cursor.Decode(&dbMark)
				if decodeErr != nil {
					httpInternalError(w, decodeErr.Error())
					return
				}
				teacherCarnet, getTeacherErr := search.GetTeacherCarnetById(dbMark.TeacherId)
				if getTeacherErr != nil {
					httpNotFoundError(w, teacherCarnet)
					return
				}
				mark := dbMark.ToGetRequest(carnet, teacherCarnet)
				marks = append(marks, mark)
			}
			cursorCloseErr := cursor.Close(context.TODO())
			if cursorCloseErr != nil {
				httpInternalError(w, cursorCloseErr.Error())
				return
			}

			marksCh <- marks
		}
	}(studentCarnet)

	marks := <-marksCh
	wg.Wait()
	utilities.WriteJson(w, http.StatusOK, marks)
}

func getMark(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	anyMark := anyMark(id)
	if !anyMark {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("id").Msg)
	} else {
		filter := bson.M{"id": id}
		var mark models.Mark
		err := dbContext.Marks.FindOne(context.TODO(), filter).Decode(&mark)
		if err != nil {
			httpNotFoundError(w, customErrors.NewNotFoundMongoError("id").Msg)
		} else {
			utilities.WriteJson(w, http.StatusOK, mark)
		}
	}
}

func deleteMark(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	anyMark := anyMark(id)
	if !anyMark {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("id").Msg)
	} else {
		filter := bson.M{"id": id}
		_, err := dbContext.Marks.DeleteOne(context.TODO(), filter)
		if err != nil {
			httpInternalError(w, err.Error())
		}
	}

}

func updateMark(w http.ResponseWriter, r *http.Request) {

	carnet := r.URL.Query().Get("Carnet")
	anyMark := anyMark(carnet)
	if !anyMark {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
	} else {
		var mark models.Mark
		err := utilities.ReadJson(w, r, &mark)
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
