package repository_test

import (
	"github.com/globalsign/mgo/bson"
	"payment-service/model"
	"payment-service/repository"
	"payment-service/test"
	"testing"
)

func TestMongoRepository_InsertShoulBeSuccessful(t *testing.T) {

	t.Logf("Given the DB is up and running")
	{
		t.Logf("\tWhen Inserting objct into DB")
		{
			payment := model.Payment{Type: "Payment", ID: bson.NewObjectId()}
			err := repository.RepositoryUnderTest.Insert("paymentDb", "payments", payment)
			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}
		}
	}
}

func TestMongoRepository_DeleteShoulBeSuccessful(t *testing.T) {

	t.Logf("Given the DB is up and running")
	{
		t.Logf("\tWhen Inserting objct into DB")
		{
			// Insert payment
			obi := bson.NewObjectId()
			payment := model.Payment{Type: "Payment", ID: obi}
			err := repository.RepositoryUnderTest.Insert("paymentDb", "payments", payment)
			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}

			// delete
			err = repository.RepositoryUnderTest.Delete("paymentDb", "payments", obi)
			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}
		}
	}
}

func TestMongoRepository_FindShouldBeSuccessful(t *testing.T) {

	t.Logf("Given the DB is up and running")
	{
		t.Logf("\tWhen Inserting objct into DB")
		{
			// Insert payment
			obi := bson.NewObjectId()
			payment := model.Payment{Type: "Payment", ID: obi}
			err := repository.RepositoryUnderTest.Insert("paymentDb", "payments", payment)
			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}

			// find
			res, err := repository.RepositoryUnderTest.Find("paymentDb", "payments", obi)

			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}

			if res.Data[0].ID == obi {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}
		}
	}
}

func TestMongoRepository_FindAllBeSuccessful(t *testing.T) {

	t.Logf("Given the DB is up and running")
	{
		t.Logf("\tWhen Inserting objct into DB")
		{
			// Insert payment
			obi := bson.NewObjectId()
			payment := model.Payment{Type: "Payment", ID: obi}
			err := repository.RepositoryUnderTest.Insert("paymentDb", "payments", payment)
			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}

			// find all
			_, err = repository.RepositoryUnderTest.FindAll("paymentDb", "payments")

			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}
		}
	}
}

func TestMongoRepository_UpdateBeSuccessful(t *testing.T) {

	t.Logf("Given the DB is up and running")
	{
		t.Logf("\tWhen Inserting objct into DB")
		{
			// Insert payment
			obi := bson.NewObjectId()
			payment := model.Payment{Type: "Payment", ID: obi, OrganisationId: "org1"}
			err := repository.RepositoryUnderTest.Insert("paymentDb", "payments", payment)
			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}

			// update
			updated := model.Payment{Type: "Payment", ID: obi, OrganisationId: "org2"}
			err = repository.RepositoryUnderTest.Update("paymentDb", "payments", obi, updated)

			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}

			// find
			res, err := repository.RepositoryUnderTest.Find("paymentDb", "payments", obi)

			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}

			if res.Data[0].OrganisationId == "org2" {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}
		}
	}
}
