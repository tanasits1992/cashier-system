package services

import (
	"cashier/app/models"
	"cashier/connectors"
	"cashier/fake_database"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

type CashService struct {
	eWalletConnector          connectors.IEWalletConnector
	creditCardConnector       connectors.ICreditCardConnector
	notificationConnector     connectors.INotificationConnector
	voucherService            VoucherService
	cashStoreRepository       fake_database.ICashStoreRepository
	voucherRepository         fake_database.IVoucherRepository
	saleTransactionRepository fake_database.ISaleTransactionRepository
}

func NewCashService(
	eWalletConnector connectors.IEWalletConnector,
	creditCardConnector connectors.ICreditCardConnector,
	notificationConnector connectors.INotificationConnector,
	voucherService VoucherService,
	cashStoreRepository fake_database.ICashStoreRepository,
	voucherRepository fake_database.IVoucherRepository,
	saleTransactionRepository fake_database.ISaleTransactionRepository,
) *CashService {
	return &CashService{
		eWalletConnector:          eWalletConnector,
		creditCardConnector:       creditCardConnector,
		notificationConnector:     notificationConnector,
		voucherService:            voucherService,
		cashStoreRepository:       cashStoreRepository,
		voucherRepository:         voucherRepository,
		saleTransactionRepository: saleTransactionRepository,
	}
}

// function
func (s *CashService) getStoreLimit() map[float64]int {
	limitString := os.Getenv("CASH_LIIMT")

	cashLimit := make(map[string]int)
	json.Unmarshal([]byte(limitString), &cashLimit)

	var limit = make(map[float64]int)
	for key, val := range cashLimit {
		keyNumber, _ := strconv.ParseFloat(key, 64)
		limit[keyNumber] = val
	}

	return limit
}

func (s *CashService) calculateFinalPriceNet(products []models.ProductItem, discounts []models.DiscountItem) ([]models.VoucherResponseBody, float64) {
	priceNet := 0.0
	discountNet := 0.0

	for _, productItem := range products {
		priceNet += productItem.PricePerUnit * float64(productItem.Amount)
	}

	vouchers := []models.VoucherResponseBody{}
	for _, discount := range discounts {
		if discount.Type == models.Voucher {
			voucher, err := s.voucherRepository.GetByCode(discount.Code)
			if err != nil {
				continue
			}

			pass, _, err := s.voucherService.Validate(discount.Code)
			if !pass || err != nil {
				continue
			}

			discountNet += voucher.Discount
			vouchers = append(vouchers, models.VoucherResponseBody{
				Name:     voucher.Name,
				Code:     discount.Code,
				Discount: voucher.Discount,
			})
		}
	}
	return vouchers, priceNet - discountNet
}

func (s *CashService) round(number float64, roundBy float64) float64 {
	mod := math.Mod(number, roundBy)

	if mod >= (roundBy / 2) {
		return number + roundBy - mod
	}

	return number - mod
}

func (s *CashService) pay(requestModel models.PayRequestBody, priceNet float64) (*models.PayResponseBody, error) {
	exchangeDetail := models.PayResponseBody{
		PriceNet: priceNet,
	}

	if requestModel.PaymentType == models.EWallet {
		err := s.payWithEWallet(priceNet, requestModel)
		if err != nil {
			return nil, err
		}
	} else if requestModel.PaymentType == models.CreditCard {
		err := s.payWithCreditCard(priceNet, requestModel)
		if err != nil {
			return nil, err
		}
	} else if requestModel.PaymentType == models.Cash {
		detail, err := s.payWithCash(priceNet, requestModel)
		if err != nil {
			return nil, err
		}

		exchangeDetail.Receive = detail.Receive
		exchangeDetail.Return = detail.Return
		exchangeDetail.ReturnDetails = detail.ReturnDetails
	}

	return &exchangeDetail, nil
}

func (s *CashService) sumCashDetail(receiveDetails []models.CashStoreBody) float64 {
	sum := 0.0
	for _, receiveDetail := range receiveDetails {
		sum += receiveDetail.CashUnit * float64(receiveDetail.Amount)
	}
	return sum
}

func (s *CashService) checkStoreLimit(cashStoreMap map[float64]int) bool {
	limit := s.getStoreLimit()

	for base, amount := range cashStoreMap {
		if amount > limit[base] {
			return true
		}
	}

	return false
}

func (s *CashService) updateStore(storeMap map[float64]int) error {
	request := []models.CashStoreBody{}

	for base, amount := range storeMap {
		request = append(request, models.CashStoreBody{
			CashUnit: base,
			Amount:   amount,
		})
	}

	err := s.ReplaceStore(request)
	if err != nil {
		return err
	}
	return nil
}

func (s *CashService) payWithEWallet(priceNet float64, requestModel models.PayRequestBody) error {
	err := s.eWalletConnector.PayWithEWallet(priceNet)
	return err
}

func (s *CashService) payWithCreditCard(priceNet float64, requestModel models.PayRequestBody) error {
	err := s.creditCardConnector.PayWithCreditCard(priceNet)
	return err
}

func (s *CashService) constructSaleTransactionBody(requestModel models.PayRequestBody, vouchers []models.VoucherResponseBody) models.SaleTransaction {
	vatString := os.Getenv("VAT_RATE")
	vatRate, _ := strconv.ParseFloat(vatString, 64)

	products := []models.SaleTransactionProductItem{}
	priceNetBeforeVat := 0.0
	vatPrice := 0.0
	priceNetAfterVat := 0.0

	for _, productItem := range requestModel.ProductItems {
		priceBeforeVat := productItem.PricePerUnit
		vat := 0.0

		if productItem.IsVat {
			// 2 decimal
			priceBeforeVat = math.Round((productItem.PricePerUnit/((100+vatRate)/100))*100) / 100
			vat = math.Round((productItem.PricePerUnit-priceBeforeVat)*100) / 100
		}

		products = append(products, models.SaleTransactionProductItem{
			Barcode:               productItem.Barcode,
			ProductName:           productItem.ProductName,
			IsVat:                 productItem.IsVat,
			UnitName:              productItem.UnitName,
			VatPrice:              vat,
			PricePerUnitBeforeVat: priceBeforeVat,
			PricePerUnit:          productItem.PricePerUnit,
			Amount:                productItem.Amount,
		})

		priceNetBeforeVat += math.Round(priceBeforeVat*float64(productItem.Amount)*100) / 100
		vatPrice += math.Round(vat*float64(productItem.Amount)*100) / 100
		priceNetAfterVat += math.Round(productItem.PricePerUnit*float64(productItem.Amount)*100) / 100
	}

	discounts := []models.SaleTransactionDiscountItem{}
	discountNet := 0.0

	for _, voucher := range vouchers {
		discounts = append(discounts, models.SaleTransactionDiscountItem{
			Type:     models.Voucher,
			Code:     voucher.Code,
			Discount: voucher.Discount,
		})
		discountNet += voucher.Discount
	}

	return models.SaleTransaction{
		PaymentType:       requestModel.PaymentType,
		Products:          products,
		Discounts:         discounts,
		PriceNetBeforeVat: priceNetBeforeVat,
		VatPrice:          vatPrice,
		PriceNetAfterVat:  priceNetAfterVat,
		DiscountNet:       discountNet,
	}
}

func (s *CashService) notificationLogic() error {
	reachLimitBase := []float64{}

	limit := s.getStoreLimit()
	store, err := s.cashStoreRepository.Get()
	if err != nil {
		return err
	}

	for base, amount := range store {
		if amount == limit[base] {
			reachLimitBase = append(reachLimitBase, base)
		}
	}

	if len(reachLimitBase) > 0 {
		err := s.sendNotification(reachLimitBase)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *CashService) sendNotification(reachLimitBase []float64) error {
	return s.notificationConnector.Notification(reachLimitBase)
}

func (s *CashService) validateCashInvalid(receiveDetails []models.CashStoreBody) error {
	limit := s.getStoreLimit()
	invalidCashes := []float64{}

	for _, receiveDetail := range receiveDetails {
		_, ok := limit[receiveDetail.CashUnit]
		if !ok {
			invalidCashes = append(invalidCashes, receiveDetail.CashUnit)
		}
	}

	if len(invalidCashes) > 0 {
		return fmt.Errorf("%v have not slot", invalidCashes)
	}
	return nil
}

func (s *CashService) payWithCash(price float64, requestModel models.PayRequestBody) (*models.PayResponseBody, error) {
	if err := s.validateCashInvalid(requestModel.ReceiveDetails); err != nil {
		return nil, err
	}

	priceNet := s.round(price, 0.25)

	// validate price net and store
	store, err := s.GetStore()
	if err != nil {
		return nil, err
	}
	if store.PriceStorable < priceNet {
		return nil, errors.New("product price more than store")
	}

	// check payment is complete
	receivePrice := s.sumCashDetail(requestModel.ReceiveDetails)
	if receivePrice < priceNet {
		return nil, fmt.Errorf("incomplete payment (%v)", priceNet-receivePrice)
	}

	// convert store to map
	storeMap := make(map[float64]int)
	for _, storeDetail := range store.StoreDetails {
		storeMap[storeDetail.CashUnit] = storeDetail.Amount
	}

	// add receive cash to store
	for _, receiveDetail := range requestModel.ReceiveDetails {
		storeMap[receiveDetail.CashUnit] += receiveDetail.Amount
	}

	// check limit
	exceedLimit := s.checkStoreLimit(storeMap)

	// not change
	if receivePrice == priceNet {
		if exceedLimit {
			return nil, errors.New("store not enough")
		}

		err := s.updateStore(storeMap)
		if err != nil {
			return nil, err
		}

		return &models.PayResponseBody{
			PriceNet: priceNet,
			Receive:  receivePrice,
		}, nil
	} else {
		change := receivePrice - priceNet

		returnDetails, err := s.changeProcess(change, storeMap, exceedLimit)
		if err != nil {
			return nil, err
		}

		return &models.PayResponseBody{
			PriceNet:      priceNet,
			Receive:       receivePrice,
			Return:        change,
			ReturnDetails: returnDetails,
		}, nil
	}
}

func (s *CashService) changeProcess(change float64, storeMap map[float64]int, exceedLimit bool) ([]models.CashStoreBody, error) {
	returnDetail := make(map[float64]int)
	limit := s.getStoreLimit()

	// change by exceed limti
	if exceedLimit {
		sum := 0.0
		for base, amount := range storeMap {
			if amount > limit[base] {
				returnDetail[base] = amount - limit[base]
				sum += float64(amount-limit[base]) * base
				storeMap[base] = limit[base]
			}
		}

		if sum > change {
			return nil, errors.New("store not enough")
		}
		change -= sum
	}

	// sort by base
	keys := []float64{}
	for k := range storeMap {
		keys = append(keys, float64(k))
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

	for _, base := range keys {
		if change >= base {
			amount := math.Floor(float64(change / base))

			if storeMap[base] < int(amount) {
				continue
			}

			storeMap[base] -= int(amount)
			returnDetail[base] += int(amount)
			change -= amount * base
		}
	}

	if change > 0 {
		return nil, errors.New("change not enough")
	}

	err := s.updateStore(storeMap)
	if err != nil {
		return nil, err
	}

	// convert return detail to slice
	result := []models.CashStoreBody{}
	for base, amount := range returnDetail {
		result = append(result, models.CashStoreBody{
			CashUnit: base,
			Amount:   amount,
		})
	}

	return result, nil
}

// service
func (s *CashService) GetStore() (*models.CashStoreResponseBody, error) {
	store, err := s.cashStoreRepository.Get()
	if err != nil {
		return nil, err
	}
	limit := s.getStoreLimit()

	priceStorable := 0.0
	storeDetails := []models.StoreDetail{}

	for base, amount := range store {
		remainingAmount := limit[base] - amount

		storeDetails = append(storeDetails, models.StoreDetail{
			CashUnit:        base,
			Amount:          amount,
			RemainingAmount: remainingAmount,
		})
		priceStorable += float64(remainingAmount) * base
	}

	return &models.CashStoreResponseBody{
		PriceStorable: priceStorable,
		StoreDetails:  storeDetails,
	}, nil
}

func (s *CashService) ReplaceStore(requestModel []models.CashStoreBody) error {
	if err := s.validateCashInvalid(requestModel); err != nil {
		return err
	}

	request := make(map[float64]int)
	for _, requestModel := range requestModel {
		request[requestModel.CashUnit] = requestModel.Amount
	}

	store := make(map[float64]int)
	exceedLimit := []float64{}
	limit := s.getStoreLimit()

	for base, amount := range limit {
		if request[base] > amount {
			exceedLimit = append(exceedLimit, base)
		} else {
			store[base] = request[base]
		}
	}

	if len(exceedLimit) > 0 {
		return fmt.Errorf("%v exceed limit", exceedLimit)
	}

	err := s.cashStoreRepository.Update(store)
	if err != nil {
		return err
	}
	return nil
}

func (s *CashService) Pay(requestModel models.PayRequestBody) (*models.PayResponseBody, error) {
	exchangeDetail := &models.PayResponseBody{}
	var err error

	vouchers, priceNet := s.calculateFinalPriceNet(requestModel.ProductItems, requestModel.DiscountItems)

	if priceNet > 0 {
		exchangeDetail, err = s.pay(requestModel, priceNet)
		if err != nil {
			return nil, err
		}
	}

	// Inactivate voucher
	for _, voucher := range vouchers {
		s.voucherService.Inactivate(voucher.Code)
	}

	// save sale transaction for sale invoice
	saleTransaction := s.constructSaleTransactionBody(requestModel, vouchers)
	billNo, _ := s.saleTransactionRepository.Insert(saleTransaction)
	exchangeDetail.BillNo = billNo

	// send noti when reach to limit at least 1
	s.notificationLogic()

	return exchangeDetail, nil
}
