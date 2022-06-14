package tests

import (
	"cashier/app/models"
	"cashier/app/services"
	"log"
	"math"
	"testing"

	"github.com/joho/godotenv"
)

func mockTestCashStoreServiceSetup(m *testing.T) {
	cashStoreRepo = NewMockCashStoreRepo()
	eWalletConnector = NewMockEWalletConnector()
	creditCardConnector = NewMockCreditCardConnector()
	notificationConnector = NewMockNotificationConnector()
	voucherRepository = NewMockVoucherRepo()
	saleTransactionRepository = NewMockSaleTransactionRepo()

	voucherService := services.NewVoucherService(voucherRepository)
	mockCashService = *services.NewCashService(
		eWalletConnector,
		creditCardConnector,
		notificationConnector,
		*voucherService,
		cashStoreRepo,
		voucherRepository,
		saleTransactionRepository,
	)

	err := godotenv.Load("../env.local")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}

func mockTestCashStoreServiceShutdown(m *testing.T) {
	cashStoreRepo = nil
	eWalletConnector = nil
	creditCardConnector = nil
	notificationConnector = nil
	voucherRepository = nil
	saleTransactionRepository = nil
}

func TestGetStore(t *testing.T) {
	cases := []struct {
		name              string
		priceStorable     float64
		storeDetailLength int
	}{
		{
			name:              "get store success",
			priceStorable:     20412.5,
			storeDetailLength: 9,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestCashStoreServiceSetup(t)
			defer mockTestCashStoreServiceShutdown(t)

			store, _ := mockCashService.GetStore()

			Assert(t, "get store", c.priceStorable, store.PriceStorable)
			Assert(t, "get store", c.storeDetailLength, len(store.StoreDetails))
		})
	}
}

var ReplaceStoreUnusualCash = "replaceStoreUnusualCash"

func TestReplaceStore(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		errorMessage string
	}{
		{
			name:  "replace store success",
			input: ReplaceStoreSuccess,
		},
		{
			name:         "replace store unusual cash",
			input:        ReplaceStoreUnusualCash,
			errorMessage: "[0.5] have not slot",
		},
		{
			name:         "replace store exceed limit",
			input:        ReplaceStoreExceedLimit,
			errorMessage: "[20] exceed limit",
		},
		{
			name:         "replace store update failed",
			input:        ReplaceStoreUpdateFailed,
			errorMessage: "update failed",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestCashStoreServiceSetup(t)
			defer mockTestCashStoreServiceShutdown(t)

			request := mockCashStoreBody(c.input)
			err := mockCashService.ReplaceStore(request)
			if err != nil {
				Assert(t, "replace store", c.errorMessage, err.Error())
			}
		})
	}
}

func TestPay(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		productPrice float64
		recievePrice float64
		returnPrice  float64
		errorMessage string
	}{
		{
			name:         "pay dicount more than product price",
			input:        PayDisountMoreThanPrice,
			productPrice: 0,
			recievePrice: 0,
			returnPrice:  0,
		},
		{
			name:         "pay with e-wallet: success",
			input:        PayWithEWalletSuccess,
			productPrice: 600,
			recievePrice: 0,
			returnPrice:  0,
		},
		{
			name:         "pay with e-wallet: error",
			input:        PayWithEWalletError,
			errorMessage: "now, can not pay with e-wallet",
		},
		{
			name:         "pay with credit card: success",
			input:        PayWithCreditCardSuccess,
			productPrice: 800,
			recievePrice: 0,
			returnPrice:  0,
		},
		{
			name:         "pay with credit card: error",
			input:        PayWithCreditCardError,
			errorMessage: "now, can not pay with credit card",
		},
		{
			name:         "pay with cash: error (unusual cash)",
			input:        PayWithCashErrorUnusualCash,
			errorMessage: "[0.5] have not slot",
		},
		{
			name:         "pay with cash: error (product price more than store)",
			input:        PayWithCashErrorOverProductPrice,
			errorMessage: "product price more than store",
		},
		{
			name:         "pay with cash: error (incomplete payment)",
			input:        PayWithCashErrorIncompletePayment,
			errorMessage: "incomplete payment (100)",
		},
		{
			name:         "pay with cash: error (no change and store not enough)",
			input:        PayWithCashErrorNoChangeStoreNotEnough,
			errorMessage: "store not enough",
		},
		{
			name:         "pay with cash: error (no change but update store failed)",
			input:        PayWithCashErrorNoChangeUpdateStoreFailed,
			errorMessage: "update store error",
		},
		{
			name:         "pay with cash: success (no change)",
			input:        PayWithCashSuccessNoChange,
			productPrice: 2000,
			recievePrice: 2000,
			returnPrice:  0,
		},
		{
			name:         "pay with cash: error (change store not enough)",
			input:        PayWithCashErrorChangeStoreNotEnough,
			errorMessage: "store not enough",
		},
		{
			name:         "pay with cash: error (change not enough)",
			input:        PayWithCashErrorChangeNotEnough,
			errorMessage: "change not enough",
		},
		{
			name:         "pay with cash: success (change)",
			input:        PayWithCashSuccessChange,
			productPrice: 8700,
			recievePrice: 9000,
			returnPrice:  300,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestCashStoreServiceSetup(t)
			defer mockTestCashStoreServiceShutdown(t)

			request := mockPayRequest(c.input)
			response, err := mockCashService.Pay(request)
			if err != nil {
				Assert(t, "pay error", c.errorMessage, err.Error())
			} else {
				Assert(t, "pay product price success", c.productPrice, response.PriceNet)
				Assert(t, "pay receive price success", c.recievePrice, response.Receive)
				Assert(t, "pay return price success", c.returnPrice, response.Return)
			}
		})
	}
}

