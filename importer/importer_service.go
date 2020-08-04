package importer

import (
	"errors"
	"fmt"
	"star-wars/entity"
	"star-wars/planet"
	"star-wars/swapi"
)

type Service interface {
	Process(planet entity.Planet, errchan chan<- string) error
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

func (i service) Process(planet entity.Planet, errchan chan<- string) error {
	defer close(errchan)
	exists, err := i.planetSrv.Exists(planet.Name)

	if err != nil {
		errchan <- fmt.Sprintf("%s: %s", planet.Name, err.Error())
		return err
	}

	if exists {
		err := errors.New("planet already registered")
		errchan <- fmt.Sprintf("%s: %s", planet.Name, err.Error())
		return err
	}

	adapter, err := i.swapiSrv.GetPlanetExternally(planet.Name)

	if err != nil {
		errchan <- fmt.Sprintf("%s: %s", planet.Name, err.Error())
		return err
	}

	total, err := planet.TotalAppearances(adapter.Results)
	if err != nil {
		errchan <- fmt.Sprintf("%s: %s", planet.Name, err.Error())
		return err
	}

	planet.TotalFilms = total

	err = i.planetSrv.Save(planet)

	if err != nil {
		errchan <- fmt.Sprintf("%s: %s", planet.Name, err.Error())
		return err
	}

	return nil
}
