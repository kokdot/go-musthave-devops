package store

import (
	"errors"
	"fmt"
	"os"
	"log"
	"time"

	"github.com/kokdot/go-musthave-devops/internal/metrics_server"
	// _ "github.com/kokdot/go-musthave-devops/internal/repo"
)

type FileStorage struct {
	StoreMap   *StoreMap
	storeFile string
	restoreFile string
	storeInterval time.Duration
	restore bool
	url string
	key string
	dataBaseDSN string
}

func NewFileStorage(storeInterval time.Duration, storeFile string, restore bool, url string, key string, dataBaseDSN string) (*FileStorage, error) {
	tmpfile, err := os.CreateTemp("/tmp/", "devops-metrics-db")
	if err != nil {
        log.Fatal(err)
    }
	file := tmpfile.Name()
	tmpfile.Close()
	sm := make(StoreMap)
	return &FileStorage{
		StoreMap : &sm, 
		storeFile: storeFile,
		restoreFile: file,
		storeInterval: storeInterval, 
		restore: restore,
		url: url,
		key: key,
		dataBaseDSN: dataBaseDSN,
	}, nil
}
func (f FileStorage) GetDataBaseDSN() string {
	return f.dataBaseDSN
}
func (f FileStorage) GetURL() string {
	return f.url
}

func (f FileStorage) GetStoreFile() string {
	return f.storeFile
}

func (f FileStorage) GetRestore() bool {
	return f.restore
}
func (f FileStorage) GetKey() string {
	return f.key
}
func (f FileStorage) GetStoreInterval() time.Duration {
	return f.storeInterval
}

func (f FileStorage) Save(mtxNew *Metrics) (*Metrics, error) {
	_, err := f.ReadStorageSelf()
	if err != nil {
		return nil, err
	}
	switch mtxNew.MType {
	case "gauge":
		(*f.StoreMap)[mtxNew.ID] = *mtxNew
		err := f.WriteStorageSelf()
		if err != nil {
			return nil, err
		} 
		return mtxNew, nil
	case "Gauge":
		(*f.StoreMap)[mtxNew.ID] = *mtxNew
		err := f.WriteStorageSelf()
		if err != nil {
			return nil, err
		} 
		return mtxNew, nil
	case "counter":
		mtxOld, ok := (*f.StoreMap)[mtxNew.ID]
		if !ok {
			(*f.StoreMap)[mtxNew.ID] = *mtxNew
			err := f.WriteStorage()
			if err != nil {
				return nil, err
			} 
			return mtxNew, nil
		}
		*mtxOld.Delta += *mtxNew.Delta
		err := f.WriteStorage()
		if err != nil {
			return nil, err
		} 
		return &mtxOld, nil
		case "Counter":
		mtxOld, ok := (*f.StoreMap)[mtxNew.ID]
		if !ok {
			(*f.StoreMap)[mtxNew.ID] = *mtxNew
			err := f.WriteStorage()              
			if err != nil {
				return nil, err
			} 
			return mtxNew, nil
		}
		*mtxOld.Delta += *mtxNew.Delta
		err := f.WriteStorage()
		if err != nil {
			return nil, err
		} 
		return &mtxOld, nil
	}
	return nil, errors.New("MType is wrong")
}

func (f FileStorage) Get(id string) (*Metrics, error){
	
	sm, err := f.ReadStorageSelf()
	if err != nil {
		return nil, err
	}
	mtxOld, ok := (*sm)[id]
	if !ok {
		return nil, errors.New("metrics not found")
	}
	return &mtxOld, nil
}

