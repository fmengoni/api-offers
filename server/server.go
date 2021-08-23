package server

import (
	"fmt"
	"net/http"

	"github.com/basset-la/api-geo/conf"
	"github.com/basset-la/api-geo/repository"
	"github.com/basset-la/api-geo/service"
	utils "github.com/basset-la/utils/v4/http"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/swag/example/basic/docs"
)

func Start() {
	// Configure Swagger Docs
	docs.SwaggerInfo.Title = "API Customers"
	docs.SwaggerInfo.Description = "API dedicated to manage customers"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = conf.GetProps().App.Host
	docs.SwaggerInfo.BasePath = conf.GetProps().App.Path
	docs.SwaggerInfo.Version = conf.Version

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(conf.GetProps().NewRelic.AppName),
		newrelic.ConfigLicense(conf.GetProps().NewRelic.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		logrus.Fatalf("Failed to start new relic app %v", err)
		panic(err)
	}

	repo, err := repository.NewMongoRepository(conf.GetProps().Mongo.URI, conf.GetProps().Mongo.DB, conf.GetProps().Mongo.AirportsTable, conf.GetProps().Mongo.GeoCoordinatesTable)

	if err != nil {
		panic(fmt.Errorf("failed to create mongo repository. %w", err))
	}

	defer repo.Close()

	repoV1, err := repository.NewMongoRepositoryV1(conf.GetProps().Mongo.URI)

	if err != nil {
		panic(fmt.Errorf("failed to create mongo repository. %w", err))
	}

	defer repoV1.Close()

	geoService := service.NewGeoService(repoV1)

	env = AppEnv{
		geoRepository:   repo,
		geoRepositoryV1: repoV1,
		geoService:      geoService,
	}

	logrus.Info("Application listen in port 8080")
	logrus.Fatal(http.ListenAndServe(":8080", utils.NewRouterWithNewRelic(conf.GetProps().App.Path, routes, nrApp)))
}

var env AppEnv

type AppEnv struct {
	geoRepository   *repository.MongoRepository
	geoRepositoryV1 *repository.MongoRepositoryV1
	geoService      *service.GeoService
}
