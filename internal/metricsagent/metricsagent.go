package metricsagent

import (
 	"github.com/go-resty/resty/v2"
	// "github.com/kokdot/go-musthave-devops/internal/onboardingagent"
	"github.com/kokdot/go-musthave-devops/internal/def"
	// "github.com/kokdot/go-musthave-devops/internal/monitor"
	"encoding/json"
	"fmt"
	"crypto/sha256"
	"crypto/hmac"
)

type Gauge = def.Gauge 
type Counter = def.Counter
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
	bodyBytes []byte
	strURL string
	key string
}
type StoreMap map[string] Metrics

func GetStoreMap(smPtr *StoreMap, mPtr *def.MonitorMap, url string, key string, mtxs... *Metrics) (*StoreMap, error) {
	for key, val := range *mPtr {
		mtx, err := NewMetricsGauge(key, &val, url, key)
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}
		(*smPtr)[key] = *mtx
	}
	for _, mtx := range mtxs {
		(*smPtr)[mtx.ID] = *mtx
	}
	return smPtr, nil
}

func UpdateByBatch(smPtr *StoreMap, mPtr *def.MonitorMap, pollCount Counter, randomValue Gauge, url string, key string) error {
	mtxCounter, err := NewMetricsCounter("PollCount", &pollCount, url, key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	mtxRandomValue, err := NewMetricsGauge("RandomValue", &randomValue, url, key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	smPtr, err = GetStoreMap(smPtr, mPtr, url, key, mtxCounter, mtxRandomValue)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = UpdateStoreMap(smPtr, url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func UpdateStoreMap(smPtrnew *StoreMap, url string) error {
	var err error
	var smOld = make(StoreMap, len(*smPtrnew))
	url = "http://" + url
	strURL := fmt.Sprintf("%s/updates/", url)
	bodyBytes, err := json.Marshal(smPtrnew)
	if err != nil {
		fmt.Printf("Failed marshal json for batch: %s", err)
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
		fmt.Printf("Failed unmarshall response by batch %#v: %s\n", smOld, err)
		return err
	}
	fmt.Printf("Result of requets is: %#v\n", smOld)
	return nil
}

func UpdateAll (m *def.MonitorMap, c Counter, g Gauge, url string, key string) error {
	mtxCounter, err := NewMetricsCounter("PollCount", &c, url, key)
	// fmt.Printf("mtxRandomValue:    %#v\n", mtxCounter)
	if err != nil {
		fmt.Println(err)
	}
	mtxCounter.Update(url)

	mtxRandomValue, err := NewMetricsGauge("RandomValue", &g, url, key)
	// fmt.Printf("mtxRandomValue:    %#v\n", mtxRandomValue)
	if err != nil {
		return err
	}
	mtxRandomValue.Update(url)
	// n := 0
	for k, v := range *m {
		// n++
		// if n > 4 {
		// 	break
		// }
		mtx, err := NewMetricsGauge(k, &v, url, key) 
		// fmt.Printf("mtx.Hash:    %#v\n", mtx.Hash)
		if err != nil {
			return err
		}
		mtx.Update(url)
	}
	return nil
}

func NewMetricsCounter(id string,  counterPtr *Counter, urlReal string, keyReal string) (*Metrics, error) {
	key := []byte(keyReal)
	urlReal1 := "http://" + urlReal
	if keyReal == "" {

		var varMetrics = Metrics{
				ID: id,
				MType: "counter",
				Delta: counterPtr,
				key: keyReal,
			}
		bodyBytes, err := json.Marshal(varMetrics)
		if err != nil {
			fmt.Printf("Failed marshal json counter:  %s\n", err)
			return nil, err
		}
		varMetrics.bodyBytes = bodyBytes
		strURL := fmt.Sprintf("%s/update/", urlReal1)
		varMetrics.strURL = strURL
		return &varMetrics, nil
	}
		src := []byte((fmt.Sprintf("%s:counter:%d", id, *counterPtr)))
	h := hmac.New(sha256.New, key)
    h.Write(src)
    dst := h.Sum(nil)
	var varMetrics = Metrics{
			ID: id,
			MType: "counter",
			Delta: counterPtr,
			Hash: fmt.Sprintf("%x", dst),
			key: keyReal,
		}
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		fmt.Printf("Failed marshal json: %s", err)
		return nil, err
	}
	varMetrics.bodyBytes = bodyBytes
	strURL := fmt.Sprintf("%s/update/", urlReal1)
	varMetrics.strURL = strURL
	return &varMetrics, nil
}

func NewMetricsGauge(id string, gaugePtr *Gauge,  urlReal string, keyReal string) (*Metrics, error) {
	// key := []byte(keyReal)
	// urlReal = "http://" + urlReal
	if keyReal == "" {

		var varMetrics = Metrics{
			ID: id,
			MType: "gauge",
			Value: gaugePtr,
			key: keyReal,
		}
		bodyBytes, err := json.Marshal(varMetrics)
		if err != nil {
			fmt.Printf("Failed marshal json gauge: %s\n", err)
			return nil, err
		}
		varMetrics.bodyBytes = bodyBytes
		// strURL := fmt.Sprintf("%s/update/", urlReal)
		// varMetrics.strURL = strURL
		return &varMetrics, nil
	}
	// src := []byte((fmt.Sprintf("%s:gauge:%f", id, float64(*gaugePtr))))
	// h := hmac.New(sha256.New, key)
	// h.Write(src)
	// dst := h.Sum(nil)
	var varMetrics = Metrics{
		ID: id,
		MType: "gauge",
		Value: gaugePtr,
		key: keyReal,
		// Hash: fmt.Sprintf("%x", dst),
	}
	// hash := Hash(&varMetrics, keyReal)
	varMetrics.Hash = Hash(&varMetrics, keyReal)
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		fmt.Printf("Failed marshal json: %s", err)
		return nil, err
	}
	varMetrics.bodyBytes = bodyBytes
	// strURL := fmt.Sprintf("%s/update/", urlReal)
	// varMetrics.strURL = strURL
	return &varMetrics, nil
}


func NewMetricsGet(id, mType, urlReal string, key string) (*Metrics, error) {
	urlReal = "http://" + urlReal
	var varMetrics = Metrics{
			ID: id,
			MType: mType,
		}
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		fmt.Printf("Failed marshal json for get: %s\n", err)
		return nil, err
	}
	varMetrics.bodyBytes = bodyBytes
	strURL := fmt.Sprintf("%s/value/", urlReal)
	varMetrics.strURL = strURL
	return &varMetrics, nil
}

func (mtx *Metrics) Update(url string) error{
	fmt.Println("--------------------------------------Update------------------------start-------------------------------------")
	fmt.Println("mtx.ID =   ", mtx.ID)
	fmt.Println("mtx.MType =   ", mtx.MType)
	fmt.Println("mtx.key =   ", mtx.key)
	if mtx.Value != nil {
		fmt.Println("mtx.Value =  ", *mtx.Value)
	}
	fmt.Println("mtx.Hash =   ", mtx.Hash)
	strURL := fmt.Sprintf("%s/update/", "http://" + url)
	var err error
	mtxOld := Metrics{}
	client := resty.New()
	_, err = client.R().
	SetBody(mtx.bodyBytes).
	SetResult(&mtxOld).
	Post(strURL)
	if err != nil {
		fmt.Printf("Failed unmarshall response %s: %s\n", mtxOld.MType, err)
		return err
	}
	fmt.Printf("Result of requets is: %#v\n", mtxOld)
	
	return nil
}

func (mtx *Metrics) GetValue() {
	var err error
	client := resty.New()
	_, err = client.R().
	SetResult(mtx).
	SetBody(mtx.bodyBytes).
	Get(mtx.strURL)
	if err != nil {
		fmt.Printf("Failed unmarshall response %s: %s\n", mtx.MType, err)
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