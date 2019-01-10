package main

import (
	"github.com/koestler/go-mqtt-to-influxdb/config"
	"github.com/koestler/go-mqtt-to-influxdb/influxDbClient"
	"github.com/koestler/go-mqtt-to-influxdb/statistics"
	"github.com/pkg/errors"
	"log"
)

func runInfluxClient(
	cfg *config.Config,
	statisticsInstance *statistics.Statistics,
	initiateShutdown chan<- error,
) *influxDbClient.ClientPool {
	influxDbClientPoolInstance := influxDbClient.RunPool()

	countStarted := 0

	for _, influxDbClientConfig := range cfg.InfluxDbClients {
		if cfg.LogWorkerStart {
			log.Printf(
				"main: start InfluxDB[%s]; Address='%s'",
				influxDbClientConfig.Name(),
				influxDbClientConfig.Address(),
			)
		}

		if client, err := influxDbClient.RunClient(influxDbClientConfig, statisticsInstance); err == nil {
			influxDbClientPoolInstance.AddClient(client)
			countStarted += 1
			if cfg.LogWorkerStart {
				log.Printf(
					"main: InfluxDbClient[%s] started; serverVersion='%s'",
					influxDbClientConfig.Name(), client.ServerVersion(),
				)
			}
		} else {
			log.Printf("main: InfluxDbClient[%s] start failed: %s", influxDbClientConfig.Name(), err)
		}
	}

	if countStarted < 1 {
		initiateShutdown <- errors.New("no InfluxDb client was started")
	}

	return influxDbClientPoolInstance
}