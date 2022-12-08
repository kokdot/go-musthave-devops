package main

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    // "io"
    "log"
    "net/http"
    "fmt"
    "strings"
    "strconv"
    "errors"
	"github.com/kokdot/go-musthave-devops/internal/handler"
)
type GaugeMap map[string]float64
type CounterSlice []int
type MapCounterSlice map[string]CounterSlice

var errNotFound  = errors.New("not found")
var gaugeMap = GaugeMap{}
// var counterSlice = CounterSlice{}
var mapCounterSlice = MapCounterSlice{}




func main() {
    gaugeMap["Alloc"] = 1234.6
    // определяем роутер chi
    r := chi.NewRouter()
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Get("/", handler.GetAllValues)
    r.Route("/update", func(r chi.Router) {
            r.Post("/*", handler.HandlerPostUpdate)})

    r.Route("/value", func(r chi.Router) {
		r.Route("/counter", func(r chi.Router){
            r.Get("/*", handler.getCount)
		})
        r.Route("/gauge", func(r chi.Router){
            r.Get("/*", handler.getGauge)				
        })
			
	})

    log.Fatal(http.ListenAndServe(":8080", r))
}