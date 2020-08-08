package planet

import (
	"star-wars/api/handler"
	"star-wars/entity"
	"star-wars/swapi"
)

// Service contract
type Service interface {
	Exists(name string) (bool, error)
	Save(planet *entity.Planet) error
	FindAll(limit int64, skip int64) ([]entity.Planet, error)
	FindByName(name string) (entity.Planet, error)
	FindByID(id string) (entity.Planet, error)
	Delete(id string) error
}

type srv struct {
	repo Repository
	swapi swapi.Service
}

// NewService returns a planet service instance
func NewService(r Repository, s swapi.Service) Service {
	return &srv{
		repo: r,
		swapi: s,
	}
}

// Exists search planet in database
func (s srv) Exists(name string) (bool, error) {
	planetDb, err := s.repo.FindByName(name)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return false, err
	}

	if planetDb.IsEmpty([]string{"ID"}) {
		return false, nil
	}

	return true, nil
}

// FindByName get planet
func (s srv) FindByName(name string) (entity.Planet, error) {
	var planet entity.Planet

	if name == "" {
		return planet, handler.BadRequest{Message: "name is invalid"}
	}

	planet, err := s.repo.FindByName(name)

	if err != nil {
		var newError error
		if err.Error() == "mongo: no documents in result" {
			newError = handler.NotFound{ Message: "planet not found" }
		} else {
			newError = handler.InternalServer{ Message: err.Error() }
		}
		return planet, newError
	}

	return planet, nil
}

// FindByID get planet
func (s srv) FindByID(id string) (entity.Planet, error) {
	var planet entity.Planet

	if id == "" {
		return planet, handler.BadRequest{Message: "id is invalid"}
	}

	planet, err := s.repo.FindByID(id)

	if err != nil {
		var newError error
		if err.Error() == "mongo: no documents in result" {
			newError = handler.NotFound{ Message: "planet not found" }
		} else {
			newError = handler.InternalServer{ Message: err.Error() }
		}
		return planet, newError
	}

	return planet, nil
}

// Delete planet
func (s srv) Delete(id string) error {
	if id == "" {
		return handler.BadRequest{ Message: "id is invalid" }
	}

	if err := s.repo.Delete(id); err != nil {
		if err.Error() == "the provided hex string is not a valid ObjectID" {
			return handler.BadRequest{ Message: "id is invalid" }
		}
		return handler.InternalServer{ Message: err.Error() }
	}
	return nil
}

// FindAll get planets
func (s srv) FindAll(limit int64, skip int64) ([]entity.Planet, error) {
	planets, err := s.repo.FindAll(limit, skip)
	if err != nil {
		return planets, handler.InternalServer{ Message: err.Error() }
	}
	return planets, nil
}

// Save planet
func (s srv) Save(planet *entity.Planet) error {
	name := planet.Name
	exists, err :=  s.Exists(name)

	if err != nil {
		return err
	}

	if exists {
		err := handler.BadRequest{Message: "planet already registered"}
		return err
	}

	adapter, err := s.swapi.GetPlanet(planet.Name)

	if err != nil {
		return handler.InternalServer{Message: err.Error()}
	}

	if adapter.Count == 0 {
		return handler.BadRequest{Message: "non-existent planet"}
	}

	total, err := planet.TotalAppearances(adapter.Results)
	if err != nil {
		return err
	}

	planet.TotalFilms = total

	err = s.repo.Save(planet)

	if err != nil {
		return err
	}

	return nil
}
