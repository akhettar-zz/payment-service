package repository_test

import (
	"github.com/payment-service/repository"
	"github.com/payment-service/test"
	"testing"
)

func TestMongoRepository_InsertShoulBeSuccessful(t *testing.T) {

	t.Logf("Given the DB is up and running")
	{
		t.Logf("\tWhen Inserting objct into DB")
		{
			payment := test.CreatePaymentRequest("31926819", "dfa8888434348", "USD")
			err := repository.RepositoryUnderTest.Insert("paymentDb", "payments", payment)
			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}
		}
	}
}
