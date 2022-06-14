package models

type PaymentType string
type DiscountType string

const (
	Cash       PaymentType = "CASH"
	EWallet    PaymentType = "E_WALLET"
	CreditCard PaymentType = "CREDIT_CARD"

	Voucher DiscountType = "VOUCHER"
)

type CashStoreBody struct {
	CashUnit float64 `json:"cashUnit"`
	Amount   int     `json:"amount"`
}

type CashStoreResponseBody struct {
	PriceStorable float64       `json:"priceStorable"`
	StoreDetails  []StoreDetail `json:"storeDetails"`
}

type StoreDetail struct {
	CashUnit        float64 `json:"cashUnit"`
	Amount          int     `json:"amount"`
	RemainingAmount int     `json:"remainingAmount"`
}

type PayRequestBody struct {
	PaymentType    PaymentType     `json:"paymentType"`
	ProductItems   []ProductItem   `json:"productItems"`
	DiscountItems  []DiscountItem  `json:"discountItems"`
	ReceiveDetails []CashStoreBody `json:"receiveDetails"`
}

type PayResponseBody struct {
	BillNo        string          `json:"billNo"`
	PriceNet      float64         `json:"priceNet"`
	Receive       float64         `json:"receive"`
	Return        float64         `json:"return"`
	ReturnDetails []CashStoreBody `json:"returnDetails"`
}

type ProductItem struct {
	Barcode      string  `json:"barcode"`
	ProductName  string  `json:"productName"`
	IsVat        bool    `json:"isVat"`
	UnitName     string  `json:"unitName"`
	PricePerUnit float64 `json:"pricePerUnit"`
	Amount       int     `json:"amount"`
}

type DiscountItem struct {
	Type DiscountType `json:"type"`
	Code string       `json:"code"`
}
