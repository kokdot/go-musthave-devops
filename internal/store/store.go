package store

import (
	"fmt"
	"errors"
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
	return &MemStorage{
		StoreMap : &sm, 
		file: filename,
		// ProducerPtr: p,
		// ConsumerPtr: c,

	}, nil
}
func NewMemStorage() *MemStorage {
	sm := make(StoreMap)
	return &MemStorage{
		StoreMap : &sm, 
	}
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

func (m MemStorage) Save(mtxNew *Metrics) (*Metrics, error) {
	if m.StoreMap == nil {
		return nil, errors.New("memStorage is nil")
	}
	fmt.Println("---------Save-------------------")
	if mtxNew == nil {
	fmt.Println("---------mtx.new is nil-------------------")

	}
	fmt.Println("---------mtxNew-------------------:   " , mtxNew)
	if mtxNew.Value != nil {
		fmt.Println("---------mtxNew-------------------:   " , *mtxNew.Value)

	}
	switch mtxNew.MType {
	case "Gauge":
        // fmt.Println("---------Gauge-------------------")
        // fmt.Println("---------Save-------------------")
        // fmt.Println("---------Save-------------------")
		
		(*m.StoreMap)[mtxNew.ID] = *mtxNew
		return mtxNew, nil
	case "Counter":
        fmt.Println("---------Counter-------------------")
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
	var counter1 Counter
	if m.StoreMap == nil {
		return counter1, errors.New("memStorage is nil")
	}
	mtxOld, ok := (*m.StoreMap)[id]
	if !ok {
		mtxNew := NewMetrics(id, "Counter")
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
		mtxNew := NewMetrics(id, "Gauge")
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
	if m.StoreMap == nil {
		return "", errors.New("storeMap is nil")
	}
	for key, val := range (*m.StoreMap) {
		if val.Delta != nil {
			d = *val.Delta
		}
		if val.Value == nil {
			v = *val.Value
		}
		str += fmt.Sprintf("%s: %v %v\n", key, v, d)

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