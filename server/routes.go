package server

import (
	"github.com/basset-la/swagger-go/docs"
	"github.com/basset-la/utils/v4/http"
)

var routes = http.Routes{
	{
		Name:        "Health Check",
		Method:      "GET",
		Pattern:     "/health-check",
		HandlerFunc: healthCheckHandler,
		ShouldLog:   false,
	},
	{
		Name:        "Get Swagger Docs",
		Method:      "GET",
		Pattern:     "/docs/{rest:.*}",
		HandlerFunc: docs.DocsHandler,
		ShouldLog:   false,
	},

	// V2
	{
		Name:        "Find Countries by ID V2",
		Method:      "GET",
		Pattern:     "/v2/countries/{id}",
		HandlerFunc: getCountryByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find cities by ID V2",
		Method:      "GET",
		Pattern:     "/v2/cities/{id}",
		HandlerFunc: getCityByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find high level region by ID V2",
		Method:      "GET",
		Pattern:     "/v2/high-level-regions/{id}",
		HandlerFunc: getHighLevelRegionByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find continents by ID V2",
		Method:      "GET",
		Pattern:     "/v2/continents/{id}",
		HandlerFunc: getContinentByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find multi city vicinity by ID V2",
		Method:      "GET",
		Pattern:     "/v2/multi-city-vicinities/{id}",
		HandlerFunc: getMultiCityVicinityByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find train stations by ID V2",
		Method:      "GET",
		Pattern:     "/v2/train-stations/{id}",
		HandlerFunc: getTrainStationByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find metro stations by ID V2",
		Method:      "GET",
		Pattern:     "/v2/metro-stations/{id}",
		HandlerFunc: getMetroStationByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find province state by ID V2",
		Method:      "GET",
		Pattern:     "/v2/province-states/{id}",
		HandlerFunc: getProvinceStateByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find pois by ID V2",
		Method:      "GET",
		Pattern:     "/v2/points-of-interest/{id}",
		HandlerFunc: getPOIsByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find neighborhood by ID V2",
		Method:      "GET",
		Pattern:     "/v2/neighborhoods/{id}",
		HandlerFunc: getNeighborhoodByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find airports by ID V2",
		Method:      "GET",
		Pattern:     "/v2/airports/{iata_code}",
		HandlerFunc: getAirportByIATACode,
		ShouldLog:   true,
	},

	{
		Name:        "Find polygons by ID V2",
		Method:      "GET",
		Pattern:     "/v2/polygons/{id}",
		HandlerFunc: getGeoRegion,
		ShouldLog:   true,
	},

	{
		Name:        "Save region V2",
		Method:      "POST",
		Pattern:     "/v2/regions",
		HandlerFunc: saveRegion,
		ShouldLog:   true,
	},

	{
		Name:        "Update region V2",
		Method:      "PUT",
		Pattern:     "/v2/regions",
		HandlerFunc: updateRegion,
		ShouldLog:   true,
	},

	{
		Name:        "Save airport V2",
		Method:      "POST",
		Pattern:     "/v2/airports",
		HandlerFunc: saveAirport,
		ShouldLog:   true,
	},

	{
		Name:        "Update airport V2",
		Method:      "PUT",
		Pattern:     "/v2/airports",
		HandlerFunc: updateAirport,
		ShouldLog:   true,
	},

	{
		Name:        "Save geo region V2",
		Method:      "POST",
		Pattern:     "/v2/polygons",
		HandlerFunc: saveGeoRegion,
		ShouldLog:   true,
	},

	{
		Name:        "Update geo region V2",
		Method:      "PUT",
		Pattern:     "/v2/polygons",
		HandlerFunc: updateGeoRegion,
		ShouldLog:   true,
	},

	// By Query

	{
		Name:        "Find Countries V2",
		Method:      "GET",
		Pattern:     "/v2/countries",
		HandlerFunc: getCountryByQueryHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find cities V2",
		Method:      "GET",
		Pattern:     "/v2/cities",
		HandlerFunc: getCityByQueryHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find high level region V2",
		Method:      "GET",
		Pattern:     "/v2/high-level-regions",
		HandlerFunc: getContinentByIDHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find continents V2",
		Method:      "GET",
		Pattern:     "/v2/continents",
		HandlerFunc: getContinentByQueryHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find multi city vicinity V2",
		Method:      "GET",
		Pattern:     "/v2/multi-city-vicinities",
		HandlerFunc: getMultiCityVicinityByQueryHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find train stations V2",
		Method:      "GET",
		Pattern:     "/v2/train-stations",
		HandlerFunc: getTrainStationByQueryHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find metro stations V2",
		Method:      "GET",
		Pattern:     "/v2/metro-stations",
		HandlerFunc: getMetroStationByQueryHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find province state V2",
		Method:      "GET",
		Pattern:     "/v2/province-states",
		HandlerFunc: getProvinceStateByQueryHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find pois V2",
		Method:      "GET",
		Pattern:     "/v2/points-of-interest",
		HandlerFunc: getPOIsByQueryHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find neighborhood V2",
		Method:      "GET",
		Pattern:     "/v2/neighborhoods",
		HandlerFunc: getNeighborhoodByQueryHandlerV2,
		ShouldLog:   true,
	},

	{
		Name:        "Find airports V2",
		Method:      "GET",
		Pattern:     "/v2/airports",
		HandlerFunc: getAirportByQuery,
		ShouldLog:   true,
	},

	{
		Name:        "Find intersections V2",
		Method:      "GET",
		Pattern:     "/v2/intersections",
		HandlerFunc: getIntersections,
		ShouldLog:   true,
	},

	{
		Name:        "Get Regions Nearby",
		Method:      "GET",
		Pattern:     "/v2/regions/nearby",
		HandlerFunc: getNearbyRegions,
		ShouldLog:   true,
	},

	{
		Name:        "Insert accommodations V2",
		Method:      "POST",
		Pattern:     "/v2/accommodations",
		HandlerFunc: insertAccommodation,
		ShouldLog:   true,
	},

	// V1

	{
		Name:        "Find cities",
		Method:      "GET",
		Pattern:     "/cities",
		HandlerFunc: getCitiesHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find city by ID",
		Method:      "GET",
		Pattern:     "/cities/{id}",
		HandlerFunc: getCityByIDHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Update city by ID",
		Method:      "PUT",
		Pattern:     "/cities/{id}",
		HandlerFunc: updateCityByIDHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find entities by iata_code/s",
		Method:      "GET",
		Pattern:     "/entities",
		HandlerFunc: getEntitiesByIataCodeHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Intersect geo locations",
		Method:      "GET",
		Pattern:     "/intersections",
		HandlerFunc: intersectsGeoEntities,
		ShouldLog:   true,
	},

	{
		Name:        "Insert accommodation geometry",
		Method:      "POST",
		Pattern:     "/accommodations",
		HandlerFunc: insertAccommodationGeometry,
		ShouldLog:   true,
	},

	{
		Name:        "Find accommodation by id",
		Method:      "GET",
		Pattern:     "/accommodations/{id}",
		HandlerFunc: getAccommodationByIDHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find accommodations",
		Method:      "GET",
		Pattern:     "/accommodations",
		HandlerFunc: getAccommodationsHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find countries",
		Method:      "GET",
		Pattern:     "/countries",
		HandlerFunc: getCountriesHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find country by ID",
		Method:      "GET",
		Pattern:     "/countries/{id}",
		HandlerFunc: getCountryByIDHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find states",
		Method:      "GET",
		Pattern:     "/states",
		HandlerFunc: getStatesHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find state by id",
		Method:      "GET",
		Pattern:     "/states/{id}",
		HandlerFunc: getStateByIDHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find neighbourhood by id",
		Method:      "GET",
		Pattern:     "/neighbourhoods/{id}",
		HandlerFunc: getNeighbourhoodsByIDHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find neighbourhoods",
		Method:      "GET",
		Pattern:     "/neighbourhoods",
		HandlerFunc: getNeighbourhoodsHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find airport by id",
		Method:      "GET",
		Pattern:     "/airports/{id}",
		HandlerFunc: getAirportByIDHandler,
		ShouldLog:   true,
	},

	{
		Name:        "Find airports",
		Method:      "GET",
		Pattern:     "/airports",
		HandlerFunc: getAirportsHandler,
		ShouldLog:   true,
	},
}
