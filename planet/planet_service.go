package planet

import (
	"star-wars/api/handler"
	"star-wars/entity"
)

// Service contract
type Service interface {
	Exists(name string) (bool, error)
	Save(planet entity.Planet) (error)
	FindByName(name string) (entity.Planet, error)
	FindByID(id string) (entity.Planet, error)
}

type srv struct {
	repo Repository
}

// NewService returns a planet service instance
func NewService(r Repository) Service {
	return &srv{
		repo: r,
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

// Save planet
func (s srv) Save(planet entity.Planet) (error) {
	err := s.repo.Save(planet)

	if err != nil {
		return err
	}

	return nil
}

// FindByName get planet
func (s srv) FindByName(name string) (entity.Planet, error) {
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

// FindByName get planet
func (s srv) FindByID(id string) (entity.Planet, error) {
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
