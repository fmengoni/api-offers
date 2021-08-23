package main

import (
	"github.com/basset-la/api-geo/conf"
	"github.com/basset-la/api-geo/server"
	"github.com/basset-la/logrus-logzio-hook/logzio"
	"github.com/sirupsen/logrus"
)

func main() {
	fields := logrus.Fields{"version": conf.Version}

	hook, err := logzio.NewHook(conf.GetProps().Logger.Token, conf.GetProps().Logger.AppName, fields)
	if err != nil {
		panic(err)
	}

	logrus.AddHook(hook)
	logrus.SetLevel(logrus.InfoLevel)

	server.Start()
}
