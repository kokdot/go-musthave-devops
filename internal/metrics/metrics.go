package metrics

import (
 	"github.com/go-resty/resty/v2"
	"encoding/json"
	"fmt"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
	bodyBytes []byte
	strURL string
}

type Gauge float64
type Counter int64

func NewMetricsCounter(id,  counterPtr *Counter, urlReal string) (*Metrics, error) {
	var varMetrics Metrics = Metrics{
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
	strURL := fmt.Sprintf("%s/update/", urlReal)
	varMetrics.strURL = strURL
	return &varMetrics, nil
}

func NewMetricsGauge(id string, gaugePtr *Gauge,  urlReal string) (*Metrics, error) {
	var varMetrics Metrics = Metrics{
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
	strURL := fmt.Sprintf("%s/update/", urlReal)
	varMetrics.strURL = strURL
	return &varMetrics, nil
}

func NewMetricsGet(id, mType, urlReal string) (*Metrics, error) {
	var varMetrics Metrics = Metrics{
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