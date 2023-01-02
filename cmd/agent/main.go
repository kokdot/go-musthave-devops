package main

import (
	// "fmt"
	// "encoding/json"
	"fmt"
	"math/rand"
	// "runtime"
	"sync"
	"time"

	"github.com/kokdot/go-musthave-devops/internal/metrics"
	"github.com/kokdot/go-musthave-devops/internal/def"
	"github.com/kokdot/go-musthave-devops/internal/monitor"
	"github.com/kokdot/go-musthave-devops/internal/onboarding"
	// "github.com/go-resty/resty/v2"
)

type Gauge = def.Gauge
type Counter = def.Counter
type MonitorMap = monitor.MonitorMap
// type MonitorMap map[string] Gauge

var (
	pollCount Counter
	randomValue Gauge 
	m monitor.MonitorMap
	wg sync.WaitGroup 
	// urlReal string
) 


func main() {
	wg.Add(2)
	onboarding.Onboarding()
	m = monitor.GetData()
	go func(m *MonitorMap) {
		defer wg.Done()
		for {
			<-time.After(onboarding.PollIntervalReal)
			//, mutex)
			pollCount++
			randomValue = Gauge(rand.Float64())
		}
	}(&m)
	
	go func() {
		defer wg.Done()
		// var err error
		for {
			<-time.After(onboarding.ReportInterval)
			mtxCounter, err := metrics.NewMetricsCounter("PollCount", &pollCount, onboarding.UrlReal)
			if err != nil {
				fmt.Println(err)
			}
			mtxCounter.Update()
			mtxRandomValue, err := metrics.NewMetricsGauge("RandomValue", &randomValue, onboarding.UrlReal)
			if err != nil {
				fmt.Println(err)
			}
			mtxRandomValue.Update()
			for key, val := range m {
				// val1 := metrics.Gauge(val)
				mtx, err := metrics.NewMetricsGauge(key, &val, onboarding.UrlReal) 
				if err != nil {
					fmt.Println(err)
				}
				mtx.Update()
			}
		}
	}()
	wg.Wait()
}
