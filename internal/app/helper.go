package app

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	portEnv                    = "PORT"
	defaultPort                = ":8080"
	useInMemory                = "USE_IN_MEMORY"
	defaultInMemory            = true
	weatherBaseURLENV          = "WEATHER_URL"
	defaultWeatherBaseURLValue = "http://localhost:3000"
)

func mustReadPort() string {
	port := os.Getenv(portEnv)
	if len(port) == 0 {
		log.Printf("PORT config is empty, using the default value %s", defaultPort)
		port = defaultPort
	}

	return fixPrefix(port, ":")
}

func mustReadWeatherBaseURL() string {
	weatherURL := os.Getenv(weatherBaseURLENV)
	if len(weatherURL) == 0 {
		log.Printf("WEATHER_URL config is empty, using the default value %s", defaultWeatherBaseURLValue)
		os.Setenv(weatherBaseURLENV, defaultWeatherBaseURLValue)
		weatherURL = defaultWeatherBaseURLValue
	}

	return weatherURL
}

func fixPrefix(val, prefix string) string {
	if !strings.HasPrefix(val, ":") {
		val = ":" + val
	}

	return val
}

func mustReadInMemoryMode() bool {
	var inMemory bool
	m := os.Getenv(useInMemory)
	if len(m) == 0 {
		log.Printf("USE_IN_MEMORY config is empty, using the default value %t", defaultInMemory)
		inMemory = defaultInMemory
	} else {
		modeVal, err := strconv.ParseBool(m)
		if err != nil {
			inMemory = defaultInMemory
			log.Printf("unsupported value for USE_IN_MEMORY config: %s, using the default value: %t", m, inMemory)
		} else {
			inMemory = modeVal
			log.Printf("USE_IN_MEMORY: %t", inMemory)
		}
	}
	return inMemory
}
