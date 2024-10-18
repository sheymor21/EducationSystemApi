package services

import (
	"calificationApi/internal/dto"
	"calificationApi/internal/mappers"
	"calificationApi/internal/models"
	"calificationApi/internal/server/customErrors"
	"calificationApi/internal/services/search"
	"calificationApi/internal/utilities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	studentId, studentErr := search.GetStudentIdByCarnet(input.StudentCarnet)
	if teacherErr != nil {
		httpInternalError(w, teacherErr.Error())
		return
	}
	if studentErr != nil {
		httpInternalError(w, studentErr.Error())
		return
	}

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
		return
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
		studentId, err := search.GetStudentIdByCarnet(studentCarnet)
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
			var marks []models.Mark
			filter := bson.D{{"student_id", studentId}}
			cursor, markFindErr := dbContext.Marks.Find(context.TODO(), filter)
			if markFindErr != nil {
				httpNotFoundError(w, customErrors.NewNotFoundMongoError("carnet").Msg)
				return
			}
			err := cursor.All(context.TODO(), &marks)
			if err != nil {
				httpInternalError(w, err.Error())
				return
			}
			cursorErr := cursor.Close(context.TODO())
			if cursorErr != nil {
				httpInternalError(w, cursorErr.Error())
				return
			}
			markDto := mappers.MarkListToGetDto(marks)
			marksCh <- markDto
		}
	}(studentCarnet)

	marks := <-marksCh
	wg.Wait()
	utilities.WriteJson(w, http.StatusOK, marks)
}

func getMark(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	ch := make(chan bool)
	id := r.URL.Query().Get("id")
	wg.Add(1)
	go anyMarkAtStudents(id, &wg, ch)
	if <-ch {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("id").Msg)
		return
	}
	close(ch)
	wg.Wait()
	bsonId := utilities.BsonIdFormat(id)
	filter := bson.M{"_id": bsonId}
	var mark models.Mark
	err := dbContext.Marks.FindOne(context.TODO(), filter).Decode(&mark)
	if err != nil {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("id").Msg)
		return
	}
	studentCarnet, studentErr := search.GetStudentCarnetById(mark.StudentId)
	if studentErr != nil {
	markDto, mapperErr := mappers.MarkToGetDto(mark)
	if mapperErr != nil {
		httpNotFoundError(w, mapperErr.Error())
		return
	}
	utilities.WriteJson(w, http.StatusOK, markDto)
}

func deleteMark(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	ch := make(chan bool)
	id := r.URL.Query().Get("id")
	wg.Add(1)
	go anyMarkAtStudents(id, &wg, ch)
	if <-ch {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("id").Msg)
		return
	}
	wg.Wait()
	filter := bson.M{"id": id}
	_, err := dbContext.Marks.DeleteOne(context.TODO(), filter)
	if err != nil {
		httpInternalError(w, err.Error())
		return
	}

}

func updateMark(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	ch := make(chan bool)
	id := r.URL.Query().Get("id")
	wg.Add(1)
	go anyMarks(id, &wg, ch)
	if <-ch {
		httpNotFoundError(w, customErrors.NewNotFoundMongoError("id").Msg)
		return
	}
	close(ch)
	wg.Wait()
	var markDto dto.MarksUpdateRequest
	err := utilities.ReadJson(w, r, &markDto)
	if err != nil {
		httpInternalError(w, err.Error())
		log.Println(err)
		return
	}
	mark, mapperErr := mappers.UpdateDtoToMark(markDto, id)
	if mapperErr != nil {
		httpNotFoundError(w, mapperErr.Error())
		return
	}
	bsonId := utilities.BsonIdFormat(id)
	filter := bson.M{"_id": bsonId}
	update := bson.M{"$set": mark}
	_, err = dbContext.Marks.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		httpInternalError(w, err.Error())
	}

}

func anyMarks(id string, wg *sync.WaitGroup, ch chan bool) {
	defer wg.Done()
	bsonId := utilities.BsonIdFormat(id)
	filter := bson.D{{"_id", bsonId}}
	err := dbContext.Student.FindOne(context.TODO(), filter).Decode(&models.Mark{})
	if errors.Is(err, mongo.ErrNoDocuments) {
	}
}

func anyMarkAtStudents(carnet string, wg *sync.WaitGroup, ch chan bool) {
	defer wg.Done()
	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Student.FindOne(context.TODO(), filter).Decode(&models.Mark{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		ch <- false
	} else if err != nil {
		log.Println(err)
		ch <- false
	} else {
		ch <- true
	}
}
