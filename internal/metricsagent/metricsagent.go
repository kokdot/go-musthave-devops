package metricsagent

import (
 	"github.com/go-resty/resty/v2"
	// "github.com/kokdot/go-musthave-devops/internal/onboardingagent"
	"github.com/kokdot/go-musthave-devops/internal/def"
	"github.com/rs/zerolog"
	// "github.com/kokdot/go-musthave-devops/internal/monitor"
	"encoding/json"
	"fmt"
	// "os"
	"crypto/sha256"
	"crypto/hmac"
    // "github.com/rs/zerolog/log"

)

type Gauge = def.Gauge 
type Counter = def.Counter
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
	BodyBytes []byte `json:"-"`
	StrURL string `json:"-"`
	Key string `json:"-"`
}
type SliceMetrics []Metrics
type StoreMap map[string] Metrics
var StoreMetrics = make(map[string]interface{})
var logg zerolog.Logger// = onboardingagent.GetLogg()
// var metrics = Metrics{}
// var sliceMetrics =  make(SliceMetrics, 0)
// var logStoreMetrics = zerolog.New(os.Stdout).With().
// 		Str("agent", "metricsagent").
// 		Fields(StoreMetrics).
// 		Logger()

// var logStoreMetrics = zerolog.New(os.Stdout).With().
// 		Str("foo", "bar").
// 		Array("maetrics", u).
// 		Logger()

// func (mtx Metrics) MarshalZerologObject(e *zerolog.Event) {
// 	var value Gauge
// 	var delta Counter
// 	if mtx.Value != nil {
// 		value = *(mtx.Value)
// 		// fmt.Println("value: ", value, "  int64(value): ", float64(value))
// 	}
// 	if mtx.Delta != nil {
// 		delta = *(mtx.Delta)
// 		// fmt.Println("delta: ", delta, "  int64(delta): ", int64(delta))
// 	}
// 	// if mtx.ID == "" {
// 		// 	return
// 		// }
// 		e.Str("ID", mtx.ID).
// 		Str("MType", mtx.MType).
// 		Float64("Value", float64(value)).
// 		Int64("Delta", int64(delta)).
// 		Str("URL", mtx.StrURL).
// 		Str("Key", mtx.Key)
		
		

// 	// Str("MType", mtx.Hash).
// 	// Int64("MType", int64(*mtx.Delta))
// }


// func (sliceMetrics SliceMetrics) MarshalZerologArray(a *zerolog.Array) {
// 	for _, m := range sliceMetrics {
// 		a.Object(m)
// 	}
// }
// var logStoreMetrics = zerolog.New(os.Stdout).With().
// 		Str("agent", "metricsagent").
// 		Array("SliceMetrics", sliceMetrics).
// 		Logger()
func GetLogg(loggReal zerolog.Logger)  {
	logg = loggReal
	// return logg
}
// logg.Error().Err(err).Send()
func GetStoreMap(mPtr *def.MonitorMap, url string, key string, mtxPollCountAndRandomValue... Metrics) (*StoreMap, error) {
	logg.Print("mPtr:  ", *mPtr)
	sm := make(StoreMap, 0)
	var mtx = Metrics{}
	var err  error
	for k, v := range *mPtr {
		mtx, err = NewMetricsGauge(k, v, url, key)
		if err != nil {
			logg.Error().Err(err).Send()
			return nil, fmt.Errorf("%s", err)
		}
		logg.Print("mtx.ID: ", mtx.ID)
		logg.Print("mtx.Value: ", *mtx.Value)
		sm[k] = mtx
	}

	for _, mtx1 := range mtxPollCountAndRandomValue {
		sm[mtx1.ID] = mtx1
		logg.Print("mtx.ID: ", mtx1.ID)
	}

	// sliceMetrics := make(SliceMetrics, 0)
	// for _, v := range sm {
	// 	sliceMetrics = append(sliceMetrics, v)
	// }
	
	// logStoreMetrics := zerolog.New(os.Stdout).With().
	// Str("--------------------agent", "metricsagent").
	// Array("SliceMetrics", sliceMetrics).
	// Logger()
	// logStoreMetrics.Log().Msg("GetStoreMap is finish")
	return &sm, nil
}

