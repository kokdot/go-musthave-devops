package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
	"io"
)

const (
	url            = "http://127.0.0.1:8080"
	pollInterval   = 2
	reportInterval = 4
)

// var mutex *sync.RWMutex
var wg sync.WaitGroup 

type Gauge float64
type Couter int64
type MonitorMap map[string]Gauge
var PollCount int


func NewMonitor(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
	runtime.ReadMemStats(&rtm)
	fmt.Println(rtm)
	// mutex.Lock()
	(*m)["Alloc"] = Gauge(rtm.Alloc)
	(*m)["BuckHashSys"] = Gauge(rtm.BuckHashSys)
	(*m)["TotalAlloc"] = Gauge(rtm.TotalAlloc)
	(*m)["Sys"] = Gauge(rtm.Sys)
	(*m)["Mallocs"] = Gauge(rtm.Mallocs)
	(*m)["Frees"] = Gauge(rtm.Frees)
	(*m)["PauseTotalNs"] = Gauge(rtm.PauseTotalNs)
	(*m)["NumGC"] = Gauge(rtm.NumGC)
	(*m)["GCCPUFraction"] = Gauge(rtm.GCCPUFraction)
	(*m)["GCSys"] = Gauge(rtm.GCSys)
	(*m)["HeapInuse"] = Gauge(rtm.HeapInuse)
	(*m)["HeapObjects"] = Gauge(rtm.HeapObjects)
	(*m)["HeapReleased"] = Gauge(rtm.HeapReleased)
	(*m)["HeapSys"] = Gauge(rtm.HeapSys)
	(*m)["LastGC"] = Gauge(rtm.LastGC)
	(*m)["MSpanInuse"] = Gauge(rtm.MSpanInuse)
	(*m)["MCacheSys"] = Gauge(rtm.MCacheSys)
	(*m)["MCacheInuse"] = Gauge(rtm.MCacheInuse)
	(*m)["MSpanSys"] = Gauge(rtm.MSpanSys)
	(*m)["NextGC"] = Gauge(rtm.NextGC)
	(*m)["NumForcedGC"] = Gauge(rtm.NumForcedGC)
	(*m)["OtherSys"] = Gauge(rtm.OtherSys)
	(*m)["StackSys"] = Gauge(rtm.StackSys)
	(*m)["StackInuse"] = Gauge(rtm.StackInuse)
	(*m)["TotalAlloc"] = Gauge(rtm.TotalAlloc)
	// mutex.Unlock()
}

func main() {
	wg.Add(2)
	var rtm runtime.MemStats
	var m = make(MonitorMap)
	// NewMonitor(&m, rtm)
	go func(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
		defer wg.Done()
		var interval = time.Duration(pollInterval) * time.Second
		for {
			<-time.After(interval)
			NewMonitor(m, rtm)//, mutex)
			fmt.Println(m)
		}
	}(&m, rtm)//, mutex)
	go func() {//mutex *sync.RWMutex) {
		defer wg.Done()
		var interval = time.Duration(reportInterval) * time.Second
		for {
			<-time.After(interval)

			client := &http.Client{}
			strUrl := fmt.Sprintf("%s/update/counter/%s/%v", url, "PollCount", PollCount)
			fmt.Println("strUrl:  --  ", strUrl)
			response, err := client.Post(strUrl, "text/plain", bytes.NewBufferString(""))
			if err != nil {
				log.Fatalf("Failed sent request: %s", err)
			}
			bodyBytes, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			response.Body.Close()
			bodyString := string(bodyBytes)
			fmt.Println("response.Body: ", bodyString)

			for key, val := range m {
				client := &http.Client{}
				// mutex.RLock()
				strUrl := fmt.Sprintf("%s/update/Gauge/%s/%v", url, key, val)
				// mutex.Unlock()
				fmt.Println("strUrl:  --  ", strUrl)
				response, err := client.Post(strUrl, "text/plain", bytes.NewBufferString(""))
				if err != nil {
					log.Fatalf("Failed sent request: %s", err)
				}
				bodyBytes, err := io.ReadAll(response.Body)
				if err != nil {
					log.Fatal(err)
				}
				response.Body.Close()
				bodyString := string(bodyBytes)
				fmt.Println("response.Body: ", bodyString)

			}
		}
	}()//mutex)
	wg.Wait()
}
