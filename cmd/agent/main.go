package main

import (
	// "bytes"
	"fmt"
	// "log"
	"encoding/json"
	"github.com/caarlos0/env/v6"
	// "net/http"
	"runtime"
	"flag"
	"sync"
	"time"
	"math/rand"
	// "io"
	"github.com/go-resty/resty/v2"
)

const (
	Url            = "127.0.0.1:8080"
	PollInterval time.Duration  = 2
	ReportInterval time.Duration = 10
)
type Config struct {
    Address  string 		`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
    ReportInterval time.Duration	 `env:"REPORT_INTERVAL" envDefault:"10s"`
    PollInterval time.Duration	 `env:"POLL_INTERVAL" envDefault:"2s"`
}

var wg sync.WaitGroup 

type Gauge float64
type Counter int64
type MonitorMap map[string]Gauge
var PollCount Counter
var RandomValue Gauge

var( 
	pollIntervalReal = PollInterval
	reportIntervalReal = ReportInterval
	urlReal = Url
	cfg Config

	// urlReal = "http://" + Url
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
func NewMonitor(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
	// runtime.ReadMemStats(&rtm)
	// fmt.Println(rtm)
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
func onboarding() {

    err := env.Parse(&cfg)
    if err != nil {
        fmt.Println(err)
    }
	

	urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    reportIntervalRealPtr := flag.Duration("r", 10, "interval of perort")
    pollIntervalRealPtr := flag.Duration("p", 2, "interval of poll")

    flag.Parse()
	urlReal = *urlRealPtr
	reportIntervalReal = *reportIntervalRealPtr
	pollIntervalReal = *pollIntervalRealPtr

	urlReal	= cfg.Address
	reportIntervalReal	= cfg.ReportInterval
	pollIntervalReal	= cfg.PollInterval
	urlReal = "http://" + urlReal
	fmt.Println("urlReal:     ", urlReal)
	fmt.Println("reportIntervalReal:     ", reportIntervalReal)
	fmt.Println("pollIntervalReal:     ", pollIntervalReal)

}
func mtxCounterSet(id string, counterPtr *Counter) ([]byte, error) {
	// fmt.Println("---------mtxCounterSet-----------------id--", id, "*counterPtr", *counterPtr)

	var varMetrics Metrics = Metrics{
			ID: id,
			MType: "Counter",
			Delta: counterPtr,
		}
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		fmt.Printf("Failed marshal json: %s\n", err)
		return nil, err
	}
	return bodyBytes, nil
}
func mtxGaugeSet(id string, gaugePtr *Gauge) ([]byte, error) {
	// fmt.Println("---------mtxGaugeSet-----------------id--", id, "*gaugePtr", *gaugePtr)

	var varMetrics Metrics = Metrics{
			ID: id,
			MType: "Gauge",
			Value: gaugePtr,
		}
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		fmt.Printf("Failed marshal json: %s\n", err)
		return nil, err
	}
	return bodyBytes, nil
}
// func mtxGaugeGet(id string) ([]byte, error) {
// 	var varMetrics Metrics = Metrics{
// 			ID: id,
// 			MType: "Gauge",
// 		}
// 	bodyBytes, err := json.Marshal(varMetrics)
// 	if err != nil {
// 		fmt.Printf("Failed marshal json: %s\n", err)
// 		return nil, err
// 	}
// 	return bodyBytes, nil
// }
func main() {
	wg.Add(2)
	onboarding()

	var rtm runtime.MemStats
	var m = make(MonitorMap)
	go func(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
		defer wg.Done()

		// var interval = pollIntervalReal
		// var interval = time.Duration(pollIntervalReal) * time.Second
		// fmt.Println("interval MonitorMap:        ", interval)

		for {
			<-time.After(pollIntervalReal)
			runtime.ReadMemStats(&rtm)
			NewMonitor(m, rtm)//, mutex)
			PollCount++
			RandomValue = Gauge(rand.Float64())
			// fmt.Println(m)
		}
	}(&m, rtm)
	
	
	go func() {
		defer wg.Done()
		// var interval = reportIntervalReal
		// fmt.Println("interval PollCount:        ", interval)
		for {

			<-time.After(reportIntervalReal) 
			
			//PollCount----------------------------------------------------------
			strURL := fmt.Sprintf("%s/update/", urlReal)
			// strURLGet := fmt.Sprintf("%s/value/", urlReal)
			var varMetrics Metrics
			var bodyBytes []byte
			var err error
			bodyBytes, err = mtxCounterSet("PollCount", &PollCount)
			if err != nil {
				fmt.Println(err)
			}
			client := resty.New()
			_, err = client.R().
			SetResult(&varMetrics).
			SetBody(bodyBytes).
			Post(strURL)
			if err != nil {
				fmt.Printf("Failed unmarshall response PollCount: %s\n", err)
			}
			// fmt.Println("varMetrics: ", varMetrics) 
			// fmt.Println("PollCount: ", *varMetrics.Delta) 

			//RandomValue----------------------------------------------------------
			bodyBytes, err = mtxGaugeSet("RandomValue", &RandomValue)
			if err != nil {
				fmt.Println(err)
			}
			client = resty.New()
			_, err = client.R().
			SetResult(&varMetrics).
			SetBody(bodyBytes).
			Post(strURL)
			if err != nil {
				fmt.Printf("Failed unmarshall response RandomValue: %s\n", err)
			}
			// fmt.Println("varMetrics: ", varMetrics) 
			// fmt.Println("RandomValue: ", *varMetrics.Value) 

				//RandomValueGet---------------------------------------------------
			// strURLGet := fmt.Sprintf("%s/value/", urlReal)
			// var metricsStructGet Metrics
			// // randomValue := float64(RandomValue)
			// // varMetrics = Metrics{
			// 	// 	ID: "RandomValue",
			// 	// 	MType: "Gauge",
			// 	// }
			// 	bodyBytes, err := mtxGaugeGet("RandomValue")
			// 	// bodyBytes, err = json.Marshal(varMetrics)
			// 	// if err != nil {
			// 	// 	fmt.Printf("Failed marshal json: %s\n", err)
			// 	// }
			// 	// var varMetrics1 Metrics
			// 	client := resty.New()
			// 	_, err = client.R().
			// SetResult(&metricsStructGet).
			// SetBody(bodyBytes).
			// Post(strURLGet)
			// if err != nil {
			// 	log.Printf("Failed unmarshall response: %s\n", err)
			// }
			// fmt.Println("RandomValueGet:  ", *metricsStructGet.Value) 

			// Gauge ----------------------------------------------------------
			//n := 0
			for key, val := range m {
				// n++
				// if n > 1 {
				// 	break
				// }
				// fmt.Println("key: ", key, ";  val: ", val)
				bodyBytes, err = mtxGaugeSet(key, &val)
				if err != nil {
					fmt.Println(err)
				}
				client := resty.New()
				_, err = client.R().
				SetResult(&varMetrics).
				SetBody(bodyBytes).
				Post(strURL)
				if err != nil {
					fmt.Printf("Failed unmarshall response Monitor, ID: %s; error: %s\n", key, err)
				}
				// fmt.Println("-----Update;------- Id: ", varMetrics.ID, "Value: ", *varMetrics.Value) 
			}
			// Gauge ------Get----------------------------------------------------
			// n := 0
			// fmt.Println("---------------------------------------------")
			// for key, _ := range m {
			// 	// n++
			// 	// if n > 3 {
			// 	// 	break
			// 	// }
			// 	// fmt.Println("key: ", key, ";  val: ", val)
			// 	bodyBytes, err = mtxGaugeGet(key)
			// 	if err != nil {
			// 		fmt.Println(err)
			// 	}
			// 	// _ = bodyBytes
			// 	client := resty.New()
			// 	_, err = client.R().
			// 	SetResult(&varMetrics).
			// 	SetBody(bodyBytes).
			// 	Post(strURLGet)
			// 	if err != nil {
			// 		fmt.Printf("Failed unmarshall response: %s\n", err)
			// 	}
			// 	if varMetrics.Value == nil {
			// 		fmt.Println("-----Get;------- Id: ", varMetrics.ID, "Value is nil: ", varMetrics) 

			// 	} else {
			// 		fmt.Println("-----Get;------- Id: ", varMetrics.ID, "Value: ", *varMetrics.Value) 

			// 	}
			// 	// fmt.Println("-----Get;------- Id: ", varMetrics.ID, "Value: ", *varMetrics.Value) 
			// }

			
		}
	}()
	wg.Wait()
}
