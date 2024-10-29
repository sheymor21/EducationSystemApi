package services

import (
	"SchoolManagerApi/internal/dto"
	"SchoolManagerApi/internal/models"
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
// @Failure 500 {object} string "Internal Server Error"
// @Router /login [post]
// @Tags login
func Login(w http.ResponseWriter, r *http.Request) {

	var loginRequest dto.UserLoginRequest
	jsonErr := utilities.ReadJson(w, r, &loginRequest)
	if jsonErr != nil {
		httpInternalError(w, jsonErr.Error())
		utilities.Log.Errorln(jsonErr)
		return
	}

	var userDb models.User
	dbErr := dbContext.Users.FindOne(context.TODO(), bson.D{{"carnet", loginRequest.Carnet}}).Decode(&userDb)
	hashComparison := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(loginRequest.Password))
	if hashComparison != nil {
		httpNotFoundError(w, "Incorrect Username or Password")
		return
	}

	if dbErr != nil {
		if errors.Is(dbErr, mongo.ErrNoDocuments) {
			httpNotFoundError(w, "Incorrect Username or Password")
			return
		}
		httpInternalError(w, dbErr.Error())
		return
	}
	userRol := validations.Rol(userDb.Rol)
	jwtUser := &validations.JWTUser{
		Carnet: userDb.Carnet,
		Rol:    userRol,
	}
	jwt, jwtErr := validations.CreateJWT(*jwtUser)
	if jwtErr != nil {
		httpInternalError(w, jwtErr.Error())
		return
	}
	utilities.WriteJson(w, http.StatusOK, jwt)

}
