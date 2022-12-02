package main

import (
    "fmt"
    "net/http"
	"log"
	"strings"
	"strconv"
)
// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
type GaugeMap map[string]float64

var guageMap = GaugeMap{}
func main() {

	http.HandleFunc("/update/", Handler)
    fmt.Println("Server started at port 8080")

    log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	sliceURLPath := strings.Split(urlPath, "/")
	if n, err := strconv.ParseFloat(sliceURLPath[4], 64); err == nil {
        guageMap[sliceURLPath[3]] = n
        }
	// fmt.Fprintf(w, "Hello, %v, %T", guageMap, guageMap)
	w.WriteHeader(http.StatusOK)
}

