package swapi

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"star-wars/swapi/adapter"
)

// Service contract
type Service interface {
	GetPlanet(name string) (adapter.Planets, error)
}

type swapi struct{}

// New returns a swapi service instance
func New() Service {
	return &swapi{}
}

func (s swapi) GetPlanet(name string) (adapter.Planets, error) {
	var adapter adapter.Planets

	resp, err := http.Get(os.Getenv("SWAPI_URL") + "/planets/?search=" + url.QueryEscape(name))

	if err != nil {
		log.Print(err)
		return adapter, err
	}

	err = json.NewDecoder(resp.Body).Decode(&adapter)

	if err != nil {
		log.Print(err)
		return adapter, err
	}

	return adapter, nil
}
