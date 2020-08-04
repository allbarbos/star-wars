package entity

import (
	"errors"
	"star-wars/entity/adapter"
)

// Planet entity
type Planet struct {
	ID 				 string `json:"id" bson:"_id,omitempty"`
	Name 			 string `json:"name" bson:"name,omitempty"`
	Climate 	 string `json:"climate" bson:"climate,omitempty"`
	Terrain 	 string `json:"terrain" bson:"terrain,omitempty"`
	TotalFilms int `json:"totalFilms" bson:"totalFilms,omitempty"`
}

func (p Planet) Valid() bool {
	if p.ID == "" {
		return false
	}

	if p.Name == "" {
		return false
	}

	if p.Climate == "" {
		return false
	}

	if p.Terrain == "" {
		return false
	}

	return true
}

func (p Planet) TotalAppearances(adapter []adapter.Planet) (int, error) {
	if len(adapter) != 1 {
		return 0, errors.New("search did not return the planet")
	}

	return len(adapter[0].Films), nil
}
