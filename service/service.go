package service

import "payment-service/model"

// FXService the foreign exchange service
type FXService struct {
	url string
}

// ChargesService the charges service
type ChargesService struct {
	url string
}

// Mocking Foreign exchange response
func (fxService FXService) GetExchangeRate(base, currency string, amount float64) (error, model.ForeignExchange) {
	fx := model.ForeignExchange{ContactReference: "FX123", ExchangeRate: 2.00000, OriginalAmount: amount, OriginalCurrency: base}
	return nil, fx
}

// Mocking the Charges service response
func (chService ChargesService) GetCharges(exRate float64, bearerCode string, senderCurrency string, receiverCurrency string) (error, model.ChargesInformation) {
	senderChargesAmount := 10.0
	senderCharges := []model.Charge{{Amount: senderChargesAmount, Currency: senderCurrency}, {Amount: senderChargesAmount / exRate, Currency: receiverCurrency}}
	return nil, model.ChargesInformation{BearerCode: bearerCode, SenderCharges: senderCharges, ReceiverChargesAmount: 1.0, ReceiverChargesCurrency: receiverCurrency}
}

// Creates an instance of FX service
func NewFxService(url string) FXService {
	return FXService{url: url}
}

// Creates an instance Charges service
func NewChargesService(url string) ChargesService {
	return ChargesService{url: url}
}
