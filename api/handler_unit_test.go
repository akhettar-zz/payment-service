package api_test

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"payment-service/api"
	"payment-service/mocks"
	"payment-service/model"
	"payment-service/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Handle DB insertion failure
func TestCratePayment_DBFailureShouldReturn500(t *testing.T) {
	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Create TagDAO request to endpoint:  \"%s\"", "\\tags")
		{
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mocks.NewMockRepository(mockCtrl)

			expectedErrorMessage := "Failed to create payment"
			body := test.CreatePaymentRequest("31926819", "GB29XABC10161234567801", "USD")
			err := errors.New(expectedErrorMessage)

			// set mock expectation
			mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any(), gomock.Any()).Return(err).Times(1)

			handler := api.NewPaymentHandler(mockRepo, "urlFX", "urlCF")
			router := handler.NewRouter()

			req, err := test.HttpRequest(body, "/payment", http.MethodPost)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert response code status
			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusInternalServerError)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)

			expectedResponse := model.ErrorResponse{Code: http.StatusInternalServerError, Message: expectedErrorMessage}

			// check body response matches the expected response
			test.CheckResponseMessage(response, expectedResponse, t, w)
		}
	}
}

// Handle DB insertion failure
func TestDeletePayment_DBFailureShouldReturn500(t *testing.T) {
	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Create TagDAO request to endpoint:  \"%s\"", "\\tags")
		{
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mocks.NewMockRepository(mockCtrl)

			expectedErrorMessage := "Failed to delete payment"
			body := test.CreatePaymentRequest("31926819", "GB29XABC10161234567801", "USD")
			err := errors.New(expectedErrorMessage)

			// set mock expectation
			mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(err).Times(1)
			mockRepo.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.PaymentResponse{}, nil).Times(1)

			handler := api.NewPaymentHandler(mockRepo, "urlFX", "urlCF")
			router := handler.NewRouter()

			req, err := test.HttpRequest(body, "/payment/5bd7506a9900b30008edf576", http.MethodDelete)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert response code status
			test.AssertForCallErrorAndHttpStatusCode(err, t, w.Code, http.StatusInternalServerError)

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)

			expectedResponse := model.ErrorResponse{Code: http.StatusInternalServerError, Message: expectedErrorMessage}

			// check body response matches the expected response
			test.CheckResponseMessage(response, expectedResponse, t, w)
		}
	}
}
