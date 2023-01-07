package handler

import (
	"encoding/json"
	"fmt"
	"io"

	"context"
	// "log"
	"net/http"
	"strconv"
	// "time"

	"github.com/go-chi/chi/v5"
	// "github.com/kokdot/go-musthave-devops/internal/interface_init"
	"github.com/kokdot/go-musthave-devops/internal/metrics_server"
	"github.com/kokdot/go-musthave-devops/internal/store"
	"github.com/kokdot/go-musthave-devops/internal/repo"
)
 
type keyData int

const (
	nameDataKey keyData = iota
	valueDataKey
)

var m  repo.Repo

func PutM(M repo.Repo) {
	m = M
}

// var	key = m.GetKey()




func PostUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------PostUpdate-------------------------start-------------------------------")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var mtxNew metrics_server.Metrics
	err = json.Unmarshal(bodyBytes, &mtxNew)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Printf("\n----------PostUpdate------mtxNew.----:   %#v\n", mtxNew)
	if m.GetKey() != "" {
		fmt.Println("----------------------------if store.Key != ampty string-------------------------------------")
		if !metrics_server.MtxValid(&mtxNew, m.GetKey()) {
			fmt.Printf("\n-------if !store.MtxValid(&mtxNew).----:   %#v\n", mtxNew)
			
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
    }
	// fmt.Println("----------------------------------------------------------------------------")
	// fmt.Printf("\n----------PostUpdate------mtxNew.----:   %#v\n", mtxNew)
	// fmt.Println("-126-handler--PostUpdate-----m/txNew-------------------", mtxNew)
    // if mtxNew.Delta != nil {
    // 	fmt.Println("--199--handler---PostUpdate-----mtxNew-------------------", mtxNew, "delta:  ", *mtxNew.Delta)
    // }
    // if mtxNew.Value != nil {
    //     fmt.Println("--202---handler----PostUpdate-----mtxNew-------------------", mtxNew, "Value:  ", *mtxNew.Value)
    // }
    // sm, _ := M.GetAllValues()
    // fmt.Println("----205---handler-------------/-------------------------------------------")
	// fmt.Println("--206--handler-----PostUpdate----sm-!!!!!!!!-----------------\n:  ", sm)
    // fmt.Println("---207---handler---------------------------------------------------------")
	fmt.Printf("m:   %#v", m)
	mtxOld, err := m.Save(&mtxNew)//----------------------------------------------------------------------------Save---
	// fmt.Println("---210----handler--------M.Save----------------")
    // fmt.Println("--211---handler-----PostUpdate-----mtxOld-------------------", mtxNew)
    // if mtxNew.Delta != nil {
    // fmt.Println("--213----handler------PostUpdate-----mtxOld-------------------", mtxNew, "delta:  ", *mtxNew.Delta)
    // }
    // if mtxNew.Value != nil {
    // fmt.Println("--216----handler------PostUpdate-----mtxOld-------------------", mtxNew, "Value:  ", *mtxNew.Value)
    // }
    // sm, _ = M.GetAllValues()
	// fmt.Println("--219-------handler----PostUpdate----sm-------------------", sm)
    // fmt.Println("---220-------handler---------------------------------------------------------------------------")
    // fmt.Println("----221--------handler-------------------------------------------------------------------------")
	// fmt.Printf("\n----------PostUpdate------mtxOld.----:   %#v\n", mtxOld)
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

	mtxOLd, err := m.Get(mtxNew.ID) 
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
	storeMap, err := m.GetAll()
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
	str, err := m.GetAllValues()
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
	counter, err := m.SaveCounterValue(nameData, store.Counter(valueData))//------430
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
	err := m.SaveGaugeValue(nameData, repo.Gauge(valueData))
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
	n, err := m.GetCounterValue(nameData)
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
	n, err := m.GetGaugeValue(nameData)
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