func UpdateByBatch(mPtr *def.MonitorMap, pollCount Counter, randomValue Gauge, url string, key string) error {
	mtxCounter, err := NewMetricsCounter("PollCount", pollCount, url, key)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	mtxRandomValue, err := NewMetricsGauge("RandomValue", randomValue, url, key)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	smPtr, err := GetStoreMap(mPtr, url, key, mtxCounter, mtxRandomValue)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	// for k, v := range *smPtr {
	// 	fmt.Println("-----------------------------------------   ", k)
	// 	if v.Value != nil {
	// 		fmt.Println("mtx.ID: ", v.ID)
	// 		fmt.Println("mtx.Value: ", *v.Value)
	// 	} else {
	// 		fmt.Println("mtx.ID: ", v.ID)
	// 		fmt.Println("mtx.Delta: ", *v.Delta)
	// 	}
	// }
	err = UpdateStoreMap(smPtr, url)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	return nil
}

func UpdateStoreMap(smPtrnew *StoreMap, url string) error {
	var err error
	var smOld = make(StoreMap, len(*smPtrnew))
	url = "http://" + url
	strURL := fmt.Sprintf("%s/updates1/", url)
	bodyBytes, err := json.Marshal(smPtrnew)
	if err != nil {
		logg.Error().Err(err).Msg("Failed marshal json for batch: ")
		return err
	}
	client := resty.New()
	_, err = client.R().
	SetHeader("Accept-Encoding", "gzip").
	SetHeader("Content-Type", "application/json").
	SetBody(bodyBytes).
	SetResult(&smOld).
	Post(strURL)
	if err != nil {
		logg.Error().Err(err).Msg("Failed unmarshall response by batch: ")
		// fmt.Printf("Failed unmarshall response by batch %#v: %s\n", smOld, err)
		return err
	}
	// fmt.Printf("Result of requets is: %#v\n", smOld)
	return nil
}

func UpdateAll (m *def.MonitorMap, c Counter, g Gauge, url string, key string) error {
	mtxCounter, err := NewMetricsCounter("PollCount", c, url, key)
	// fmt.Printf("mtxRandomValue:    %#v\n", mtxCounter)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	mtxCounter.Update(url)

	mtxRandomValue, err := NewMetricsGauge("RandomValue", g, url, key)
	// fmt.Printf("mtxRandomValue:    %#v\n", mtxRandomValue)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	mtxRandomValue.Update(url)
	// n := 0
	for k, v := range *m {
		// n++
		// if n > 1 {
		// 	break
		// }
		mtx, err := NewMetricsGauge(k, v, url, key) 
		logg.Print("------------------------------------UpdateAll---------------------------start---------------------")
		logg.Printf("mtx: %v; Value =  %v", mtx.ID, *mtx.Value)
		if err != nil {
			logg.Error().Err(err).Send()
			return err
		}
		mtx.Update(url)
	}
	return nil
}

func NewMetricsCounter(id string,  counter Counter, urlReal string, keyReal string) (Metrics, error) {
	key := []byte(keyReal)
	urlReal1 := "http://" + urlReal
	if keyReal == "" {

		var varMetrics = Metrics{
				ID: id,
				MType: "counter",
				Delta: &counter,
				Key: keyReal,
			}
		bodyBytes, err := json.Marshal(varMetrics)
		if err != nil {
			logg.Error().Err(err).Msg("Failed marshal json counter: ")
			// fmt.Printf("Failed marshal json counter:  %s\n", err)
			return Metrics{}, err
		}
		varMetrics.BodyBytes = bodyBytes
		strURL := fmt.Sprintf("%s/update/", urlReal1)
		varMetrics.StrURL = strURL
		return varMetrics, nil
	}
		src := []byte((fmt.Sprintf("%s:counter:%d", id, counter)))
	h := hmac.New(sha256.New, key)
    h.Write(src)
    dst := h.Sum(nil)
	var varMetrics = Metrics{
			ID: id,
			MType: "counter",
			Delta: &counter,
			Hash: fmt.Sprintf("%x", dst),
			Key: keyReal,
		}
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		fmt.Printf("Failed marshal json: %s", err)
		return Metrics{}, err
	}
	varMetrics.BodyBytes = bodyBytes
	strURL := fmt.Sprintf("%s/update/", urlReal1)
	varMetrics.StrURL = strURL
	return varMetrics, nil
}

