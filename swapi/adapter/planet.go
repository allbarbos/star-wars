package adapter

// Planets adapter
type Planets struct {
	Count    int32    `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Planet `json:"results"`
}

// Planet adapter
type Planet struct {
	Name           string   `json:"name"`
	RotationPeriod string   `json:"rotation_period"`
	OrbitalPeriod  string   `json:"orbital_period"`
	Diameter       string   `json:"diameter"`
	Climate        string   `json:"climate"`
	Gravity        string   `json:"resulgravityts"`
	Terrain        string   `json:"terrain"`
	SurfaceWater   string   `json:"surface_water"`
	Population     string   `json:"population"`
	Residents      []string `json:"residents"`
	Films          []string `json:"films"`
	Created        string   `json:"created"`
	Edited         string   `json:"edited"`
	URL            string   `json:"url"`
}
