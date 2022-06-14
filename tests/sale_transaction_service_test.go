package tests

import (
	"cashier/app/models"
	"cashier/app/services"
	"testing"
)

func mockTestSaleTransactionServiceSetup(m *testing.T) {
	saleTransactionRepo = NewMockSaleTransactionRepo()
	mockSaleTransactionService = *services.NewSaleTransactionService(saleTransactionRepo)
}

func mockTestSaleTransactionServiceShutdown(m *testing.T) {
	saleTransactionRepo = nil
}

func TestInsertSaleTransaction(t *testing.T) {
	cases := []struct {
		name         string
		input        models.PaymentType
		billNo       string
		errorMessage string
	}{
		{
			name:   "insert sale transaction success",
			input:  InsertSaleTransactionSuccess,
			billNo: "1",
		},
		{
			name:         "insert sale transaction error",
			input:        InsertSaleTransactionError,
			errorMessage: "insert failed",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestSaleTransactionServiceSetup(t)
			defer mockTestSaleTransactionServiceShutdown(t)

			request := mockSaleTransactionRequest(c.input)
			billNo, err := mockSaleTransactionService.Insert(request)

			if err != nil {
				Assert(t, "insert sale transaction error", c.errorMessage, err.Error())
			} else {
				Assert(t, "insert sale transaction", c.billNo, billNo)
			}
		})
	}
}
func TestGetByBillNoSaleTransaction(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		paymentType  models.PaymentType
		errorMessage string
	}{
		{
			name:        "get by code success",
			input:       GetSaleTransactionSuccess,
			paymentType: models.Cash,
		},
		{
			name:         "get by code error",
			input:        GetSaleTransactionError,
			errorMessage: "code is invalid",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestSaleTransactionServiceSetup(t)
			defer mockTestSaleTransactionServiceShutdown(t)

			sale, err := mockSaleTransactionService.GetByBillNo(c.input)

			if err != nil {
				Assert(t, "get sale transaction error", c.errorMessage, err.Error())
			} else {
				Assert(t, "get sale transaction", c.paymentType, sale.PaymentType)
			}
		})
	}
}
