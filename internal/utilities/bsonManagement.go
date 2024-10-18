package utilities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func BsonIdFormat(id string) primitive.ObjectID {
	hex, hexErr := primitive.ObjectIDFromHex(id)
	if hexErr != nil {
		log.Println(hexErr)
	}
	return hex
}
