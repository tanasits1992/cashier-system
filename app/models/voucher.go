package models

import "time"

type VoucherCollection struct {
	Name     string    `json:"name"`
	Discount float64   `json:"discount"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Active   bool      `json:"active"`
}

type CreateVoucherRequestBody struct {
	Name     string    `json:"name"`
	Discount float64   `json:"discount"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
}

type InactivateVoucherRequestBody struct {
	Active *bool `json:"active"`
}

type VoucherResponseBody struct {
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Discount float64   `json:"discount"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Active   bool      `json:"active"`
}
