package handler

import (
	"net/http"
	"strings"
	"fmt"
	"strconv"
	"github.com/kokdot/go-musthave-devops/internal/store"
)
var m = new(stroe.MemStorage)
func HandlerPostUpdate(w http.ResponseWriter, r *http.Request) {
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
            w.WriteHeader(http.StatusNotFound)
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
func getGaugeValue(name string) (store.Gauge, error) {
    n, err := m.getGauge(name)
    if err != nil {
        return 0, err
    }
    return n, nil
}
func getCountValue(name string) (store.Counter, error) {
    n, err := m.getCounter(name)
    if err != nil {
        return 0, err
    }
    return n, nil
}
func GetAllValues(w http.ResponseWriter, r *http.Request){
    w.Header().Set("content-type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
	w.Write(m.GetAll)
	// fmt.Fprintf(w, "Poll: %v", counterSlice[len(counterSlice) - 1])
}