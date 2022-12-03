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
		fmt.Fprint(w, "")

	case sliceURLPath[2] == "Guage":
		n, err := strconv.ParseFloat(sliceURLPath[4], 64)
		if  err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "")
		} else {
			guageMap[sliceURLPath[3]] = n
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "")
		}

	case sliceURLPath[2] == "Counter":
		n, err := strconv.Atoi(sliceURLPath[4])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "")
		} else {
			counterMap[sliceURLPath[3]] = append(counterMap[sliceURLPath[3]], n)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "")
		}
	}
	// if sliceURLPath[2] == "Guage" {
		
	// 	if n, err := strconv.ParseFloat(sliceURLPath[4], 64); err == nil {
	// 		guageMap[sliceURLPath[3]] = n
    //     }
	// }
	// if sliceURLPath[2] == "Counter" {
	// 	if n, err := strconv.Atoi(sliceURLPath[4]); err == nil {
	// 		counterMap[sliceURLPath[3]] = append(counterMap[sliceURLPath[3]], n)
	// 	}
	// }
	// fmt.Fprintf(w, "Hello, %v, %T", guageMap, guageMap)
	// w.WriteHeader(http.StatusOK)
}

