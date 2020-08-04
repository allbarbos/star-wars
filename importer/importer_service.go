package importer

import (
	"log"
	"star-wars/entity"
	"star-wars/planet"
	"star-wars/swapi"
)

type Service interface {
	Process(planet entity.Planet) error
}

type service struct {
	planetSrv planet.Service
	swapiSrv swapi.Service
}

func NewImporter(s planet.Service, swapi swapi.Service) Service {
	return &service{
		planetSrv: s,
		swapiSrv: swapi,
	}
}

func (i service) Process(planet entity.Planet) error {
	exists, err := i.planetSrv.Exists(planet.Name)

	if err != nil {
		log.Print(err)
		return err
	}

	if exists {
		log.Print("planet already registered")
		return nil
	}

	adapter, err := i.swapiSrv.GetPlanetExternally(planet.Name)

	if err != nil {
		log.Print(err)
		return err
	}

	total, err := planet.TotalAppearances(adapter.Results)
	if err != nil {
		log.Print(err)
		return err
	}

	planet.TotalFilms = total

	err = i.planetSrv.Save(planet)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
