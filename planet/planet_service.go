package planet

import (
	"context"
	"star-wars/api/handler"
	"star-wars/entity"
	"star-wars/swapi"
)

// Service contract
type Service interface {
	Exists(ctx context.Context, name string) (bool, error)
	Save(ctx context.Context, planet *entity.Planet) error
	FindAll(ctx context.Context, limit int64, skip int64) (*[]entity.Planet, error)
	FindByName(ctx context.Context, name string) (*entity.Planet, error)
	FindByID(ctx context.Context, id string) (*entity.Planet, error)
	Delete(ctx context.Context, id string) error
}

type srv struct {
	repo  Repository
	swapi swapi.Service
}

// NewService returns a planet service instance
func NewService(r Repository, s swapi.Service) Service {
	return &srv{
		repo:  r,
		swapi: s,
	}
}

// Exists search planet in database
func (s srv) Exists(ctx context.Context, name string) (bool, error) {
	planet := entity.Planet{Name: name}

	if planet.IsEmpty([]string{"Name"}) {
		return false, handler.BadRequest{Message: "name is invalid"}
	}

	_, err := s.repo.FindByName(ctx, planet.Name)

	if err != nil {
		return false, err
	}

	return true, nil
}

// FindByName get planet
func (s srv) FindByName(ctx context.Context, name string) (*entity.Planet, error) {
	planet := &entity.Planet{Name: name}

	if planet.IsEmpty([]string{"Name"}) {
		return nil, handler.BadRequest{Message: "name is invalid"}
	}

	planet, err := s.repo.FindByName(ctx, name)

	if err != nil {
		var newError error
		if err.Error() == "mongo: no documents in result" {
			newError = handler.NotFound{Message: "planet not found"}
		} else {
			newError = handler.InternalServer{Message: err.Error()}
		}
		return nil, newError
	}

	return planet, nil
}

// FindByID get planet
func (s srv) FindByID(ctx context.Context, id string) (*entity.Planet, error) {
	planet := &entity.Planet{ID: id}

	if planet.IsEmpty([]string{"ID"}) {
		return nil, handler.BadRequest{Message: "id is invalid"}
	}

	planet, err := s.repo.FindByID(ctx, id)

	if err != nil {
		var newError error
		if err.Error() == "mongo: no documents in result" {
			newError = handler.NotFound{Message: "planet not found"}
		} else {
			newError = handler.InternalServer{Message: err.Error()}
		}
		return planet, newError
	}

	return planet, nil
}

// Delete planet
func (s srv) Delete(ctx context.Context, id string) error {
	if id == "" {
		return handler.BadRequest{Message: "id is invalid"}
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if err.Error() == "the provided hex string is not a valid ObjectID" {
			return handler.BadRequest{Message: "id is invalid"}
		}
		return handler.InternalServer{Message: err.Error()}
	}
	return nil
}

// FindAll get planets
func (s srv) FindAll(ctx context.Context, limit int64, skip int64) (*[]entity.Planet, error) {
	planets, err := s.repo.FindAll(ctx, limit, skip)
	if err != nil {
		return nil, handler.InternalServer{Message: err.Error()}
	}
	return planets, nil
}

// Save planet
func (s srv) Save(ctx context.Context, planet *entity.Planet) error {
	name := planet.Name
	exists, err := s.Exists(ctx, name)

	if err != nil && err.Error() != "mongo: no documents in result" {
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

	err = s.repo.Save(ctx, planet)

	if err != nil {
		return err
	}

	return nil
}
