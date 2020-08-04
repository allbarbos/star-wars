package planet

import (
	"log"
	"star-wars/entity"
)

// Service contract
type Service interface {
	Exists(name string) (bool, error)
	Save(planet entity.Planet) (error)
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
		log.Print(err)
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
		log.Print(err)
		return err
	}

	return nil
}
