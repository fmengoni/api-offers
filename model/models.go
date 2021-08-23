package model

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
)

// GeometryType serves to enumerate different geometry types
// For more info: http://geojson.org/
type GeometryType string

// Geometry types from GeoJSON used in this app
// For more info: http://geojson.org/
const (
	GeometryPoint        GeometryType = "Point"
	GeometryMultiPolygon GeometryType = "MultiPolygon"
	GeometryPolygon      GeometryType = "Polygon"
	GeometryMultiPoint   GeometryType = "MultiPoint"
)

// Geometry represents a GeoJSON Geometry
// For more info: http://geojson.org/
type Geometry struct {
	Type         GeometryType    `json:"type" bson:"type"`
	Coordinates  []interface{}   `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
	MultiPolygon [][][][]float64 `json:"multipolygon,omitempty" bson:"-"`
	Polygon      [][][]float64   `json:"polygon,omitempty" bson:"-"`
	Point        []float64       `json:"point,omitempty" bson:"-"`
	MultiPoint   [][]float64     `json:"multipoint,omitempty" bson:"-"`
}

// MarshalJSON converts the geometry object into the correct JSON.
// This fulfills the json.Marshaler interface.
func (g *Geometry) MarshalJSON() ([]byte, error) {
	// defining a struct here lets us define the order of the JSON elements.
	type geometry struct {
		Type         GeometryType    `json:"type"`
		MultiPolygon [][][][]float64 `json:"multipolygon,omitempty" bson:"-"`
		Polygon      [][][]float64   `json:"polygon,omitempty" bson:"-"`
		Point        []float64       `json:"point,omitempty" bson:"-"`
		MultiPoint   [][]float64     `json:"multipoint,omitempty" bson:"-"`
		Coordinates  []interface{}   `json:"coordinates,omitempty"`
	}

	geo := &geometry{
		Type:        g.Type,
		Coordinates: g.Coordinates,
	}

	switch g.Type {
	case GeometryPoint:
		geo.Point = encodePoint(g.Coordinates)
	case GeometryPolygon:
		geo.Polygon = encodePolygon(g.Coordinates)
	case GeometryMultiPoint:
		geo.MultiPoint = encodeMultiPoint(g.Coordinates)
	case GeometryMultiPolygon:
		geo.MultiPolygon = encodeMultiPolygon(g.Coordinates)
	}

	return json.Marshal(geo)
}

func (g *GeoRegion) MapCoordinates() {
	switch g.Geometry.Type {
	case GeometryPoint:
		g.Geometry.Point = encodePoint(g.Geometry.Coordinates)
	case GeometryPolygon:
		g.Geometry.Polygon = encodePolygon(g.Geometry.Coordinates)
	case GeometryMultiPoint:
		g.Geometry.MultiPoint = encodeMultiPoint(g.Geometry.Coordinates)
	case GeometryMultiPolygon:
		g.Geometry.MultiPolygon = encodeMultiPolygon(g.Geometry.Coordinates)
	}
}

func encodePoint(co interface{}) []float64 {
	pt := make([]float64, 0, len(co.([]interface{})))

	for _, v := range co.([]interface{}) {
		pt = append(pt, v.(float64))
	}

	return pt
}

func encodeMultiPoint(co []interface{}) [][]float64 {
	pt := make([][]float64, 0, len(co))

	for _, c := range co {
		p := encodePoint(c)

		pt = append(pt, p)
	}

	return pt
}

func encodePolygon(co []interface{}) [][][]float64 {
	pt := make([][][]float64, 0, len(co))

	for _, c := range co {
		p := encodeMultiPoint(c.([]interface{}))

		pt = append(pt, p)
	}

	return pt
}

func encodeMultiPolygon(co []interface{}) [][][][]float64 {
	pt := make([][][][]float64, 0, len(co))

	for _, c := range co {
		p := encodePolygon(c.([]interface{}))

		pt = append(pt, p)
	}

	return pt
}

// NewMultiPointGeometry creates and initializes a multipoint geometry with the give coordinate.
func NewMultiPointGeometry(coordinate [][]float64) *Geometry {
	return &Geometry{
		Type:       GeometryMultiPoint,
		MultiPoint: coordinate,
	}
}

// NewPointGeometry creates and initializes a point geometry with the give coordinate.
func NewPointGeometry(coordinate []interface{}) *Geometry {
	return &Geometry{
		Type:        GeometryPoint,
		Coordinates: coordinate,
	}
}

// NewMultiPolygonGeometry creates and initializes a multi-polygon geometry with the given polygons.
func NewMultiPolygonGeometry(polygons ...[][][]float64) *Geometry {
	return &Geometry{
		Type:         GeometryMultiPolygon,
		MultiPolygon: polygons,
	}
}

// UnmarshalJSON decodes the data into a GeoJSON geometry.
// This fulfills the json.Unmarshaler interface.
/*func (g *Geometry) UnmarshalJSON(data []byte) error {
	var object map[string]interface{}

	err := json.Unmarshal(data, &object)
	if err != nil {
		return fmt.Errorf("failed to unmarshal geometry. %w", err)
	}

	return decodeGeometry(g, object)
}


// GetBSON implements bson.Getter.
func (g Geometry) GetBSON() (interface{}, error) {
	type geometry struct {
		Type        string      `bson:"type"`
		Coordinates interface{} `bson:"coordinates"`
	}

	geo := &geometry{
		Type: string(g.Type),
	}

	switch g.Type {
	case GeometryPoint:
		geo.Coordinates = g.Point
	case GeometryMultiPolygon:
		geo.Coordinates = g.MultiPolygon
	case GeometryMultiPoint:
		geo.Coordinates = g.MultiPoint
	}

	return geo, nil
}
*/
// SetBSON implements bson.Setter.
/*func (g *Geometry) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		Type        GeometryType  `bson:"type"`
		Coordinates []interface{} `bson:"coordinates"`
	})

	err := raw.Unmarshal(decoded)

	if err == nil {
		g.Type = decoded.Type
		g.Coordinates = decoded.Coordinates

		switch g.Type {
		case GeometryPoint:
			g.Point, err = decodePosition(decoded.Coordinates)
		case GeometryMultiPolygon:
			g.MultiPolygon, err = decodePolygonSet(decoded.Coordinates)
		case GeometryMultiPoint:
			g.MultiPoint, err = decodePositionSet(decoded.Coordinates)
		}
	}

	return err
}

func decodeGeometry(g *Geometry, object map[string]interface{}) error {
	t, ok := object["type"]
	if !ok {
		return errors.New("type property not defined")
	}

	if s, ok := t.(string); ok {
		g.Type = GeometryType(s)
	} else {
		return errors.New("type property not string")
	}

	var err error

	switch g.Type {
	case GeometryPoint:
		g.Point, err = decodePosition(object["coordinates"])
	case GeometryMultiPolygon:
		g.MultiPolygon, err = decodePolygonSet(object["coordinates"])
	case GeometryMultiPoint:
		g.MultiPoint, err = decodePositionSet(object["coordinates"])
	}

	return err
}

func decodePosition(data interface{}) ([]float64, error) {
	coords, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid position, got %v", data)
	}

	result := make([]float64, 0, len(coords))

	for _, coord := range coords {
		if f, ok := coord.(float64); ok {
			result = append(result, f)
		} else {
			return nil, fmt.Errorf("not a valid coordinate, got %v", coord)
		}
	}

	return result, nil
}

func decodePolygonSet(data interface{}) ([][][][]float64, error) {
	polygons, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid polygon, got %v", data)
	}

	result := make([][][][]float64, 0, len(polygons))

	for _, polygon := range polygons {
		if p, err := decodePathSet(polygon); err == nil {
			result = append(result, p)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodePathSet(data interface{}) ([][][]float64, error) {
	sets, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid path, got %v", data)
	}

	result := make([][][]float64, 0, len(sets))

	for _, set := range sets {
		if s, err := decodePositionSet(set); err == nil {
			result = append(result, s)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodePositionSet(data interface{}) ([][]float64, error) {
	points, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid set of positions, got %v", data)
	}

	result := make([][]float64, 0, len(points))

	for _, point := range points {
		if p, err := decodePosition(point); err == nil {
			result = append(result, p)
		} else {
			return nil, err
		}
	}

	return result, nil
}
*/
// Language is a simple representation of language in ISO 639-1
type Language string

func (l Language) String() string {
	return string(l)
}

// RegionType is the type of entity
type RegionType string

// types of entities
const (
	RegionTypeCity              RegionType = "city"
	RegionTypeCountry           RegionType = "country"
	RegionTypeContinent         RegionType = "continent"
	RegionTypeHighLevelRegion   RegionType = "high_level_region"
	RegionTypeMetroStation      RegionType = "metro_station"
	RegionTypeProvinceState     RegionType = "province_state"
	RegionTypeMultiCityVicinity RegionType = "multi_city_vicinity"
	RegionTypePOI               RegionType = "point_of_interest"
	RegionTypeNeighborhood      RegionType = "neighborhood"
	RegionTypeTrainStation      RegionType = "train_station"
	RegionTypeAccommodation     RegionType = "accommodation"
	RegionTypeGeo               RegionType = "geo_coordinates"
)

// BaseRegion is a basic region data
type BaseRegion struct {
	ID    bson.ObjectId `json:"_id" bson:"_id"`
	GeoID string        `json:"id" bson:"geo_id"`
	Type  RegionType    `json:"type" bson:"type"`
}

// Region is the new version of geo entity
type Region struct {
	BaseRegion  `bson:",inline"`
	Name        map[Language]string `json:"name" bson:"name"`
	CountryCode string              `json:"country_code,omitempty" bson:"country_code,omitempty"`
	Center      Center              `json:"center" bson:"coordinates"`
	Ancestors   []Ancestor          `json:"ancestors" bson:"ancestors"`
	Descendants Descendants         `json:"descendants" bson:"descendants"`
}

// GeoRegion is a baseRegion + a geometry
type GeoRegion struct {
	BaseRegion `bson:",inline"`
	Geometry   Geometry `json:"geometry" bson:"bounding_polygon"`
}

// Ancestor reprecents a container/father region
type Ancestor struct {
	ID   string     `json:"id" bson:"geo_id"`
	Type RegionType `json:"type" bson:"type"`
}

// Descendants represents a child region
type Descendants struct {
	Cities              []string `json:"cities,omitempty" bson:"city,omitempty"`
	Countries           []string `json:"countries,omitempty" bson:"country,omitempty"`
	POIs                []string `json:"points_of_interest,omitempty" bson:"point_of_interest,omitempty"`
	HighLevelRegions    []string `json:"high_level_regions,omitempty" bson:"high_level_region,omitempty"`
	TrainStations       []string `json:"train_stations,omitempty" bson:"train_station,omitempty"`
	MetroStations       []string `json:"metro_stations,omitempty" bson:"metro_station,omitempty"`
	Neighbourhoods      []string `json:"neighborhoods,omitempty" bson:"neighborhood,omitempty"`
	MultiCityVicinities []string `json:"multi_city_vicinities,omitempty" bson:"multi_city_vicinity,omitempty"`
	ProvinceStates      []string `json:"province_states,omitempty" bson:"province_state,omitempty"`
	Accommodations      []string `json:"accommodations,omitempty" bson:"accommodation,omitempty"`
}

// AirportV2 is the new version of airports
type AirportV2 struct {
	ID          bson.ObjectId       `json:"id" bson:"_id"`
	IataCode    string              `json:"iata_code" bson:"iata"`
	Name        map[Language]string `json:"name" bson:"fullname"`
	CountryCode string              `json:"country_code" bson:"countrycode"`
	Coordinates Coordinates         `json:"coordinates" bson:"coordinates"`
	Region      AirportRegion       `json:"region" bson:"region"`
}

// AirportRegion represents a region of an airport
type AirportRegion struct {
	ID   string              `json:"id" bson:"id"`
	Type string              `json:"type" bson:"regiontype"`
	Name map[Language]string `json:"name" bson:"fullname"`
}

// Coordinates represents a coordinates in degrees
type Coordinates struct {
	Longitude float64 `json:"longitude" bson:"longitude"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
}

// Center represents the center coordinates of a region
type Center struct {
	Longitude float64 `json:"longitude" bson:"center_longitude"`
	Latitude  float64 `json:"latitude" bson:"center_latitude"`
}
