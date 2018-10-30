package api_test

import (
	"encoding/json"
	"github.com/payment-service/api"
	"github.com/payment-service/logger"
	"github.com/payment-service/model"
	"github.com/payment-service/test"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	urlFx                 = "url1" // dummy url for FX service
	urlCh                 = "url2" // dummy url for Charges service
	beneficiaryCurrency   = "USD"
	debtorAccountNumb     = "GB29XABC10161234567801"
	beneficiaryAccountNum = "31926819"
)

func TestCarPaymentHandler_Health(t *testing.T) {
	t.Logf("Given the need to use the health endpoint to query container status")
	{
		t.Logf("\tWhen checking \"%s\" for status code \"%d\"", "\\health", http.StatusOK)
		{
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)
			router := handler.NewRouter()
			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusOK)
		}
	}
}

func TestPaymentHandler_CreatePayment(t *testing.T) {
	t.Logf("Given the need to create a payment")
	{
		t.Logf("\tWhen sending Create Payment request to endpoint %s", "\\payment")
		{
			body := test.CreatePaymentRequest(beneficiaryAccountNum, debtorAccountNumb, beneficiaryCurrency)

			bytes, _ := json.Marshal(body)
			logger.Info.Println(string(bytes))
			w := httptest.NewRecorder()

			req, err := test.HttpRequest(body, "/payment", http.MethodPost)
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)
			router := handler.NewRouter()
			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusCreated)

			var response model.CreatePaymentResponse
			json.NewDecoder(w.Body).Decode(&response)
			if response.Id != "" {
				t.Logf("\t\tThe response should contain payment id. %v %v", response.Id, test.CheckMark)
			} else {
				t.Errorf("\t\tThe response should contain payment id. %v %v", response.Id, test.BallotX)
			}
		}
	}
}

func TestPaymentHandler_CreatePaymentForeignExchangeNotRequired(t *testing.T) {
	t.Logf("Given the need to create a payment")
	{
		t.Logf("\tWhen sending Create Payment request to endpoint %s", "\\payment")
		{
			body := test.CreatePaymentRequest(beneficiaryAccountNum, debtorAccountNumb, "GBP")

			bytes, _ := json.Marshal(body)
			logger.Info.Println(string(bytes))
			w := httptest.NewRecorder()

			req, err := test.HttpRequest(body, "/payment", http.MethodPost)
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)
			router := handler.NewRouter()
			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusCreated)

			var response model.CreatePaymentResponse
			json.NewDecoder(w.Body).Decode(&response)
			if response.Id != "" {
				t.Logf("\t\tThe response should contain payment id. %v %v", response.Id, test.CheckMark)
			} else {
				t.Errorf("\t\tThe response should contain payment id. %v %v", response.Id, test.BallotX)
			}
		}
	}
}

func TestPaymentHandler_CreatePaymentWithEmptyBodyShouldReturnBadRequest(t *testing.T) {
	t.Logf("Given the need to create a payment")
	{
		t.Logf("\tWhen sending an invalid Create Request to endpoint %s", "\\payment")
		{
			body := model.CreatePaymentRequest{}

			bytes, _ := json.Marshal(body)
			logger.Info.Println(string(bytes))
			w := httptest.NewRecorder()

			req, err := test.HttpRequest(body, "/payment", http.MethodPost)
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)
			router := handler.NewRouter()
			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusBadRequest)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)
			expectedMessage := "Failed to parse payment request"
			if response.Message == expectedMessage {
				t.Logf("\t\tThe response should be: %v %v", expectedMessage, test.CheckMark)
			} else {
				t.Errorf("\t\tThe response should be: %v %v %v", expectedMessage, test.BallotX, response.Message)
			}
		}
	}
}

func TestPaymentHandler_QueryAllShouldReturnAllPayment(t *testing.T) {
	t.Logf("Given the need to Query all payment")
	{
		t.Logf("\tWhen sending Query All Payment request to endpoint %s", "\\payment")
		{
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)

			// create first payment
			test.CreatePaymentAndAssertResponse(t, handler)

			// create second payment
			test.CreatePaymentAndAssertResponse(t, handler)

			req, err := http.NewRequest(http.MethodGet, "/payment", nil)
			router := handler.NewRouter()
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusOK)

			var response model.PaymentResponse
			json.NewDecoder(w.Body).Decode(&response)
			if len(response.Data) >= 2 {
				t.Logf("\t\tThe response should contain %v payments %v", 2, test.CheckMark)
			} else {
				t.Errorf("\t\tThe response should contain %v payment %v %v", 2, test.BallotX, len(response.Data))
			}
		}
	}
}

