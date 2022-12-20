package store

import (
	"fmt"
	"errors"
)


type Counter int64
type Gauge float64
// type GaugeMap map[string]Gauge
// type CounterMap map[string]Counter
type StoreMap map[string]Metrics

type MemStorage struct {
	StoreMap   StoreMap
}

// type Repo interface {
// 	SaveCounterValue(name string, counter Counter) Counter
// 	SaveGaugeValue(name string, gauge Gauge)
// 	GetCounterValue(name string) (Counter, error)
// 	GetGaugeValue(name string) (Gauge, error)
// 	GetAllValues() string
// 	GetAllValuesJson() (GaugeMap, CounterMap)
// }
type Repo interface {
	Save(mtx *Metrics) *Metrics
	Get(id string) (*Metrics, error)
	GetAll() (StoreMap)
	SaveCounterValue(name string, counter Counter) Counter
	SaveGaugeValue(name string, gauge Gauge)
	GetCounterValue(name string) (Counter, error)
	GetGaugeValue(name string) (Gauge, error)
	GetAllValues() string
}
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
var zeroG Gauge
var zeroC Counter
func NewMetrics(id string, mType string) Metrics {
	if mType == "Gauge" {
		return Metrics{
		ID: id,
		MType: "Gauge",
		Value: &zeroG,
		}
	} else {
		return Metrics{
			ID: id,
			MType: "Counter",
			Delta: &zeroC,
		}
	}
}


func (m MemStorage) Save(mtxNew *Metrics) *Metrics {
	switch mtxNew.MType {
	case "Gauge":
		m.StoreMap[mtxNew.ID] = *mtxNew
		return mtxNew
	case "Counter":
		mtxOld, ok := m.StoreMap[mtxNew.ID]
		if !ok {
			m.StoreMap[mtxNew.ID] = *mtxNew
			return mtxNew
		}
		*mtxOld.Delta += *mtxNew.Delta
		return &mtxOld
	}
	return mtxNew
}

func (m MemStorage) Get(id string) (*Metrics, error) {
	mtxOld, ok := m.StoreMap[id]
	if !ok {
		return nil, errors.New("metrics not found")
	}
	return &mtxOld, nil
}

func (m MemStorage) GetAll() StoreMap {
	return m.StoreMap
}

func (m *MemStorage) SaveCounterValue(id string, counter Counter) Counter {
// func (m *MemStorage) SaveCounterValue(name string, counter Counter) Counter {
	mtxOld, ok := m.StoreMap[id]
	if !ok {
		mtxNew := NewMetrics(id, "Counter")
		mtxNew.Delta = &counter
		m.StoreMap[id] = mtxNew
		return counter
	}
	*mtxOld.Delta += counter
	return *mtxOld.Delta
}

func (m *MemStorage) SaveGaugeValue(id string, gauge Gauge) {
	mtxOld, ok := m.StoreMap[id]
	if !ok {
		mtxNew := NewMetrics(id, "Gauge")
		mtxNew.Value = &gauge
		m.StoreMap[id] = mtxNew
	}
	*mtxOld.Value = gauge
}

func (m *MemStorage) GetCounterValue(id string) (Counter, error) {
	mtxOld, ok := m.StoreMap[id]
	if !ok {
		return 0, errors.New("this counter don't find")
	}
	return *mtxOld.Delta, nil
}

func (m *MemStorage) GetGaugeValue(id string) (Gauge, error) {
	mtxOld, ok := m.StoreMap[id]
	if !ok {
		return 0, errors.New("this gauge don't find")
	}
	return *mtxOld.Value, nil
}

func (m *MemStorage) GetAllValues() string {
	var str string
	for key, val := range m.StoreMap {
		str += fmt.Sprintf("%s: %v %v\n", key, val.Value, val.Delta)

	}
	// mapAll := make(map[string]string)
	// for key, val := range m.CounterMap {
	// 	mapAll[key] = fmt.Sprintf("%v", val)
	// }
	// for key, val := range m.GaugeMap {
	// 	mapAll[key] = fmt.Sprintf("%v", val)
	// }
	// for key, val := range mapAll{
	// 	str += fmt.Sprintf("%s: %s\n", key, val)
	// }
	return str
}
// func (m *MemStorage) GetAllValuesJson() (GaugeMap, CounterMap) {
// 	gaugeMap := m.GaugeMap
// 	counterMap := m.CounterMap
// 	return gaugeMap, counterMap 
// }