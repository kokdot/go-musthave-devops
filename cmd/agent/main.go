package main

import (
	// "fmt"
	// "os"
	"math/rand"
	"sync"
	"time"

	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"

	"github.com/kokdot/go-musthave-devops/internal/def"
	"github.com/kokdot/go-musthave-devops/internal/metricsagent"
	"github.com/kokdot/go-musthave-devops/internal/monitor"
	"github.com/kokdot/go-musthave-devops/internal/onboardingagent"
	"github.com/kokdot/go-musthave-devops/internal/virtualmemory"
)

type Gauge = def.Gauge
type Counter = def.Counter
type GaugeMap = def.GaugeMap

var (
	// pollInterval time.Duration
	// reportInterval time.Duration
	// url string
	// key string
	// batch bool
	pollCount Counter
	randomValue Gauge 
	// m MonitorMap
	// vm def.VirtualMemoryMap
	wg sync.WaitGroup 
	
) 
// var conf *def.Conf
// var logg = log.Logger
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
	wg.Add(3)
	conf := onboardingagent.OnboardingAgent()
	logg := conf.Logg
	// pollInterval, reportInterval, url, key, batch, logg = onboardingagent.OnboardingAgent()
	metricsagent.GetConf(conf)
	gm := make(def.GaugeMap)
	vm := make(def.GaugeMap)
	// sm := make(metricsagent.StoreMap)
	// logMetrics := zerolog.New(os.Stdout).With().
	// 	Str("foo", "bar").
	// 	Object("user", u).
	// 	Logger()

	// log.Log().Msg("hello world")
	go func (vm *def.GaugeMap) {
		defer wg.Done()
		for {
			<-time.After(conf.PollInterval)
			vm = virtualmemory.GetData(vm)
		}
	}(&vm)

	go func(gm *GaugeMap) {
		defer wg.Done()
		for {
			<-time.After(conf.PollInterval)
			gm = monitor.GetData(gm)
			pollCount++
			randomValue = Gauge(rand.Float64())
			// sm = *metricsagent.GetStoreMap(&sm)
		}
	}(&gm)
	
	go func() {
		logg.Print("main is going to send the report.--------------------")
		defer wg.Done()
		for {
			<-time.After(conf.ReportInterval)
			for k, v := range vm {
				gm[k] = v
			}
			gm["randomValue"] = randomValue 
			if conf.Batch {
				err := metricsagent.UpdateByBatch(&gm, pollCount)
				// err := metricsagent.UpdateByBatch(&m, &vm, pollCount, randomValue, url, key)
				if err != nil {
					logg.Error().Err(err).Send()
				} else {
					logg.Print("Get response by batch request.")
				}
			} else {
				err := metricsagent.UpdateAll(&gm, pollCount)
				// err := metricsagent.UpdateAll(&m, &vm, pollCount, randomValue, url, key)
				if err != nil {
					logg.Error().Err(err).Send()
				} else {
					logg.Print("Get response by common request.")
				}
			}
		}
	}()
	wg.Wait()
}