func mockPayRequest(input string) models.PayRequestBody {
	if input == PayDisountMoreThanPrice {
		return mockPayRequestNormal(models.EWallet, 600, true)
	} else if input == PayWithEWalletSuccess {
		return mockPayRequestNormal(models.EWallet, 600, false)
	} else if input == PayWithCreditCardSuccess {
		return mockPayRequestNormal(models.CreditCard, 800, false)
	} else if input == PayWithEWalletError {
		return mockPayRequestNormal(models.EWallet, PriceNetError, false)
	} else if input == PayWithCreditCardError {
		return mockPayRequestNormal(models.CreditCard, PriceNetError, false)
	} else if input == PayWithCashErrorUnusualCash {
		payment := mockPayRequestNormal(models.Cash, 1000.50, false)

		receiveDetails := []models.CashStoreBody{}
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 1000,
			Amount:   1,
		})
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 0.50,
			Amount:   1,
		})
		payment.ReceiveDetails = receiveDetails

		return payment
	} else if input == PayWithCashErrorOverProductPrice {
		payment := mockPayRequestNormal(models.Cash, 30000, false)

		receiveDetails := []models.CashStoreBody{}
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 1000,
			Amount:   30,
		})
		payment.ReceiveDetails = receiveDetails

		return payment
	} else if input == PayWithCashErrorIncompletePayment {
		payment := mockPayRequestNormal(models.Cash, 1000, false)

		receiveDetails := []models.CashStoreBody{}
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 500,
			Amount:   1,
		})
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 100,
			Amount:   4,
		})
		payment.ReceiveDetails = receiveDetails

		return payment
	} else if input == PayWithCashErrorNoChangeStoreNotEnough {
		payment := mockPayRequestNormal(models.Cash, 10000, false)

		receiveDetails := []models.CashStoreBody{}
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 500,
			Amount:   20,
		})
		payment.ReceiveDetails = receiveDetails

		return payment
	} else if input == PayWithCashErrorNoChangeUpdateStoreFailed {
		payment := mockPayRequestNormal(models.Cash, 300, false)

		receiveDetails := []models.CashStoreBody{}
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 100,
			Amount:   3,
		})
		payment.ReceiveDetails = receiveDetails

		return payment
	} else if input == PayWithCashSuccessNoChange {
		payment := mockPayRequestNormal(models.Cash, 2000, false)

		receiveDetails := []models.CashStoreBody{}
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 1000,
			Amount:   2,
		})
		payment.ReceiveDetails = receiveDetails

		return payment
	} else if input == PayWithCashErrorChangeStoreNotEnough {
		payment := mockPayRequestNormal(models.Cash, 19350, false)

		receiveDetails := []models.CashStoreBody{}
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 1000,
			Amount:   20,
		})
		payment.ReceiveDetails = receiveDetails

		return payment
	} else if input == PayWithCashErrorChangeNotEnough {
		payment := mockPayRequestNormal(models.Cash, 1005, false)

		receiveDetails := []models.CashStoreBody{}
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 1000,
			Amount:   1,
		})
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 10,
			Amount:   1,
		})
		payment.ReceiveDetails = receiveDetails

		return payment
	} else if input == PayWithCashSuccessChange {
		payment := mockPayRequestNormal(models.Cash, 8700, false)

		receiveDetails := []models.CashStoreBody{}
		receiveDetails = append(receiveDetails, models.CashStoreBody{
			CashUnit: 1000,
			Amount:   9,
		})
		payment.ReceiveDetails = receiveDetails

		return payment
	}
	return mockPayRequestNormal(models.EWallet, 100, false)
}

func mockPayRequestNormal(paymentType models.PaymentType, productPrice float64, discountF bool) models.PayRequestBody {
	price1 := math.Floor(productPrice*1*100/3/2) / 100
	price2 := productPrice - (price1 * 2)

	productItems := []models.ProductItem{}
	productItems = append(productItems, models.ProductItem{
		Barcode:      "001",
		ProductName:  "A01",
		IsVat:        true,
		UnitName:     "ชิ้น",
		PricePerUnit: price1,
		Amount:       2,
	})
	productItems = append(productItems, models.ProductItem{
		Barcode:      "002",
		ProductName:  "A02",
		IsVat:        false,
		UnitName:     "ชิ้น",
		PricePerUnit: price2,
		Amount:       1,
	})

	discountItems := []models.DiscountItem{}
	if discountF {
		discountItems = append(discountItems, models.DiscountItem{
			Type: models.Voucher,
			Code: VoucherCode,
		})
	}

	return models.PayRequestBody{
		PaymentType:   paymentType,
		ProductItems:  productItems,
		DiscountItems: discountItems,
	}
}
