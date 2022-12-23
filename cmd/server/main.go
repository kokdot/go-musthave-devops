package main

import (
	"log"
	"net/http"
	// "strconv"
    "time"
    "flag"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"fmt"

	"github.com/kokdot/go-musthave-devops/internal/handler"
)
const (
    url = "127.0.0.1:8080"
    StoreInterval = 300
    StoreFile = "/tmp/devops-metrics-db.json"
    Restore = true
)
type Config struct {
    Address  string 		`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
    StoreInterval  int 		`env:"STORE_INTERVAL" envDefault:"300"`
    StoreFile  string 		`env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
    Restore  bool 		`env:"RESTORE" envDefault:"true"`
}
var (
    urlReal = url
	storeInterval = StoreInterval
	storeFile = StoreFile
	restore = Restore
    // syncDownload = false
    cfg Config
)

func onboarding() {
    err := env.Parse(&cfg)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("main:  %+v\n", cfg)
	urlReal	= cfg.Address
    storeInterval = cfg.StoreInterval
    storeFile = cfg.StoreFile
    restore = cfg.Restore

    urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    restorePtr := flag.Bool("r", true, "restore Metrics(Bool)")
    storeFilePtr := flag.String("f", "/tmp/devops-metrics-db.json", "file name")
    storeIntervalPtr := flag.Int("i", 300, "interval of download")

    flag.Parse()
    if urlReal == url {
        urlReal = *urlRealPtr
    }
    if storeInterval == StoreInterval {
        storeInterval = *storeIntervalPtr
    }
    if storeFile == StoreFile {
        storeFile = *storeFilePtr
    }
    if !Restore {
        restore = *restorePtr
    }
    
    if storeFile != "" {
        if storeInterval > 0 {
            DownloadToFile(storeFile)
        } else {
            // syncDownload = true
            handler.CheckSyncDownload(storeFile)
        }
    }
    if restore{
        handler.UpdateMemStorageFromFile(storeFile)
    }
   
    fmt.Println("urlRealPrt:", *urlRealPtr)
    fmt.Println("restorePtr:", *restorePtr)
    fmt.Println("storeFilePtr:", *storeFilePtr)
    fmt.Println("storeIntervalPtr:", *storeIntervalPtr)
}
func DownloadToFile(file string) {
    go func() {
        var interval = time.Duration(storeInterval) * time.Second
        for {
            <-time.After(interval) 
            fmt.Println("main; line: 67; DownloadToFile", ";  file:  ", file)
            handler.DownloadMemStorageToFile(file)
        }
    }()
}

func main() {
	onboarding()

    // определяем роутер chi
    r := chi.NewRouter()
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
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

    log.Fatal(http.ListenAndServe(urlReal, r))
    // log.Fatal(http.ListenAndServe(":8080", r))
}
