package importer

import (
	"context"
	"star-wars/entity"
	"star-wars/planet"
	"star-wars/swapi"
	"sync"
)

var wg sync.WaitGroup
var mutex = &sync.Mutex{}

// Service contract
type Service interface {
	Import(ctx context.Context, planets []entity.Planet) []error
}

type service struct {
	planetSrv planet.Service
	swapiSrv  swapi.Service
}

// NewImporter returns a importer service instance
func NewImporter(s planet.Service, swapi swapi.Service) Service {
	return &service{
		planetSrv: s,
		swapiSrv:  swapi,
	}
}

func saveTask(ctx context.Context, planet entity.Planet, srv planet.Service, errs *[]error) {
	defer wg.Done()
	err := srv.Save(ctx, &planet)
	if err != nil {
		mutex.Lock()
		*errs = append(*errs, err)
		mutex.Unlock()
	}
}

// Import data import
func (i service) Import(ctx context.Context, planets []entity.Planet) []error {
	var errs []error

	for _, planet := range planets {
		wg.Add(1)
		go saveTask(ctx, planet, i.planetSrv, &errs)
	}
	wg.Wait()

	return errs
}
