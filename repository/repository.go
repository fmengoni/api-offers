package repository

import (
	"errors"
	"fmt"

	pkgErrors "github.com/basset-la/api-geo/errors"
	geoModel "github.com/basset-la/api-geo/model"
	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const earthRadius float64 = 6378.1

type Repository interface {
	GetRegionByTypeAndGeoID(regionType geoModel.RegionType, geoID string, r *geoModel.Region) error
	GetRegions(query QueryRegion, r *[]geoModel.Region) error
	GetCountryByCountryCode(countryCode string, r *geoModel.Region) error
	Count(regionType geoModel.RegionType) (int, error)
	SaveRegion(e *geoModel.Region) error
	UpdateRegion(e *geoModel.Region) error
	SaveAirport(e *geoModel.AirportV2) error
	UpdateAirport(e *geoModel.AirportV2) error
	GetAirportByIATACode(iataCode string, a *geoModel.AirportV2) error
	GetAirportByQuery(q QueryAirport, a *[]geoModel.AirportV2) error
	GetIntersectedRegions(geometry geoModel.Geometry, regionTypes []geoModel.RegionType, r *[]geoModel.GeoRegion) error
	InsertAccommodation(accommodation *geoModel.GeoRegion) ([]geoModel.GeoRegion, error)
	GetNearByRegions(latitude float64, longitude float64, regionTypes []geoModel.RegionType, radius float64) ([]geoModel.GeoRegion, error)
}

// QueryRegion for regions
type QueryRegion struct {
	RegionType          geoModel.RegionType
	CountryCode         string
	GeoIDs              []string
	Descendants         []string
	Basic               bool
	Page                int
	Limit               int
	Ancestors           []string
	AncestorsRegionType geoModel.RegionType
}

// QueryAirport for airports
type QueryAirport struct {
	CountryCode string
	IataCodes   []string
}

// MongoRepository handles all requests to MongoDB
type MongoRepository struct {
	Session             *mgo.Session
	airportTable        string
	geoCoordinatesTable string
	db                  string
}

// NewMongoRepository creates a new mongo repository
func NewMongoRepository(url, db, airportTable, geoCoordinatesTable string) (*MongoRepository, error) {
	s, err := mgo.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connecto to mongodb. %w", err)
	}

	r := &MongoRepository{
		Session:             s,
		airportTable:        airportTable,
		geoCoordinatesTable: geoCoordinatesTable,
		db:                  db,
	}

	return r, nil
}

func (repo *MongoRepository) Close() {
	repo.Session.Close()
}

// GetRegionByTypeAndGeoID returns a region by type and id
func (repo *MongoRepository) GetRegionByTypeAndGeoID(regionType geoModel.RegionType, geoID string, r *geoModel.Region) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(string(regionType))

	err := col.Find(bson.M{"geo_id": geoID}).One(r)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return fmt.Errorf("failed to get region by type %s and geo id %s. %w", regionType, geoID, pkgErrors.ErrEntityNotFound)
		}

		return fmt.Errorf("failed to get region by type %s and geo id %s. %w", regionType, geoID, err)
	}

	return nil
}

// SaveRegion saves a region in mongoDB
func (repo *MongoRepository) SaveRegion(r *geoModel.Region) error {
	s := repo.Session.Copy()
	defer s.Close()

	r.ID = bson.NewObjectId()
	col := s.DB(repo.db).C(string(r.Type))

	if err := col.Insert(r); err != nil {
		return fmt.Errorf("failed to save region. %w", err)
	}

	return nil
}

// UpdateRegion saves a region in mongoDB
func (repo *MongoRepository) UpdateRegion(e *geoModel.Region) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(string(e.Type))
	err := col.Update(bson.M{"geo_id": e.GeoID}, e)

	if err != nil {
		return fmt.Errorf("failed to update region. %w", err)
	}

	return nil
}

// SaveAirport saves a airport in mongoDB
func (repo *MongoRepository) SaveAirport(a *geoModel.AirportV2) error {
	s := repo.Session.Copy()
	defer s.Close()

	a.ID = bson.NewObjectId()
	col := s.DB(repo.db).C(repo.airportTable)

	if err := col.Insert(a); err != nil {
		return fmt.Errorf("failed to save airport. %w", err)
	}

	return nil
}

