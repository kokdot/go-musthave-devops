package handler

import (
	"fmt"
	"encoding/json"
    "io"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"net/http"
	"context"
	"github.com/kokdot/go-musthave-devops/internal/store"
	"strconv"
    "log"
    // "strconv"
)

type key int

const (
    nameDataKey key = iota
    valueDataKey
)
type Config struct {
    StoreInterval  int	`env:"STORE_INTERVAL" envDefault:"300"`
    StoreFile  string 		`env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
}

var m store.Repo
var ms = new(store.MemStorage)
var SyncDownload bool
var cfg Config
var StoreFile string  = "/tmp/devops-metrics-db.json"
var StoreInterval int

func init() {
	ms.StoreMap = make(store.StoreMap)
    m = ms
    err := env.Parse(&cfg)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("handler init():  %+v\n", cfg)

    if cfg.StoreFile != "" &&  cfg.StoreInterval == 0 {
        SyncDownload = true
    }
}

func CheckSyncDownload(file string) {
    SyncDownload = true
    StoreFile = file
    fmt.Println("SyncDownload:  ", SyncDownload, "; StoreFile:  ", StoreFile)
}

func DownloadMemStorageToFile(file string) {
    fmt.Println("handler; line: 53; DownloadMemStorageToFile: ", m,"; SyncDownload:  ", SyncDownload, "   ;file: ", file, "; m.GetAllValues:  ", m.GetAllValues())
    m.DownloadMemStorage(file)
}
func UpdateMemStorageFromFile(file string) {
    fmt.Println("handler; line: 57; UpdateMemStorageFromFile", m, "   ;file: ", file)
    m = m.UpdateMemStorage(file)
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        // fmt.Fprint(w, "http.StatusBadRequest")
        return
    }
    var mtxNew store.Metrics
    err = json.Unmarshal(bodyBytes, &mtxNew)
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        // fmt.Fprint(w, "http.StatusBadRequest")
        return
    }
    mtxOld := m.Save(&mtxNew)
    bodyBytes, err = json.Marshal(mtxOld)
     if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        // fmt.Fprint(w, "http.StatusBadRequest")
        return
    }
    w.Header().Set("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    // fmt.Fprintf(w, "%v", bodyBytes) 
    w.Write(bodyBytes)
    if SyncDownload {
        DownloadMemStorageToFile(StoreFile)
    }
}

func GetValue(w http.ResponseWriter, r *http.Request) {
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        // fmt.Fprint(w, "http.StatusBadRequest")
        return
    }
    var mtxNew store.Metrics
    err = json.Unmarshal(bodyBytes, &mtxNew)
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        // fmt.Fprint(w, "http.StatusBadRequest")
        return
    }
    mtxOLd, err := m.Get(mtxNew.ID)
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        // fmt.Fprint(w, "http.StatusBadRequest")
        return
    }
    bodyBytes, err = json.Marshal(mtxOLd)
     if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        // fmt.Fprint(w, "http.StatusBadRequest")
        return
    }
    w.Header().Set("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    // fmt.Fprintf(w, "%v", bodyBytes) 
    w.Write(bodyBytes)
}


func GetAllJson(w http.ResponseWriter, r *http.Request) {
    storeMap := m.GetAll()
    bodyBytes, err := json.Marshal(storeMap)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    fmt.Println(string(bodyBytes))
    w.Header().Set("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(bodyBytes)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
    storeMap := m.GetAll()
    bodyBytes, err := json.Marshal(storeMap)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    fmt.Println(string(bodyBytes))
    w.Header().Set("content-type", "text/plain; charset=utf-8")
    // w.Header().Set("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "%s", bodyBytes) 
}


func PostCounterCtx(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var nameData string
        var valueData int

		nameDataStr := chi.URLParam(r, "nameData")
		valueDataStr := chi.URLParam(r, "valueData")

        if nameDataStr == "" || valueDataStr == "" {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
		    w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "http.StatusNotFound")
            return
        }
        nameData = nameDataStr
        valueData, err := strconv.Atoi(valueDataStr)
        if err != nil {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
            fmt.Fprint(w, "http.StatusBadRequest")
            return
        }

		ctx := context.WithValue(r.Context(), nameDataKey, nameData)
		ctx = context.WithValue(ctx, valueDataKey, valueData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetCtx(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var nameData string
		nameDataStr := chi.URLParam(r, "nameData")
        if nameDataStr == "" {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
		    w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "line: 115; http.StatusNotFound")
            return
        }
        nameData = nameDataStr

		ctx := context.WithValue(r.Context(), nameDataKey, nameData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PostGaugeCtx(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var nameData string
        var valueData float64

		nameDataStr := chi.URLParam(r, "nameData")
		valueDataStr := chi.URLParam(r, "valueData")

        if nameDataStr == "" || valueDataStr == "" {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
		    w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "http.StatusNotFound")
            return
        }
        nameData = nameDataStr
        valueData, err := strconv.ParseFloat(valueDataStr, 64)
        if err != nil {
            w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
            fmt.Fprint(w, "http.StatusBadRequest")
            return
        }

		ctx := context.WithValue(r.Context(), nameDataKey, nameData)
		ctx = context.WithValue(ctx, valueDataKey, valueData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func PostUpdateCounter(w http.ResponseWriter, r *http.Request) {
	valueData := r.Context().Value(valueDataKey).(int)
	nameData := r.Context().Value(nameDataKey).(string)
	// fmt.Println("__________________________________", m)
    m.SaveCounterValue(nameData, store.Counter(valueData))
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "http.StatusOK")
}
func PostUpdateGauge(w http.ResponseWriter, r *http.Request) {
	valueData := r.Context().Value(valueDataKey).(float64)
	nameData := r.Context().Value(nameDataKey).(string)
    m.SaveGaugeValue(nameData, store.Gauge(valueData))
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "http.StatusOK")
}
func GetCounter(w http.ResponseWriter, r *http.Request) {
    nameData := r.Context().Value(nameDataKey).(string)
    n, err := m.GetCounterValue(nameData)
    if err != nil {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, "line: 175; http.StatusNotFound")
    } else {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "%v", n)
    }
}
func GetGauge(w http.ResponseWriter, r *http.Request) {
    nameData := r.Context().Value(nameDataKey).(string)
    n, err := m.GetGaugeValue(nameData)
    if err != nil {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, "line: 188; http.StatusNotFound, error: ", err)
    } else {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "%v", n)
    }    
}

// func GetAll(w http.ResponseWriter, r *http.Request) {
//     w.Header().Set("content-type", "text/plain; charset=utf-8")
//     w.WriteHeader(http.StatusOK)
//     fmt.Fprintf(w, "%v", m.GetAllValues()) 
// }

// func PostUpdate(w http.ResponseWriter, r *http.Request) {
//     var metrics Metrics
//         // fmt.Println(string(r, "-----------------------------------------------------------")

//     bodyBytes, err := io.ReadAll(r.Body)
//     if err != nil {
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusNotFound)
//         // fmt.Fprint(w, "http.StatusBadRequest")
//         return
//     }
//     err = json.Unmarshal(bodyBytes, &metrics)
//     if err != nil {
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusNotFound)
//         // fmt.Fprint(w, "http.StatusBadRequest")
//         return
//     }
//     switch metrics.MType  {
//     case "Gauge":
//         m.SaveGaugeValue(metrics.ID, store.Gauge(*metrics.Value))
         
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusOK)
//         // fmt.Fprintf(w, "%v", bodyBytes) 
//         w.Write(bodyBytes)
//     case "Counter":
//         delta := m.SaveCounterValue(metrics.ID, store.Counter(*metrics.Delta))
//         *metrics.Delta = int64(delta)
//         bodyBytes, err := json.Marshal(metrics)
//         if err != nil {
//             w.Header().Set("content-type", "application/json")
//             w.WriteHeader(http.StatusNotFound)
//             // fmt.Fprint(w, "http.StatusBadRequest")
//             return
//         }
//         // fmt.Println(string(bodyBytes), "-----------------------------------------------------------")
//         // fmt.Println(bodyBytes, "-----------------------------------------------------------")
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusOK)
//         // fmt.Fprintf(w, "%v", bodyBytes) 
//         w.Write(bodyBytes)
//     default:
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusNotFound)
//         // fmt.Fprint(w, "http.StatusBadRequest")
//         return
//     }
// }
// func GetValue(w http.ResponseWriter, r *http.Request) {
//     bodyBytes, err := io.ReadAll(r.Body)
//     if err != nil {
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusNotFound)
//         // fmt.Fprint(w, "http.StatusBadRequest")
//         return
//     }
//     var metrics Metrics
//     err = json.Unmarshal(bodyBytes, &metrics)
//     if err != nil {
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusNotFound)
//         // fmt.Fprint(w, "http.StatusBadRequest")
//         return
//     }
//     switch metrics.MType  {
//     case "Gauge":
//         gaugeValue, err := m.GetGaugeValue(metrics.ID)
//         if err != nil {
//             w.Header().Set("content-type", "application/json")
//             w.WriteHeader(http.StatusNotFound)
//             // fmt.Fprint(w, "http.StatusBadRequest")
//             return
//         }
//         fmt.Println("gaugeValue: ", gaugeValue, "   ;float64(gaugeValue): ", float64(gaugeValue), "    ;metrics: ", metrics)
//         gaugeValue1 := float64(gaugeValue)
//         metrics1 := Metrics{
// 				ID: metrics.ID,
// 				MType: metrics.MType,
// 				Value: &gaugeValue1,
// 			}
//         // *metrics.Value = gaugeValue1
//         bodyBytes, err := json.Marshal(metrics1)
//          if err != nil {
//             w.Header().Set("content-type", "application/json")
//             w.WriteHeader(http.StatusNotFound)
//             // fmt.Fprint(w, "http.StatusBadRequest")
//             return
//         }
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusOK)
//         w.Write(bodyBytes)
//         // fmt.Fprintf(w, "%v", bodyBytes) 
//     case "Counter":
//         delta, err := m.GetCounterValue(metrics.ID)
//          if err != nil {
//             w.Header().Set("content-type", "application/json")
//             w.WriteHeader(http.StatusNotFound)
//             // fmt.Fprint(w, "http.StatusBadRequest")
//             return
//         }
//         // fmt.Println("gaugeValue: ", gaugeValue, "   ;float64(gaugeValue): ", float64(gaugeValue), "    ;metrics: ", metrics)
//         counterValue1 := int64(delta)
//         metrics1 := Metrics{
// 				ID: metrics.ID,
// 				MType: metrics.MType,
// 				Delta: &counterValue1,
// 			}
//         // *metrics.Value = gaugeValue1
//         bodyBytes, err := json.Marshal(metrics1)
//          if err != nil {
//             w.Header().Set("content-type", "application/json")
//             w.WriteHeader(http.StatusNotFound)
//             // fmt.Fprint(w, "http.StatusBadRequest")
//             return
//         }
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusOK)
//         w.Write(bodyBytes) 
//     default:
//         w.Header().Set("content-type", "application/json")
//         w.WriteHeader(http.StatusNotFound)
//         // fmt.Fprint(w, "http.StatusBadRequest")
//         return
//     }
// }
