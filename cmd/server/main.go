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

var gaugeMap = GaugeMap{}
var counterMap = CounterMap{}

func main() {

	http.HandleFunc("/update/", Handler)
	fmt.Println("Server started at port 8080")

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	sliceURLPath := strings.Split(urlPath, "/")

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
			counterMap[sliceURLPath[3]] = append(counterMap[sliceURLPath[3]], n)
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "n, err := strconv.Atoi(sliceURLPath[4]) err == nil; http.StatusOK: %v; sliceURLPath: %v; method: %v", http.StatusOK, sliceURLPath, r.Method)
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