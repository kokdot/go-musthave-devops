package main

import (


	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"

	"fmt"
	"github.com/kokdot/go-musthave-devops/internal/handler"
)
const (
    url = "127.0.0.1:8080"
)
type Config struct {
    Address  string 		`env:"ADDRESS"`
}

func main() {
    var cfg Config
	var urlReal = url
    err := env.Parse(&cfg)
    if err != nil {
        log.Fatal(err)
    }
	if cfg.Address != ""{
		urlReal	= cfg.Address
	} 
    // определяем роутер chi
    r := chi.NewRouter()
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Get("/", handler.GetAll)
    r.Route("/update", func(r chi.Router) {
        r.Post("/", handler.PostUpdate)
        r.Route("/counter", func(r chi.Router) {
            r.Route("/{nameData}/{valueData}", func(r chi.Router) {
                r.Use(handler.PostCounterCtx)
                r.Post("/", handler.PostUpdateCounter)
            })
        })
        r.Route("/gauge", func(r chi.Router) {
            r.Route("/{nameData}/{valueData}", func(r chi.Router) {
                r.Use(handler.PostGaugeCtx)
                r.Post("/", handler.PostUpdateGauge)
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
        r.Post("/", handler.GetValue)
		r.Route("/counter", func(r chi.Router){
            r.Route("/{nameData}", func(r chi.Router) {
                r.Use(handler.GetCtx)
                r.Get("/", handler.GetCounter)
            })
        })
       	r.Route("/gauge", func(r chi.Router){
            r.Route("/{nameData}", func(r chi.Router) {
                r.Use(handler.GetCtx)
                r.Get("/", handler.GetGauge)
            })
        })
	})

    log.Fatal(http.ListenAndServe(urlReal, r))
    // log.Fatal(http.ListenAndServe(":8080", r))
}
