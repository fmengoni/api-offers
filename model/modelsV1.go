package model

import (
	"gopkg.in/mgo.v2/bson"
)

// GeoEntityType serves to enumerate different GeoEntity types
type GeoEntityType string

// Different GeoEntity types
const (
	GeoEntityTypeCountry       GeoEntityType = "COUNTRY"
	GeoEntityTypeCity          GeoEntityType = "CITY"
	GeoEntityTypeState         GeoEntityType = "STATE"
	GeoEntityTypeNeighbourhood GeoEntityType = "NEIGHBOURHOOD"
	GeoEntityTypeAirport       GeoEntityType = "AIRPORT"
	GeoEntityTypeAccommodation GeoEntityType = "ACCOMMODATION"
)

// BaseEntity base fields for entity
type BaseEntity struct {
	ID   bson.ObjectId       `json:"id" bson:"_id"`
	Name map[Language]string `json:"name" bson:"name"`
}

// AccommodationGeometry : accommodation with geometry information
type AccommodationGeometry struct {
	Accommodation
	Location Geometry `json:"location"`
}

// Accommodation : representation of accommodation
type Accommodation struct {
	BaseEntity `bson:",inline"`
	CatalogID  string `json:"catalog_id" bson:"catalog_id"`
	InternalID string `json:"internal_id" bson:"internal_id"`
	ProviderID string `json:"provider_id" bson:"provider_id"`
}

// Airport is a basic representation of a Airport
type Airport struct {
	BaseEntity  `bson:",inline"`
	IataCode    string        `json:"iata_code" bson:"iata_code"`
	CityID      bson.ObjectId `json:"city_id,omitempty" bson:"city_id,omitempty"`
	CountryID   bson.ObjectId `json:"country_id,omitempty" bson:"country_id,omitempty"`
	CountryCode string        `json:"country_code" bson:"country_code"`
	StateID     bson.ObjectId `json:"state_id,omitempty" bson:"state_id,omitempty"`
	EanID       int           `json:"-" bson:"ean_id"`
}

// Country is a basic representation of a country
type Country struct {
	BaseEntity   `bson:",inline"`
	OfficialName map[Language]string `json:"official_name,omitempty" bson:"official_name"`
	Alpha2Code   string              `json:"alpha2_code" bson:"alpha2_code"`
	Alpha3Code   string              `json:"alpha3_code" bson:"alpha3_code"`
	GeoID        string              `json:"-" bson:"geo_id"`
}

// City is a basic representation of a city
type City struct {
	BaseEntity `bson:",inline"`
	IataCode   string        `json:"iata_code" bson:"iata_code,omitempty"`
	GeoID      string        `json:"-" bson:"geo_id,omitempty"`
	CountryID  bson.ObjectId `json:"country_id,omitempty" bson:"country_id,omitempty"`
	StateID    bson.ObjectId `json:"state_id,omitempty" bson:"state_id,omitempty"`
}

// Neighbourhood is a basic representation of a neighbourhood
type Neighbourhood struct {
	BaseEntity `bson:",inline"`
	GeoID      string        `json:"-" bson:"geo_id"`
	CountryID  bson.ObjectId `json:"country_id,omitempty" bson:"country_id,omitempty"`
	StateID    bson.ObjectId `json:"state_id,omitempty" bson:"state_id,omitempty"`
	CityID     bson.ObjectId `json:"city_id,omitempty" bson:"city_id,omitempty"`
}

// State is a basic representation of a state
type State struct {
	BaseEntity   `bson:",inline"`
	Timezone     string        `json:"timezone,omitempty" bson:"timezone"`
	GeoID        string        `json:"-" bson:"geo_id"`
	CountryID    bson.ObjectId `json:"country_id,omitempty" bson:"country_id,omitempty"`
	Abbreviation string        `json:"abbreviation" bson:"abbreviation"`
}

// GeoEntity is used to be stored in a MongoDB to make different geospatial queries
type GeoEntity struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Type     GeoEntityType `json:"type" bson:"type"`
	Geometry Geometry      `json:"geometry,omitempty" bson:"geometry"`
}

// NewGeoEntityAccommodation returns a new GeoEntity of type ACCOMMODATION
func NewGeoEntityAccommodation(id bson.ObjectId, geometry Geometry) *GeoEntity {
	return &GeoEntity{
		ID:       id,
		Type:     GeoEntityTypeAccommodation,
		Geometry: geometry,
	}
}

// Entity geo entity
type Entity struct {
	ID        string        `json:"id"`
	IataCode  string        `json:"iata_code"`
	Name      string        `json:"name"`
	Type      GeoEntityType `json:"type"`
	CountryID string        `json:"country_id"`
}
