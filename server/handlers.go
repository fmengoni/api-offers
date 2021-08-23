package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/basset-la/api-geo/conf"
	pkgErrors "github.com/basset-la/api-geo/errors"
	"github.com/basset-la/api-geo/model"
	"github.com/basset-la/api-geo/repository"
	"github.com/basset-la/utils/v4/api"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	log "github.com/sirupsen/logrus"
)

// healthCheckHandler godoc
// @Summary Health Check
// @Description Method used by the application load balancer to check the status of the application
// @Produce json
// @Tags Internal
// @Success 200 {string} string version
// @Router /health-check [get]
func healthCheckHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	txn.Ignore()

	mongoDBConection := "Alive"

	isAlive := env.geoRepository.CheckIfRepositoryIsActive()

	if !isAlive {
		mongoDBConection = "Dead"
	}

	return api.DataJSON(http.StatusOK, map[string]string{"version": conf.Version, "mongo-db": mongoDBConection}, nil)
}

func getRegionByTypeAndID(regionType model.RegionType, r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	id := mux.Vars(r)["id"]

	var region model.Region

	err := env.geoRepository.GetRegionByTypeAndGeoID(regionType, id, &region)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusNotFound, fmt.Errorf("%s not found", regionType), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, region, nil)
}

func getRegionsByQuery(regionType model.RegionType, r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	qp := r.URL.Query()

	var err error

	qpGeoIDs := qp.Get("ids")
	qpDesc := qp.Get("descendants")
	basicQuery := qp.Get("basic")
	qspage := qp.Get("page")
	qslimit := qp.Get("limit")

	var geoIds []string

	var descendants []string

	if len(qpGeoIDs) > 0 {
		geoIds = strings.Split(qpGeoIDs, ",")
	}

	if len(qpDesc) > 0 {
		descendants = strings.Split(qpDesc, ",")
	}

	q, err := checkQueryParams(basicQuery, qslimit, qspage)

	if err != nil {
		return api.ErrJSON(http.StatusBadRequest, err, nil)
	}

	q.GeoIDs = geoIds
	q.RegionType = regionType
	q.Descendants = descendants

	if len(qp.Get("country_code")) > 0 && regionType != model.RegionTypeCountry {
		country := model.Region{}
		err = env.geoRepository.GetCountryByCountryCode(qp.Get("country_code"), &country)

		if err != nil {
			if errors.Is(err, pkgErrors.ErrEntityNotFound) {
				return api.ErrJSON(http.StatusNotFound, fmt.Errorf("%s not found", model.RegionTypeCountry), nil)
			}

			txn.NoticeError(err)

			return api.ErrJSON(http.StatusInternalServerError, err, nil)
		}

		q.Ancestors = []string{country.GeoID}
		q.AncestorsRegionType = model.RegionTypeCountry
	} else {
		q.CountryCode = qp.Get("country_code")
	}

	return getRegionsQueryResponse(q, regionType, txn)
}

func checkQueryParams(basicQuery, qslimit, qspage string) (*repository.QueryRegion, error) {
	basic := false

	var err error

	if len(basicQuery) > 0 {
		basic, err = strconv.ParseBool(basicQuery)

		if err != nil {
			return nil, fmt.Errorf("[basic] must be true or false")
		}
	}

	var limit int
	if len(qslimit) > 0 {
		limit, err = strconv.Atoi(qslimit)

		if err != nil {
			return nil, fmt.Errorf("[limit] must be a valid number")
		}
	}

	var page int
	if len(qspage) > 0 {
		page, err = strconv.Atoi(qspage)

		if err != nil {
			return nil, fmt.Errorf("[page] must be a valid number")
		}

		if page <= 0 {
			return nil, fmt.Errorf("[page] must be a number greater than 0")
		}
	}

	q := repository.QueryRegion{
		Basic: basic,
		Page:  page,
		Limit: limit,
	}

	return &q, nil
}

