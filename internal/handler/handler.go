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
    "time"
    "flag"
    // "strconv"
)

type key int

const (
    nameDataKey key = iota
    valueDataKey
)
const (
    url = "127.0.0.1:8080"
    StoreInterval = 200
    StoreFile = "/tmp/devops-metrics-db.json"
    Restore = false
)
type Config struct {
    Address  string 		`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
    StoreInterval  int 		`env:"STORE_INTERVAL" envDefault:"30"`
    StoreFile  string 		`env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
    Restore  bool 		`env:"RESTORE" envDefault:"true"`
}
var (
    M store.Repo 
    ms = new(store.MemStorage)
    UrlReal = url
	storeInterval = StoreInterval
	storeFile = StoreFile
	restore = Restore
    cfg Config
)

func init() {
	fmt.Println("---------init-------------------")
    onboarding()
}
func InterfaceInit() {
     if storeInterval > 0 {
        //---------save to MemStorage-------------------
        // fmt.Println("---------save to MemStorage-------------------")
        if storeFile == ""{
        // fmt.Println("---------storeFile == \"\"-------------------")

            ms := store.NewMemStorage()
            M = ms
        } else {
        // fmt.Println("---------storeFile == \"\" else------------------")

            ms, err := store.NewMemStorageWithFile(storeFile)
            // fmt.Println("---------ms-------------------", ms)

            if err != nil {
                log.Fatalf("failed to create MemStorage, err: %s", err)
            }
            M = ms
            // fmt.Println("---------m-------------------", M, "restore:  :", restore)

            if restore {
                // fmt.Println("---------restore-------------------", M)

                _ , err := M.ReadStorage()
                fmt.Println("---------M.ReadStorage()-------------------", M)

                if err != nil {
                    log.Printf("Can't to read data froM file, err: %s", err)
                }
            }
            // DownloadingToFile()
            M = ms
        }
    } else {
        //---------save to FileStorage------------------
        fmt.Println("---------save to FileStorage-------------------")

        if storeFile == ""{
            ms, err := store.NewFileStorage()
            if err != nil {
                log.Fatalf("failed to create FileStorage, err: %s", err)
            }
            M = ms
        } else {
            ms, err := store.NewFileStorageWithFile(storeFile)
            if err != nil {
                log.Fatalf("failed to create FileStorage, err: %s", err)
            }
            M = ms
            if restore {
                _ , err := M.ReadStorage()
                if err != nil {
                    log.Printf("Can't to read data from file, err: %s", err)
                }
            }
            // DownloadingToFile()
        }   
    }
    // fmt.Printf("---------init exit %+v: -------------------", M)
}

func onboarding() {
	fmt.Println("---------onboarding-------------------")
	fmt.Println("------1---storeInterval-------------------", storeInterval)

    err := env.Parse(&cfg)
    if err != nil {
        log.Print(err)
    }
    fmt.Printf("main:  %+v\n", cfg)

    urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    restorePtr := flag.Bool("r", true, "restore Metrics(Bool)")
    storeFilePtr := flag.String("f", "/tmp/devops-metrics-db.json", "file name")
    storeIntervalPtr := flag.Int("i", 300, "interval of download")

    flag.Parse()
	fmt.Println("-----2----storeInterval-------------------", storeInterval)

    fmt.Println("urlRealPrt:", *urlRealPtr)
    fmt.Println("restorePtr:", *restorePtr)
    fmt.Println("storeFilePtr:", *storeFilePtr)
    fmt.Println("storeIntervalPtr:", *storeIntervalPtr)
    UrlReal = *urlRealPtr
    storeInterval = *storeIntervalPtr
    storeFile = *storeFilePtr
    restore = *restorePtr
	fmt.Println("----3-----storeInterval-------------------", storeInterval)

    UrlReal	= cfg.Address
    storeInterval = cfg.StoreInterval
    storeFile = cfg.StoreFile
    restore = cfg.Restore
	fmt.Println("---4------storeInterval-------------------", storeInterval)

    // fmt.Println("---------onboarding m: -------------------", M)
    // fmt.Println("---------onboarding m: -------------------", ms)
	fmt.Println("---5------cfg.StoreInterval-------------------", cfg.StoreInterval)
    
    
   
    
}

func DownloadingToFile() {
    // fmt.Println("---------DownloadingToFile m: -------------------", M)
    // fmt.Println("---------DownloadingToFile-------------------", storeInterval)

    go func() {
        var interval = time.Duration(storeInterval) * time.Second
        for {
            <-time.After(interval) 
            fmt.Println("main; line: 67; DownloadToFile", ";  file:  ")
            err := M.WriteStorage()
            if err != nil {
                log.Printf("StoreMap did not been saved in file, err: %s", err)
            }
        }
    }()
}


func PostUpdate(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("------PostUpdate---m-------------------", M)

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
	// fmt.Println("---------mtxNew-------------------", mtxNew, "delta:  ", mtxNew)
	// fmt.Println("-----PostUpdate----m-------------------", M)


    mtxOld, err := M.Save(&mtxNew)
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        // fmt.Fprint(w, "http.StatusBadRequest")
        return
    }
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
    // if syncDownload {
    //     DownloadMemStorageToFile(StoreFile)
    // }
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
    mtxOLd, err := M.Get(mtxNew.ID)
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
    storeMap, err := M.GetAll()
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        return
    }
    bodyBytes, err := json.Marshal(storeMap)
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        return
    }
    fmt.Println(string(bodyBytes))
    w.Header().Set("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(bodyBytes)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
    storeMap, err := M.GetAll()
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        return
    }
    bodyBytes, err := json.Marshal(storeMap)
    if err != nil {
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusNotFound)
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
    M.SaveCounterValue(nameData, store.Counter(valueData))
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "http.StatusOK")
}
func PostUpdateGauge(w http.ResponseWriter, r *http.Request) {
	valueData := r.Context().Value(valueDataKey).(float64)
	nameData := r.Context().Value(nameDataKey).(string)
    M.SaveGaugeValue(nameData, store.Gauge(valueData))
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "http.StatusOK")
}
func GetCounter(w http.ResponseWriter, r *http.Request) {
    nameData := r.Context().Value(nameDataKey).(string)
    n, err := M.GetCounterValue(nameData)
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
    n, err := M.GetGaugeValue(nameData)
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
