package tests

import (
	"cashier/app/models"
	"cashier/app/services"
	"cashier/connectors"
	"cashier/fake_database"
	"errors"
	"testing"
	"time"
)

// variable
var VoucherCode = "voucherForTestPay"
var VoucherDiscount = 1000
var PriceNetError = 1000000.0
var StoreNetError = 3320.0

// case
var InsertVoucherSuccess = "insertVoucherSuccess"
var InsertVoucherError = "insertVoucherError"

var InactiveVoucherSuccess = "inactiveVoucherSuccess"
var InactiveVoucherGetByCodeError = "inactiveVoucherGetByCodeError"
var InactiveVoucherAlreadyInactivated = "inactiveVoucherAlreadyInactivated"
var InactiveVoucherUpdateError = "inactiveVoucherUpdateError"

var GetVoucherSuccess = "getVoucherSuccess"
var GetVoucherError = "getVoucherError"

var ValidateVoucherSuccess = "ValidateVoucherSuccess"
var ValidateVoucherGetByCodeError = "ValidateVoucherGetByCodeError"
var ValidateVoucherNotEffect = "validateVoucherNotEffect"
var ValidateVoucherExpired = "validateVoucherExpired"
var ValidateVoucherAlreadyUsed = "validateVoucherAlreadyUsed"

var InsertSaleTransactionSuccess = models.Cash
var InsertSaleTransactionError = models.CreditCard

var GetSaleTransactionSuccess = "getSaleTransactionSuccess"
var GetSaleTransactionError = "getSaleTransactionError"

var ReplaceStoreSuccess = "replaceStoreSuccess"
var ReplaceStoreExceedLimit = "replaceStoreExceedLimit"
var ReplaceStoreUpdateFailed = "replaceStoreUpdateFailed"

var PayDisountMoreThanPrice = "payDisountMoreThanPrice"
var PayWithEWalletSuccess = "payWithEWalletSuccess"
var PayWithEWalletError = "payWithEWalletError"
var PayWithCreditCardSuccess = "payWithCreditCardSuccess"
var PayWithCreditCardError = "payWithCreditCardError"
var PayWithCashErrorUnusualCash = "payErrorUnusualCash"
var PayWithCashErrorOverProductPrice = "payWithCashErrorOverProductPrice"
var PayWithCashErrorIncompletePayment = "payWithCashErrorIncompletePayment"
var PayWithCashErrorNoChangeStoreNotEnough = "payWithCashErrorNoChangeStoreNotEnough"
var PayWithCashErrorNoChangeUpdateStoreFailed = "payWithCashErrorNoChangeUpdateStoreFailed"
var PayWithCashSuccessNoChange = "payWithCashSuccessNoChange"
var PayWithCashErrorChangeStoreNotEnough = "payWithCashErrorChangeStoreNotEnough"
var PayWithCashErrorChangeNotEnough = "payWithCashErrorChangeNotEnough"
var PayWithCashSuccessChange = "payWithCashSuccessChange"

// mock repository
var mockVoucherService services.VoucherService
var voucherRepo fake_database.IVoucherRepository

var mockSaleTransactionService services.SaleTransactionService
var saleTransactionRepo fake_database.ISaleTransactionRepository

var mockCashService services.CashService
var cashStoreRepo fake_database.ICashStoreRepository

var eWalletConnector connectors.IEWalletConnector
var creditCardConnector connectors.ICreditCardConnector
var notificationConnector connectors.INotificationConnector
var voucherRepository fake_database.IVoucherRepository
var saleTransactionRepository fake_database.ISaleTransactionRepository

type mockVoucherRepo struct{}
type mockSaleTransactionRepo struct{}
type mockCashStoreRepo struct{}
type mockEWalletConnector struct{}
type mockCreditCardConnector struct{}
type mockNotificationConnector struct{}

func NewMockVoucherRepo() *mockVoucherRepo {
	return &mockVoucherRepo{}
}
func NewMockSaleTransactionRepo() *mockSaleTransactionRepo {
	return &mockSaleTransactionRepo{}
}
func NewMockCashStoreRepo() *mockCashStoreRepo {
	return &mockCashStoreRepo{}
}
func NewMockEWalletConnector() *mockEWalletConnector {
	return &mockEWalletConnector{}
}
func NewMockCreditCardConnector() *mockCreditCardConnector {
	return &mockCreditCardConnector{}
}
func NewMockNotificationConnector() *mockNotificationConnector {
	return &mockNotificationConnector{}
}

