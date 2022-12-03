package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

)

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
type GaugeMap map[string]float64
type CounterMap map[string][]int

var guageMap = GaugeMap{}
var counterMap = CounterMap{}

func main() {

	http.HandleFunc("/update/", Handler)
    fmt.Println("Server started at port 8080")

    log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	sliceURLPath := strings.Split(urlPath, "/")

	switch  {
	case len(sliceURLPath) != 5:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "http.StatusNotFound")

	case sliceURLPath[2] == "guage":
		n, err := strconv.ParseFloat(sliceURLPath[4], 64)
		if  err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "http.StatusBadRequest")
		} else {
			guageMap[sliceURLPath[3]] = n
			w.WriteHeader(http.StatusOK)
			// fmt.Fprint(w, "http.StatusOK")
			fmt.Fprintf(w, "http.StatusOK: %v; sliceURLPath: %v; method: %v",http.StatusOK, sliceURLPath, r.Method)

		}
	case sliceURLPath[2] == "counter":
		n, err := strconv.Atoi(sliceURLPath[4])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "http.StatusBadRequest: %v; sliceURLPath: %v; method: %v",http.StatusBadRequest, sliceURLPath, r.Method)
		} else {
			counterMap[sliceURLPath[3]] = append(counterMap[sliceURLPath[3]], n)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "http.StatusOK")
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "http.StatusNotFound")
	}
}

