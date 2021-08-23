package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/basset-la/api-geo/conf"
	pkgErrors "github.com/basset-la/api-geo/errors"
	"github.com/basset-la/api-geo/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Repository is a basic interface to a database
type V1 interface {
	FindCountryByID(id string) (*model.Country, error)
	FindCountries(q CountryQuery) ([]model.Country, error)
	FindCityByID(id string) (*model.City, error)
	UpdateCityByID(id string, city model.City) error
	FindCities(q CityQuery) ([]model.City, error)
	FindStateByID(id string) (*model.State, error)
	FindStates(q StateQuery) ([]model.State, error)
	FindNeighbourhoodByID(id string) (*model.Neighbourhood, error)
	FindNeighbourhoods(q NeighbourhoodQuery) ([]model.Neighbourhood, error)
	UpdateNeighbourhoodByID(id string, neighbourhood model.Neighbourhood) error
	FindAirportByID(id string) (*model.Airport, error)
	UpdateAirportByID(id string, airport model.Airport) error
	FindAirports(q AirportQuery) ([]model.Airport, error)
	FindGeoEntityByID(id string) (*model.GeoEntity, error)
	FindGeoEntities(query bson.M, resultsPerPage int, page int) ([]model.GeoEntity, error)
	InsertGeoEntity(geometry model.Geometry, id string) error
	IntersectsGeoEntities(geometry model.Geometry, excludedEntityTypes []model.GeoEntityType) ([]model.GeoEntity, error)
	InsertAccommodation(accommodation model.Accommodation) (*string, error)
	FindAccommodationByID(id string) (*model.Accommodation, error)
	FindAccommodations(q AccommodationQuery) ([]model.Accommodation, error)
}

// MongoRepositoryV1 handles all requests to MongoDB
type MongoRepositoryV1 struct {
	Session *mgo.Session
}

func (r *MongoRepositoryV1) Close() {
	r.Session.Close()
}

// NewMongoRepositoryV1 creates a new mongo repository
func NewMongoRepositoryV1(mongoURL string) (*MongoRepositoryV1, error) {
	s, err := mgo.Dial(mongoURL)

	if err != nil {
		return nil, fmt.Errorf("error creating mongo repository %w", err)
	}

	r := &MongoRepositoryV1{
		Session: s,
	}

	return r, nil
}

// AirportQuery is used as query to get airports
type AirportQuery struct {
	CountryID string
	StateID   string
	CityID    string
	IataCode  string
}

// CityQuery is used as query to get cities
type CityQuery struct {
	CountryID string
	StateID   string
	IataCode  string
}

// CountryQuery is used as query to get countries
type CountryQuery struct {
	Alpha2Code string
	Alpha3Code string
}

// StateQuery is used as query to get states
type StateQuery struct {
	CountryID string
}

// NeighbourhoodQuery is used as query to get neighbourhoods
type NeighbourhoodQuery struct {
	CountryID string
	StateID   string
	CityID    string
}

// AccommodationQuery is used as query to get accommodations
type AccommodationQuery struct {
}

func newIntersectsQuery(geometry model.Geometry, excludedEntityTypes []model.GeoEntityType) bson.M {
	query := bson.M{}
	query["geometry"] = map[string]map[string]model.Geometry{"$geoIntersects": {"$geometry": geometry}}
	query["type"] = map[string][]model.GeoEntityType{"$nin": excludedEntityTypes}

	return query
}

// InsertGeoEntity : insert geo entity with geometry
func (r *MongoRepositoryV1) InsertGeoEntity(geometry model.Geometry, id string) error {
	geoEntity := model.NewGeoEntityAccommodation(bson.ObjectIdHex(id), geometry)
	err := r.insertEntity(geoEntity, conf.GetProps().Mongo.GeoEntitiesTable)

	if err != nil {
		return fmt.Errorf("failed to save geo entity")
	}

	return nil
}

// FindAccommodations searches accommodations by query
func (r *MongoRepositoryV1) FindAccommodations(q AccommodationQuery) ([]model.Accommodation, error) {
	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(conf.GetProps().Mongo.AccommodationTable)
	accommodation := make([]model.Accommodation, 0)

	query := bson.M{}

	err := col.Find(query).All(&accommodation)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("error finding accommodations %w", err)
	}

	return accommodation, nil
}

// FindAccommodationByID : find accommodation by id
func (r *MongoRepositoryV1) FindAccommodationByID(id string) (*model.Accommodation, error) {
	accommodation := &model.Accommodation{}
	err := r.findEntityByID(id, conf.GetProps().Mongo.AccommodationTable, accommodation)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("error finding accommodation %w", err)
	}

	return accommodation, nil
}

