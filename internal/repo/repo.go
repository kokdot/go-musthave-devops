package repo

import "time"

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
	GetURL() string
	GetKey() string
	GetStoreFile() string
	GetRestore() bool
	GetStoreInterval() time.Duration
}

type Consumer interface {
    ReadStorage() (*StoreMap, error) // для чтения события
    Close() error               // для закрытия ресурса (файла)
}

type Producer interface {
    WriteStorage() error // для записи события
    Close() error            // для закрытия ресурса (файла)
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  []byte   `json:"hash,omitempty"`  // значение хеш-функции
}
type Counter int64
type Gauge float64
type StoreMap map[string]Metrics