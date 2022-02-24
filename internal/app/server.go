package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/shahulsonhal/store-service/internal/data"
)

// Config holds the attributes of the Server
type Config struct {
	Name,
	Port,
	WeatherURL string
	UseInMemory bool
}

func newConfig() *Config {
	c := new(Config)

	c.Name = "store-location"
	c.Port = mustReadPort()
	c.UseInMemory = mustReadInMemoryMode()
	c.WeatherURL = mustReadWeatherBaseURL()

	return c
}

// Server is the top level store location server application object.
type Server struct {
	// TODO : db repo
	*Config

	repo data.Repo
}

// NewServer creates server object
func NewServer() *Server {
	s := new(Server)

	s.Config = newConfig()
	s.repo = data.NewRepo(s.UseInMemory)

	if s.UseInMemory {
		s.repo = readStoreData()
	}

	return s
}

func (s *Server) Start() {
	server := http.Server{
		Addr:    s.Port,
		Handler: s.InitRouter(),
	}

	go func() {
		log.Printf("listening on port %s...\n", server.Addr)
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal("failed to start server: ", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

// Country wise store data map
var record = map[string]string{
	"DE": "store_data_de.json",
	"FR": "store_data_fr.json",
}

func readStoreData() *data.St {
	var (
		st []data.StoreDetails
		fx data.St
	)

	for k, v := range record {
		st = nil
		file, err := ioutil.ReadFile(v)
		if err != nil {
			log.Fatalf("Store data read error: %v", err)
		}

		err = json.Unmarshal([]byte(file), &st)
		if err != nil {
			log.Fatalf("Store data read error: %v", err)
		}

		fx.StoreData.Store(k, st)
	}

	return &fx
}
