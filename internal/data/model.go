package data

type LocationData struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}

type StoreDetails struct {
	StoreID     string       `json:"store_id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Country     string       `json:"country,omitempty"`
	CountryCode string       `json:"country_code,omitempty"`
	Location    LocationData `json:"location,omitempty"`
	SlowService bool         `json:"slow_service,omitempty"`
}

type Accuweather struct {
	Temp               float64
	PrecipitationLevel string
	PrecipitationType  string
}

type Aerisweather struct {
	Temp                 float64
	PrecipitationLast24h float64
	Precipitation        []string
}
