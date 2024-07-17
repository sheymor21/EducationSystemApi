package models

type Teacher struct {
	Carnet    string `bson:"carnet"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Age       uint8  `bson:"age"`
	Classroom string `bson:"classroom"`
}
