package services

import (
	"SchoolManagerApi/internal/models"
	"SchoolManagerApi/internal/utilities"
	"SchoolManagerApi/internal/validations"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func addUser(firstName, lastname, carnet string, rol validations.Rol) error {
	userName := fmt.Sprintf("%s %s", firstName, lastname)
	password := fmt.Sprintf("%s-%s", carnet[len(carnet)-3:], strings.ToLower(lastname))
	hashPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), 4)
	if hashErr != nil {
		utilities.Log.Errorln(hashErr)
		return hashErr
	}
	user := models.User{
		Carnet:   carnet,
		Username: userName,
		Password: string(hashPassword),
		Rol:      rol.EnumIndex(),
	}
	_, insertUserErr := dbContext.Users.InsertOne(context.TODO(), user)
	if insertUserErr != nil {
		utilities.Log.Errorln(insertUserErr)
		_, deleteErr := dbContext.Student.DeleteOne(context.TODO(), bson.D{{"carnet", carnet}})
		if deleteErr != nil {
			utilities.Log.Errorln(deleteErr)
			return deleteErr
		}
	}
	return nil
}
