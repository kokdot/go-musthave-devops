package handler

import (
	"encoding/json"
	"fmt"
	"io"

	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kokdot/go-musthave-devops/internal/store"
)
 
type key int

const (
	nameDataKey key = iota
	valueDataKey
)

var (
	M  store.Repo
	StoreInterval time.Duration// = StoreInterval
)

func InterfaceInit(storeInterval time.Duration, storeFile string, restore bool) {
    StoreInterval = storeInterval
	if storeInterval > 0 {
		if storeFile == "" {
			ms, err := store.NewMemStorage()
            if err != nil {
				log.Fatalf("failed to create MemStorage, err: %s", err)
			}
			M = ms
		} else {
			ms, err := store.NewMemStorageWithFile(storeFile)
			if err != nil {
				log.Fatalf("failed to create MemStorage, err: %s", err)
			}
			M = ms
			if restore {
				_, err := M.ReadStorage()
				if err != nil {
					log.Printf("Can't to read data froM file, err: %s", err)
				}
			}
			DownloadingToFile()
			M = ms
		}
	} else {
		if storeFile == "" {
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
				_, err := M.ReadStorage()
				if err != nil {
					log.Printf("Can't to read data from file, err: %s", err)
				}
			}
		}
	}
}



func DownloadingToFile() {
	// fmt.Println("---------DownloadingToFile m: -------------------", M)
	// fmt.Println("---------DownloadingToFile-------------------", storeInterval)

	go func() {
		// var interval = StoreInterval
		// var interval = time.Duration(storeInterval) * time.Second
		for {
			<-time.After(StoreInterval)
			fmt.Println("main; line: 67; DownloadToFile", ";  file:  ")
			err := M.WriteStorage()
			if err != nil {
				log.Printf("StoreMap did not been saved in file, err: %s", err)
			}
		}
	}()
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------PostUpdate-------------------------start-------------------------------")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var mtxNew store.Metrics
	err = json.Unmarshal(bodyBytes, &mtxNew)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Printf("\n----------PostUpdate------mtxNew.----:   %#v\n", mtxNew)
	if store.Key != "" {
		fmt.Println("----------------------------if store.Key != ampty string-------------------------------------")
		if !store.MtxValid(&mtxNew) {
			fmt.Printf("\n-------if !store.MtxValid(&mtxNew).----:   %#v\n", mtxNew)
			
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
    }
	// fmt.Println("----------------------------------------------------------------------------")
	// fmt.Printf("\n----------PostUpdate------mtxNew.----:   %#v\n", mtxNew)
	mtxOld, err := M.Save(&mtxNew)
	fmt.Printf("\n----------PostUpdate------mtxOld.----:   %#v\n", mtxOld)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyBytes, err = json.Marshal(mtxOld)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}

func GetValue(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------GetValue-------------------------start-------------------------------")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var mtxNew store.Metrics
	err = json.Unmarshal(bodyBytes, &mtxNew)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("\n----------GetValue------mtxNew.----:   %#v\n", mtxNew)

	mtxOLd, err := M.Get(mtxNew.ID) 
	// fmt.Println("----------------------------------------------------------------------------")
	fmt.Printf("\n----------GetValue------mtxOLd.----:   %#v\n", mtxOLd)
	if err != nil {
        fmt.Println("-----------------------------------err line 274, err:  ", err)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bodyBytes, err = json.Marshal(mtxOLd)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}

func GetAllJSON(w http.ResponseWriter, r *http.Request) {
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
	str, err := M.GetAllValues()
	if err != nil {
		w.Header().Set("content-type", "test/html")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(str))
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
			return
		}
		nameData = nameDataStr
		valueData, err := strconv.Atoi(valueDataStr)
		if err != nil {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
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
			return
		}
		nameData = nameDataStr
		valueData, err := strconv.ParseFloat(valueDataStr, 64)
		if err != nil {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
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
	counter, err := M.SaveCounterValue(nameData, store.Counter(valueData))//------430
    if err != nil {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, counter)
}
func PostUpdateGauge(w http.ResponseWriter, r *http.Request) {
	valueData := r.Context().Value(valueDataKey).(float64)
	nameData := r.Context().Value(nameDataKey).(string)
	err := M.SaveGaugeValue(nameData, store.Gauge(valueData))
    if err != nil {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, valueData)
}
func GetCounter(w http.ResponseWriter, r *http.Request) {
	nameData := r.Context().Value(nameDataKey).(string)
	n, err := M.GetCounterValue(nameData)
	if err != nil {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
	} else {
	    w.Header().Set("content-type", "text/html")
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
	} else {
	    w.Header().Set("content-type", "text/html")
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", n)
	}
}
