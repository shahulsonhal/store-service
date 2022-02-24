package data

import (
	"fmt"
	"strings"
	"sync"
)

func s() {
	Store = &St{}
	Store.StoreData = sync.Map{}
}

var (
	Store               *St
	ErrResourceNotFound = fmt.Errorf("Not Found")
)

type St struct {
	StoreData sync.Map
}

const (
	countryGermany         = "DE"
	countryFrance          = "FR"
	precipitationLevelHigh = "HIGH"
)

// GetStore lists the store details.
func (f *St) GetStore(max int, country string) ([]StoreDetails, error) {
	val, ok := f.StoreData.Load(country)
	if !ok {
		return nil, ErrResourceNotFound
	}

	store := val.([]StoreDetails)

	if max > len(store) || max == 0 {
		max = len(store)
	}

	store = store[0:max]

	if err := checkServiceType(store); err != nil {
		return nil, err
	}

	return store, nil
}

func checkServiceType(store []StoreDetails) error {
	for i, item := range store {
		switch item.CountryCode {
		case countryGermany:
			weather, err := getAccuweather(countryGermany, item.Location.Lat, item.Location.Lng)
			if err != nil {
				return err
			}

			store[i].SlowService = weather.PrecipitationLevel == precipitationLevelHigh
		case countryFrance:
			weather, err := getAerisweather(strings.ToLower(item.Country), item.Location.Lat, item.Location.Lng)
			if err != nil {
				return err
			}

			// if precipitation is over 40mm  means there is SNOW or ICE precipitation.
			// The service will be slow.
			// We are only considering the value precipitation level for last 24 hours to calculte the slow service.
			// Consideration of Precipitation types may cause ambiguity.
			store[i].SlowService = weather.PrecipitationLast24h > 40
		}
	}

	return nil
}