// UpdateAirport update a airport in mongodb
func (repo *MongoRepository) UpdateAirport(a *geoModel.AirportV2) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(repo.airportTable)
	err := col.Update(bson.M{"iata": a.IataCode}, a)

	if err != nil {
		return fmt.Errorf("failed to update airport. %w", err)
	}

	return nil
}

// GetAirportByIATACode returns an airport by iataCode
func (repo *MongoRepository) GetAirportByIATACode(iataCode string, a *geoModel.AirportV2) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(repo.airportTable)

	err := col.Find(bson.M{"iata": iataCode}).One(a)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return pkgErrors.ErrEntityNotFound
		}

		return fmt.Errorf("failed to get airport by iata code %s. %w", iataCode, err)
	}

	return nil
}

// GetGeoRegion returns a geo region by id
func (repo *MongoRepository) GetGeoRegion(geoID string, r *geoModel.GeoRegion) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(repo.geoCoordinatesTable)

	err := col.Find(bson.M{"geo_id": geoID}).One(r)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return fmt.Errorf("geo region %s not found. %w", geoID, pkgErrors.ErrEntityNotFound)
		}

		return fmt.Errorf("failed to get geo %s region %w", geoID, err)
	}

	r.MapCoordinates()

	return nil
}

// UpdateGeoRegion update a geo region in mongodb
func (repo *MongoRepository) UpdateGeoRegion(r *geoModel.GeoRegion) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(repo.geoCoordinatesTable)
	err := col.Update(bson.M{"geo_id": r.GeoID}, r)

	if err != nil {
		return fmt.Errorf("failed to update geo region. %w", err)
	}

	return nil
}

func (repo *MongoRepository) Count(regionType geoModel.RegionType) (count int, err error) {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(string(regionType))

	return col.Count()
}

// GetRegions returns a slice of regions
func (repo *MongoRepository) GetRegions(query QueryRegion, r *[]geoModel.Region) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(string(query.RegionType))

	dbQuery := bson.M{}

	if len(query.GeoIDs) > 0 {
		dbQuery["geo_id"] = bson.M{
			"$in": query.GeoIDs,
		}
	}

	if len(query.Descendants) > 0 {
		for _, desc := range query.Descendants {
			dbQuery["descendants."+desc] = bson.M{
				"$exists": 1,
			}
		}
	}

	if query.CountryCode != "" {
		dbQuery["country_code"] = query.CountryCode
	}

	if len(query.Ancestors) > 0 && len(query.AncestorsRegionType) > 0 {
		dbQuery["ancestors.type"] = bson.M{
			"$eq": query.AncestorsRegionType,
			"$ne": query.RegionType,
		}
		dbQuery["ancestors.geo_id"] = bson.M{
			"$in": query.Ancestors,
		}
	}

	q := col.Find(dbQuery)
	if query.Page > 0 && query.Limit > 0 {
		q = q.Skip((query.Page - 1) * query.Limit).Limit(query.Limit)
	}

	var err error
	if query.Basic {
		err = q.Select(bson.M{"geo_id": 1, "name": 1}).All(r)
	} else {
		err = q.All(r)
	}

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return pkgErrors.ErrEntityNotFound
		}

		return fmt.Errorf("failed to get regions %w", err)
	}

	return nil
}

// GetCountryByCountryCode returns a country for a given country code
func (repo *MongoRepository) GetCountryByCountryCode(countryCode string, r *geoModel.Region) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(string(geoModel.RegionTypeCountry))

	dbQuery := bson.M{}

	if countryCode == "" {
		return fmt.Errorf("invalid country code")
	}

	dbQuery["country_code"] = countryCode

	err := col.Find(dbQuery).One(r)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return pkgErrors.ErrEntityNotFound
		}

		return fmt.Errorf("failed to get country %w", err)
	}

	return nil
}

// GetAirportByQuery returns a slice of regions
func (repo *MongoRepository) GetAirportByQuery(q QueryAirport, a *[]geoModel.AirportV2) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(repo.airportTable)

	dbQuery := bson.M{}

	if len(q.IataCodes) > 0 {
		dbQuery["iata"] = bson.M{
			"$in": q.IataCodes,
		}
	}

	if q.CountryCode != "" {
		dbQuery["countrycode"] = q.CountryCode
	}

	err := col.Find(dbQuery).All(a)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return pkgErrors.ErrEntityNotFound
		}

		return fmt.Errorf("failed to get airport %w", err)
	}

	return nil
}

