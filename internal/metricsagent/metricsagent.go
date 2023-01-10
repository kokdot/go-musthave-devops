package metricsagent

import (
 	"github.com/go-resty/resty/v2"
	"github.com/kokdot/go-musthave-devops/internal/onboardingagent"
	"github.com/kokdot/go-musthave-devops/internal/def"
	"encoding/json"
	"fmt"
	"crypto/sha256"
	"crypto/hmac"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
	// Hash  []byte   `json:"hash,omitempty"`  // значение хеш-функции
	bodyBytes []byte
	strURL string
}

type Gauge = def.Gauge
type Counter = def.Counter

func NewMetricsCounter(id string,  counterPtr *Counter, urlReal string) (*Metrics, error) {
	keyReal := onboardingagent.KeyReal
	key := []byte(keyReal)
	urlReal1 := "http://" + urlReal
	if keyReal == "" {

		var varMetrics = Metrics{
				ID: id,
				MType: "counter",
				Delta: counterPtr,
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

func NewMetricsGauge(id string, gaugePtr *Gauge,  urlReal string) (*Metrics, error) {
	keyReal := onboardingagent.KeyReal
	key := []byte(keyReal)
	urlReal1 := "http://" + urlReal
	if keyReal == "" {

		var varMetrics = Metrics{
			ID: id,
			MType: "gauge",
			Value: gaugePtr,
		}
		bodyBytes, err := json.Marshal(varMetrics)
		if err != nil {
			fmt.Printf("Failed marshal json gauge: %s\n", err)
			return nil, err
		}
		varMetrics.bodyBytes = bodyBytes
		strURL := fmt.Sprintf("%s/update/", urlReal1)
		varMetrics.strURL = strURL
		return &varMetrics, nil
	}
	src := []byte((fmt.Sprintf("%s:gauge:%f", id, float64(*gaugePtr))))
	h := hmac.New(sha256.New, key)
	h.Write(src)
	dst := h.Sum(nil)
	var varMetrics = Metrics{
			ID: id,
			MType: "gauge",
			Value: gaugePtr,
			Hash: fmt.Sprintf("%x", dst),
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

func NewMetricsGet(id, mType, urlReal string) (*Metrics, error) {
	urlReal1 := "http://" + urlReal
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
	strURL := fmt.Sprintf("%s/value/", urlReal1)
	varMetrics.strURL = strURL
	return &varMetrics, nil
}
func (mtx *Metrics) Update() {
	var err error
	client := resty.New()
	_, err = client.R().
	SetResult(mtx).
	SetBody(mtx.bodyBytes).
	Post(mtx.strURL)
	if err != nil {
		fmt.Printf("Failed unmarshall response %s: %s\n", mtx.MType, err)
	}
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