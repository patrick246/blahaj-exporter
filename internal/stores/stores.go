package stores

import (
	_ "embed"
	"encoding/json"
)

//go:embed stores.json
var stores []byte

//nolint:gochecknoglobals // cache for the also global embed variable
var storeCache []Store

type Store struct {
	ID          string    `json:"buCode"`
	Name        string    `json:"name"`
	Coordinates []float64 `json:"coordinates"`
	CountryCode string    `json:"countryCode"`
}

func GetStores() ([]Store, error) {
	if storeCache != nil {
		return storeCache, nil
	}

	err := json.Unmarshal(stores, &storeCache)
	return storeCache, err
}

func GetCountryCodes() ([]string, error) {
	stores, err := GetStores()
	if err != nil {
		return nil, err
	}

	codes := map[string]struct{}{}
	for _, store := range stores {
		codes[store.CountryCode] = struct{}{}
	}

	countryCodes := make([]string, 0, len(codes))
	for code := range codes {
		countryCodes = append(countryCodes, code)
	}

	return countryCodes, nil
}
