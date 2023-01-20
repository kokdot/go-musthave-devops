package main

import (
	"fmt"
	// "os"
	"math/rand"
	"sync"
	"time"
	// "github.com/rs/zerolog"

	"github.com/kokdot/go-musthave-devops/internal/metricsagent"
	"github.com/kokdot/go-musthave-devops/internal/def"
	"github.com/kokdot/go-musthave-devops/internal/monitor"
	"github.com/kokdot/go-musthave-devops/internal/onboardingagent"
)

type Gauge = def.Gauge
type Counter = def.Counter
type MonitorMap = def.MonitorMap

var (
	pollInterval time.Duration
	reportInterval time.Duration
	url string
	key string
	batch bool
	pollCount Counter
	randomValue Gauge 
	m MonitorMap
	wg sync.WaitGroup 
) 
// func (mtx metricsagent.Metrics) MarshalZerologObject(e *zerolog.Event) {
// 	e.Str("ID", mtx.ID).
// 	Str("MType", mtx.MType).
// 	Str("MType", mtx.Hash).
// 	Str("MType", mtx.Key).
// 	Str("MType", mtx.StrURL).
// 	Float64("MType", mtx.Value).
// 	int64("MType", mtx.Delta)
// }

func main() {
	wg.Add(2)
	pollInterval, reportInterval, url, key, batch = onboardingagent.OnboardingAgent()
	m = make(def.MonitorMap)
	// sm := make(metricsagent.StoreMap)
	// logMetrics := zerolog.New(os.Stdout).With().
	// 	Str("foo", "bar").
	// 	Object("user", u).
	// 	Logger()

	// log.Log().Msg("hello world")
	go func(m *MonitorMap) {
		defer wg.Done()
		for {
			<-time.After(pollInterval)
			m = monitor.GetData(m)
			pollCount++
			randomValue = Gauge(rand.Float64())
			// sm = *metricsagent.GetStoreMap(&sm)
		}
	}(&m)
	
	go func() {
		defer wg.Done()
		for {
			<-time.After(reportInterval)
			if batch {
				err := metricsagent.UpdateByBatch(&m, pollCount, randomValue, url, key)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Get response by batch request.")
				}
			} else {
				err := metricsagent.UpdateAll(&m, pollCount, randomValue, url, key)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Get response by common request.")
				}
			}
		}
	}()
	wg.Wait()
}
