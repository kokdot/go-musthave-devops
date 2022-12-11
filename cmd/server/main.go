package main

import (
	"context"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"log"
	"net/http"

	"fmt"
	"github.com/kokdot/go-musthave-devops/internal/store"
)

var m store.Repo

func main() {
	var ms = new(store.MemStorage)
	ms.GaugeMap = make(store.GaugeMap)
	ms.CounterMap = make(store.CounterMap)
	m = ms
    // определяем роутер chi
    r := chi.NewRouter()
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Get("/", GetAll)
    r.Route("/update", func(r chi.Router) {
        r.Route("/counter", func(r chi.Router) {
            r.Route("/{nameData}/{valueData}", func(r chi.Router) {
                r.Use(PostCounterCtx)
                r.Post("/", PostUpdateCounter)
            })
        })
        r.Route("/gauge", func(r chi.Router) {
            r.Route("/{nameData}/{valueData}", func(r chi.Router) {
                r.Use(PostGaugeCtx)
                r.Post("/", PostUpdateGauge)
            })
        })
        r.Route("/",func(r chi.Router) {
            r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
		        w.Header().Set("content-type", "text/plain; charset=utf-8")
                w.WriteHeader(http.StatusNotImplemented)
                fmt.Fprint(w, "line: 52; http.StatusNotImplemented")
	        })
        })
    })

    r.Route("/value", func(r chi.Router) {
		r.Route("/counter", func(r chi.Router){
            r.Route("/{nameData}", func(r chi.Router) {
                r.Use(GetCtx)
                r.Get("/", GetCounter)
            })
        })
       	r.Route("/gauge", func(r chi.Router){
            r.Route("/{nameData}", func(r chi.Router) {
                r.Use(GetCtx)
                r.Get("/", GetGauge)
            })
        })
	})

    log.Fatal(http.ListenAndServe(":8080", r))
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

		ctx := context.WithValue(r.Context(), "nameData", nameData)
		ctx = context.WithValue(ctx, "valueData", valueData)
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

		ctx := context.WithValue(r.Context(), "nameData", nameData)
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

		ctx := context.WithValue(r.Context(), "nameData", nameData)
		ctx = context.WithValue(ctx, "valueData", valueData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func PostUpdateCounter(w http.ResponseWriter, r *http.Request) {
	valueData := r.Context().Value("valueData").(int)
	nameData := r.Context().Value("nameData").(string)
    m.SaveCounterValue(nameData, store.Counter(valueData))
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "http.StatusOK")
}
func PostUpdateGauge(w http.ResponseWriter, r *http.Request) {
	valueData := r.Context().Value("valueData").(float64)
	nameData := r.Context().Value("nameData").(string)
    m.SaveGaugeValue(nameData, store.Gauge(valueData))
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "http.StatusOK")
}
func GetCounter(w http.ResponseWriter, r *http.Request) {
    nameData := r.Context().Value("nameData").(string)
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
    nameData := r.Context().Value("nameData").(string)
    n, err := m.GetGaugeValue(nameData)
    if err != nil {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, "line: 188; http.StatusNotFound")
    } else {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "%v", n)
    }    
}
func GetAll(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "%v", m.GetAllValues()) 
}