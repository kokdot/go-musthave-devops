package handler

import (
	"encoding/json"
	"fmt"
	"io"

	"context"
	"net/http"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/go-chi/chi/v5"
	"github.com/kokdot/go-musthave-devops/internal/metricsserver"
	"github.com/kokdot/go-musthave-devops/internal/store"
	"github.com/kokdot/go-musthave-devops/internal/repo"
)
 
type keyData int

const (
	nameDataKey keyData = iota
	valueDataKey
)

var m  repo.Repo
var logg zerolog.Logger

func PutM(M repo.Repo) {
	m = M
}
func GetLogg(loggReal zerolog.Logger)  {
	logg = loggReal
}
func PostUpdateByBatch(w http.ResponseWriter, r *http.Request) {
	logg.Print("--------------------PostUpdateByBatch------------1-------------start-------------------------------")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
	logg.Print("--------------------PostUpdateByBatch------------2-------------start-------------------------------")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	slm := make([]repo.Metrics, 0)
	err = json.Unmarshal(bodyBytes, &slm)
	if err != nil {
	logg.Print("--------------------PostUpdateByBatch--------------3-----------start-------------------------------")
	logg.Error().Err(err).Send()	
	w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logg.Printf("Getting of requets is: %#v\n", slm)

	
	smOld, err := m.SaveByBatch(slm)
	
	logg.Printf("Answer to requets is: %#v\n", smOld)
	if err != nil {
	logg.Print("--------------------PostUpdateByBatch-------------4------------start-------------------------------")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logg.Print("--------------------PostUpdateByBatch-------------5------------start-------------------------------")
	mtxOld := repo.Metrics{}
	bodyBytes, err = json.Marshal(&mtxOld)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}
func PostUpdateByBatch1(w http.ResponseWriter, r *http.Request) {
	logg.Print("--------------------PostUpdateByBatch------------1-------------start-------------------------------")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
	logg.Print("--------------------PostUpdateByBatch------------2-------------start-------------------------------")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	smNew := make(repo.StoreMap)
	err = json.Unmarshal(bodyBytes, &smNew)
	if err != nil {
	logg.Print("--------------------PostUpdateByBatch--------------3-----------start-------------------------------")
	logg.Print(err)	
	w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logg.Printf("Getting of requets is: %#v\n", smNew)

	smOld, err := m.SaveByBatch1(&smNew)
	
	logg.Printf("Answer to requets is: %#v\n", smOld)
	if err != nil {
	logg.Print("--------------------PostUpdateByBatch-------------4------------start-------------------------------")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logg.Print("--------------------PostUpdateByBatch-------------5------------start-------------------------------")
	bodyBytes, err = json.Marshal(smOld)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}
func GetPing(w http.ResponseWriter, r *http.Request) {
	ok, err := m.GetPing()
 	if !ok {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		logg.Printf("%s", err)
		return
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}
}
func PostUpdate(w http.ResponseWriter, r *http.Request) {
	logg.Print("--------------------PostUpdate-------------------------start-------------------------------")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var mtxNew metricsserver.Metrics
	err = json.Unmarshal(bodyBytes, &mtxNew)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	logg.Printf("----------PostUpdate------mtxNew.----:   %#v", mtxNew)
	if m.GetKey() != "" {
		logg.Print("----------------------------if store.Key != ampty string-------------------------------------")
		if !metricsserver.MtxValid(&mtxNew, m.GetKey()) {
			logg.Printf("\n-------if !store.MtxValid(&mtxNew).----:   %#v\n", mtxNew)
			
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
    }
	if mtxNew.Delta != nil {
		logg.Print(" Delta = ", *mtxNew.Delta)
	}
	if mtxNew.Value != nil {
		logg.Print(" Value = ", *mtxNew.Value)
	}
	mtxOld, err := m.Save(&mtxNew)//----------------------------------------------------------------------------Save---

	if err != nil {
		logg.Print("-------after--Save-------err:   ", err)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if mtxOld.Delta != nil {
		logg.Print(" Delta = ", *mtxOld.Delta)
	}
	if mtxNew.Value != nil {
		logg.Print(" Value = ", *mtxNew.Value)
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
	logg.Print("--------------------GetValue-------------------------start-------------------------------")

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
	logg.Printf("\n----------GetValue------mtxNew.----:   %#v\n", mtxNew)

	mtxOLd, err := m.Get(mtxNew.ID) 
	logg.Printf("\n----------GetValue------mtxOLd.----:   %#v\n", mtxOLd)
	if err != nil {
        logg.Print("-----------------------------------err line 274, err:  ", err)
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
	logg.Print(string(bodyBytes))
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}
func GetAll(w http.ResponseWriter, r *http.Request) {
	str := m.GetAllValues()
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
	logg.Print("-----------------------------------------------------------------------------PostUpdateCounter-----------------")
	valueData := r.Context().Value(valueDataKey).(int)
	nameData := r.Context().Value(nameDataKey).(string)
	logg.Debug().Str("nameData", nameData).Int("ValueData", valueData).Send()
	counter, err := m.SaveCounterValue(nameData, store.Counter(valueData))
    if err != nil {
		logg.Error().Err(err).Send()
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
	logg.Debug().Str("nameData", nameData).Int("ValueData", int(counter)).Send()
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
