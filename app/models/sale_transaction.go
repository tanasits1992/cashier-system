package models

type SaleTransaction struct {
	PaymentType       PaymentType                   `json:"paymentType"`
	Products          []SaleTransactionProductItem  `json:"products"`
	Discounts         []SaleTransactionDiscountItem `json:"discounts"`
	PriceNetBeforeVat float64                       `json:"priceNetBeforeVat"`
	VatPrice          float64                       `json:"vatPrice"`
	PriceNetAfterVat  float64                       `json:"priceNetAfterVat"`
	DiscountNet       float64                       `json:"discountNet"`
}

type SaleTransactionProductItem struct {
	Barcode               string  `json:"barcode"`
	ProductName           string  `json:"productName"`
	IsVat                 bool    `json:"isVat"`
	UnitName              string  `json:"unitName"`
	PricePerUnitBeforeVat float64 `json:"pricePerUnitBeforeVat"`
	VatPrice              float64 `json:"vatPrice"`
	PricePerUnit          float64 `json:"pricePerUnit"`
	Amount                int     `json:"amount"`
}

type SaleTransactionDiscountItem struct {
	Type     DiscountType `json:"type"`
	Code     string       `json:"code"`
	Discount float64      `json:"discount"`
}
