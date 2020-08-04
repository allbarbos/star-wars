package entity

import (
	"errors"
	"reflect"
	"star-wars/swapi/adapter"
)

// Planet entity
type Planet struct {
	ID 				 string `json:"id" bson:"_id,omitempty"`
	Name 			 string `json:"name" bson:"name,omitempty"`
	Climate 	 string `json:"climate" bson:"climate,omitempty"`
	Terrain 	 string `json:"terrain" bson:"terrain,omitempty"`
	TotalFilms int `json:"totalFilms" bson:"totalFilms,omitempty"`
}

// IsEmpty validate fields
func (p Planet) IsEmpty(fields []string) bool {
	for _, key := range fields {
		value := reflect.ValueOf(p).FieldByName(key)
		if value.Interface() == "" || value.Interface() == 0  {
			return true
		}
	}

	return false
}

// TotalAppearances counts film appearances
func (p Planet) TotalAppearances(adapter []adapter.Planet) (int, error) {
	if len(adapter) != 1 {
		return 0, errors.New("search did not return the planet")
	}

	return len(adapter[0].Films), nil
}
