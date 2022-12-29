package store

import (
	"errors"
	"fmt"
	"os"
	"log"

)
type FileStorage struct {
	StoreMap   *StoreMap
	file string
	fileRestore string
	// ProducerPtr *producer
	// ProducerFilePtr *producer
	// ConsumerPtr *consumer
	// ConsumerFilePtr *consumer
}
func NewFileStorage() (*FileStorage, error) {
	tmpfile, err := os.CreateTemp("/tmp/", "devops-metrics-db")
	if err != nil {
        log.Fatal(err)
    }
	file := tmpfile.Name()
	tmpfile.Close()
	// p1, err := NewProducer(tmpfile.Name())
	// if err != nil {
	// 	return nil, err
	// }
	// c1, err := NewConsumer(tmpfile.Name())
	// if err != nil {
	// 	return nil, err
	// }
	sm := make(StoreMap)
	return &FileStorage{
		StoreMap : &sm, 
		file: file,
		// ProducerPtr: p1,
		// ConsumerPtr: c1,

	}, nil
}
func NewFileStorageWithFile(fileRestore string) (*FileStorage, error) {
	// p, err := NewProducer(filename)
	// if err != nil {
	// 	return nil, err
	// }
	// c, err := NewConsumer(filename)
	// if err != nil {
	// 	return nil, err
	// }
	tmpfile, err := os.CreateTemp("/tmp/", "devops-metrics-db")
	if err != nil {
        log.Fatal(err)
    }
	file := tmpfile.Name()
	tmpfile.Close()
	// p1, err := NewProducer(tmpfile.Name())
	// if err != nil {
	// 	return nil, err
	// }
	// c1, err := NewConsumer(tmpfile.Name())
	// if err != nil {
	// 	return nil, err
	// }
	sm := make(StoreMap)
	return &FileStorage{
		StoreMap : &sm,
		file: file,
		fileRestore: fileRestore, 
		// ProducerFilePtr: p,
		// ConsumerFilePtr: c,
		// ProducerPtr: p1,
		// ConsumerPtr: c1,

	}, nil
}
// func (f FileStorage) Read () error {
// 	sm, err := f.ConsumerPtr.ReadStorage()
// 	if err != nil {
// 		return nil
// 	}
// 	f.StoreMap = sm
// 	return nil
// }
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
	p, err := NewProducer(f.fileRestore)
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
	c, err := NewConsumer(f.fileRestore)
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
	p, err := NewProducer(f.file)
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
	c, err := NewConsumer(f.file)
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
		mtxNew := NewMetrics(id, "counter")
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
		mtxNew := NewMetrics(id, "gauge")
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
// func (m *FileStorage) Close() error {
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