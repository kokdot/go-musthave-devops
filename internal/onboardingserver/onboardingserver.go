package onboardingserver

import (
	// "fmt"
    "time"
	// "github.com/kokdot/go-musthave-devops/internal/store"
	"strconv"
	"os"
	"flag"
    "github.com/caarlos0/env/v6"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)
const (
    URL = "127.0.0.1:8080"
    StoreInterval time.Duration = time.Second * 200
    StoreFile = "/tmp/devops-metrics-db.json"
    Key = ""
    Restore = false
    DataBaseDSN = ""
    Debug = false
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
    logg zerolog.Logger
)

func GetLogg() zerolog.Logger {
	return logg
}

func OnboardingServer() (string, string, string, bool, time.Duration, string, zerolog.Logger) {
	logg.Print("---------onboarding-------------------")
    err := env.Parse(&cfg)
    if err != nil {
        logg.Print(err)
    }

    urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    restorePtr := flag.Bool("r", false, "restore Metrics(Bool)")
    storeFilePtr := flag.String("f", "/tmp/devops-metrics-db.json", "file name")
    storeIntervalPtr := flag.Duration("i", 300000000000, "interval of download")
    keyPtr := flag.String("k", "", "secret key")
    DataBaseDSNPtr := flag.String("d", "", "Data Base DSN")
    debug := flag.Bool("debug", false, "sets log level to debug")

    flag.Parse()
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    if *debug {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// log.Logger = log.With().Caller().Logger()
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
	logg = log.Logger

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
    logg.Print("--------------------------const-------------server------------------")
    logg.Print("URL:  ", URL)
    logg.Print("StoreInterval:  ", StoreInterval)
    logg.Print("StoreFile:  ", StoreFile)
    logg.Print("Restore:  ", Restore)
    logg.Print("Key:  ", Key)
    logg.Print("DataBaseDSN:  ", DataBaseDSN)
    logg.Print("---------------------------flag------------------------------")
    logg.Print("URLRealPtr:", *urlRealPtr)
    logg.Print("restorePtr:", *restorePtr)
    logg.Print("storeFilePtr:", *storeFilePtr)
    logg.Print("storeIntervalPtr:", *storeIntervalPtr)
    logg.Print("keyPtr:", *keyPtr)
    logg.Print("DataBaseDSNPtr:", *DataBaseDSNPtr)
    logg.Print("debug:     ", *debug)
    logg.Print("---------------------------cfg------------------------------")
    logg.Print("cfg.Address:", cfg.Address)
    logg.Print("cfg.Restore:", cfg.Restore)
    logg.Print("cfg.StoreFile:", cfg.StoreFile)
    logg.Print("cfg.StoreInterval:", cfg.StoreInterval)
    logg.Print("cfg.Key:", cfg.Key)
    logg.Print("cfg.DataBaseDSN:", cfg.DataBaseDSN)
    logg.Print("------------------------real---------------------------------")
    logg.Print("urlReal:", urlReal)
    logg.Print("restoreReal:", restoreReal)
    logg.Print("storeFileReal:", storeFileReal)
    logg.Print("storeIntervalReal:", storeIntervalReal)
    logg.Print("keyReal:", keyReal)
    logg.Print("dataBaseDSNReal:", dataBaseDSNReal)
    logg.Print("------------------------Ok---------------------------------")
    return urlReal, storeFileReal, keyReal, restoreReal, storeIntervalReal, dataBaseDSNReal, log.Logger
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