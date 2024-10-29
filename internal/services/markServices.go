package services

import (
	"SchoolManagerApi/internal/dto"
	"SchoolManagerApi/internal/mappers"
	"SchoolManagerApi/internal/models"
	"SchoolManagerApi/internal/server/customErrors"
	"SchoolManagerApi/internal/services/search"
	"SchoolManagerApi/internal/utilities"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// @Summary Add a new mark
// @Description Creates a new mark entry with student and teacher details
// @Param request body MarkAddRequest true "Mark Add Request"
// @Success 200
// @Failure 500 string error
// @Router /mark [post]
// @Tags mark
func addMark(w http.ResponseWriter, r *http.Request) {
	defer utilities.Recover()
	var markDto dto.MarkAddRequest
	jsonErr := utilities.ReadJson(w, r, &markDto)
	customErrors.ThrowHttpError(jsonErr, w, "", http.StatusInternalServerError)
	teacherId, teacherErr := search.GetTeacherIdByCarnet(markDto.TeacherCarnet)
	studentId, studentErr := search.GetStudentIdByCarnet(markDto.StudentCarnet)
	customErrors.ThrowHttpError(teacherErr, w, "", http.StatusNotFound)
	customErrors.ThrowHttpError(studentErr, w, "", http.StatusNotFound)

	mark := mappers.MarkAddToModel(markDto, teacherId, studentId)
	_, dbErr := dbContext.Marks.InsertOne(context.TODO(), mark)
	customErrors.ThrowHttpError(dbErr, w, "", http.StatusInternalServerError)

}

// @Summary Retrieves student's marks by carnet
// @Description Finds and returns the marks of a student using their carnet number.
// @Tags marks
// @Param Carnet query string true "Student Carnet"
// @Success 200 {array} MarksGetRequest
// @Failure 404 {string} string "Mark Not Found"
// @Failure 500 string error
// @Router /marks [get]
// @Produce json
func getMarksByStudentCarnet(w http.ResponseWriter, r *http.Request) {
	defer utilities.Recover()
	studentCarnet := r.URL.Query().Get("Carnet")
	studentId, err := search.GetStudentIdByCarnet(studentCarnet)
	customErrors.ThrowHttpError(err, w, "Not Found this studentCarnet", http.StatusNotFound)
	var marks []models.Mark
	filter := bson.D{{"student_id", studentId}}
	cursor, markFindErr := dbContext.Marks.Find(context.TODO(), filter)
	customErrors.ThrowHttpError(markFindErr, w, "Not Found this carnet", http.StatusNotFound)
	cursorErr := cursor.All(context.TODO(), &marks)
	customErrors.ThrowHttpError(cursorErr, w, "", http.StatusNotFound)
	cursorCloseErr := cursor.Close(context.TODO())
	customErrors.ThrowHttpError(cursorCloseErr, w, "", http.StatusNotFound)

	markDto := mappers.MarkListToGetRequest(marks)
	utilities.WriteJson(w, http.StatusOK, markDto)
}

// @Summary Retrieve a mark
// @Description Fetches a mark object based on the provided student ID
// @Param id query string true "Student ID"
// @Success 200 {object} MarksGetRequest
// @Failure 500 string error
// @Router /mark [get]
// @Tags marks
func getMark(w http.ResponseWriter, r *http.Request) {
	defer utilities.Recover()
	id := r.URL.Query().Get("id")
	markExist := anyMarkAtStudents(id)
	if markExist {
		http.Error(w, "Not found this id", http.StatusNotFound)
		return
	}
	bsonId := utilities.BsonIdFormat(id)
	filter := bson.M{"_id": bsonId}
	var mark models.Mark
	err := dbContext.Marks.FindOne(context.TODO(), filter).Decode(&mark)
	customErrors.ThrowHttpError(err, w, "Not found this id", http.StatusNotFound)
	markDto, mapperErr := mappers.MarkToGetRequest(mark)
	customErrors.ThrowHttpError(mapperErr, w, "", http.StatusNotFound)
	utilities.WriteJson(w, http.StatusOK, markDto)
}

// @Summary Delete a mark
// @Description Deletes a mark entry from the database using the provided mark ID
// @Param id query string true "Mark ID"
// @Success 204
// @Failure 404 {object} string "Mark Not Found"
// @Failure 500 string error
// @Router /mark [delete]
// @Tags mark
func deleteMark(w http.ResponseWriter, r *http.Request) {
	defer utilities.Recover()
	id := r.URL.Query().Get("id")
	markExist := anyMarkAtStudents(id)
	if markExist {
		http.Error(w, "Not found this id", http.StatusNotFound)
		return
	}
	filter := bson.M{"id": id}
	_, err := dbContext.Marks.DeleteOne(context.TODO(), filter)
	customErrors.ThrowHttpError(err, w, "", http.StatusInternalServerError)
	utilities.WriteJson(w, http.StatusNoContent, nil)
}

// @Summary Update a mark
// @Description Modifies an existing mark entry using the supplied ID and mark details
// @Param id query string true "Mark ID"
// @Param request body MarksUpdateRequest true "Marks Update Request"
// @Success 200
// @Failure 404 {object} string "Mark Not Found"
// @Failure 500 string error
// @Router /marks [put]
// @Tags mark
func updateMark(w http.ResponseWriter, r *http.Request) {
	defer utilities.Recover()
	id := r.URL.Query().Get("id")
	markExist := anyMarks(id)
	if markExist {
		http.Error(w, "Not found this id", http.StatusNotFound)
		return
	}
	var markDto dto.MarksUpdateRequest
	jsonErr := utilities.ReadJson(w, r, &markDto)
	customErrors.ThrowHttpError(jsonErr, w, "", http.StatusInternalServerError)
	mark, mapperErr := mappers.MarkUpdateToModel(markDto, id)
	customErrors.ThrowHttpError(mapperErr, w, "", http.StatusNotFound)
	bsonId := utilities.BsonIdFormat(id)
	filter := bson.M{"_id": bsonId}
	update := bson.M{"$set": mark}
	_, dbErr := dbContext.Marks.UpdateOne(context.TODO(), filter, update)
	customErrors.ThrowHttpError(dbErr, w, "", http.StatusInternalServerError)

}

func anyMarks(id string) bool {
	bsonId := utilities.BsonIdFormat(id)
	filter := bson.D{{"_id", bsonId}}
	err := dbContext.Student.FindOne(context.TODO(), filter).Decode(&models.Mark{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false
	} else if err != nil {
		utilities.Log.Errorln(err)
		return false
	} else {
		return true
	}
}

func anyMarkAtStudents(carnet string) bool {
	filter := bson.D{{"carnet", carnet}}
	err := dbContext.Student.FindOne(context.TODO(), filter).Decode(&models.Mark{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false
	} else if err != nil {
		utilities.Log.Errorln(err)
		return false
	} else {
		return true
	}
}
