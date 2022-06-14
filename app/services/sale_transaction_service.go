package services

import (
	"cashier/app/models"
	"cashier/fake_database"
)

type SaleTransactionService struct {
	saleTransationRepository fake_database.ISaleTransactionRepository
}

func NewSaleTransactionService(saleTransationRepository fake_database.ISaleTransactionRepository) *SaleTransactionService {
	return &SaleTransactionService{
		saleTransationRepository: saleTransationRepository,
	}
}

func (s *SaleTransactionService) Insert(requestModel models.SaleTransaction) (string, error) {
	billNo, err := s.saleTransationRepository.Insert(requestModel)
	if err != nil {
		return "", nil
	}

	return billNo, nil
}

func (s *SaleTransactionService) GetByBillNo(billNo string) (models.SaleTransaction, error) {
	saleTransaction, err := s.saleTransationRepository.GetByBillNo(billNo)
	if err != nil {
		return models.SaleTransaction{}, err
	}

	return saleTransaction, nil
}
