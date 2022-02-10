package exporter

import (
	"context"
	"fmt"
	"github.com/patrick246/blahaj-exporter/internal/client"
	"github.com/patrick246/blahaj-exporter/internal/stores"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"sync"
	"time"
)

const blahajArtNo = "30373588"

type Exporter struct {
	stockClient *client.StockLevelClient
	descs       []*prometheus.Desc
	log         *zap.SugaredLogger
}

func New(stockClient *client.StockLevelClient, log *zap.SugaredLogger) *Exporter {
	return &Exporter{
		stockClient: stockClient,
		descs: []*prometheus.Desc{
			prometheus.NewDesc("ikea_info", "Informational metric about IKEA stores", []string{"store", "name", "lat", "lon", "country"}, nil),
			prometheus.NewDesc("ikea_blahaj_count", "Number of Blahajs in stock", []string{"store"}, nil),
		},
		log: log,
	}
}

func (e *Exporter) Describe(descs chan<- *prometheus.Desc) {
	for _, desc := range e.descs {
		descs <- desc
	}
}

func (e *Exporter) Collect(metrics chan<- prometheus.Metric) {
	storeList, err := stores.GetStores()
	if err != nil {
		e.log.Warnw("store error", "error", err)
		metrics <- prometheus.NewInvalidMetric(e.descs[0], err)
	} else {
		const numCoordinates = 2
		for _, store := range storeList {
			if len(store.Coordinates) != numCoordinates {
				metrics <- prometheus.MustNewConstMetric(
					e.descs[0],
					prometheus.GaugeValue,
					1,
					store.ID,
					store.Name,
					"n/a",
					"n/a",
					store.CountryCode,
				)
			} else {
				metrics <- prometheus.MustNewConstMetric(
					e.descs[0],
					prometheus.GaugeValue,
					1,
					store.ID,
					store.Name,
					fmt.Sprintf("%f", store.Coordinates[0]),
					fmt.Sprintf("%f", store.Coordinates[1]),
					store.CountryCode,
				)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) //nolint:gomnd // half a minute
	defer cancel()

	countryCodes, err := stores.GetCountryCodes()
	if err != nil {
		e.log.Warnw("country code error", "error", err)
		metrics <- prometheus.NewInvalidMetric(e.descs[1], err)
		return
	}

	const parallelWorkers = 4

	semaphore := make(chan struct{}, parallelWorkers)
	wg := sync.WaitGroup{}
	for _, country := range countryCodes {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(country string) {
			response, err := e.stockClient.StockLevels(ctx, country, blahajArtNo)
			if err != nil {
				e.log.Warnw("stock error", "error", err)
			} else {
				for i := range response.Availabilities {
					metrics <- prometheus.MustNewConstMetric(
						e.descs[1],
						prometheus.GaugeValue,
						float64(response.Availabilities[i].BuyingOption.CashCarry.Availability.Quantity),
						response.Availabilities[i].ClassUnitKey.ClassUnitCode,
					)
				}
			}
			<-semaphore
			wg.Done()
		}(country)
	}
	wg.Wait()
}
