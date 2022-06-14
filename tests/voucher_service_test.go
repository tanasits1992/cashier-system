package tests

import (
	"cashier/app/services"
	"testing"
)

func mockTestVoucherServiceSetup(m *testing.T) {
	voucherRepo = NewMockVoucherRepo()
	mockVoucherService = *services.NewVoucherService(voucherRepo)
}

func mockTestVoucherServiceShutdown(m *testing.T) {
	voucherRepo = nil
}

func TestInsertVoucher(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		voucherCode  string
		errorMessage string
	}{
		{
			name:        "insert success",
			input:       InsertVoucherSuccess,
			voucherCode: "1",
		},
		{
			name:         "insert error",
			input:        InsertVoucherError,
			errorMessage: "insert failed",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestVoucherServiceSetup(t)
			defer mockTestVoucherServiceShutdown(t)

			request := mockVoucherRequest(c.input)
			voucherCode, err := mockVoucherService.Insert(request)

			if err != nil {
				Assert(t, "insert voucher error", c.errorMessage, err.Error())
			} else {
				Assert(t, "insert voucher", c.voucherCode, voucherCode)
			}

		})
	}
}

func TestInactivateVoucher(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		errorMessage string
	}{
		{
			name:         "inactivate success",
			input:        InactiveVoucherSuccess,
			errorMessage: "",
		},
		{
			name:         "inactivate get by code error",
			input:        InactiveVoucherGetByCodeError,
			errorMessage: "code is invalid",
		},
		{
			name:         "inactivate already inactivate",
			input:        InactiveVoucherAlreadyInactivated,
			errorMessage: "voucher has already been inactivated",
		},
		{
			name:         "inactivate update error",
			input:        InactiveVoucherUpdateError,
			errorMessage: "update voucher error",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestVoucherServiceSetup(t)
			defer mockTestVoucherServiceShutdown(t)

			err := mockVoucherService.Inactivate(c.input)

			if err != nil {
				Assert(t, "inactivate voucher error", c.errorMessage, err.Error())
			}
		})
	}
}

func TestListVoucher(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		length int
	}{
		{
			name:   "list success",
			input:  InactiveVoucherSuccess,
			length: 3,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestVoucherServiceSetup(t)
			defer mockTestVoucherServiceShutdown(t)

			vouchers, _ := mockVoucherService.List()

			Assert(t, "list voucher error", c.length, len(vouchers))

		})
	}
}

func TestGetByCodeVoucher(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		voucherName  string
		errorMessage string
	}{
		{
			name:        "get by code success",
			input:       GetVoucherSuccess,
			voucherName: "1",
		},
		{
			name:         "get by code error",
			input:        GetVoucherError,
			errorMessage: "code is invalid",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestVoucherServiceSetup(t)
			defer mockTestVoucherServiceShutdown(t)

			voucher, err := mockVoucherService.GetByCode(c.input)
			if err != nil {
				Assert(t, "insert voucher error", c.errorMessage, err.Error())
			} else {
				Assert(t, "insert voucher", c.voucherName, voucher.Name)
			}

		})
	}
}

func TestValidateVoucher(t *testing.T) {
	cases := []struct {
		name         string
		input        string
		pass         bool
		detail       string
		errorMessage string
	}{
		{
			name:  "validate voucher success",
			input: ValidateVoucherSuccess,
			pass:  true,
		},
		{
			name:         "validate voucher get by code error",
			input:        ValidateVoucherGetByCodeError,
			errorMessage: "code is invalid",
		},
		{
			name:   "validate voucher not effect",
			input:  ValidateVoucherNotEffect,
			pass:   false,
			detail: "voucher is not effect",
		},
		{
			name:   "validate voucher expired",
			input:  ValidateVoucherExpired,
			pass:   false,
			detail: "voucher has expired",
		},
		{
			name:   "validate voucher already used",
			input:  ValidateVoucherAlreadyUsed,
			pass:   false,
			detail: "voucher has already been used",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockTestVoucherServiceSetup(t)
			defer mockTestVoucherServiceShutdown(t)

			pass, detail, err := mockVoucherService.Validate(c.input)
			if err != nil {
				Assert(t, "validate voucher error", c.errorMessage, err.Error())
			} else {
				Assert(t, "validate voucher", c.pass, pass)
				Assert(t, "validate voucher", c.detail, detail)
			}

		})
	}
}
