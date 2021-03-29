package data

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Here is the definition of the database structure.
type MongoDB struct {
	db     *mongo.Database
	client *mongo.Client
}

func NewMongoDB() *MongoDB {
	return &MongoDB{
		db:     nil,
		client: nil,
	}
}

// InitDatabase function initialize the connection to the database.
func (m *MongoDB) InitDatabase() error {
	opt := options.Client().ApplyURI("mongodb://49.234.35.49:27017")
	// Change the port and connection method for connecting to the database according to the situation.
	client, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	db := client.Database("labnote")
	// Init these variable of mongodb.
	m.db = db
	m.client = client
	return nil
}

// CloseDatabase is a function close the connection with mongodb.
func (m *MongoDB) CloseDatabase() error {
	err := m.client.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

// CheckUserAuth is a check of if user permissions are correct.
func (m *MongoDB) CheckUserAuth(user *User) (bool, error) {
	var res User
	err := m.db.Collection("user").FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&res)
	if err != nil {
		// Unable to find the result, the login information is wrong.
		return false, err
	}
	if res.Password == user.Password {
		return true, nil
	}
	return false, nil
}

// InsertOneNote insert a new note into the database's note collection.
func (m *MongoDB) InsertOneNote(note *Note) error {
	_, err := m.db.Collection("note").InsertOne(context.TODO(), note)
	if err != nil {
		return err
	}
	return nil
}

// GetAllNotes function return all the notes in the database.
func (m *MongoDB) GetAllNotes() (*[]Note, error) {
	var all = make([]Note, 0)
	cur, err := m.db.Collection("note").Find(context.TODO(), bson.D{})
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