// GetIntersectedRegions finds all regions that intersects with a polygon
func (repo *MongoRepository) GetIntersectedRegions(geometry geoModel.Geometry, regionTypes []geoModel.RegionType, r *[]geoModel.GeoRegion) error {
	s := repo.Session.Copy()
	defer s.Close()

	col := s.DB(repo.db).C(repo.geoCoordinatesTable)

	dbQuery := bson.M{
		"bounding_polygon": map[string]map[string]geoModel.Geometry{"$geoIntersects": {"$geometry": geometry}},
	}

	if len(regionTypes) > 0 {
		dbQuery["type"] = map[string][]geoModel.RegionType{"$in": regionTypes}
	}

	err := col.Find(dbQuery).Select(bson.M{"type": 1, "geo_id": 1}).All(r)

	if err != nil {
		if errors.Is(err, mgo.ErrNotFound) {
			return pkgErrors.ErrEntityNotFound
		}

		return fmt.Errorf("failed to get regions %w", err)
	}

	return nil
}

// InsertAccommodation inserts an acommodation in the GeoRegion collection and adds all ancestor's relations
func (repo *MongoRepository) InsertAccommodation(accommodation *geoModel.GeoRegion) ([]geoModel.GeoRegion, error) {
	intersectedRegions := make([]geoModel.GeoRegion, 0)
	err := repo.GetIntersectedRegions(accommodation.Geometry, nil, &intersectedRegions)

	accommodation.Type = geoModel.RegionTypeAccommodation

	if err != nil {
		return nil, err
	}

	for _, e := range intersectedRegions {
		var r geoModel.Region
		err1 := repo.GetRegionByTypeAndGeoID(e.Type, e.GeoID, &r)

		if err1 != nil {
			log.Error(fmt.Errorf("region id: %s, Type: %s", e.Type, e.GeoID))
		}

		if r.Type == geoModel.RegionTypeCity || r.Type == geoModel.RegionTypeMultiCityVicinity || r.Type == geoModel.RegionTypeNeighborhood {
			if len(r.Descendants.Accommodations) == 0 {
				r.Descendants.Accommodations = make([]string, 0)
			}

			r.Descendants.Accommodations = append(r.Descendants.Accommodations, accommodation.GeoID)
			err1 = repo.UpdateRegion(&r)

			if err1 != nil {
				log.Error(fmt.Errorf("region id: %s, Type: %s", e.Type, e.GeoID))
			}
		}
	}

	err = repo.SaveGeoRegion(accommodation)

	if err != nil {
		return nil, err
	}

	return intersectedRegions, nil
}

// SaveGeoRegion saves a GeoRegion in mongoDB
func (repo *MongoRepository) SaveGeoRegion(r *geoModel.GeoRegion) error {
	s := repo.Session.Copy()
	defer s.Close()

	r.ID = bson.NewObjectId()
	col := s.DB(repo.db).C(repo.geoCoordinatesTable)

	if err := col.Insert(r); err != nil {
		return fmt.Errorf("failed to save region")
	}

	return nil
}

// GetNearByRegions returns regions in a radius
func (repo *MongoRepository) GetNearByRegions(latitude float64, longitude float64, regionTypes []geoModel.RegionType, radius float64) ([]geoModel.GeoRegion, error) {
	s := repo.Session.Copy()
	defer s.Close()

	query := bson.M{"bounding_polygon": bson.M{
		"$geoWithin": bson.M{
			"$centerSphere": []interface{}{
				[]interface{}{longitude, latitude}, radius / earthRadius,
			},
		},
	}}

	regions := make([]geoModel.GeoRegion, 0)

	err := s.DB(repo.db).C(repo.geoCoordinatesTable).Find(query).All(&regions)

	if err != nil {
		return []geoModel.GeoRegion{}, nil
	}

	return regions, nil
}

func (repo *MongoRepository) CheckIfRepositoryIsActive() bool {
	if err := repo.Session.Ping(); err != nil {
		return false
	}

	return true
}