func NewMetricsGauge(id string, gauge Gauge,  urlReal string, keyReal string) (Metrics, error) {
	// key := []byte(keyReal)
	// urlReal = "http://" + urlReal
	if keyReal == "" {

		var varMetrics = Metrics{
			ID: id,
			MType: "gauge",
			Value: &gauge,
			Key: keyReal,
		}
		bodyBytes, err := json.Marshal(varMetrics)
		if err != nil {
			logg.Error().Err(err).Msg("Failed marshal json gauge: ")
			// fmt.Printf("Failed marshal json gauge: %s\n", err)
			return Metrics{}, err
		}
		varMetrics.BodyBytes = bodyBytes
		// strURL := fmt.Sprintf("%s/update/", urlReal)
		// varMetrics.strURL = strURL
		return varMetrics, nil
	}
	// src := []byte((fmt.Sprintf("%s:gauge:%f", id, float64(*gaugePtr))))
	// h := hmac.New(sha256.New, key)
	// h.Write(src)
	// dst := h.Sum(nil)
	var varMetrics = Metrics{
		ID: id,
		MType: "gauge",
		Value: &gauge,
		Key: keyReal,
		// Hash: fmt.Sprintf("%x", dst),
	}
	// hash := Hash(&varMetrics, keyReal)
	varMetrics.Hash = Hash(&varMetrics, keyReal)
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		logg.Error().Err(err).Msg("Failed marshal json gauge: ")
		// fmt.Printf("Failed marshal json: %s", err)
		return Metrics{}, err
	}
	varMetrics.BodyBytes = bodyBytes
	// strURL := fmt.Sprintf("%s/update/", urlReal)
	// varMetrics.strURL = strURL
	return varMetrics, nil
}


func NewMetricsGet(id, mType, urlReal string, key string) (*Metrics, error) {
	urlReal = "http://" + urlReal
	var varMetrics = Metrics{
			ID: id,
			MType: mType,
		}
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		logg.Error().Err(err).Msg("Failed marshal json for get: ")
		// fmt.Printf("Failed marshal json for get: %s\n", err)
		return nil, err
	}
	varMetrics.BodyBytes = bodyBytes
	strURL := fmt.Sprintf("%s/value/", urlReal)
	varMetrics.StrURL = strURL
	return &varMetrics, nil
}

func (mtx *Metrics) Update(url string) error{
	logg.Print("--------------------------------------Update------------------------start-------------------------------------")
	// logg.Print("logg: ", logg)
	logg.Printf("mtx.ID =   %v", mtx.ID)
	logg.Printf("mtx.MType =   %v", mtx.MType)
	logg.Printf("mtx.key =   %v", mtx.Key)
	if mtx.Value != nil {
		logg.Printf("mtx.Value =  %v", *mtx.Value)
	}
	logg.Printf("mtx.Hash =  %v", mtx.Hash)
	strURL := fmt.Sprintf("%s/update/", "http://" + url)
	var err error
	mtxOld := Metrics{}
	client := resty.New()
	_, err = client.R().
	SetBody(mtx.BodyBytes).
	SetResult(&mtxOld).
	Post(strURL)
	logg.Print("--------------------------------------Get perponse-------------------------------------------------------------")
	if err != nil {
		// logg.Printf("Failed unmarshall response %s: %s\n", mtxOld.MType, err)
		logg.Error().Err(err).Msg("Failed unmarshall response")
		// logg.Error().Err(err.Printf("Failed unmarshall response %s: %s\n", mtxOld.MType, err).Send()

		return err
	}
	// logg.Printf("Result of requets is: %#v", mtxOld)
	logg.Printf("mtxOld.ID =   %v", mtxOld.ID)
	logg.Printf("mtxOld.MType =   %v", mtxOld.MType)
	logg.Printf("mtxOld.key =   %v", mtxOld.Key)
	if mtxOld.Value != nil {
		logg.Printf("mtxOld.Value =  %v", *mtxOld.Value)
	}
	if mtxOld.Delta != nil {
		logg.Printf("mtxOld.Delta =  %v", *mtxOld.Delta)
	}
	logg.Printf("mtxOld.Hash =  %v\n ", mtxOld.Hash)
	
	return nil
}

func (mtx *Metrics) GetValue() {
	var err error
	client := resty.New()
	_, err = client.R().
	SetResult(mtx).
	SetBody(mtx.BodyBytes).
	Get(mtx.StrURL)
	if err != nil {
		logg.Error().Err(err).Msg("Failed unmarshall response: ")
		// fmt.Printf("Failed unmarshall response %s: %s\n", mtx.MType, err)
	}
}

func Hash(m *Metrics, key string) string {
	var data string
	switch m.MType {
	case "counter":
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case "gauge":
		data = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}