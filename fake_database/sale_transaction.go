package fake_database

import (
	"cashier/app/models"
	"errors"

	"github.com/google/uuid"
)

type SaleTransactionRepository struct{}

func NewSaleTransactionRepository() ISaleTransactionRepository {
	return &SaleTransactionRepository{}
}

type ISaleTransactionRepository interface {
	Insert(requestModel models.SaleTransaction) (string, error)
	GetByBillNo(billNo string) (models.SaleTransaction, error)
}

var saleTransactionData = make(map[string]models.SaleTransaction)

func (r *SaleTransactionRepository) Insert(requestModel models.SaleTransaction) (string, error) {
	biilNo := uuid.New().String()
	saleTransactionData[biilNo] = requestModel

	return biilNo, nil
}

func (r *SaleTransactionRepository) GetByBillNo(billNo string) (models.SaleTransaction, error) {
	saleTransaction, ok := saleTransactionData[billNo]
	if !ok {
		return models.SaleTransaction{}, errors.New("bill no is invalid")
	}

	return saleTransaction, nil

}