func getRegionsQueryResponse(q *repository.QueryRegion, regionType model.RegionType, txn *newrelic.Transaction) *api.Response {
	regions := make([]model.Region, 0)

	err := env.geoRepository.GetRegions(*q, &regions)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusNotFound, fmt.Errorf("%s not found", regionType), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, regions, nil)
}

func getCountryByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypeCountry, r)
}

func getCityByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypeCity, r)
}

func getHighLevelRegionByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypeHighLevelRegion, r)
}

func getContinentByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypeContinent, r)
}

func getMultiCityVicinityByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypeMultiCityVicinity, r)
}

func getTrainStationByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypeTrainStation, r)
}

func getMetroStationByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypeMetroStation, r)
}

func getProvinceStateByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypeProvinceState, r)
}

func getPOIsByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypePOI, r)
}

func getNeighborhoodByIDHandlerV2(r *http.Request) *api.Response {
	return getRegionByTypeAndID(model.RegionTypeNeighborhood, r)
}

func getAirportByIATACode(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	iataCode := mux.Vars(r)["iata_code"]

	var airport model.AirportV2

	err := env.geoRepository.GetAirportByIATACode(iataCode, &airport)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusNotFound, fmt.Errorf("airport not found"), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, airport, nil)
}

// Get regions by query

func getCountryByQueryHandlerV2(r *http.Request) *api.Response {
	return getRegionsByQuery(model.RegionTypeCountry, r)
}

func getCityByQueryHandlerV2(r *http.Request) *api.Response {
	return getRegionsByQuery(model.RegionTypeCity, r)
}

// func getHighLevelRegionByQueryHandlerV2(env *Env, r *http.Request, txn newrelic.Transaction) *api.Response {
//	return getRegionsByQuery(RegionTypeHighLevelRegion, env, r, txn)
//}

func getContinentByQueryHandlerV2(r *http.Request) *api.Response {
	return getRegionsByQuery(model.RegionTypeContinent, r)
}

func getMultiCityVicinityByQueryHandlerV2(r *http.Request) *api.Response {
	return getRegionsByQuery(model.RegionTypeMultiCityVicinity, r)
}

func getTrainStationByQueryHandlerV2(r *http.Request) *api.Response {
	return getRegionsByQuery(model.RegionTypeTrainStation, r)
}

func getMetroStationByQueryHandlerV2(r *http.Request) *api.Response {
	return getRegionsByQuery(model.RegionTypeMetroStation, r)
}

func getProvinceStateByQueryHandlerV2(r *http.Request) *api.Response {
	return getRegionsByQuery(model.RegionTypeProvinceState, r)
}

func getPOIsByQueryHandlerV2(r *http.Request) *api.Response {
	return getRegionsByQuery(model.RegionTypePOI, r)
}

func getNeighborhoodByQueryHandlerV2(r *http.Request) *api.Response {
	return getRegionsByQuery(model.RegionTypeNeighborhood, r)
}

func getAirportByQuery(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	qp := r.URL.Query()

	qpIataCodes := qp.Get("iata_codes")

	var iataCodes []string

	if len(qpIataCodes) > 0 {
		iataCodes = strings.Split(qpIataCodes, ",")
	}

	q := repository.QueryAirport{
		CountryCode: qp.Get("country_code"),
		IataCodes:   iataCodes,
	}

	airports := make([]model.AirportV2, 0)

	err := env.geoRepository.GetAirportByQuery(q, &airports)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusNotFound, fmt.Errorf("airport not found"), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, airports, nil)
}

