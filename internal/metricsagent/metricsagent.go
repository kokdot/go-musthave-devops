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
// type SliceMetrics []Metrics
type MetricsSlice [] Metrics
var logg zerolog.Logger// = onboardingagent.GetLogg()
var conf *def.Conf
func GetConf(config *def.Conf) {
	conf = config
	logg = conf.Logg
}

func UpdateByBatch(gm *def.GaugeMap, counter Counter) error {
	var m = make(MetricsSlice, 0)

	mtxCounter, err := NewMetricsCounter("PollCount", counter)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	err = GetMetricsSlice(gm, &m)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	m = append(m, mtxCounter)
	
	err = UpdateStoreMap(&m)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	return nil
}

func GetMetricsSlice(gm *def.GaugeMap, m *MetricsSlice) error {
	for k, v := range *gm {
		mtx, err := NewMetricsGauge(k, v)
		if err != nil {
			logg.Error().Err(err).Send()
			return fmt.Errorf("%s", err)
		}
		logg.Print("mtx.ID: ", mtx.ID)
		logg.Print("mtx.Value: ", *mtx.Value)
		*m = append(*m, mtx)
	}
	return nil
}

func UpdateStoreMap(m *MetricsSlice) error {
	var err error
	url := "http://" + conf.URL
	strURL := fmt.Sprintf("%s/updates/", url)
	bodyBytes, err := json.Marshal(m)
	if err != nil {
		logg.Error().Err(err).Msg("Failed marshal json for batch: ")
		return err
	}
	client := resty.New()
	_, err = client.R().
	SetHeader("Accept-Encoding", "gzip").
	SetHeader("Content-Type", "application/json").
	SetBody(bodyBytes).
	SetResult(&m).
	Post(strURL)
	if err != nil {
		logg.Error().Err(err).Msg("Failed unmarshall response by batch: ")
		return err
	}
	return nil
}

func UpdateAll (gm *def.GaugeMap, counter Counter) error {
	logg.Print("------------------------------------UpdateAll---------------------------start---------------------")
	mtxCounter, err := NewMetricsCounter("PollCount", counter)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	err = Update(mtxCounter)
	if err != nil {
		logg.Error().Err(err).Send()
		return err
	}
	// n := 0
	for k, v := range *gm {
		// n++
		// if n > 1 {
		// 	break
		// }
		mtx, err := NewMetricsGauge(k, v) 
		logg.Printf("mtx: %v; Value =  %v", mtx.ID, *mtx.Value)
		if err != nil {
			logg.Error().Err(err).Send()
			return err
		}
		Update(mtx)
	}
	return nil
}

func NewMetricsCounter(id string, counter Counter) (Metrics, error) {
	if conf.Key == "" {
		var varMetrics = Metrics{
				ID: id,
				MType: "counter",
				Delta: &counter,
			}
		bodyBytes, err := json.Marshal(varMetrics)
		if err != nil {
			logg.Error().Err(err).Msg("Failed marshal json counter: ")
			return Metrics{}, err
		}
		varMetrics.BodyBytes = bodyBytes
		return varMetrics, nil
	}
	var varMetrics = Metrics{
			ID: id,
			MType: "counter",
			Delta: &counter,
		}
	varMetrics.Hash = Hash(&varMetrics, conf.Key)
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		fmt.Printf("Failed marshal json: %s", err)
		return Metrics{}, err
	}
	varMetrics.BodyBytes = bodyBytes
	return varMetrics, nil
}

func NewMetricsGauge(id string, gauge Gauge) (Metrics, error) {
	if conf.Key == "" {
		var varMetrics = Metrics{
			ID: id,
			MType: "gauge",
			Value: &gauge,
		}
		bodyBytes, err := json.Marshal(varMetrics)
		if err != nil {
			logg.Error().Err(err).Msg("Failed marshal json gauge: ")
			return Metrics{}, err
		}
		varMetrics.BodyBytes = bodyBytes
		return varMetrics, nil
	}
	var varMetrics = Metrics{
		ID: id,
		MType: "gauge",
		Value: &gauge,
	}
	varMetrics.Hash = Hash(&varMetrics, conf.Key)
	bodyBytes, err := json.Marshal(varMetrics)
	if err != nil {
		logg.Error().Err(err).Msg("Failed marshal json gauge: ")
		return Metrics{}, err
	}
	varMetrics.BodyBytes = bodyBytes
	return varMetrics, nil
}


// func NewMetricsGet(id, mType string) (*Metrics, error) {
// 	var varMetrics = Metrics{
// 			ID: id,
// 			MType: mType,
// 		}
// 	bodyBytes, err := json.Marshal(varMetrics)
// 	if err != nil {
// 		logg.Error().Err(err).Msg("Failed marshal json for get: ")
// 		return nil, err
// 	}
// 	varMetrics.BodyBytes = bodyBytes
// 	return &varMetrics, nil
// }

func Update(mtx Metrics) error{
	logg.Print("--------------------------------------Update------------------------start-------------------------------------")
	logg.Printf("mtx.ID =   %v", mtx.ID)
	logg.Printf("mtx.MType =   %v", mtx.MType)
	logg.Printf("mtx.key =   %v", mtx.Key)
	if mtx.Value != nil {
		logg.Printf("mtx.Value =  %v", *mtx.Value)
	}
	logg.Printf("mtx.Hash =  %v", mtx.Hash)

	strURL := fmt.Sprintf("%s/update/", "http://" + conf.URL)
	var err error
	mtxOld := Metrics{}
	client := resty.New()
	_, err = client.R().
	SetBody(mtx.BodyBytes).
	SetResult(&mtxOld).
	Post(strURL)
	logg.Print("--------------------------------------Get perponse-------------------------------------------------------------")
	if err != nil {
		logg.Error().Err(err).Msg("Failed unmarshall response")
		return err
	}

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