// InsertAccommodation : insert new accommodation and return id or error
func (r *MongoRepositoryV1) InsertAccommodationV1(accommodation model.Accommodation) (*string, error) {
	id := bson.NewObjectId().Hex()
	accommodation.ID = bson.ObjectIdHex(id)
	err := r.insertEntity(accommodation, conf.GetProps().Mongo.AccommodationTable)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("error saving accommodation %w", err)
	}

	return &id, nil
}

// UpdateCityByID update a city by id
func (r *MongoRepositoryV1) UpdateCityByID(id string, city model.City) error {
	err := r.updateEntityByID(id, city, conf.GetProps().Mongo.CitiesTable)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return pkgErrors.ErrEntityNotFound
		}

		return fmt.Errorf("failed to save update city %w", err)
	}

	return nil
}

// FindGeoEntityByID searches a geoEntity by id
func (r *MongoRepositoryV1) FindGeoEntityByID(id string) (*model.GeoEntity, error) {
	geoEntity := &model.GeoEntity{}
	err := r.findEntityByID(id, conf.GetProps().Mongo.GeoEntitiesTable, geoEntity)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get entity %w", err)
	}

	return geoEntity, nil
}

// FindGeoEntities searches geoEntities by query
func (r *MongoRepositoryV1) FindGeoEntities(query bson.M, resultsPerPage int, page int) ([]model.GeoEntity, error) {
	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(conf.GetProps().Mongo.GeoEntitiesTable)
	geoEntities := make([]model.GeoEntity, 0)

	err := col.Find(query).Sort("_id").Skip(page * resultsPerPage).Limit(resultsPerPage).All(&geoEntities)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get entity %w", err)
	}

	return geoEntities, nil
}

// IntersectsGeoEntities intersects GeoEntities by geometry
func (r *MongoRepositoryV1) IntersectsGeoEntities(geometry model.Geometry, excludedEntityTypes []model.GeoEntityType) ([]model.GeoEntity, error) {
	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(conf.GetProps().Mongo.GeoEntitiesTable)
	geoEntities := make([]model.GeoEntity, 0)

	q := newIntersectsQuery(geometry, excludedEntityTypes)
	err := col.Find(q).Select(bson.M{"type": 1}).All(&geoEntities) // Just retrieve fields _id and type

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to intersect entities %w", err)
	}

	return geoEntities, nil
}

// FindCityByID searches a city by id
func (r *MongoRepositoryV1) FindCityByID(id string) (*model.City, error) {
	city := &model.City{}
	err := r.findEntityByID(id, conf.GetProps().Mongo.CitiesTable, city)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get city %w", err)
	}

	return city, nil
}

// FindCountryByID searches a city by id
func (r *MongoRepositoryV1) FindCountryByID(id string) (*model.Country, error) {
	country := &model.Country{}
	err := r.findEntityByID(id, conf.GetProps().Mongo.CountryTable, country)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get country %w", err)
	}

	return country, nil
}

func (r *MongoRepositoryV1) findEntityByID(id string, collection string, target interface{}) error {
	err := validateObjectID(id)

	if err != nil {
		return fmt.Errorf("id is not in the correct format")
	}

	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(collection)
	err = col.FindId(bson.ObjectIdHex(id)).One(target)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return pkgErrors.ErrEntityNotFound
		}

		return fmt.Errorf("failed to get entity %w", err)
	}

	return nil
}

func (r *MongoRepositoryV1) updateEntityByID(id string, entity interface{}, collection string) error {
	err := validateObjectID(id)
	if err != nil {
		return err
	}

	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(collection)
	err = col.UpdateId(bson.ObjectIdHex(id), entity)

	return fmt.Errorf("error updating entity %w", err)
}

func (r *MongoRepositoryV1) insertEntity(entity interface{}, collection string) error {
	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(collection)
	err := col.Insert(entity)

	return fmt.Errorf("error inserting entity %w", err)
}

// FindCities searches cities by query
func (r *MongoRepositoryV1) FindCities(q CityQuery) ([]model.City, error) {
	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(conf.GetProps().Mongo.CitiesTable)
	cities := make([]model.City, 0)

	query := bson.M{}

	if q.CountryID != "" {
		query["country_id"] = bson.ObjectIdHex(q.CountryID)
	}

	if q.StateID != "" {
		query["state_id"] = bson.ObjectIdHex(q.StateID)
	}

	if q.IataCode != "" {
		iataCodes := strings.Split(q.IataCode, ",")
		query["iata_code"] = bson.M{"$in": iataCodes}
	}

	err := col.Find(query).All(&cities)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get cities %w", err)
	}

	return cities, nil
}