func TestPaymentHandler_QueryForAGivenPaymentIdShouldReturn200(t *testing.T) {
	t.Logf("Given the need to Query a payment for given payment Id")
	{
		t.Logf("\tWhen sending Query Payment request to endpoint %s", "\\payment")
		{
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)

			// create first payment
			res := test.CreatePaymentAndAssertResponse(t, handler)

			req, err := test.HttpRequest(nil, "/payment/"+res.Id, http.MethodGet)
			router := handler.NewRouter()
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusOK)

			var response model.PaymentResponse
			json.NewDecoder(w.Body).Decode(&response)
			if response.Data[0].Id.Hex() == res.Id {
				t.Logf("\t\tThe payment Id should be %v %v", res.Id, test.CheckMark)
			} else {
				t.Errorf("\t\tThe payment Id should be %v %v %v", res.Id, test.BallotX, response.Data[0].Id)
			}
		}
	}
}

func TestPaymentHandler_QueryForAGivenDummyID_ShouldReturnNotFound(t *testing.T) {
	t.Logf("Given the need to Query a payment for given payment Id")
	{
		t.Logf("\tWhen sending Query Payment request to endpoint %s", "\\payment")
		{
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)

			req, err := test.HttpRequest(nil, "/payment/"+bson.NewObjectId().Hex(), http.MethodGet)
			router := handler.NewRouter()
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusNotFound)
		}
	}
}

func TestPaymentHandler_DeleteShouldDeleteTheGivenPayment(t *testing.T) {
	t.Logf("Given the need to Query a payment for given payment Id")
	{
		t.Logf("\tWhen sending Query Payment request to endpoint %s", "\\payment")
		{
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)

			// create first payment
			res := test.CreatePaymentAndAssertResponse(t, handler)

			// Assert the record exist
			req, err := test.HttpRequest(nil, "/payment/"+res.Id, http.MethodGet)
			router := handler.NewRouter()
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusOK)

			// delete the resource
			w = httptest.NewRecorder()
			req, err = test.HttpRequest(nil, "/payment/"+res.Id, http.MethodDelete)
			router.ServeHTTP(w, req)

			// assert successful delete
			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusNoContent)
		}
	}
}

func TestPaymentHandler_ForPaymentNotFoundDeleteShouldReturn404(t *testing.T) {
	t.Logf("Given the need to Query a payment for given payment Id")
	{
		t.Logf("\tWhen sending Query Payment request to endpoint %s", "\\payment")
		{
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)

			// create first payment
			res := test.CreatePaymentAndAssertResponse(t, handler)

			// Assert the record exist
			req, err := test.HttpRequest(nil, "/payment/"+res.Id, http.MethodGet)
			router := handler.NewRouter()
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusOK)

			// delete the resource
			w = httptest.NewRecorder()
			req, err = test.HttpRequest(nil, "/payment/"+bson.NewObjectId().Hex(), http.MethodDelete)
			router.ServeHTTP(w, req)

			// assert successful delete
			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusNotFound)
		}
	}
}

func TestPaymentHandler_SuccessfulUpdatePaymentShouldReturn204(t *testing.T) {
	t.Logf("Given the need to update a payment")
	{
		t.Logf("\tWhen sending Update Payment request to endpoint %s", "\\payment")
		{
			// Create payment
			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)
			res := test.CreatePaymentAndAssertResponse(t, handler)

			// update payment
			newDebtorAccNum := "GB29XABC101613434343"
			update := test.CreatePaymentRequest(beneficiaryAccountNum, newDebtorAccNum, beneficiaryCurrency)
			req, err := test.HttpRequest(update, "/payment/"+res.Id, http.MethodPut)

			w := httptest.NewRecorder()
			router := handler.NewRouter()
			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusNoContent)

			req, err = test.HttpRequest(nil, "/payment/"+res.Id, http.MethodGet)
			w = httptest.NewRecorder()

			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusOK)

			var response model.PaymentResponse
			json.NewDecoder(w.Body).Decode(&response)
			if response.Data[0].DebtorParty.AccountNumber == newDebtorAccNum {
				t.Logf("\t\tThe debtor account number should be %v %v", newDebtorAccNum, test.CheckMark)
			} else {
				t.Errorf("\t\tThe debtor account number should be %v %v %v", newDebtorAccNum, test.BallotX,
					response.Data[0].DebtorParty.AccountNumber)
			}
		}
	}
}

func TestPaymentHandler_AttemptUpdatePaymentNotFoundShouldReturn404(t *testing.T) {
	t.Logf("Given the need to update a payment")
	{
		t.Logf("\tWhen sending Update Payment request to endpoint %s", "\\payment")
		{

			handler := api.NewPaymentHandler(Repository, urlFx, urlCh)

			// update payment
			newDebtorAccNum := "GB29XABC101613434343"
			update := test.CreatePaymentRequest(beneficiaryAccountNum, newDebtorAccNum, beneficiaryCurrency)
			dummyId := bson.NewObjectId()
			req, err := test.HttpRequest(update, "/payment/"+dummyId.Hex(), http.MethodPut)

			w := httptest.NewRecorder()
			router := handler.NewRouter()
			router.ServeHTTP(w, req)

			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusNotFound)
		}
	}
}
