package service

import (
	"fmt"
	"strings"
	"sync"

	pkgErrors "github.com/basset-la/api-geo/errors"
	"github.com/basset-la/api-geo/model"
	"github.com/basset-la/api-geo/repository"
)

const ErrorsMaxCount = 2

type Service interface {
	MapGeoEntitiesToBaseEntities(geoEntities []model.GeoEntity) map[model.GeoEntityType]*model.BaseEntity
	MapGeoEntityToBaseEntity(geoEntity model.GeoEntity) (*model.BaseEntity, error)
	MapGeoEntitiesToAccommodation(geoEntities []model.GeoEntity) map[string]*model.Accommodation
	MapGeoEntityToAccommodation(geoEntity model.GeoEntity) (*model.Accommodation, error)
	GetEntitiesByIataCode(iataCode string, language string) ([]model.Entity, error)
}

type GeoService struct {
	repo *repository.MongoRepositoryV1
}

func NewGeoService(r *repository.MongoRepositoryV1) *GeoService {
	return &GeoService{
		repo: r,
	}
}

func (s *GeoService) MapGeoEntitiesToBaseEntities(geoEntities []model.GeoEntity) map[model.GeoEntityType]*model.BaseEntity {
	result := map[model.GeoEntityType]*model.BaseEntity{}

	for _, geoEntity := range geoEntities {
		entity, _ := s.MapGeoEntityToBaseEntity(geoEntity)
		result[geoEntity.Type] = entity
	}

	return result
}

func (s *GeoService) MapGeoEntityToBaseEntity(geoEntity model.GeoEntity) (*model.BaseEntity, error) {
	switch geoEntity.Type {
	case model.GeoEntityTypeCity:
		city, _ := s.repo.FindCityByID(geoEntity.ID.Hex())

		return &city.BaseEntity, nil

	case model.GeoEntityTypeCountry:
		country, _ := s.repo.FindCountryByID(geoEntity.ID.Hex())

		return &country.BaseEntity, nil

	case model.GeoEntityTypeState:
		state, _ := s.repo.FindStateByID(geoEntity.ID.Hex())

		return &state.BaseEntity, nil

	case model.GeoEntityTypeAccommodation:
		accommodation, _ := s.repo.FindAccommodationByID(geoEntity.ID.Hex())

		return &accommodation.BaseEntity, nil

	case model.GeoEntityTypeAirport:
		airport, _ := s.repo.FindAirportByID(geoEntity.ID.Hex())

		return &airport.BaseEntity, nil

	case model.GeoEntityTypeNeighbourhood:
		neighborhood, _ := s.repo.FindNeighbourhoodByID(geoEntity.ID.Hex())

		return &neighborhood.BaseEntity, nil
	}

	return nil, fmt.Errorf("not found geo entity type %s", geoEntity.Type)
}

func (s *GeoService) MapGeoEntitiesToAccommodation(geoEntities []model.GeoEntity) map[string]*model.Accommodation {
	result := map[string]*model.Accommodation{}

	for _, geoEntity := range geoEntities {
		entity, err := s.MapGeoEntityToAccommodation(geoEntity)
		if err == nil {
			result[geoEntity.ID.Hex()] = entity
		}
	}

	return result
}

func (s *GeoService) MapGeoEntityToAccommodation(geoEntity model.GeoEntity) (*model.Accommodation, error) {
	if geoEntity.Type == model.GeoEntityTypeAccommodation {
		accommodation, err := s.repo.FindAccommodationByID(geoEntity.ID.Hex())

		return accommodation, fmt.Errorf("error mapping geo entity, %w", err)
	}

	return nil, fmt.Errorf("geo entity %s must be accommodation type", geoEntity.Type)
}

func (s *GeoService) GetEntitiesByIataCode(iataCode string, language string) ([]model.Entity, error) {
	err := s.validateRequest(iataCode, language)
	if err != nil {
		return nil, err
	}

	var result []model.Entity

	var errs []string

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		cities, err := s.findCities(iataCode, language)

		if err != nil {
			errs = append(errs, err.Error())
		} else {
			result = append(result, cities...)
		}

		wg.Done()
	}()

	wg.Add(1)

	go func() {
		airports, err := s.findAirports(iataCode, language)

		if err != nil {
			errs = append(errs, err.Error())
		} else {
			result = append(result, airports...)
		}

		wg.Done()
	}()

	wg.Wait()

	if len(errs) == ErrorsMaxCount {
		return nil, pkgErrors.ErrEntityNotFound
	}

	return result, nil
}

func (s *GeoService) validateRequest(iataCode string, language string) error {
	if iataCode == "" || language == "" {
		return pkgErrors.ErrMissingParameters
	}

	return nil
}

func (s *GeoService) findAirports(iataCode string, language string) ([]model.Entity, error) {
	aq := repository.AirportQuery{
		IataCode: iataCode,
	}

	airports, err := s.repo.FindAirports(aq)

	if err != nil {
		return nil, fmt.Errorf("error finding airport %w", err)
	}

	result := make([]model.Entity, 0, len(airports))

	for _, airport := range airports {
		result = append(result, model.Entity{
			ID:        airport.ID.Hex(),
			IataCode:  airport.IataCode,
			Name:      airport.Name[model.Language(strings.ToLower(language))],
			Type:      model.GeoEntityTypeAirport,
			CountryID: airport.CountryID.Hex(),
		})
	}

	return result, nil
}

func (s *GeoService) findCities(iataCode string, language string) ([]model.Entity, error) {
	cq := repository.CityQuery{
		IataCode: iataCode,
	}

	cities, err := s.repo.FindCities(cq)
	if err != nil {
		return nil, fmt.Errorf("error finding city %w", err)
	}

	result := make([]model.Entity, 0, len(cities))

	for _, city := range cities {
		result = append(result, model.Entity{
			ID:        city.ID.Hex(),
			IataCode:  city.IataCode,
			Name:      city.Name[model.Language(strings.ToLower(language))],
			Type:      model.GeoEntityTypeCity,
			CountryID: city.CountryID.Hex(),
		})
	}

	return result, nil
}
