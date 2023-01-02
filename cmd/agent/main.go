package main

import (
	"fmt"
	"encoding/json"
	"runtime"
	"sync"
	"time"
	"math/rand"
	"github.com/kokdot/go-musthave-devops/internal/metrics"
 	"github.com/go-resty/resty/v2"
)




var wg sync.WaitGroup 


type MonitorMap map[string]Gauge
var PollCount Counter
var RandomValue Gauge 


func main() {
	wg.Add(2)
	onboarding()
	var rtm runtime.MemStats
	var m = make(MonitorMap)
	go func(m *MonitorMap, rtm runtime.MemStats) {
		defer wg.Done()
		for {
			<-time.After(pollIntervalReal)
			runtime.ReadMemStats(&rtm)
			NewMonitor(m, rtm)//, mutex)
			PollCount++
			RandomValue = Gauge(rand.Float64())
		}
	}(&m, rtm)
	
	go func() {
		defer wg.Done()
		mtxCounter := NewMetricsCounter("pollCount", pollCount, urlReal)
		mtxCounter.Update()
		mtxRandomValue := NewMetricsGauge("randomValue", randomValue, urlReal)
		mtxRandomValue.Update()
		for key, val := range m {
			mtx := NewMetricsGauge(key, &val, urlReal) 
			mtx.Update()
		}
		
	}()
	wg.Wait()
}

// func (mtx *Metrics) CounterUpdate() error {
// 	// var varMetrics Metrics
// 	// var bodyBytes []byte
// 	var err error
// 	// bodyBytes, err = mtxCounterSet("PollCount", &PollCount)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	client := resty.New()
// 	_, err = client.R().
// 	SetHeader("Accept-Encoding", "gzip").
// 	SetResult(mtx).
// 	SetBody(mtx.bodyBytes).
// 	Post(mtx.strURL)

// 	// client := resty.New()
// 	// _, err = client.R().
// 	// SetHeader("Accept-Encoding", "gzip").
// 	// SetResult(&varMetrics).
// 	// SetBody(bodyBytes).
// 	// Post(strURL)
// 	fmt.Println("--------------------------new--------------------------------------")
// 	if err != nil {
// 		fmt.Printf("Failed unmarshall response PollCount: %s\n", err)
// 	}

// }
// func (mtx *Metrics) GaugeUpdate() {
// 	var varMetrics Metrics
// 	var bodyBytes []byte
// 	var err error
// 	for {
// 			<-time.After(reportIntervalReal) 

			
			
			

// 			for key, val := range m {
// 				bodyBytes, err = mtxGaugeSet(key, &val)
// 				if err != nil {
// 					fmt.Println(err)
// 				}
// 				client := resty.New()
// 				_, err = client.R().
// 				SetResult(&varMetrics).
// 				SetBody(bodyBytes).
// 				Post(strURL)
// 				if err != nil {
// 					fmt.Printf("Failed unmarshall response Monitor, ID: %s; error: %s\n", key, err)
// 				}
// 			}
// 		}
// }