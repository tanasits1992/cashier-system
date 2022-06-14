package connectors

import "fmt"

type EWalletConnector struct{}

func NewEWalletConnector() *EWalletConnector {
	return &EWalletConnector{}
}

type IEWalletConnector interface {
	PayWithEWallet(priceNet float64) error
}

func (c *EWalletConnector) PayWithEWallet(priceNet float64) error {
	fmt.Println("pay with e-wallet")
	return nil
}
