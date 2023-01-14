package main

import (
	"log"
	"net/http"

	// "strconv"
	// "flag"

	// "github.com/caarlos0/env/v6"

	// "github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"fmt"
	// "time"

	"github.com/kokdot/go-musthave-devops/internal/handler"
	"github.com/kokdot/go-musthave-devops/internal/interfaceinit"
	"github.com/kokdot/go-musthave-devops/internal/onboardingserver"
	// "github.com/kokdot/go-musthave-devops/internal/repo"
	// "github.com/kokdot/go-musthave-devops/internal/store"
	// "github.com/kokdot/go-musthave-devops/internal/downloadingtofile"
)

//:PATH="$PATH:/mnt/c/Users/user/devopstest
// devopstest -test.v -test.run=^TestIteration2[b]*$ -source-path=. -binary-path=cmd/server/server
// devopstest -test.v -test.run=^TestIteration4$ -source-path=. -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent
// SERVER_PORT=$(random unused-port) ADDRESS="localhost:${SERVER_PORT}" TEMP_FILE=$(random tempfile) devopstest -test.v -test.run=^TestIteration7$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=$SERVER_PORT -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' -file-storage-path=$TEMP_FILE
// SERVER_PORT=$(random unused-port) ADDRESS="localhost:${SERVER_PORT}" TEMP_FILE=$(random tempfile) devopstest -test.v -test.run=^TestIteration6$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=$SERVER_PORT -database-dsn='postgres://postgres:postgrespw@localhost:49164?sslmode=disable' -file-storage-path=$TEMP_FILE
//devopstest -test.v -test.run=^TestIteration8 -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=8080 -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' -file-storage-path=azxs123
//SERVER_PORT=$(random unused-port)
//devopstest -test.v -test.run=^TestIteration9$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=8080 -file-storage-path=/tmp/wert123 -database-dsn='postgres://postgres:postgrespw@localhost:49164?sslmode=disable' -key=/tmp/wert1234

// SERVER_PORT="33658" ADDRESS="localhost:33658" TEMP_FILE="/tmp/tgy785"  devopstest -test.v -test.run=^TestIteration6$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=33658 -database-dsn='postgres://postgres:postgrespw@localhost:49164?sslmode=disable' -file-storage-path=/tmp/tgy785

// SERVER_PORT=33658 ADDRESS="localhost:33658" TEMP_FILE=jkr678 devopstest -test.v -test.run=^TestIteration10[b]*$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=33658 -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' -key="jkr678"
//func init() {
//     onboarding_server.OnboardingServer()

// }
var (
	// ms  store.MemStorage
	// m  repo.Repo
    // url = onboarding_server.GetURL()
	// storeInterval time.Duration = onboarding_server.GetStoreInterval()
	// storeFile = onboarding_server.GetStoreFile()
	// restore = onboarding_server.GetRestore()
	// key = onboarding_server.GetKey()
)
func main() { 
    url, storeFile, key, restore, storeInterval, dataBaseDSNReal  := onboardingserver.OnboardingServer()
    fmt.Println("--------------------main-------------------------------------------")
    fmt.Println("url:  ", url)
    fmt.Println("storeInterval:  ", storeInterval)
    fmt.Println("storeFile:  ", storeFile)
    fmt.Println("restore:  ", restore)
    fmt.Println("key:  ", key)
    fmt.Println("dataBaseDSNReal:  ", dataBaseDSNReal)

    m, err := interfaceinit.InterfaceInit(storeInterval, storeFile, restore, url, key, dataBaseDSNReal)
    if err != nil {
        fmt.Printf("\nthere in error in starting interface and restore data: %s", err)
    }
    handler.PutM(m)
    fmt.Printf("\nm:   %#v\n", m)
    fmt.Println("--------------------main--started-----------------------------------------")
    
    // if m.GetDataBaseDSN() != "" {
    //     downloading_to_file.DownloadingToFile(m)
    // }

    // определяем роутер chi
    r := chi.NewRouter()
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    // r.Use(middleware.Compress(5, "gzip"))
    r.Use(middleware.Compress(5))
    r.Get("/", handler.GetAll)
    r.Get("/ping", handler.GetPing)
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

    log.Fatal(http.ListenAndServe(url, r))
    // log.Fatal(http.ListenAndServe(":8080", r))
}
