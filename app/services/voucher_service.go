package services

import (
	"cashier/app/models"
	"cashier/fake_database"
	"errors"
	"time"
)

type VoucherService struct {
	voucherRepository fake_database.IVoucherRepository
}

func NewVoucherService(voucherRepository fake_database.IVoucherRepository) *VoucherService {
	return &VoucherService{
		voucherRepository: voucherRepository,
	}
}

func (v *VoucherService) Insert(requestModel models.CreateVoucherRequestBody) (string, error) {
	voucher := models.VoucherCollection{
		Name:     requestModel.Name,
		Start:    requestModel.Start,
		End:      requestModel.End,
		Discount: requestModel.Discount,
		Active:   true,
	}

	voucherCode, err := v.voucherRepository.Insert(voucher)
	if err != nil {
		return "", err
	}
	return voucherCode, nil
}

func (v *VoucherService) Inactivate(code string) error {
	voucher, err := v.voucherRepository.GetByCode(code)
	if err != nil {
		return err
	}

	if !voucher.Active {
		return errors.New("voucher has already been inactivated")
	}

	voucher.Active = false

	err = v.voucherRepository.Update(code, voucher)
	if err != nil {
		return err
	}

	return nil
}

func (v *VoucherService) List() ([]models.VoucherResponseBody, error) {
	vouchersDto := []models.VoucherResponseBody{}

	vouchers, err := v.voucherRepository.List()
	if err != nil {
		return vouchersDto, err
	}

	for key, voucher := range vouchers {
		vouchersDto = append(vouchersDto, models.VoucherResponseBody{
			Name:     voucher.Name,
			Code:     key,
			Start:    voucher.Start,
			End:      voucher.End,
			Discount: voucher.Discount,
			Active:   voucher.Active,
		})
	}

	return vouchersDto, nil
}

func (v *VoucherService) GetByCode(code string) (models.VoucherResponseBody, error) {
	voucher, err := v.voucherRepository.GetByCode(code)
	if err != nil {
		return models.VoucherResponseBody{}, err
	}

	return models.VoucherResponseBody{
		Name:     voucher.Name,
		Code:     code,
		Start:    voucher.Start,
		End:      voucher.End,
		Discount: voucher.Discount,
		Active:   voucher.Active,
	}, nil
}

func (v *VoucherService) Validate(code string) (bool, string, error) {
	voucher, err := v.voucherRepository.GetByCode(code)
	if err != nil {
		return false, "", err
	}

	pass := true
	detail := ""

	if voucher.Start.After(time.Time(time.Now())) {
		pass = false
		detail = "voucher is not effect"
	} else if voucher.End.Before(time.Time(time.Now())) {
		pass = false
		detail = "voucher has expired"
	} else if !voucher.Active {
		pass = false
		detail = "voucher has already been used"
	}

	return pass, detail, nil
}
