package onboarding_server

import (
	"fmt"
    "time"
	"github.com/kokdot/go-musthave-devops/internal/store"
	"flag"
    "github.com/caarlos0/env/v6"
)
const (
    URL = "127.0.0.1:8080"
    StoreInterval time.Duration = time.Second * 200
    StoreFile = "/tmp/devops-metrics-db.json"
    Key = ""
    Restore = false
)

type Config struct {
    Address  string 		`env:"ADDRESS"`// envDefault:"127.0.0.1:8080"`
    StoreInterval  time.Duration `env:"STORE_INTERVAL"`// envDefault:"30s"`
    StoreFile  string 		`env:"STORE_FILE"`// envDefault:"/tmp/devops-metrics-db.json"`
    Restore  bool 		`env:"RESTORE"`// envDefault:"true"`
    Key string 			`env:"KEY"`
}
var (
    M store.Repo 
    // ms = new(store.MemStorage)
    URLReal = URL
	StoreIntervalReal = StoreInterval
	StoreFileReal = StoreFile
	RestoreReal = Restore
    KeyReal = Key
    cfg Config
)
func OnboardingServer() {
	fmt.Println("---------onboarding-------------------")
    err := env.Parse(&cfg)
    if err != nil {
        fmt.Println(err)
    }

    URLRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    restorePtr := flag.Bool("r", true, "restore Metrics(Bool)")
    storeFilePtr := flag.String("f", "/tmp/devops-metrics-db.json", "file name")
    storeIntervalPtr := flag.Duration("i", 300000000000, "interval of download")
    keyPtr := flag.String("k", "", "secret key")

    flag.Parse()
    URLReal = *URLRealPtr
    StoreIntervalReal = *storeIntervalPtr
    StoreFileReal = *storeFilePtr
    RestoreReal = *restorePtr
    KeyReal = *keyPtr

    if cfg.Address != "" {
        URLReal	= cfg.Address
    }
    if cfg.StoreInterval != 0 {
        StoreIntervalReal = cfg.StoreInterval
    }
    if cfg.StoreFile != "" {
        StoreFileReal = cfg.StoreFile
    }
    if cfg.Restore {
        RestoreReal = cfg.Restore
    } 
    if cfg.Key != "" {
        KeyReal	= cfg.Key
    }
    fmt.Println("--------------------------const-------------server------------------")
    fmt.Println("URL:  ", URL)
    fmt.Println("StoreInterval:  ", StoreInterval)
    fmt.Println("StoreFile:  ", StoreFile)
    fmt.Println("Restore:  ", Restore)
    fmt.Println("Key:  ", Key)
    fmt.Println("---------------------------flag------------------------------")
    fmt.Println("URLRealPtr:", *URLRealPtr)
    fmt.Println("restorePtr:", *restorePtr)
    fmt.Println("storeFilePtr:", *storeFilePtr)
    fmt.Println("storeIntervalPtr:", *storeIntervalPtr)
    fmt.Println("keyPtr:", *keyPtr)
    fmt.Println("---------------------------cfg------------------------------")
    fmt.Println("cfg.Address:", cfg.Address)
    fmt.Println("cfg.Restore:", cfg.Restore)
    fmt.Println("cfg.StoreFile:", cfg.StoreFile)
    fmt.Println("cfg.StoreInterval:", cfg.StoreInterval)
    fmt.Println("cfg.Key:", cfg.Key)
    fmt.Println("------------------------real---------------------------------")
    fmt.Println("URLReal:", URLReal)
    fmt.Println("RestoreReal:", RestoreReal)
    fmt.Println("StoreFileReal:", StoreFileReal)
    fmt.Println("StoreIntervalReal:", StoreIntervalReal)
    fmt.Println("KeyReal:", KeyReal)
    fmt.Println("------------------------Ok---------------------------------")
}