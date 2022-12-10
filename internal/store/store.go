package store

import (
	"fmt"
	"errors"
)


type Counter int
type Gauge float64
type GaugeMap map[string]Gauge
type CounterMap map[string]Counter

type MemStorage struct {
	GaugeMap   GaugeMap
	CounterMap CounterMap
}

type Repo interface {
	SaveCounterValue(name string, counter Counter)
	SaveGaugeValue(name string, gauge Gauge)
	GetCounterValue(name string) (Counter, error)
	GetGaugeValue(name string) (Gauge, error)
	GetAllValues() string
}

func (m *MemStorage) SaveCounterValue(name string, counter Counter) {
	n, ok := m.CounterMap[name]
	if !ok {
		m.CounterMap[name] = counter
		return
	}
	m.CounterMap[name] = n + counter
}

func (m *MemStorage) SaveGaugeValue(name string, gauge Gauge) {
	m.GaugeMap[name] = gauge
}

func (m *MemStorage) GetCounterValue(name string) (Counter, error) {
	n, ok := m.CounterMap[name]
	if !ok {
		return 0, errors.New("this counter don't find")
	}
	return n, nil
}

func (m *MemStorage) GetGaugeValue(name string) (Gauge, error) {
	n, ok := m.GaugeMap[name]
	if !ok {
		return 0, errors.New("this gauge don't find")
	}
	return n, nil
}

func (m *MemStorage) GetAllValues() string {
	mapAll := make(map[string]string)
	for key, val := range m.CounterMap {
		mapAll[key] = fmt.Sprintf("%v", val)
	}
	for key, val := range m.GaugeMap {
		mapAll[key] = fmt.Sprintf("%v", val)
	}
	var str string
	for key, val := range mapAll{
		str += fmt.Sprintf("%s: %s\n", key, val)
	}
	return str
}
