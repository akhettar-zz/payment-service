package repository

import (
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/payment-service/model"
)

// MongoRepository type
type MongoRepository struct {
	Session *mgo.Session
}

// Repository interface
type Repository interface {

	// Insert content in the given db and collection
	Insert(db, col string, content interface{}) error

	// Find all the notes
	FindAll(db, col string) (model.PaymentResponse, error)

	// Find a payment for a given ID
	Find(db, col string, oid bson.ObjectId) (model.PaymentResponse, error)

	// Delete a payment for a given ID
	Delete(db, col string, oid bson.ObjectId) error

	// Update a payment for given ID
	Update(db, col string, oid bson.ObjectId, content interface{}) error
}

// Insert content into db
func (repo *MongoRepository) Insert(db string, col string, content interface{}) error {
	return repo.Session.DB(db).C(col).Insert(&content)
}

// Find query tag for a given id
func (repo *MongoRepository) Find(db string, collection string, oid bson.ObjectId) (model.PaymentResponse, error) {
	var result model.Payment
	err := repo.Session.DB(db).C(collection).FindId(oid).One(&result)
	return model.PaymentResponse{Data: []model.Payment{result}, Links: model.Links{Self: "https://api.test.form3.tech/v1/payments"}}, err
}

// FindAll query all the
func (repo *MongoRepository) FindAll(db string, col string) (model.PaymentResponse, error) {
	var result []model.Payment
	err := repo.Session.DB(db).C(col).Find(nil).All(&result)
	return model.PaymentResponse{Data: result, Links: model.Links{Self: "https://api.test.form3.tech/v1/payments"}}, err
}

// Delete payment
func (repo *MongoRepository) Delete(db, col string, oid bson.ObjectId) error {
	return repo.Session.DB(db).C(col).RemoveId(oid)
}

// Update Given Payment
func (repo *MongoRepository) Update(db string, collection string, oid bson.ObjectId, content interface{}) error {
	return repo.Session.DB(db).C(collection).UpdateId(oid, content)
}

// NewRepository creates a Repository type
func NewRepository(uri string) Repository {
	dialInfo, err := mgo.ParseURL(uri)
	if err != nil {
		log.Panicf("Failed to parse Mongo URI")
	}

	// get session
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		panic(err)
	}
	repository := &MongoRepository{session}
	return repository

}
