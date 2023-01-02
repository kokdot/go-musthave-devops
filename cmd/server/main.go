package main

import (
	"log"
	"net/http"
	// "strconv"
	"github.com/caarlos0/env/v6"
    "flag"

	// "github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"time"
	"fmt"

	"github.com/kokdot/go-musthave-devops/internal/handler"
	"github.com/kokdot/go-musthave-devops/internal/store"
)
// func SaveToFile() {
//     fmt.Println("---------SaveToFile  m: -------------------")

//     handler.DownloadingToFile()
// }

//:PATH="$PATH:/mnt/c/Users/user/devopstest
// devopstest -test.v -test.run=^TestIteration2[b]*$ -source-path=. -binary-path=cmd/server/server
// devopstest -test.v -test.run=^TestIteration4$ -source-path=. -binary-path=cmd/server/server -agent-binary-path=cmd/agent/agent
// SERVER_PORT=$(random unused-port) ADDRESS="localhost:${SERVER_PORT}" TEMP_FILE=$(random tempfile) devopstest -test.v -test.run=^TestIteration7$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=$SERVER_PORT -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' -file-storage-path=$TEMP_FILE
// SERVER_PORT=$(random unused-port) ADDRESS="localhost:${SERVER_PORT}" TEMP_FILE=$(random tempfile) devopstest -test.v -test.run=^TestIteration6$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=$SERVER_PORT -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' -file-storage-path=$TEMP_FILE
//devopstest -test.v -test.run=^TestIteration8 -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port=8080 -database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' -file-storage-path=azxs123
const (
    url = "127.0.0.1:8080"
    StoreInterval time.Duration = time.Second * 200
    StoreFile = "/tmp/devops-metrics-db.json"
    Restore = false
)

type Config struct {
    Address  string 		`env:"ADDRESS"`// envDefault:"127.0.0.1:8080"`
    StoreInterval  time.Duration `env:"STORE_INTERVAL"`// envDefault:"30s"`
    StoreFile  string 		`env:"STORE_FILE"`// envDefault:"/tmp/devops-metrics-db.json"`
    Restore  bool 		`env:"RESTORE"`// envDefault:"true"`
}
var (
    M store.Repo 
    // ms = new(store.MemStorage)
    UrlReal = url
	storeInterval = StoreInterval
	storeFile = StoreFile
	restore = Restore
    cfg Config
)


// func init() {
//     onboarding()


// }
func onboarding() {
	fmt.Println("---------onboarding-------------------")
	// fmt.Println("------1---storeInterval-------------------", storeInterval)

    err := env.Parse(&cfg)
    if err != nil {
        fmt.Println(err)
    }
    // fmt.Printf("main:  %+v\n", cfg)

    urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    restorePtr := flag.Bool("r", true, "restore Metrics(Bool)")
    storeFilePtr := flag.String("f", "/tmp/devops-metrics-db.json", "file name")
    storeIntervalPtr := flag.Duration("i", 300000000000, "interval of download")

    flag.Parse()
	// fmt.Println("-----2----storeInterval-------------------", storeInterval)

    
    UrlReal = *urlRealPtr
    storeInterval = *storeIntervalPtr
    storeFile = *storeFilePtr
    restore = *restorePtr
	// fmt.Println("----3-----storeInterval-------------------", storeInterval)
    if cfg.Address != "" {
        UrlReal	= cfg.Address
    }
    if cfg.StoreInterval != 0 {
        storeInterval = cfg.StoreInterval
    }
    if cfg.StoreFile != "" {
        storeFile = cfg.StoreFile
    }
    if cfg.Restore {
        restore = cfg.Restore
    } 
	// fmt.Println("---4------storeInterval-------------------", storeInterval)

    // fmt.Println("---------onboarding m: -------------------", M)
    // fmt.Println("---------onboarding m: -------------------", ms)
	// fmt.Println("---5------cfg.StoreInterval-------------------", cfg.StoreInterval)
    fmt.Println("--------------------------const-------------------------------")
    fmt.Println("url:  ", url)
    fmt.Println("StoreInterval:  ", StoreInterval)
    fmt.Println("StoreFile:  ", StoreFile)
    fmt.Println("Restore:  ", Restore)
    fmt.Println("---------------------------flag------------------------------")
    fmt.Println("urlRealPtr:", *urlRealPtr)
    fmt.Println("restorePtr:", *restorePtr)
    fmt.Println("storeFilePtr:", *storeFilePtr)
    fmt.Println("storeIntervalPtr:", *storeIntervalPtr)
    fmt.Println("---------------------------cfg------------------------------")
    fmt.Println("cfg.Address:", cfg.Address)
    fmt.Println("cfg.Restore:", cfg.Restore)
    fmt.Println("cfg.StoreFile:", cfg.StoreFile)
    fmt.Println("cfg.StoreInterval:", cfg.StoreInterval)
    fmt.Println("------------------------real---------------------------------")
    fmt.Println("UrlReal:", UrlReal)
    fmt.Println("restore:", restore)
    fmt.Println("storeFile:", storeFile)
    fmt.Println("storeInterval:", storeInterval)
    fmt.Println("------------------------Ok---------------------------------")
}

func main() {
    onboarding()

    handler.InterfaceInit(storeInterval, storeFile, restore)
    // fmt.Println("\n --------main.func ---m: ", handler.M)
    // SaveToFile()
	// onboarding()

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

    log.Fatal(http.ListenAndServe(UrlReal, r))
    // log.Fatal(http.ListenAndServe(":8080", r))
}
