package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/kokdot/go-musthave-devops/internal/metricsagent"
	"github.com/kokdot/go-musthave-devops/internal/def"
	"github.com/kokdot/go-musthave-devops/internal/monitor"
	"github.com/kokdot/go-musthave-devops/internal/onboardingagent"
)

type Gauge = def.Gauge
type Counter = def.Counter
type MonitorMap = def.MonitorMap

var (
	pollCount Counter
	randomValue Gauge 
	m MonitorMap
	wg sync.WaitGroup 
) 

func main() {
	wg.Add(2)
	onboardingagent.OnboardingAgent()
	m = make(def.MonitorMap)
	go func(m *MonitorMap) {
		defer wg.Done()
		for {
			<-time.After(onboardingagent.PollIntervalReal)
			m = monitor.GetData(m)
			pollCount++
			randomValue = Gauge(rand.Float64())
		}
	}(&m)
	
	go func() {
		defer wg.Done()
		for {
			<-time.After(onboardingagent.ReportInterval)
			mtxCounter, err := metricsagent.NewMetricsCounter("PollCount", &pollCount, onboardingagent.URLReal)
			// fmt.Printf("mtxRandomValue:    %#v\n", mtxCounter)
			if err != nil {
				fmt.Println(err)
			}
			mtxCounter.Update()
			mtxRandomValue, err := metricsagent.NewMetricsGauge("RandomValue", &randomValue, onboardingagent.URLReal)
			// fmt.Printf("mtxRandomValue:    %#v\n", mtxRandomValue)
			if err != nil {
				fmt.Println(err)
			}
			mtxRandomValue.Update()
			for key, val := range m {
				mtx, err := metricsagent.NewMetricsGauge(key, &val, onboardingagent.URLReal) 
				// fmt.Printf("mtxRandomValue:    %#v\n", mtx)
				if err != nil {
					fmt.Println(err)
				}
				mtx.Update()
			}
		}
	}()
	wg.Wait()
}
