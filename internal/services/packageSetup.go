package services

import (
	"calificationApi/internal/database"
	"net/http"
)

var dbContext = database.GetMongoContext()

func HttpStudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getStudents(w)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func HttpTeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeachers(w)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

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

func HttpMarksHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getMarksByStudentCarnet(w, r)
	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}
