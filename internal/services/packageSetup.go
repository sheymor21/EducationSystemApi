package services

import (
	"calificationApi/internal/database"
	"net/http"
	"sync"
)

var dbContext = database.GetMongoContext()

func HttpStudentHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()
	wg.Wait()

}

func HttpMarkHandler(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		switch r.Method {
		case http.MethodGet:
			go getMark(w, r)
		case http.MethodPost:
			addMark(w, r)
		case http.MethodPut:
			updateMark(w, r)
		case http.MethodDelete:
			deleteMark(w, r)
		default:
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	}()
	wg.Wait()

}

func HttpTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

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
	}()
	wg.Wait()
}

func HttpMarksHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		switch r.Method {
		case http.MethodGet:
			getMarksByStudentCarnet(w, r)
		default:
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	}()
	wg.Wait()
}
