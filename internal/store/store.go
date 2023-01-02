package store

import (
	"fmt"
	"errors"
	"sort"
	// "os"
	// "bufio"
	// "log"
	// "encoding/json"
)


type Counter int64
type Gauge float64
type StoreMap map[string]Metrics

type MemStorage struct {
	StoreMap   *StoreMap
	// ProducerPtr *producer
	// ConsumerPtr *consumer
	file string
}
func NewMemStorageWithFile(filename string) (*MemStorage, error) {
	// p, err := NewProducer(filename)
	// if err != nil {
	// 	return nil, err
	// }
	// c, err := NewConsumer(filename)
	// if err != nil {
	// 	return nil, err
	// }
	sm := make(StoreMap)
	sm["1"] = Metrics{
		ID: "1",
		MType: "1",
		
	}
	return &MemStorage{
		StoreMap : &sm, 
		file: filename,
		// ProducerPtr: p,
		// ConsumerPtr: c,

	}, nil
}
func NewMemStorage() (*MemStorage, error) {
	sm := make(StoreMap)
	sm["1"] = Metrics{
		ID: "1",
		MType: "1",

	}
	return &MemStorage{
		StoreMap : &sm,
	}, nil
}
//------------------------------------interface--------------------------------------
type Consumer interface {
    ReadStorage() (*StoreMap, error) // для чтения события
    Close() error               // для закрытия ресурса (файла)
}
type Producer interface {
    WriteStorage() error // для записи события
    Close() error            // для закрытия ресурса (файла)
}
type Repo interface {
	Save(mtx *Metrics) (*Metrics, error)
	Get(id string) (*Metrics, error)
	GetAll() (StoreMap, error)
	SaveCounterValue(name string, counter Counter) (Counter, error)
	SaveGaugeValue(name string, gauge Gauge) error
	GetCounterValue(name string) (Counter, error)
	GetGaugeValue(name string) (Gauge, error)
	GetAllValues() (string, error)
	ReadStorage() (*StoreMap, error)
	WriteStorage() error 
	// WriteStorage(file string)
	// ReadStorage(file string) *MemStorage
}
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
var zeroG Gauge = 0
var zeroC Counter = 0

func NewMetrics(id string, mType string) Metrics {
	if mType == "gauge" {
		return Metrics{
		ID: id,
		MType: "gauge",
		Value: &zeroG,
		}
	} else {
		return Metrics{
			ID: id,
			MType: "counter",
			Delta: &zeroC,
		}
	}
}

func (m MemStorage) Save(mtxNew *Metrics) (*Metrics, error) {
	if m.StoreMap == nil {
		// fmt.Println("--store---Save--err line 107-----------------------------")
		return nil, errors.New("memStorage is nil")
	}
	// fmt.Println("--111---store----Save-------------------")
	if mtxNew == nil {
	// fmt.Println("--113--store---Save line 112-----mtx.new is nil-------------------")

	}
	// fmt.Println("---116--store---Save line 115----mtxNew-------------------:   " , mtxNew)
	if mtxNew.Value != nil {
		fmt.Println("--118---store---Save line 118----mtxNew.Value-------------------:   " , *mtxNew.Value, "  ----ID----:  ", *mtxNew.Value)

	}
	switch mtxNew.MType {
	case "Gauge":
        // fmt.Println("---123---store---Save line 122---Gauge-------------------")
        // fmt.Println("---124--store---Save line----Save-------------------")
        // fmt.Println("---125--store---Save line----Save-------------------")
		
		(*m.StoreMap)[mtxNew.ID] = *mtxNew
		// sm, _ := m.GetAll()
        // fmt.Println("--129---store---Save line----Save---m.GetAll()----------------:   ", sm)

		return mtxNew, nil
	case "gauge":
        // fmt.Println("---123---store---Save line 122---Gauge-------------------")
        // fmt.Println("---124--store---Save line----Save-------------------")
        // fmt.Println("---125--store---Save line----Save-------------------")
		
		(*m.StoreMap)[mtxNew.ID] = *mtxNew
		// sm, _ := m.GetAll()
        // fmt.Println("--129---store---Save line----Save---m.GetAll()----------------:   ", sm)

		return mtxNew, nil
	case "counter":
        // fmt.Println("---133--store----Counter-------------------")
		mtxOld, ok := (*m.StoreMap)[mtxNew.ID]
		if !ok {
			(*m.StoreMap)[mtxNew.ID] = *mtxNew
			return mtxNew, nil
		}
		*mtxOld.Delta += *mtxNew.Delta
		return &mtxOld, nil
	case "Counter":
        // fmt.Println("---133--store----Counter-------------------")
		mtxOld, ok := (*m.StoreMap)[mtxNew.ID]
		if !ok {
			(*m.StoreMap)[mtxNew.ID] = *mtxNew
			return mtxNew, nil
		}
		*mtxOld.Delta += *mtxNew.Delta
		return &mtxOld, nil
	}
	return mtxNew, nil
}

