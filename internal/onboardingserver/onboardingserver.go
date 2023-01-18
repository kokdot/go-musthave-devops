package onboardingserver

import (
	"fmt"
    "time"
	// "github.com/kokdot/go-musthave-devops/internal/store"
	"flag"
    "github.com/caarlos0/env/v6"
)
const (
    URL = "127.0.0.1:8080"
    StoreInterval time.Duration = time.Second * 200
    StoreFile = "/tmp/devops-metrics-db.json"
    Key = ""
    Restore = false
    DataBaseDSN = ""
    // DataBaseDSN = "postgres://postgres:postgres@localhost:5432/postgres"
)

type Config struct {
    Address  string 		`env:"ADDRESS"`// envDefault:"127.0.0.1:8080"`
    StoreInterval  time.Duration `env:"STORE_INTERVAL"`// envDefault:"30s"`
    StoreFile  string 		`env:"STORE_FILE"`// envDefault:"/tmp/devops-metrics-db.json"`
    Restore  bool 		`env:"RESTORE" envDefault:"false"`
    Key string 			`env:"KEY"`
    DataBaseDSN string    `env:"DATABASE_DSN"`
}
var (
    urlReal = URL
	storeIntervalReal = StoreInterval
	storeFileReal = StoreFile
	restoreReal = Restore
    keyReal = Key
    cfg Config
    dataBaseDSNReal = DataBaseDSN
)
func OnboardingServer() (string, string, string, bool, time.Duration, string) {
	fmt.Println("---------onboarding-------------------")
    err := env.Parse(&cfg)
    if err != nil {
        fmt.Println(err)
    }

    urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    restorePtr := flag.Bool("r", false, "restore Metrics(Bool)")
    storeFilePtr := flag.String("f", "/tmp/devops-metrics-db.json", "file name")
    storeIntervalPtr := flag.Duration("i", 300000000000, "interval of download")
    keyPtr := flag.String("k", "", "secret key")
    DataBaseDSNPtr := flag.String("d", "", "Data Base DSN")

    flag.Parse()
    urlReal = *urlRealPtr
    storeIntervalReal = *storeIntervalPtr
    storeFileReal = *storeFilePtr
    restoreReal = *restorePtr
    keyReal = *keyPtr
    if *DataBaseDSNPtr != "" {
        dataBaseDSNReal = *DataBaseDSNPtr
    }

    if cfg.Address != "" {
        urlReal	= cfg.Address
    }
    if cfg.StoreInterval != 0 {
        storeIntervalReal = cfg.StoreInterval
    }
    if cfg.StoreFile != "" {
        storeFileReal = cfg.StoreFile
    }
    if cfg.Restore {
        restoreReal = cfg.Restore
    } 
    if cfg.Key != "" {
        keyReal	= cfg.Key
    }
    if cfg.DataBaseDSN != "" {
        dataBaseDSNReal	= cfg.DataBaseDSN
    }
    fmt.Println("--------------------------const-------------server------------------")
    fmt.Println("URL:  ", URL)
    fmt.Println("StoreInterval:  ", StoreInterval)
    fmt.Println("StoreFile:  ", StoreFile)
    fmt.Println("Restore:  ", Restore)
    fmt.Println("Key:  ", Key)
    fmt.Println("DataBaseDSN:  ", DataBaseDSN)
    fmt.Println("---------------------------flag------------------------------")
    fmt.Println("URLRealPtr:", *urlRealPtr)
    fmt.Println("restorePtr:", *restorePtr)
    fmt.Println("storeFilePtr:", *storeFilePtr)
    fmt.Println("storeIntervalPtr:", *storeIntervalPtr)
    fmt.Println("keyPtr:", *keyPtr)
    fmt.Println("DataBaseDSNPtr:", *DataBaseDSNPtr)
    fmt.Println("---------------------------cfg------------------------------")
    fmt.Println("cfg.Address:", cfg.Address)
    fmt.Println("cfg.Restore:", cfg.Restore)
    fmt.Println("cfg.StoreFile:", cfg.StoreFile)
    fmt.Println("cfg.StoreInterval:", cfg.StoreInterval)
    fmt.Println("cfg.Key:", cfg.Key)
    fmt.Println("cfg.DataBaseDSN:", cfg.DataBaseDSN)
    fmt.Println("------------------------real---------------------------------")
    fmt.Println("urlReal:", urlReal)
    fmt.Println("restoreReal:", restoreReal)
    fmt.Println("storeFileReal:", storeFileReal)
    fmt.Println("storeIntervalReal:", storeIntervalReal)
    fmt.Println("keyReal:", keyReal)
    fmt.Println("dataBaseDSNReal:", dataBaseDSNReal)
    fmt.Println("------------------------Ok---------------------------------")
    return urlReal, storeFileReal, keyReal, restoreReal, storeIntervalReal, dataBaseDSNReal
}


func GetStoreInterval() time.Duration {
    return storeIntervalReal
}

func GetStoreFile() string {
    return storeFileReal
}


func GetRestore() bool {
    return restoreReal
}

func GetURL() string {
    return urlReal
}

func GetKey() string {
    return keyReal
}