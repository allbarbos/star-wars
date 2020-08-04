package swapi

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"star-wars/entity/adapter"
)

type Service interface {
	GetPlanetExternally(name string) (adapter.Planets, error)
}

type swapi struct {}

func New() Service {
	return &swapi{}
}

func (s swapi) GetPlanetExternally(name string) (adapter.Planets, error) {
	var adapter adapter.Planets

	resp, err := http.Get(os.Getenv("SWAPI_URL") + "/planets/?search=" + name)

	if err != nil {
		log.Print(err)
		return adapter, err
	}

	json.NewDecoder(resp.Body).Decode(&adapter)
	return adapter, nil
}
