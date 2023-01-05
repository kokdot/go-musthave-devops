package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/kokdot/go-musthave-devops/internal/metrics"
	"github.com/kokdot/go-musthave-devops/internal/def"
	"github.com/kokdot/go-musthave-devops/internal/monitor"
	"github.com/kokdot/go-musthave-devops/internal/onboarding_agent"
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
	onboarding_agent.OnboardingAgent()
	m = make(def.MonitorMap)
	go func(m *MonitorMap) {
		defer wg.Done()
		for {
			<-time.After(onboarding_agent.PollIntervalReal)
			m = monitor.GetData(m)
			pollCount++
			randomValue = Gauge(rand.Float64())
		}
	}(&m)
	
	go func() {
		defer wg.Done()
		for {
			<-time.After(onboarding_agent.ReportInterval)
			mtxCounter, err := metrics.NewMetricsCounter("PollCount", &pollCount, onboarding_agent.URLReal)
			fmt.Printf("mtxRandomValue:    %#v\n", mtxCounter)
			if err != nil {
				fmt.Println(err)
			}
			mtxCounter.Update()
			mtxRandomValue, err := metrics.NewMetricsGauge("RandomValue", &randomValue, onboarding_agent.URLReal)
			fmt.Printf("mtxRandomValue:    %#v\n", mtxRandomValue)
			if err != nil {
				fmt.Println(err)
			}
			mtxRandomValue.Update()
			for key, val := range m {
				mtx, err := metrics.NewMetricsGauge(key, &val, onboarding_agent.URLReal) 
				fmt.Printf("mtxRandomValue:    %#v\n", mtx)
				if err != nil {
					fmt.Println(err)
				}
				mtx.Update()
			}
		}
	}()
	wg.Wait()
}
