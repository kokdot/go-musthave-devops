package main

import (
	// "bytes"
	"fmt"
	"log"
	// "net/http"
	"runtime"
	"sync"
	"time"
	// "io"
	"github.com/go-resty/resty/v2"
)

const (
	url            = "http://127.0.0.1:8080"
	pollInterval   = 2
	reportInterval = 4
)

// var mutex *sync.RWMutex
var wg sync.WaitGroup 

type Guage float64
type Couter int64
type MonitorMap map[string]Guage
var PollCount int
func NewMonitor(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
	runtime.ReadMemStats(&rtm)
	// fmt.Println(rtm)
	// mutex.Lock()
	(*m)["Alloc"] = Guage(rtm.Alloc)
	(*m)["BuckHashSys"] = Guage(rtm.BuckHashSys)
	(*m)["TotalAlloc"] = Guage(rtm.TotalAlloc)
	(*m)["Sys"] = Guage(rtm.Sys)
	(*m)["Mallocs"] = Guage(rtm.Mallocs)
	(*m)["Frees"] = Guage(rtm.Frees)
	(*m)["PauseTotalNs"] = Guage(rtm.PauseTotalNs)
	(*m)["NumGC"] = Guage(rtm.NumGC)
	(*m)["GCCPUFraction"] = Guage(rtm.GCCPUFraction)
	(*m)["GCSys"] = Guage(rtm.GCSys)
	(*m)["HeapInuse"] = Guage(rtm.HeapInuse)
	(*m)["HeapObjects"] = Guage(rtm.HeapObjects)
	(*m)["HeapReleased"] = Guage(rtm.HeapReleased)
	(*m)["HeapSys"] = Guage(rtm.HeapSys)
	(*m)["LastGC"] = Guage(rtm.LastGC)
	(*m)["MSpanInuse"] = Guage(rtm.MSpanInuse)
	(*m)["MCacheSys"] = Guage(rtm.MCacheSys)
	(*m)["MCacheInuse"] = Guage(rtm.MCacheInuse)
	(*m)["MSpanSys"] = Guage(rtm.MSpanSys)
	(*m)["NextGC"] = Guage(rtm.NextGC)
	(*m)["NumForcedGC"] = Guage(rtm.NumForcedGC)
	(*m)["OtherSys"] = Guage(rtm.OtherSys)
	(*m)["StackSys"] = Guage(rtm.StackSys)
	(*m)["StackInuse"] = Guage(rtm.StackInuse)
	(*m)["TotalAlloc"] = Guage(rtm.TotalAlloc)
	// mutex.Unlock()
}
func main() {
	wg.Add(2)
	var rtm runtime.MemStats
	var m = make(MonitorMap)
	go func(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
		defer wg.Done()

		var interval = time.Duration(pollInterval) * time.Second
		for {
			<-time.After(interval)
			NewMonitor(m, rtm)//, mutex)
			PollCount++
			// fmt.Println(m)
		}
	}(&m, rtm)
	
	
	go func() {
		defer wg.Done()
		var interval = time.Duration(reportInterval) * time.Second
		// PollCount35 := 35
		for {

			<-time.After(interval) 
			strURL := fmt.Sprintf("%s/update/counter/%s/%v", url, "PollCount", PollCount)
			client := resty.New()
			response, err := client.R().Post(strURL)
			if err != nil {
				log.Fatalf("Failed sent request: %s", err)
			}
			fmt.Println(response) 
			// strURL = fmt.Sprintf("%s/update/counter/%s/%v", url, "PollCount35", PollCount35)
			// client = resty.New()
			// response, err = client.R().Post(strURL)
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// PollCount35++
			// fmt.Println(response) 
			// n := 0
			for key, val := range m {
				// n++
				// if n > 1 {
				// 	break
				// }
				strURL := fmt.Sprintf("%s/update/gauge/%s/%v", url, key, val)
				client := resty.New()
				response, err := client.R().Post(strURL)
				if err != nil {
					log.Fatalf("Failed sent request: %s", err)
				}
				fmt.Println(response)
			}
		}
	}()
	wg.Wait()
}