func (f FileStorage) GetAll() (StoreMap, error) {
	sm, err := f.ReadStorageSelf()
	if err != nil {
		return nil, err
	}
	return *sm, nil
}
//------------------------WriteStorage-------------------------------
func (f FileStorage) WriteStorage() error {
	p, err := NewProducer(f.storeFile)
	if err != nil  {
		err1 := fmt.Errorf("can't to create producer: %s", err)
		return err1
    }
	err = p.WriteStorage((f.StoreMap))
	if err != nil {
        err1 := fmt.Errorf("can't to write fileStorege: %s", err)
		return err1
    }
	return nil
}
//------------------------ReadStorage-------------------------------
func (f FileStorage) ReadStorage() (*StoreMap, error) {
	c, err := NewConsumer(f.storeFile)
	if err != nil  {
		err1 := fmt.Errorf("can't to create consumer: %s", err)
		return nil, err1
    }
	sm, err := c.ReadStorage()
	err = fmt.Errorf("file for StoreMap is ampty: %s", err)
	if sm == nil && err !=  nil {
		err1 := fmt.Errorf("file for StoreMap is ampty: %s", err)
		return nil, err1
	} else if sm == nil {
		return nil, errors.New("file for StoreMap is ampty")
	}
	(*f.StoreMap) = *sm
	return sm, nil
}
//------------------------WriteStorageSelf-------------------------------
func (f FileStorage) WriteStorageSelf() error {
	p, err := NewProducer(f.restoreFile)
	if err != nil  {
		err1 := fmt.Errorf("can't to create producer: %s", err)
		return err1
    }
	err = p.WriteStorage((f.StoreMap))
	if err != nil {
        err1 := fmt.Errorf("can't to write fileStorege: %s", err)
		return err1
    }
	return nil
}
//------------------------ReadStorageSelf-------------------------------
func (f FileStorage) ReadStorageSelf() (*StoreMap, error) {
	c, err := NewConsumer(f.restoreFile)
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
	(*f.StoreMap) = *sm
	return sm, nil
}
func (f FileStorage) SaveCounterValue(id string, counter Counter) (Counter, error) {
	_, err := f.ReadStorageSelf()
	if err != nil {
		return counter, err
	}
	mtxOld, ok := (*f.StoreMap)[id]
	if !ok {
		mtxNew := metrics_server.NewMetrics(id, "counter")
		mtxNew.Delta = &counter
		(*f.StoreMap)[id] = mtxNew
		err := f.WriteStorage()
		if err != nil {
			return counter, err
		} 
		return counter, nil
	}
	*mtxOld.Delta += counter
	err = f.WriteStorage()
	if err != nil {
		return counter, err
	}
	return *mtxOld.Delta, nil
}

func (f FileStorage) SaveGaugeValue(id string, gauge Gauge) error {
	_, err := f.ReadStorageSelf()
	if err != nil {
		return err
	}
	mtxOld, ok := (*f.StoreMap)[id]
	if !ok {
		mtxNew := metrics_server.NewMetrics(id, "gauge")
		mtxNew.Value = &gauge
		(*f.StoreMap)[id] = mtxNew
		err := f.WriteStorage()
		if err != nil {
			return err
		}
	}else {
		*mtxOld.Value = gauge
		err := f.WriteStorage()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f FileStorage) GetCounterValue(id string) (Counter, error) {
	var counter Counter
	_, err := f.ReadStorageSelf()
	if err != nil {
		return counter, err
	}
	mtxOld, ok := (*f.StoreMap)[id]
	if !ok {
		return counter, errors.New("this counter don't find")
	}
	return *mtxOld.Delta, nil
}

func (f FileStorage) GetGaugeValue(id string) (Gauge, error) {
	var gauge Gauge
	_, err := f.ReadStorageSelf()
	if err != nil {
		return gauge, err
	}
	mtxOld, ok := (*f.StoreMap)[id]
	if !ok {
		return gauge, errors.New("this gauge don't find")
	}
	return *mtxOld.Value, nil
}

func (f FileStorage) GetAllValues() (string, error) {
	var str string
	var v Gauge
	var d Counter
	_, err := f.ReadStorageSelf()
	if err != nil {
		return "", err
	}
	for key, val := range *f.StoreMap {
		if val.Delta != nil {
			d = *val.Delta
		}
		if val.Value != nil {
			v = *val.Value
		}
		str += fmt.Sprintf("%s: %v %v\n", key, v, d)

	}

	return str, nil
}