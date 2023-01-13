package store

import (
	
	"errors"
	"fmt"
	"sort"
	"time"
	"github.com/kokdot/go-musthave-devops/internal/repo"
	"github.com/kokdot/go-musthave-devops/internal/metricsserver"

)

type Gauge = repo.Gauge
type Counter = repo.Counter
type StoreMap = repo.StoreMap
type Metrics = repo.Metrics

type MemStorage struct {
	StoreMap   *StoreMap
	storeFile string
	restore bool
	storeInterval time.Duration
	key string
	url string
	dataBaseDSN string
}

// var key string

func (m MemStorage) GetPing() (bool, error) {
	return false, errors.New("MemStorage not defines")
}

func (m MemStorage) GetDataBaseDSN() string {
	return m.dataBaseDSN
}
func (m MemStorage) GetURL() string {
	return m.url
}
func (m MemStorage) GetRestore() bool {
	return m.restore
}
func (m MemStorage) GetStoreFile() string {
	return m.storeFile
}
func (m MemStorage) GetKey() string {
	return m.key
}
func (m MemStorage) GetStoreInterval() time.Duration {
	return m.storeInterval
}
func NewMemStorageWithFile(filename string) (*MemStorage, error) {
	sm := make(StoreMap)
	sm["1"] = repo.Metrics{
		ID: "1",
		MType: "1",
		
	} 
	return &MemStorage{
		StoreMap : &sm, 
		storeFile: filename,
	}, nil
}
func NewMemStorage(storeInterval time.Duration, storeFile string, restore bool, url string, key string, dataBaseDSN string) (*MemStorage, error) {
	sm := make(StoreMap)
	sm["1"] = Metrics{
		ID: "1",
		MType: "1",
	}
	return &MemStorage{
		StoreMap : &sm,
		storeFile: storeFile,
		restore: restore,
		storeInterval: storeInterval,
		key: key,
		url: url,
		dataBaseDSN: dataBaseDSN,
	}, nil
}

func (m MemStorage) Save(mtxNew *Metrics) (*Metrics, error) {
	if key := m.GetKey(); key != "" {

		switch mtxNew.MType {
		case "Gauge":
			(*m.StoreMap)[mtxNew.ID] = *mtxNew
			return mtxNew, nil
		case "gauge":
			(*m.StoreMap)[mtxNew.ID] = *mtxNew
			return mtxNew, nil
		case "counter":
			mtxOld, ok := (*m.StoreMap)[mtxNew.ID]
			if !ok {
				(*m.StoreMap)[mtxNew.ID] = *mtxNew
				return mtxNew, nil
			}
			delta := *mtxNew.Delta + *mtxOld.Delta
			mtxOld = *metricsserver.NewCounterMetrics(mtxNew.ID, delta, key)
			(*m.StoreMap)[mtxOld.ID] = mtxOld
			return &mtxOld, nil
		case "Counter":
			mtxOld, ok := (*m.StoreMap)[mtxNew.ID]
			if !ok {
				(*m.StoreMap)[mtxNew.ID] = *mtxNew
				return mtxNew, nil
			}
			delta := *mtxNew.Delta + *mtxOld.Delta
			mtxOld = *metricsserver.NewCounterMetrics(mtxNew.ID, delta, key)
			(*m.StoreMap)[mtxOld.ID] = mtxOld
			return &mtxOld, nil
		}
		return mtxNew, nil
	} else {
		switch mtxNew.MType {
		case "Gauge":
			(*m.StoreMap)[mtxNew.ID] = *mtxNew
			return mtxNew, nil
		case "gauge":
			(*m.StoreMap)[mtxNew.ID] = *mtxNew
			return mtxNew, nil
		case "counter":
			mtxOld, ok := (*m.StoreMap)[mtxNew.ID]
			if !ok {
				(*m.StoreMap)[mtxNew.ID] = *mtxNew
				return mtxNew, nil
			}
			*mtxOld.Delta += *mtxNew.Delta
			return &mtxOld, nil
		case "Counter":
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
}

func (m MemStorage) Get(id string) (*Metrics, error) {
	if m.StoreMap == nil {
		return nil, errors.New("memStorage is nil")
	}
	mtxOld, ok := (*m.StoreMap)[id]
	if !ok {
		return nil, errors.New("metrics not found")
	}
	// if mtxOld.Hash != nil {
	// 	return &mtxOld, nil
	// }
	fmt.Printf("mtxOld:   %#v", mtxOld)
	if key := m.GetKey(); key != "" {

		if mtxOld.MType == "Gauge" || mtxOld.MType == "gauge" {
			mtxOld = *metricsserver.NewGaugeMetrics(mtxOld.ID, *mtxOld.Value, key)
		} else {
				mtxOld = *metricsserver.NewCounterMetrics(mtxOld.ID, *mtxOld.Delta, key) //-------------------------------------line : 216
		}
			return &mtxOld, nil
	} else {
		return &mtxOld, nil
	}
}

func (m MemStorage) GetAll() (StoreMap, error) {
	if m.StoreMap == nil {
		return nil, errors.New("memStorage is nil")
	}
	return *m.StoreMap, nil
}

func (m *MemStorage) SaveCounterValue(id string, counter Counter) (Counter, error) {
	if m.StoreMap == nil {
		return counter, errors.New("memStorage is nil")
	}
	mtxOld, ok := (*m.StoreMap)[id]
	if !ok {
		mtxNew := metricsserver.NewMetrics(id, "counter")
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
		mtxNew := metricsserver.NewMetrics(id, "gauge")
		mtxNew.Value = &gauge
		(*m.StoreMap)[id] = mtxNew
	}else {
		*mtxOld.Value = gauge
	}
	return nil
}

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
	mtxOld, ok := (*m.StoreMap)[id]
	if !ok {
		return 0, errors.New("this gauge don't find")
	}
	return *mtxOld.Value, nil
}

func (m *MemStorage) GetAllValues() string {
	var str string
	var v Gauge
	var d Counter
	var i int
	if m.StoreMap == nil {
		return ""
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
	return str
}

func (m MemStorage) ReadStorage() error {
	c, err := NewConsumer(m.storeFile)
	if err != nil  {
		err1 := fmt.Errorf("can't to create consumer: %s", err)
		return err1
    }
	sm, err := c.ReadStorage()
	if sm == nil && err !=  nil {
		err1 := fmt.Errorf("file for StoreMap is ampty: %s", err)
		return err1
	} else if sm == nil {
		return errors.New("file for StoreMap is ampty")
	}
	(*m.StoreMap) = *sm
	return nil
}

func (m MemStorage) WriteStorage() error{
	p, err := NewProducer(m.storeFile)
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