// FindCountries searches countries by query
func (r *MongoRepositoryV1) FindCountries(q CountryQuery) ([]model.Country, error) {
	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(conf.GetProps().Mongo.CountryTable)
	countries := make([]model.Country, 0)

	query := bson.M{}

	if q.Alpha2Code != "" {
		query["alpha2_code"] = q.Alpha2Code
	}

	if q.Alpha3Code != "" {
		query["alpha3_code"] = q.Alpha3Code
	}

	err := col.Find(query).All(&countries)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get countries %w", err)
	}

	return countries, nil
}

// FindStateByID searches a state by id
func (r *MongoRepositoryV1) FindStateByID(id string) (*model.State, error) {
	state := &model.State{}
	err := r.findEntityByID(id, conf.GetProps().Mongo.StatesTable, state)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get state %w", err)
	}

	return state, nil
}

// FindStates searches states by query
func (r *MongoRepositoryV1) FindStates(q StateQuery) ([]model.State, error) {
	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(conf.GetProps().Mongo.StatesTable)
	states := make([]model.State, 0)

	query := bson.M{}

	if q.CountryID != "" {
		query["country_id"] = bson.ObjectIdHex(q.CountryID)
	}

	err := col.Find(query).All(&states)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get states %w", err)
	}

	return states, nil
}

// FindNeighbourhoodByID searches a neighbourhood by id
func (r *MongoRepositoryV1) FindNeighbourhoodByID(id string) (*model.Neighbourhood, error) {
	neighbourhood := &model.Neighbourhood{}
	err := r.findEntityByID(id, conf.GetProps().Mongo.NeighbourhoodsTable, neighbourhood)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get neighborhood %w", err)
	}

	return neighbourhood, nil
}

// FindNeighbourhoods searches neighbourhoods by query
func (r *MongoRepositoryV1) FindNeighbourhoods(q NeighbourhoodQuery) ([]model.Neighbourhood, error) {
	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(conf.GetProps().Mongo.NeighbourhoodsTable)
	neighbourhoods := make([]model.Neighbourhood, 0)

	query := bson.M{}

	if q.CountryID != "" {
		query["country_id"] = bson.ObjectIdHex(q.CountryID)
	}

	if q.StateID != "" {
		query["state_id"] = bson.ObjectIdHex(q.StateID)
	}

	if q.CityID != "" {
		query["city_id"] = bson.ObjectIdHex(q.CityID)
	}

	err := col.Find(query).All(&neighbourhoods)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get neighborhoods %w", err)
	}

	return neighbourhoods, nil
}

// UpdateNeighbourhoodByID update a neighbourhood by id
func (r *MongoRepositoryV1) UpdateNeighbourhoodByID(id string, neighbourhood model.Neighbourhood) error {
	err := r.updateEntityByID(id, neighbourhood, conf.GetProps().Mongo.NeighbourhoodsTable)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return pkgErrors.ErrEntityNotFound
		}

		return fmt.Errorf("failed to update neighborhood %w", err)
	}

	return nil
}

// UpdateAirportByID update a city by id
func (r *MongoRepositoryV1) UpdateAirportByID(id string, airport model.Airport) error {
	err := r.updateEntityByID(id, airport, conf.GetProps().Mongo.AirportsTable)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return pkgErrors.ErrEntityNotFound
		}

		return fmt.Errorf("failed to update airport %w", err)
	}

	return nil
}

// FindAirportByID searches an airport by id
func (r *MongoRepositoryV1) FindAirportByID(id string) (*model.Airport, error) {
	airport := &model.Airport{}
	err := r.findEntityByID(id, conf.GetProps().Mongo.AirportsTable, airport)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get airport %w", err)
	}

	return airport, nil
}

// FindAirports searches airports by query
func (r *MongoRepositoryV1) FindAirports(q AirportQuery) ([]model.Airport, error) {
	s := r.Session.Copy()
	defer s.Close()

	col := s.DB(conf.GetProps().Mongo.DBV1).C(conf.GetProps().Mongo.AirportsTableV1)
	airports := make([]model.Airport, 0)

	query := bson.M{}

	if q.CountryID != "" {
		query["country_id"] = bson.ObjectIdHex(q.CountryID)
	}

	if q.StateID != "" {
		query["state_id"] = bson.ObjectIdHex(q.StateID)
	}

	if q.CityID != "" {
		query["city_id"] = bson.ObjectIdHex(q.CityID)
	}

	if q.IataCode != "" {
		iataCodes := strings.Split(q.IataCode, ",")
		query["iata_code"] = bson.M{"$in": iataCodes}
	}

	err := col.Find(query).All(&airports)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return nil, pkgErrors.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to get airports %w", err)
	}

	for _, airport := range airports {
		country, err := r.FindCountryByID(airport.CountryID.String())
		if err == nil {
			airport.CountryCode = country.Alpha3Code
		}
	}

	return airports, nil
}

func validateObjectID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("ID[%s] is invalid", id)
	}

	return nil
}
