package fake_database

import (
	"cashier/app/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type VoucherRepository struct{}

func NewVoucherRepository() IVoucherRepository {
	return &VoucherRepository{}
}

type IVoucherRepository interface {
	Load(filename string) error
	List() (map[string]models.VoucherCollection, error)
	GetByCode(code string) (models.VoucherCollection, error)
	Insert(voucher models.VoucherCollection) (string, error)
	Update(code string, voucher models.VoucherCollection) error
}

var voucherData = make(map[string]models.VoucherCollection)

func (v *VoucherRepository) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("can not read voucher.json: %v", err)
		return err
	}
	defer file.Close()

	var voucher map[string]models.VoucherCollection
	if err := json.NewDecoder(file).Decode(&voucher); err != nil {
		log.Fatalf("can not decode voucher.json: %v", err)
		return err
	}

	voucherData = voucher
	return nil
}

func (v *VoucherRepository) List() (map[string]models.VoucherCollection, error) {
	return voucherData, nil
}

func (v *VoucherRepository) GetByCode(code string) (models.VoucherCollection, error) {
	voucher, ok := voucherData[code]
	if !ok {
		return models.VoucherCollection{}, errors.New("code is invalid")
	}

	return voucher, nil
}

func (v *VoucherRepository) Insert(voucher models.VoucherCollection) (string, error) {
	code := uuid.New().String()
	voucherData[code] = models.VoucherCollection{
		Name:     voucher.Name,
		Discount: voucher.Discount,
		Start:    voucher.Start,
		End:      voucher.End,
		Active:   true,
	}
	return code, nil
}

func (v *VoucherRepository) Update(code string, voucher models.VoucherCollection) error {
	_, err := v.GetByCode(code)
	if err != nil {
		fmt.Println("code is invalid")
		return errors.New("code is invalid")
	}

	voucherData[code] = voucher

	return nil
}
