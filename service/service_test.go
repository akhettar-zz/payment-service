package service_test

import (
	"github.com/payment-service/service"
	"github.com/payment-service/test"
	"testing"
)

func TestGetForeignExchangeService_GetFXDetatils(t *testing.T) {
	t.Logf("Given the need to get forgein exchange details")
	{
		t.Logf("\tWhen invoking foreign exchange service")
		{
			fx := service.NewFxService("url1")
			_, res := fx.GetExchangeRate("USD", "GBP", 100.0)

			if res.ExchangeRate == 2.0 {
				t.Logf("\t\tThe exchange rate is . %v %v", 2.0, test.CheckMark)
			} else {
				t.Errorf("\t\tThe response should contain payment id. %v %v %v", 2.0, test.BallotX, res.ExchangeRate)
			}
		}
	}
}

func TestChargesService_GetChargesDetails(t *testing.T) {
	t.Logf("Given the need to get forgein exchange details")
	{
		t.Logf("\tWhen invoking foreign exchange service")
		{
			fx := service.NewChargesService("url1")
			_, res := fx.GetCharges(2.0, "SHAR", "USD", "GBP")

			if res.ReceiverChargesAmount == 1.0 {
				t.Logf("\t\tThe exchange rate is . %v %v", 1.0, test.CheckMark)
			} else {
				t.Errorf("\t\tThe response should contain payment id. %v %v %v", 1.0, test.BallotX, res.ReceiverChargesAmount)
			}
		}
	}
}
