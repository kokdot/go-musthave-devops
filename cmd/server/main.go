package main

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    // "io"
    "log"
    "net/http"
    "fmt"
    "strings"
    "strconv"
    "errors"
)
type GaugeMap map[string]float64
type CounterSlice []int
type MapCounterSlice map[string]CounterSlice

var errNotFound  = errors.New("not found")
var gaugeMap = GaugeMap{}
var counterSlice = CounterSlice{}
var mapCounterSlice = MapCounterSlice{}

func Handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	sliceURLPath := strings.Split(urlPath, "/")
    // w.Header().Set("content-type", "text/plain; charset=utf-8")
	// w.WriteHeader(http.StatusNotFound)
	// fmt.Fprintf(w, "len(sliceURLPath) != 5; http.StatusNotFound: %v; sliceURLPath: %v; method: %v", http.StatusNotFound, sliceURLPath, r.Method)

	switch {
	case len(sliceURLPath) != 5:
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "len(sliceURLPath) != 5; http.StatusNotFound: %v; sliceURLPath: %v; method: %v", http.StatusNotFound, sliceURLPath, r.Method)
		// fmt.Fprint(w, "http.StatusNotFound")

	case sliceURLPath[2] == "gauge":
		// fmt.Println("case sliceURLPath[2] == \"gauge\":")
		n, err := strconv.ParseFloat(sliceURLPath[4], 64)
		if err != nil {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			// fmt.Fprint(w, "http.StatusBadRequest")
			fmt.Fprintf(w, "n, err := strconv.ParseFloat(sliceURLPath[4], 64) err != nil; http.StatusBadRequest: %v; sliceURLPath: %v; method: %v", http.StatusBadRequest, sliceURLPath, r.Method)
		} else {
			gaugeMap[sliceURLPath[3]] = n
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			// fmt.Fprint(w, "http.StatusOK")
			fmt.Fprintf(w, "n, err := strconv.ParseFloat(sliceURLPath[4], 64) err == nil; http.StatusOK: %v; sliceURLPath: %v; method: %v", http.StatusOK, sliceURLPath, r.Method)

		}
	case sliceURLPath[2] == "counter":
		n, err := strconv.Atoi(sliceURLPath[4])
		if err != nil {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "n, err := strconv.Atoi(sliceURLPath[4]) err != nil; http.StatusBadRequest: %v; sliceURLPath: %v; method: %v", http.StatusBadRequest, sliceURLPath, r.Method)
		} else {
			mapCounterSlice[sliceURLPath[3]] = append(mapCounterSlice[sliceURLPath[3]], n)
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "line 63; mapCounterSlice: %v;  http.StatusOK: %v; sliceURLPath: %v; method: %v", mapCounterSlice, http.StatusOK, sliceURLPath, r.Method)
			// fmt.Fprint(w, "http.StatusOK")
		}
	case sliceURLPath[2] != "counter" && sliceURLPath[2] != "gauge":
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(w, "http.StatusNotImplemented")
	default:
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "default: ;sliceURLPath[2] = %v; http.StatusNotFound: %v; sliceURLPath: %v; method: %v", sliceURLPath[2], http.StatusNotFound, sliceURLPath, r.Method)
		// fmt.Fprint(w, "http.StatusNotFound")
	}

}
func getCount(w http.ResponseWriter, r *http.Request) {
    urlPath := r.URL.Path
    sliceURLPath := strings.Split(urlPath, "/")
    if len(sliceURLPath) != 4 {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "%v", "")
    } else {
        n, err := getCountValue(sliceURLPath[3])
        
        if err != nil {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "line: 92;mapCounterSlice: %v; %v",sliceURLPath[3], n)

        } else {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "%v", n)
    
        }
    }
}
func getGauge(w http.ResponseWriter, r *http.Request) {
    urlPath := r.URL.Path
    sliceURLPath := strings.Split(urlPath, "/")
	if len(sliceURLPath) != 4{
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "%v", "")
    } else {
        n, err := getGaugeValue(sliceURLPath[3])
        
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
func getGaugeValue(name string) (float64, error) {
    n, ok := gaugeMap[name]
    if !ok {
        return 0, errNotFound
    }
    return n, nil
}
func getCountValue(name string) (int, error) {
    slice, ok := mapCounterSlice[name]
    if !ok {
        return 0, errNotFound
    }
    return slice[len(slice) - 1], nil
}
func GetAll(w http.ResponseWriter, r *http.Request){
    w.Header().Set("content-type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    for key, val := range gaugeMap {
        fmt.Fprintf(w, "%v: %v\n", key, val)
    }
	 for key := range mapCounterSlice {
        fmt.Fprintf(w, "%v: %v\n", key, mapCounterSlice[key][len(mapCounterSlice[key]) - 1])
    }
}

func main() {
    gaugeMap["Alloc"] = 1234.6
    // определяем роутер chi
    r := chi.NewRouter()
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Get("/", GetAll)
    r.Route("/update", func(r chi.Router) {
            r.Post("/*", Handler)})

    r.Route("/value", func(r chi.Router) {
		r.Route("/count", func(r chi.Router){
            r.Get("/*", getCount)
		})
        r.Route("/gauge", func(r chi.Router){
            r.Get("/*", getGauge)				
        })
			
	})

    log.Fatal(http.ListenAndServe(":8080", r))
}