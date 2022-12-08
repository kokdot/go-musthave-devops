package handler

import (
	"net/http"
	"strings"
	"fmt"
	"strconv"
	"errors"
	// "github.com/kokdot/go-musthave-devops/internal/store"
)
// var m = new(store.MemStorage)
// var m store.Repo
// var m Repo
// var ms store.MemStorage
// m = ms
// m.getAll()

// func Test()




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
	PostUpdateHandler(w http.ResponseWriter, r *http.Request)
	GetCountHandler(w http.ResponseWriter, r *http.Request)
	GetGaugeHandler(w http.ResponseWriter, r *http.Request)
	GetAllHandler(w http.ResponseWriter, r *http.Request)
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



func  (m *MemStorage) PostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	sliceURLPath := strings.Split(urlPath, "/")

	switch {
	case len(sliceURLPath) != 5:
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "len(sliceURLPath) != 5; http.StatusNotFound: %v; sliceURLPath: %v; method: %v", http.StatusNotFound, sliceURLPath, r.Method)

	case sliceURLPath[2] == "gauge":
		// fmt.Println("case sliceURLPath[2] == \"gauge\":")
		n, err := strconv.ParseFloat(sliceURLPath[4], 64)
		if err != nil {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "n, err := strconv.ParseFloat(sliceURLPath[4], 64) err != nil; http.StatusBadRequest: %v; sliceURLPath: %v; method: %v", http.StatusBadRequest, sliceURLPath, r.Method)
		} else {
			m.SaveGaugeValue(sliceURLPath[3], Gauge(n))
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "n, err := strconv.ParseFloat(sliceURLPath[4], 64) err == nil; http.StatusOK: %v; sliceURLPath: %v; method: %v", http.StatusOK, sliceURLPath, r.Method)

		}
	case sliceURLPath[2] == "counter":
		n, err := strconv.Atoi(sliceURLPath[4])
		if err != nil {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "n, err := strconv.Atoi(sliceURLPath[4]) err != nil; http.StatusBadRequest: %v; sliceURLPath: %v; method: %v", http.StatusBadRequest, sliceURLPath, r.Method)
		} else {
			m.SaveCounterValue(sliceURLPath[3], Counter(n))
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "line 63; mapCounterSlice: ;  http.StatusOK: %v; sliceURLPath: %v; method: %v", http.StatusOK, sliceURLPath, r.Method)
		}
	case sliceURLPath[2] != "counter" && sliceURLPath[2] != "gauge":
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(w, "http.StatusNotImplemented")
	default:
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "default: ;sliceURLPath[2] = %v; http.StatusNotFound: %v; sliceURLPath: %v; method: %v", sliceURLPath[2], http.StatusNotFound, sliceURLPath, r.Method)
	}
}

func  (m *MemStorage) GetCountHandler(w http.ResponseWriter, r *http.Request) {
    urlPath := r.URL.Path
    sliceURLPath := strings.Split(urlPath, "/")
    if len(sliceURLPath) != 4 {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "%v", "")
    } else {
        n, err := m.GetCounterValue(sliceURLPath[3])
        if err != nil {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprintf(w, "line: 147; name: %v; %v",sliceURLPath[3], n)
        } else {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "%v", n)
        }
    }
}
func (m *MemStorage) GetGaugeHandler(w http.ResponseWriter, r *http.Request) {
    urlPath := r.URL.Path
    sliceURLPath := strings.Split(urlPath, "/")
	if len(sliceURLPath) != 4{
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "%v", "")
    } else {
        n, err := m.GetGaugeValue(sliceURLPath[3])
        if err != nil {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "%v", n)
        }
    }
}
// func (m MemStorage) getGauge(name string) (Gauge, error) {
//     n, err := m.getGauge(name)
//     if err != nil {
//         return 0, err
//     }
//     return n, nil
// }
// func (m MemStorage) getCount(name string) (Counter, error) {
//     n, err := m.getCounter(name)
//     if err != nil {
//         return 0, err
//     }
//     return n, nil
// }
func (m *MemStorage) GetAllHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("content-type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
	w.Write([]byte(m.GetAllValues()))
	// fmt.Fprintf(w, "Poll: %v", counterSlice[len(counterSlice) - 1])
}