package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	url            = "http://127.0.0.1:8080"
	pollInterval   = 2
	reportInterval = 10
)

var mutex *sync.RWMutex

type Guage float64
type Couter int64
type MonitorMap map[string]Guage

func NewMonitor(m *MonitorMap, rtm *runtime.MemStats, mutex *sync.RWMutex) {
	runtime.ReadMemStats(rtm)
	mutex.Lock()
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
	mutex.Unlock()
}

func main() {
	var rtm *runtime.MemStats
	var m = make(MonitorMap)
	go func(m *MonitorMap, rtm *runtime.MemStats, mutex *sync.RWMutex) {
		var interval = time.Duration(pollInterval) * time.Second
		for {
			<-time.After(interval)
			NewMonitor(m, rtm, mutex)
		}
	}(&m, rtm, mutex)
	go func(mutex *sync.RWMutex) {
		var interval = time.Duration(reportInterval) * time.Second
		for {
			<-time.After(interval)
			for key, val := range m {
				client := &http.Client{}
				mutex.RLock()
				strUrl := fmt.Sprintf("%s/update/Guage/%s/%v", url, key, val)
				mutex.Unlock()
				_, err := client.Post(strUrl, "text/plain", bytes.NewBufferString(""))
				if err != nil {
					log.Fatalf("Failed sent request: %s", err)
				}

			}
		}
	}(mutex)
}
