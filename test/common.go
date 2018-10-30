package test

import (
	"bytes"
	"encoding/json"
	"github.com/payment-service/api"
	"github.com/payment-service/logger"
	"github.com/payment-service/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	// CheckMark used for unit test highlight.
	CheckMark = "\u2713"

	// BallotX used for unit test highlight.
	BallotX = "\u2717"
)

// HttpRequest helper
func HttpRequest(jsonReq interface{}, endpoint string, method string) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint, RequestBody(jsonReq))
	if err != nil {
		panic("Failed to marshall json request")
	}
	req.Header.Add("Content-Type", "application/json")
	return req, err
}

// CheckStatus helper
func CheckStatus(w *httptest.ResponseRecorder, t *testing.T, status int) {
	if w.Code == status {
		t.Logf("\t\tShould receive a \"%d\" status. %v", status, CheckMark)
	} else {
		t.Errorf("\t\tShould receive a \"%d\" status. %v %v", status, BallotX, w.Code)
	}
}

// CheckResponseMessage helper.
func CheckResponseMessage(response model.ErrorResponse, expectedResponse model.ErrorResponse, t *testing.T, w *httptest.ResponseRecorder) {
	if response.Message == expectedResponse.Message {
		t.Logf("\t\t\t\tThe body response should  contain a message \"%s\" . %v", expectedResponse.Message, CheckMark)
	} else {
		t.Errorf("\t\t\t\tThe body response should contain a message \"%s\". %v %v", response.Message, BallotX, response.Message)
	}
}

// RequestBody helper
func RequestBody(req interface{}) *bytes.Buffer {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		panic("Failed to marshall json request")
	}
	return bytes.NewBuffer(jsonBytes)
}

func AssertForCallErrorAndHttpStatusCode(err error, t *testing.T, code int, expectedCode int) {
	if err != nil {
		t.Fatal("\t\tShould be able to make the Get call.",
			BallotX, err)
	}
	if code == expectedCode {
		t.Logf("\t\t\t\tShould receive a \"%d\" status. %v", expectedCode, CheckMark)
	} else {
		t.Errorf("\t\t\t\tShould receive a \"%d\" status. %v %v", expectedCode, BallotX, code)
	}
}

// Helper method to create a payment
func CreatePaymentAndAssertResponse(t *testing.T, handler *api.PaymentHandler) model.CreatePaymentResponse {
	body := CreatePaymentRequest("31926819", "GB29XABC10161234567801", "GBP")
	bytes, _ := json.Marshal(body)
	logger.Info.Println(string(bytes))
	w := httptest.NewRecorder()
	req, err := HttpRequest(body, "/payment", http.MethodPost)
	router := handler.NewRouter()
	router.ServeHTTP(w, req)
	if err != nil {
		t.Fatal("\t\tShould be able to make the Get call.",
			BallotX, err)
	}
	if w.Code == http.StatusCreated {
		t.Logf("\t\t\t\tShould receive a \"%d\" status. %v", http.StatusCreated, CheckMark)
	} else {
		t.Errorf("\t\t\t\tShould receive a \"%d\" status. %v %v", http.StatusCreated, BallotX, w.Code)
	}

	var response model.CreatePaymentResponse
	json.NewDecoder(w.Body).Decode(&response)
	return response
}

// Helper function to build CreatePaymentRequest
func CreatePaymentRequest(beneficiaryAccNum, debtorAccNum, beneficiaryCurrency string) model.CreatePaymentRequest {
	beneficiary := model.Party{AccountName: "W Owens", AccountNumber: beneficiaryAccNum, AccountNumberCode: "BBAN",
		AccountType: 0, Address: "1 The Beneficiary Localtown SE2", BankId: "403000", BankIdCode: "GBDSC",
		Name: "Wilfred Jeremiah Owens", Currency: beneficiaryCurrency}

	debtor := model.Party{AccountName: "EJ Brown Black", AccountNumber: debtorAccNum, AccountNumberCode: "IBAN",
		AccountType: 0, Address: "10 Debtor Crescent Sourcetown NE1", BankId: "203301", BankIdCode: "GBDSC",
		Name: "Emelia Jane Brown", Currency: "GBP"}

	return model.CreatePaymentRequest{OrganisationID: "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
		Amount: 200.42, BeneficiaryParty: beneficiary, DebtorParty: debtor, PaymentPurpose: "Paying for goods/services",
		PaymentScheme: "FPS", PaymentType: "Credit", Reference: "Payment for Em's piano lessons",
		SchemePaymentSubType: "InternetBanking", SchemePaymentType: "ImmediatePayment",
		SponsorParty: model.SponsorParty{AccountNumber: "56781234", BankId: "123123", BankIdCode: "GBDSC"},
		BearerCode:   "SHAR", ProcessingDate: time.Now()}
}