func (m MemStorage) Get(id string) (*Metrics, error) {
	if m.StoreMap == nil {
		return nil, errors.New("memStorage is nil")
	}
	mtxOld, ok := (*m.StoreMap)[id]
	if !ok {
		return nil, errors.New("metrics not found")
	}
	return &mtxOld, nil
}

func (m MemStorage) GetAll() (StoreMap, error) {
	if m.StoreMap == nil {
		return nil, errors.New("memStorage is nil")
	}
	return *m.StoreMap, nil
}

func (m *MemStorage) SaveCounterValue(id string, counter Counter) (Counter, error) {
	// var counter1 Counter
	if m.StoreMap == nil {
		return counter, errors.New("memStorage is nil")
	}
	mtxOld, ok := (*m.StoreMap)[id]
	if !ok {
		mtxNew := NewMetrics(id, "counter")
		mtxNew.Delta = &counter
		(*m.StoreMap)[id] = mtxNew
		return counter, nil
	}
	*mtxOld.Delta += counter
	return *mtxOld.Delta, nil
}

func (m *MemStorage) SaveGaugeValue(id string, gauge Gauge) error {
	if m.StoreMap == nil {
		return errors.New("memStorage is nil")
	}
	mtxOld, ok := (*m.StoreMap)[id]
	if !ok {
		mtxNew := NewMetrics(id, "gauge")
		mtxNew.Value = &gauge
		(*m.StoreMap)[id] = mtxNew
	}else {
		*mtxOld.Value = gauge
	}
	return nil
}
// func (m *MemStorage) Close() error {
// 	err := (*m.ConsumerPtr).file.Close()
// 	if err != nil {
// 		return err
// 	}
// 	err = (*m.ProducerPtr).file.Close()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (m *MemStorage) GetCounterValue(id string) (Counter, error) {
	var counter1 Counter
	if m.StoreMap == nil {
		return counter1, errors.New("memStorage is nil")
	}
	mtxOld, ok := (*m.StoreMap)[id]
	if !ok {
		return 0, errors.New("this counter don't find")
	}
	return *mtxOld.Delta, nil
}

func (m *MemStorage) GetGaugeValue(id string) (Gauge, error) {
	var gauge1 Gauge
	if m.StoreMap == nil {
		return gauge1, errors.New("memStorage is nil")
	}
	// fmt.Println("id:  ", id, "-----------------------((((((((((((((((((((((((((((((((((((")
	mtxOld, ok := (*m.StoreMap)[id]
	// fmt.Println("mtxOld:  ", mtxOld, ";  OK: ", ok, "MemStorage", m)
	if !ok {
		return 0, errors.New("this gauge don't find")
	}
	// fmt.Println("*mtxOld.Value:  ", *mtxOld.Value, "----------------------------------")
	return *mtxOld.Value, nil
}

func (m *MemStorage) GetAllValues() (string, error) {
	var str string
	var v Gauge
	var d Counter
	var i int
	if m.StoreMap == nil {
		return "", errors.New("storeMap is nil")
	}
	sm := *m.StoreMap
	keys := make([]string, 0, len(sm))
	for k := range sm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		i++
		if sm[key].Delta != nil {
			d = *sm[key].Delta
		}
		if sm[key].Value != nil {
			v = *sm[key].Value
		}
		str += fmt.Sprintf("%d; %s: %v %v\n",i , key, v, d)
	}
	return str, nil
}

func (m MemStorage) ReadStorage() (*StoreMap, error) {
	c, err := NewConsumer(m.file)
	if err != nil  {
		err1 := fmt.Errorf("can't to create consumer: %s", err)
		return nil, err1
    }
	sm, err := c.ReadStorage()
	if sm == nil && err !=  nil {
		err1 := fmt.Errorf("file for StoreMap is ampty: %s", err)
		return nil, err1
	} else if sm == nil {
		return nil, errors.New("file for StoreMap is ampty")
	}
	(*m.StoreMap) = *sm
	return sm, nil
}

func (m MemStorage) WriteStorage() error{
	p, err := NewProducer(m.file)
	if err != nil  {
		err1 := fmt.Errorf("can't to create producer: %s", err)
		return err1
    }
	err = p.WriteStorage((m.StoreMap))
	if err != nil {
        err1 := fmt.Errorf("can't to write memStorege: %s", err)
		return err1
    }
	return nil
}