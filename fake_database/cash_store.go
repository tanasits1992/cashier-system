package fake_database

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type CashStoreRepository struct{}

func NewCashStoreRepository() ICashStoreRepository {
	return &CashStoreRepository{}
}

type ICashStoreRepository interface {
	Load(filename string) error
	Get() (map[float64]int, error)
	Update(cash map[float64]int) error
}

var cashStoreData = make(map[float64]int)

func (r *CashStoreRepository) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("can not read %s: %v", filename, err)
		return err
	}
	defer file.Close()

	var cash map[string]int
	if err := json.NewDecoder(file).Decode(&cash); err != nil {
		log.Fatalf("can not decode %s: %v", filename, err)
		return err
	}

	for key, val := range cash {
		keyNumber, err := strconv.ParseFloat(key, 64)
		if err != nil {
			return err
		}

		cashStoreData[keyNumber] = val
	}

	return nil
}

func (r *CashStoreRepository) Get() (map[float64]int, error) {
	return cashStoreData, nil
}

func (r *CashStoreRepository) Update(cash map[float64]int) error {
	for key, value := range cash {
		cashStoreData[key] = value
	}
	return nil
}
