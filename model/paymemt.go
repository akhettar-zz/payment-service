package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type PaymentResponse struct {
	Data  []Payment `json:"data"`
	Links `json:"links"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:code`
}

// HealthResponse the health json response
type HealthResponse struct {
	Message string `json:"message"`
}

// CreatePaymentRequest
type CreatePaymentRequest struct {
	OrganisationID       string       `json:"organisation_id" binding:"required"`
	BeneficiaryParty     Party        `json:"beneficiary_party" binding:"required"`
	DebtorParty          Party        `json:"debtor_party" binding:"required"`
	PaymentPurpose       string       `json:"payment_purpose"`
	PaymentScheme        string       `json:"payment_scheme"`
	PaymentType          string       `json:"payment_type"`
	Reference            string       `json:"reference"`
	EndToEndReference    string       `json:"end_to_end_reference"`
	SchemePaymentSubType string       `json:"scheme_payment_sub_type"`
	SchemePaymentType    string       `json:"scheme_payment_type"`
	SponsorParty         SponsorParty `json:"sponsor_party"`
	NumericReference     string       `json:"numeric_reference"`
	PaymentID            string       `json:"payment_id"`
	Amount               float64      `json:"amount" binding:"required"`
	BearerCode           string       `json:"bearer_code" binding:"required"`
	ProcessingDate       time.Time    `json:"processing_date"`
}

type CreatePaymentResponse struct {
	Id             string `json:"id"`
	OrganisationId string `json:"organisation_id"`
}

type Links struct {
	Self string `json: "self"`
}

// Payment type
type Payment struct {
	Type           string        `json:"type"`
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Version        int           `json:"version"`
	OrganisationId string        `json:"organisation_id"`
	Attributes     `json: "attributes"`
}

type Attributes struct {
	Amount               float64            `json:"amount"`
	BeneficiaryParty     Party              `json:"beneficiary_party`
	ChargesInformation   ChargesInformation `json:"charges_information"`
	Currency             string             `json:"currency"`
	DebtorParty          Party              `json:"debtor_party"`
	EndToEndReference    string             `json:"end_to_end_reference"`
	Fx                   ForeignExchange    `json:"fx", omitempty`
	NumericReference     string             `json:"numeric_reference"`
	PaymentID            string             `json:"payment_id"`
	PaymentPurpose       string             `json:"payment_purpose"`
	PaymentScheme        string             `json:"payment_scheme"`
	PaymentType          string             `json:"payment_type"`
	ProcessingDate       time.Time          `json:"processing_date"`
	Reference            string             `json:"reference"`
	SchemePaymentSubType string             `json:"scheme_payment_sub_type"`
	SchemePaymentType    string             `json:"scheme_payment_type"`
	SponsorParty         SponsorParty       `json:"sponsor_party"`
}

type SponsorParty struct {
	AccountNumber string `json:"account_number"`
	BankId        string `json:"bank_id"`
	BankIdCode    string `json:"bank_id_code"`
}

type Party struct {
	AccountName       string `json:"account_name`
	AccountNumber     string `json:"account_number`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type", omitempty`
	Address           string `json:"address"`
	BankId            string `json:"bank_id:"`
	BankIdCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
	Currency          string `json:"currency", omitempty`
}

type Charge struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type ChargesInformation struct {
	BearerCode              string   `json:"bearer_code`
	SenderCharges           []Charge `json:"sender_charges`
	ReceiverChargesAmount   float64  `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string   `json:"receiver_charges_currency"`
}

type ForeignExchange struct {
	ContactReference string  `json:"contract_reference"`
	ExchangeRate     float64 `json:"exchange_rate"`
	OriginalAmount   float64 `json:"original_amount"`
	OriginalCurrency string  `json:"original_currency"`
}

type EmptyBody struct{}