func getIntersections(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	qp := r.URL.Query()

	latitude, err := strconv.ParseFloat(qp.Get("latitude"), 64)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusBadRequest, err, nil)
	}

	longitude, err := strconv.ParseFloat(qp.Get("longitude"), 64)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusBadRequest, err, nil)
	}

	regionTypes := qp.Get("region_types")

	var rts []model.RegionType

	if len(regionTypes) > 0 {
		rts = make([]model.RegionType, 0, len(regionTypes))
		for i, e := range strings.Split(regionTypes, ",") {
			rts[i] = model.RegionType(e)
		}
	}

	regions := make([]model.GeoRegion, 0)

	err = env.geoRepository.GetIntersectedRegions(*model.NewPointGeometry([]interface{}{longitude, latitude}), rts, &regions)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusBadRequest, err, nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, regions, nil)
}

func insertAccommodation(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	var accommodation model.GeoRegion

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&accommodation)

	if err != nil {
		return api.ErrJSON(http.StatusBadRequest, fmt.Errorf("failed to read body"), nil)
	}

	intersectedRegions, err := env.geoRepository.InsertAccommodation(&accommodation)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, fmt.Errorf("failed to insert accommodation id: %s", accommodation.GeoID), nil)
	}

	return api.DataJSON(http.StatusOK, intersectedRegions, nil)
}

func getNearbyRegions(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	qp := r.URL.Query()
	longitude, _ := strconv.ParseFloat(qp.Get("longitude"), 64)
	latitude, _ := strconv.ParseFloat(qp.Get("latitude"), 64)
	radius, _ := strconv.ParseFloat(qp.Get("radius"), 64)
	types := qp.Get("types")

	rt := strings.Split(types, ",")

	regionTypes := make([]model.RegionType, 0, len(rt))

	for i, e := range rt {
		regionTypes[i] = model.RegionType(e)
	}

	regions, err := env.geoRepository.GetNearByRegions(latitude, longitude, regionTypes, radius)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, regions, nil)
}

func getGeoRegion(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	id := mux.Vars(r)["id"]

	var region model.GeoRegion

	err := env.geoRepository.GetGeoRegion(id, &region)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusNotFound, fmt.Errorf("region not found"), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, region, nil)
}

func saveRegion(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	var region model.Region
	err := json.NewDecoder(r.Body).Decode(&region)

	if err != nil {
		return api.ErrJSON(http.StatusBadRequest, fmt.Errorf("failed to read body"), nil)
	}

	err = env.geoRepository.SaveRegion(&region)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, region, nil)
}

func updateRegion(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	var region model.Region

	err := json.NewDecoder(r.Body).Decode(&region)

	if err != nil {
		return api.ErrJSON(http.StatusBadRequest, fmt.Errorf("failed to read body"), nil)
	}

	err = env.geoRepository.UpdateRegion(&region)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, region, nil)
}

func saveAirport(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	var airport model.AirportV2

	err := json.NewDecoder(r.Body).Decode(&airport)

	if err != nil {
		return api.ErrJSON(http.StatusBadRequest, fmt.Errorf("failed to read body"), nil)
	}

	err = env.geoRepository.SaveAirport(&airport)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, airport, nil)
}

func updateAirport(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	var airport model.AirportV2

	err := json.NewDecoder(r.Body).Decode(&airport)

	if err != nil {
		return api.ErrJSON(http.StatusBadRequest, fmt.Errorf("failed to read body"), nil)
	}

	err = env.geoRepository.UpdateAirport(&airport)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, airport, nil)
}

func saveGeoRegion(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	var region model.GeoRegion
	err := json.NewDecoder(r.Body).Decode(&region)

	if err != nil {
		return api.ErrJSON(http.StatusBadRequest, fmt.Errorf("failed to read body"), nil)
	}

	err = env.geoRepository.SaveGeoRegion(&region)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, region, nil)
}

func updateGeoRegion(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	var region model.GeoRegion

	err := json.NewDecoder(r.Body).Decode(&region)

	if err != nil {
		return api.ErrJSON(http.StatusBadRequest, fmt.Errorf("failed to read body"), nil)
	}

	err = env.geoRepository.UpdateGeoRegion(&region)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, region, nil)
}

// V1

func getCitiesHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	qp := r.URL.Query()
	countryCode := qp.Get("country_code")
	countryID := qp.Get("country_id")
	iataCode := qp.Get("iata_code")
	stateID := qp.Get("state_id")

	if countryCode != "" && countryID == "" {
		countries, err := env.geoRepositoryV1.FindCountries(repository.CountryQuery{
			Alpha2Code: countryCode,
		})

		if err != nil {
			if errors.Is(err, pkgErrors.ErrEntityNotFound) {
				return api.DataJSON(http.StatusOK, []model.City{}, nil)
			}

			txn.NoticeError(err)

			return api.ErrJSON(http.StatusInternalServerError, err, nil)
		}

		if len(countries) == 0 {
			return api.DataJSON(http.StatusOK, []model.City{}, nil)
		}

		countryID = countries[0].ID.Hex()
	}

	q := repository.CityQuery{
		CountryID: countryID,
		StateID:   stateID,
		IataCode:  iataCode,
	}

	cities, err := env.geoRepositoryV1.FindCities(q)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusOK, []model.City{}, nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, cities, nil)
}

func getCountriesHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	qp := r.URL.Query()
	q := repository.CountryQuery{
		Alpha2Code: qp.Get("alpha2_code"),
		Alpha3Code: qp.Get("alpha3_code"),
	}

	countries, err := env.geoRepositoryV1.FindCountries(q)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusOK, []model.Country{}, nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, countries, nil)
}

func getCountryByIDHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	id := mux.Vars(r)["id"]

	country, err := env.geoRepositoryV1.FindCountryByID(id)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusNotFound, fmt.Errorf("country not found"), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, country, nil)
}

func getStateByIDHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	id := mux.Vars(r)["id"]

	state, err := env.geoRepositoryV1.FindStateByID(id)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusNotFound, fmt.Errorf("state not found"), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, state, nil)
}

func getStatesHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	qp := r.URL.Query()
	countryID := qp.Get("country_id")

	countryCode := qp.Get("country_code")

	if countryCode != "" {
		countries, err := env.geoRepositoryV1.FindCountries(repository.CountryQuery{
			Alpha2Code: countryCode,
		})
		if err != nil {
			if errors.Is(err, pkgErrors.ErrEntityNotFound) {
				return api.DataJSON(http.StatusOK, []model.State{}, nil)
			}

			txn.NoticeError(err)

			return api.ErrJSON(http.StatusInternalServerError, err, nil)
		}

		if len(countries) == 0 {
			return api.DataJSON(http.StatusOK, []model.State{}, nil)
		}

		countryID = countries[0].ID.Hex()
	}

	states, err := env.geoRepositoryV1.FindStates(repository.StateQuery{
		CountryID: countryID,
	})

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusOK, []model.State{}, nil)
		}

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, states, nil)
}

func getNeighbourhoodsByIDHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	id := mux.Vars(r)["id"]

	neighbourhood, err := env.geoRepositoryV1.FindNeighbourhoodByID(id)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusNotFound, fmt.Errorf("neighborhood not found"), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, neighbourhood, nil)
}

func getNeighbourhoodsHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	qp := r.URL.Query()
	q := repository.NeighbourhoodQuery{
		CountryID: qp.Get("country_id"),
		CityID:    qp.Get("city_id"),
		StateID:   qp.Get("state_id"),
	}

	neighbourhoods, err := env.geoRepositoryV1.FindNeighbourhoods(q)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusOK, []model.Neighbourhood{}, nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, neighbourhoods, nil)
}

func getAirportByIDHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	id := mux.Vars(r)["id"]

	airport, err := env.geoRepositoryV1.FindAirportByID(id)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusOK, fmt.Errorf("airport not found"), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, airport, nil)
}

func getAirportsHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	qp := r.URL.Query()
	q := repository.AirportQuery{
		CountryID: qp.Get("country_id"),
		CityID:    qp.Get("city_id"),
		StateID:   qp.Get("state_id"),
		IataCode:  qp.Get("iata_code"),
	}

	airports, err := env.geoRepositoryV1.FindAirports(q)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusOK, []model.Airport{}, nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, airports, nil)
}

func getAccommodationsHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	geoID := r.URL.Query().Get("geo_id")

	if geoID == "" {
		return api.ErrJSON(http.StatusBadRequest, fmt.Errorf("[geo_id] must be valid"), nil)
	}

	entity, err := env.geoRepositoryV1.FindGeoEntityByID(geoID)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusOK, []model.Accommodation{}, nil)
		}

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	// Intersect geo entities
	entities, err := env.geoRepositoryV1.IntersectsGeoEntities(entity.Geometry, []model.GeoEntityType{model.GeoEntityTypeCountry,
		model.GeoEntityTypeCity,
		model.GeoEntityTypeState,
		model.GeoEntityTypeNeighbourhood,
		model.GeoEntityTypeAirport})

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	result := env.geoService.MapGeoEntitiesToAccommodation(entities)

	return api.DataJSON(http.StatusOK, result, nil)
}

func getAccommodationByIDHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	id := mux.Vars(r)["id"]

	accommodation, err := env.geoRepositoryV1.FindAccommodationByID(id)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusNotFound, fmt.Errorf("accommodation not found"), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, accommodation, nil)
}

func insertAccommodationGeometry(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())

	var accommodation model.AccommodationGeometry

	err := json.NewDecoder(r.Body).Decode(&accommodation)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusBadRequest, err, nil)
	}

	id, err := env.geoRepositoryV1.InsertAccommodationV1(accommodation.Accommodation)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	err = env.geoRepositoryV1.InsertGeoEntity(accommodation.Location, *id)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusCreated, map[string]string{"id": *id}, nil)
}

func intersectsGeoEntities(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	qp := r.URL.Query()
	latitude, err := strconv.ParseFloat(qp.Get("latitude"), 64)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusBadRequest, err, nil)
	}

	longitude, err := strconv.ParseFloat(qp.Get("longitude"), 64)

	if err != nil {
		log.Error(err)

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusBadRequest, err, nil)
	}

	entities, err := env.geoRepositoryV1.IntersectsGeoEntities(*model.NewPointGeometry([]interface{}{longitude, latitude}), []model.GeoEntityType{})

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	result := env.geoService.MapGeoEntitiesToBaseEntities(entities)

	return api.DataJSON(http.StatusOK, result, nil)
}

func updateCityByIDHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	id := mux.Vars(r)["id"]

	var city model.City
	err := json.NewDecoder(r.Body).Decode(&city)

	if err != nil {
		txn.NoticeError(err)

		return api.ErrJSON(http.StatusBadRequest, err, nil)
	}

	err = env.geoRepositoryV1.UpdateCityByID(id, city)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusNotFound, fmt.Errorf("city not found"), nil)
		}

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, struct{}{}, nil)
}

func getCityByIDHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	id := mux.Vars(r)["id"]

	city, err := env.geoRepositoryV1.FindCityByID(id)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.ErrJSON(http.StatusNotFound, fmt.Errorf("city not found"), nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, city, nil)
}

func getEntitiesByIataCodeHandler(r *http.Request) *api.Response {
	txn := newrelic.FromContext(r.Context())
	qp := r.URL.Query()
	iataCode := qp.Get("iata_code")
	language := qp.Get("language")

	entity, err := env.geoService.GetEntitiesByIataCode(iataCode, language)

	if err != nil {
		if errors.Is(err, pkgErrors.ErrMissingParameters) {
			return api.ErrJSON(http.StatusBadRequest, err, nil)
		} else if errors.Is(err, pkgErrors.ErrEntityNotFound) {
			return api.DataJSON(http.StatusOK, []model.Entity{}, nil)
		}

		txn.NoticeError(err)

		return api.ErrJSON(http.StatusInternalServerError, err, nil)
	}

	return api.DataJSON(http.StatusOK, entity, nil)
}
