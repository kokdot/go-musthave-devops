package store

import (
	"fmt"
	"errors"
	"sort"
	"crypto/sha256"
	"crypto/hmac"
)

type Counter int64
type Gauge float64
type StoreMap map[string]Metrics

type MemStorage struct {
	StoreMap   *StoreMap
	file string
}
func NewMemStorageWithFile(filename string) (*MemStorage, error) {
	sm := make(StoreMap)
	sm["1"] = Metrics{
		ID: "1",
		MType: "1",
		
	}
	return &MemStorage{
		StoreMap : &sm, 
		file: filename,
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
}
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  []byte   `json:"hash,omitempty"`  // значение хеш-функции
}
var zeroG Gauge = 0
var zeroC Counter = 0
var Key string

func GetKey(key string) {
	Key = key
}

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
func NewCounterMetrics(id string, counter Counter) *Metrics {
	src := []byte(fmt.Sprintf("%s:counter:%d", id, counter))
	// keyCrypto := Key
	if Key == "" {
		panic("bad")
	}
	h := hmac.New(sha256.New, []byte(Key))
    h.Write(src)
    dst := h.Sum(nil)
	var varMetrics Metrics = Metrics{
			ID: id,
			MType: "counter",
			Delta: &counter,
			Hash: dst,
		}
	return &varMetrics
}

func NewGaugeMetrics(id string, gauge Gauge) *Metrics {
	src := []byte(fmt.Sprintf("%s:gauge:%f", id, float64(gauge)))
	// keyCrypto := Key
	h := hmac.New(sha256.New, []byte(Key))
    h.Write(src)
    dst := h.Sum(nil)
	var varMetrics Metrics = Metrics{
			ID: id,
			MType: "gauge",
			Value: &gauge,
			Hash: dst,
		}
	return &varMetrics
}
func MtxValid(mtx *Metrics) bool {
	// fmt.Println("-----------------------------------MtxValid-----start----")

	if Key == "" {
	fmt.Println("-----------------------------------MtxValid-------if Key == nil--")

		return true
	}
	if mtx.Hash == nil {
		fmt.Println("--------------------------------------------------------------------------------------------------mtx.Hash is ampty----")
		return false
	}
	var src []byte
	if mtx.MType == "gauge"{
	// fmt.Println("-----------------------------------MtxValid-------if mtx.MType == gauge--")

		src = []byte((fmt.Sprintf("%s:gauge:%f", mtx.ID, *mtx.Value)))
		// src = []byte((fmt.Sprintf("%s:gauge:%f", mtx.ID, float64(*mtx.Value))))
	} else if mtx.MType == "counter" {
	// fmt.Println("-----------------------------------MtxValid-------else if mtx.MType == counter--")

		src = []byte((fmt.Sprintf("%s:counter:%v", mtx.ID, *mtx.Delta)))
	} else {
		fmt.Println("-----------------------------------MtxValid-------else --false--")
		fmt.Printf("not counter not gauge: %#v:    \n", mtx)
		
		return false
	}
	
	h := hmac.New(sha256.New, []byte(Key))
    h.Write(src)
    dst := h.Sum(nil)
	// fmt.Println("-----------------------------------MtxValid-----return----")
	fmt.Println("hash old: ", dst)
	fmt.Println("hash new: ", mtx.Hash)
	fmt.Println("hmac.Equal(dst, mtx.Hash): ", hmac.Equal(dst, mtx.Hash))
	// fmt.Println("-----------------------MtxValid-----finish---------------------")
	return hmac.Equal(dst, mtx.Hash)
}

func (m MemStorage) Save(mtxNew *Metrics) (*Metrics, error) {
	if Key != "" {

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
			mtxOld = *NewCounterMetrics(mtxNew.ID, delta)
			(*m.StoreMap)[mtxOld.ID] = mtxOld
			return &mtxOld, nil
		case "Counter":
			mtxOld, ok := (*m.StoreMap)[mtxNew.ID]
			if !ok {
				(*m.StoreMap)[mtxNew.ID] = *mtxNew
				return mtxNew, nil
			}
			delta := *mtxNew.Delta + *mtxOld.Delta
			mtxOld = *NewCounterMetrics(mtxNew.ID, delta)
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
	if Key != "" {

		if mtxOld.MType == "Gauge" || mtxOld.MType == "gauge" {
			mtxOld = *NewGaugeMetrics(mtxOld.ID, *mtxOld.Value)
		} else {
				mtxOld = *NewCounterMetrics(mtxOld.ID, *mtxOld.Delta) //-------------------------------------line : 216
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