package store

import (
	"errors"
	"fmt"
)

type Counter int
type Gauge float64
type GaugeMap map[string]Gauge
type CounterMap map[string]Counter

type MemStorage struct {
	gaugeMap   GaugeMap
	counterMap CounterMap
}

type Repo interface {
	saveCounter(name string, counter Counter)
	saveGauge(name string, gauge Gauge)
	getCounter(name string) (Counter, error)
	getGauge(name string) (Gauge, error)
	GetAllValues() string
}

func (m MemStorage) saveCounter(name string, counter Counter) {
	n, ok := m.counterMap[name]
	if !ok {
		m.counterMap[name] = counter
		return
	}
	m.counterMap[name] = n + counter
}

func (m MemStorage) saveGauge(name string, gauge Gauge) {
	m.gaugeMap[name] = gauge
}

func (m MemStorage) getCounter(name string) (Counter, error) {
	n, ok := m.counterMap[name]
	if !ok {
		return 0, errors.New("this counter don't find")
	}
	return n, nil
}

func (m MemStorage) getGauge(name string) (Gauge, error) {
	n, ok := m.gaugeMap[name]
	if !ok {
		return 0, errors.New("this gauge don't find")
	}
	return n, nil
}

func (m MemStorage) GetAllValues() string {
	mapAll := make(map[string]string)
	for key, val := range m.counterMap {
		mapAll[key] = fmt.Sprintf("%v", val)
	}
	for key, val := range m.gaugeMap {
		mapAll[key] = fmt.Sprintf("%v", val)
	}
	var str string
	for key, val := range mapAll{
		str += fmt.Sprintf("%s: %s\n", key, val)
	}
	return str
}