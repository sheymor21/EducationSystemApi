package services

import (
	"SchoolManagerApi/internal/dto"
	"SchoolManagerApi/internal/models"
	"SchoolManagerApi/internal/server/customErrors"
	"SchoolManagerApi/internal/utilities"
	"SchoolManagerApi/internal/validations"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Login handles the creation of a new mark entry in the database.
// @Summary Get JWT by Login
// @Description Get JWT by Login
// @Accept json
// @Produce json
// @Param user body UserLoginRequest true "Login User"
// @Success 200 {object} string "JWT"
// @Failure 500 string error
// @Router /login [post]
// @Tags login
func Login(w http.ResponseWriter, r *http.Request) {
	defer utilities.Recover()
	var loginRequest dto.UserLoginRequest

	jsonErr := utilities.ReadJson(w, r, &loginRequest)
	customErrors.ThrowHttpError(jsonErr, w, "", http.StatusInternalServerError)

	var userDb models.User
	dbErr := dbContext.Users.FindOne(context.TODO(), bson.D{{"carnet", loginRequest.Carnet}}).Decode(&userDb)
	hashComparisonErr := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(loginRequest.Password))
	customErrors.ThrowHttpError(hashComparisonErr, w, "Incorrect Username or Password", http.StatusNotFound)

	if dbErr != nil {
		if errors.Is(dbErr, mongo.ErrNoDocuments) {
			http.Error(w, "Incorrect Username or Password", http.StatusNotFound)
			return
		}
		http.Error(w, dbErr.Error(), http.StatusNotFound)
		return
	}
	userRol := validations.Rol(userDb.Rol)
	jwtUser := validations.JWTUser{
		Carnet: userDb.Carnet,
		Rol:    userRol,
	}
	jwt, jwtErr := validations.CreateJWT(jwtUser)
	customErrors.ThrowHttpError(jwtErr, w, "", http.StatusInternalServerError)
	utilities.WriteJson(w, http.StatusOK, jwt)

}