// apply interface
func (s *mockVoucherRepo) Load(filename string) error {
	return nil
}
func (s *mockVoucherRepo) List() (map[string]models.VoucherCollection, error) {
	return map[string]models.VoucherCollection{
		"1": {},
		"2": {},
		"3": {},
	}, nil
}
func (s *mockVoucherRepo) GetByCode(code string) (models.VoucherCollection, error) {
	now := time.Now()
	if code == InactiveVoucherGetByCodeError || code == GetVoucherError || code == ValidateVoucherGetByCodeError {
		return models.VoucherCollection{}, errors.New("code is invalid")
	} else if code == InactiveVoucherAlreadyInactivated || code == ValidateVoucherAlreadyUsed {
		return mockVoucherResponse(code, false), nil
	} else if code == GetVoucherSuccess || code == VoucherCode {
		return mockVoucherResponse("1", true), nil
	} else if code == ValidateVoucherNotEffect {
		return mockVoucherResponseSpecifyRange(now.AddDate(0, 0, 2), now.AddDate(0, 0, 30)), nil
	} else if code == ValidateVoucherExpired {
		return mockVoucherResponseSpecifyRange(now.AddDate(0, 0, -2), now.AddDate(0, 0, -24)), nil
	}

	return mockVoucherResponse(code, true), nil
}
func (s *mockVoucherRepo) Insert(voucher models.VoucherCollection) (string, error) {
	if voucher.Name == InsertVoucherError {
		return "", errors.New("insert failed")
	}

	return "1", nil
}
func (s *mockVoucherRepo) Update(code string, voucher models.VoucherCollection) error {
	if code == InactiveVoucherUpdateError {
		return errors.New("update voucher error")
	}
	return nil
}

func (s *mockSaleTransactionRepo) Insert(requestModel models.SaleTransaction) (string, error) {
	if requestModel.PaymentType == InsertSaleTransactionError {
		return "", errors.New("insert failed")
	}
	return "1", nil
}
func (s *mockSaleTransactionRepo) GetByBillNo(billNo string) (models.SaleTransaction, error) {
	if billNo == GetSaleTransactionError {
		return models.SaleTransaction{}, errors.New("code is invalid")
	}
	return models.SaleTransaction{
		PaymentType: models.Cash,
	}, nil
}

func (r *mockCashStoreRepo) Load(filename string) error {
	return nil
}
func (r *mockCashStoreRepo) Get() (map[float64]int, error) {
	return mockCashStore(), nil
}
func (r *mockCashStoreRepo) Update(cash map[float64]int) error {
	store := 0.0

	for base, amount := range cash {
		store += base * float64(amount)
	}

	if store == StoreNetError {
		return errors.New("update store error")
	}

	return nil
}

func (c *mockEWalletConnector) PayWithEWallet(priceNet float64) error {
	if priceNet == PriceNetError {
		return errors.New("now, can not pay with e-wallet")
	}
	return nil
}
func (c *mockCreditCardConnector) PayWithCreditCard(priceNet float64) error {
	if priceNet == PriceNetError {
		return errors.New("now, can not pay with credit card")
	}
	return nil
}
func (c *mockNotificationConnector) Notification(reachLimitBase []float64) error {
	return nil
}

// function mock
func mockVoucherRequest(name string) models.CreateVoucherRequestBody {
	return models.CreateVoucherRequestBody{
		Name:     name,
		Discount: 1000,
		Start:    time.Now().AddDate(0, 0, -12),
		End:      time.Now().AddDate(0, 0, 12),
	}
}
func mockVoucherResponse(name string, active bool) models.VoucherCollection {
	return models.VoucherCollection{
		Name:     name,
		Discount: float64(VoucherDiscount),
		Start:    time.Now().AddDate(0, 0, -12),
		End:      time.Now().AddDate(0, 0, 12),
		Active:   active,
	}
}
func mockVoucherResponseSpecifyRange(start time.Time, end time.Time) models.VoucherCollection {
	return models.VoucherCollection{
		Name:     "test",
		Discount: float64(VoucherDiscount),
		Start:    start,
		End:      end,
		Active:   true,
	}
}

func mockSaleTransactionRequest(payment models.PaymentType) models.SaleTransaction {
	return models.SaleTransaction{
		PaymentType: payment,
	}
}

func mockCashStore() map[float64]int {
	return map[float64]int{
		1000: 1,
		500:  2,
		100:  5,
		50:   3,
		20:   14,
		10:   9,
		5:    0,
		1:    0,
		0.25: 0,
	}
}
func mockCashStoreBody(input string) []models.CashStoreBody {
	mock := []models.CashStoreBody{}
	mock = append(mock, models.CashStoreBody{
		CashUnit: 500,
		Amount:   2,
	})
	mock = append(mock, models.CashStoreBody{
		CashUnit: 100,
		Amount:   4,
	})

	if input == ReplaceStoreUnusualCash {
		mock = append(mock, models.CashStoreBody{
			CashUnit: 0.5,
			Amount:   4,
		})
	}

	return mock
}

// function for test
func Assert(t *testing.T, testCase string, expect interface{}, actual interface{}) {
	if expect != actual {
		t.Fatalf("Assert %v. \r\nexpected : %v\r\nactual   : %v", testCase, expect, actual)
	}
}
