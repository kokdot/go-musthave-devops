package store

import (
	"fmt"
	"errors"
	"os"
	"bufio"
	"log"
	"encoding/json"
)


type Counter int64
type Gauge float64
type StoreMap map[string]Metrics

type MemStorage struct {
	StoreMap   StoreMap
}

type Repo interface {
	Save(mtx *Metrics) *Metrics
	Get(id string) (*Metrics, error)
	GetAll() (StoreMap)
	SaveCounterValue(name string, counter Counter) Counter
	SaveGaugeValue(name string, gauge Gauge)
	GetCounterValue(name string) (Counter, error)
	GetGaugeValue(name string) (Gauge, error)
	GetAllValues() string
	DownloadMemStorage(file string)
	UpdateMemStorage(file string) *MemStorage
}
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
var zeroG Gauge = 0
var zeroC Counter = 0

type Producer interface {
    WriteMemStorage(memStorage *MemStorage) // для записи события
    Close() error            // для закрытия ресурса (файла)
}

type Consumer interface {
    ReadMemStorage() (*MemStorage, error) // для чтения события
    Close() error               // для закрытия ресурса (файла)
}
type producer struct {
    file *os.File
    // добавляем writer в Producer
    writer *bufio.Writer
}

func (m MemStorage) UpdateMemStorage(file string) *MemStorage {
	consumerPtr, err := NewConsumer(file)
	if err != nil {
        log.Fatal(err)
    }
	ms, err := consumerPtr.ReadMemStorage()
	if err != nil {
        log.Fatal(err)
    }
	return ms
}

func (m MemStorage) DownloadMemStorage(file string) {
	producerPtr, err := NewProducer(file)
	if err != nil {
		fmt.Println("store; line: 77; DownloadMemStorage file: ", file, "   ;err: ", err)
        log.Fatal(err)
    }
	err = producerPtr.WriteMemStorage(&m)
	if err != nil {
        log.Fatal(err)
    }
}

func NewProducer(filename string) (*producer, error) {
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	file.Truncate(0)
    // file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
    if err != nil {
        return nil, err
    }

    return &producer{
        file: file,
        // создаём новый Writer
        writer: bufio.NewWriter(file),
    }, nil
}

func (p *producer) WriteMemStorage(memStorage *MemStorage) error {
    data, err := json.Marshal(&memStorage)
    if err != nil {
        return err
    }

    // записываем событие в буфер
    if _, err := p.writer.Write(data); err != nil {
        return err
    }

    // добавляем перенос строки
    if err := p.writer.WriteByte('\n'); err != nil {
        return err
    }

    // записываем буфер в файл
    return p.writer.Flush()
}
type consumer struct {
    file *os.File
    // заменяем reader на scanner
    scanner *bufio.Scanner
}

func NewConsumer(filename string) (*consumer, error) {
    file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
    // file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
    if err != nil {
        return nil, err
    }

    return &consumer{
        file: file,
        // создаём новый scanner
        scanner: bufio.NewScanner(file),
    }, nil
}

func (c *consumer) ReadMemStorage() (*MemStorage, error) {
    // одиночное сканирование до следующей строки
    if !c.scanner.Scan() {
        return nil, c.scanner.Err()
    }
    // читаем данные из scanner
    data := c.scanner.Bytes()

    memStorage := MemStorage{}
    err := json.Unmarshal(data, &memStorage)
    if err != nil {
        return nil, err
    }

    return &memStorage, nil
}

func (c *consumer) Close() error {
    return c.file.Close()
}
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
	// fmt.Println("OK:  ", ok, "-----------------********************************************-----------------")

	if !ok {
		mtxNew := NewMetrics(id, "Gauge")
		mtxNew.Value = &gauge
		m.StoreMap[id] = mtxNew
	}else {
		// fmt.Println("*mtxOld:  ", mtxOld, "----------------------------------")
		// fmt.Println("*mtxOld.Value:  ", *mtxOld.Value, "----------------------------------")
		*mtxOld.Value = gauge
		
	}
}

func (m *MemStorage) GetCounterValue(id string) (Counter, error) {
	mtxOld, ok := m.StoreMap[id]
	if !ok {
		return 0, errors.New("this counter don't find")
	}
	return *mtxOld.Delta, nil
}

func (m *MemStorage) GetGaugeValue(id string) (Gauge, error) {
	// fmt.Println("id:  ", id, "-----------------------((((((((((((((((((((((((((((((((((((")
	mtxOld, ok := m.StoreMap[id]
	// fmt.Println("mtxOld:  ", mtxOld, ";  OK: ", ok, "MemStorage", m)
	if !ok {
		return 0, errors.New("this gauge don't find")
	}
	// fmt.Println("*mtxOld.Value:  ", *mtxOld.Value, "----------------------------------")
	return *mtxOld.Value, nil
}

func (m *MemStorage) GetAllValues() string {
	var str string
	for key, val := range m.StoreMap {
		str += fmt.Sprintf("%s: %v %v\n", key, val.Value, val.Delta)

	}
	return str
}
