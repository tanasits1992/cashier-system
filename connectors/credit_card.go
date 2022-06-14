package connectors

import "fmt"

type CreditCardConnector struct{}

func NewCerditCardConnector() ICreditCardConnector {
	return &CreditCardConnector{}
}

type ICreditCardConnector interface {
	PayWithCreditCard(priceNet float64) error
}

func (c *CreditCardConnector) PayWithCreditCard(priceNet float64) error {
	fmt.Println("pay with credit card")
	return nil
}
