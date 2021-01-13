package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// There are some common variables used in functions.
var (
	db     *mongo.Database
	client *mongo.Client
)

// InitDatabase function initialize the connection to the database.
func InitDatabase() {
	opt := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	// Change the port and connection method for connecting to the database according to the situation.
	client, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		return
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return
	}
	db = client.Database("labnote")
}

// CloseDatabase is a function close the connection with mongodb.
func CloseDatabase() {
	err := client.Disconnect(context.TODO())
	if err != nil {
		return
	}
}

// CheckUserAuth is a check of if user permissions are correct.
func CheckUserAuth(user *User) bool {
	var res User
	err := db.Collection("user").FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&res)
	if err != nil {
		// Unable to find the result, the login information is wrong.
		return false
	}
	if res.Password == user.Password {
		return true
	}
	return false
}

// InsertOneNote insert a new note into the database's note collection.
func InsertOneNote(note *Note) error {
	_, err := db.Collection("note").InsertOne(context.TODO(), note)
	if err != nil {
		return err
	}
	return nil
}

// GetAllNotes function return all the notes in the database.
func GetAllNotes() (*[]Note, error) {
	var all = make([]Note, 0)
	cur, err := db.Collection("note").Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		// Traverse all notes in the database.
		var elem Note
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		all = append(all, elem)
	}
	return &all, nil
}
