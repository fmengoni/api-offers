package conf

import "github.com/basset-la/tools/env"

var properties *Props

func GetProps() Props {
	if properties == nil {
		properties = &Props{}
		env.LoadProperties("env", properties)
	}

	return *properties
}

type Props struct {
	App struct {
		Path string `yaml:"appPath"`
		Host string `yaml:"host"`
	} `yaml:"app"`
	NewRelic struct {
		AppName    string `yaml:"appName"`
		LicenseKey string `yaml:"licenseKey"`
	} `yaml:"newRelic"`
	Logger struct {
		Token   string `yaml:"token"`
		AppName string `yaml:"appName"`
	} `yaml:"logger"`
	Mongo struct {
		URI                 string `yaml:"url"`
		DB                  string `yaml:"db"`
		DBV1                string `yaml:"dbV1"`
		AirportsTableV1     string `yaml:"airportsTableV1"`
		AirportsTable       string `yaml:"airportsTable"`
		GeoCoordinatesTable string `yaml:"geoCoordinatesTable"`
		CitiesTable         string `yaml:"citiesTable"`
		CountryTable        string `yaml:"countryTable"`
		StatesTable         string `yaml:"statesTable"`
		NeighbourhoodsTable string `yaml:"neighbourhoodsTable"`
		GeoEntitiesTable    string `yaml:"geoEntitiesTable"`
		AccommodationTable  string `yaml:"accommodationTable"`
	} `yaml:"mongo"`
}
