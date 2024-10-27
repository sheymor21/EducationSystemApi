package utilities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BsonIdFormat(id string) primitive.ObjectID {
	hex, hexErr := primitive.ObjectIDFromHex(id)
	if hexErr != nil {
		Log.Errorln(hexErr)
	}
	return hex
}
