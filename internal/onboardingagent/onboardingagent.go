package onboardingagent
import (
	"github.com/caarlos0/env/v6"
	"time"
	"flag"
	"strconv"
	"os"
	// "fmt"
	"github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

const (
	URL           				 	= "127.0.0.1:8080"
	PollInterval time.Duration  	= time.Second * 2
	ReportInterval time.Duration 	= time.Second * 10
	Key = ""
	Batch 							= false
	Debug = false
)

type Config struct {
    Address  string 		`env:"ADDRESS"`
    ReportInterval time.Duration	 `env:"REPORT_INTERVAL"`
    PollInterval time.Duration	 `env:"POLL_INTERVAL"`
	Key string 			`env:"KEY"`
	Batch bool 			`env:"BATCH"`
}

var( 
	pollIntervalReal = PollInterval
	reportIntervalReal = ReportInterval
	urlReal = URL
	keyReal string
	cfg Config
	batchReal = false
	logg zerolog.Logger
)
// func GetReportInterval () time.Duration {
// 	return reportIntervalReal
// }
// func GetPollInterval () time.Duration {
// 	return pollIntervalReal
// }
// func GetURL () string {
// 	return urlReal
// }
// func GetKey () string {
// 	return keyReal
// }
// func GetBatch () bool{
// 	return batchReal
// }
func GetLogg() zerolog.Logger {
	return logg
}

func OnboardingAgent() (time.Duration, time.Duration, string, string, bool, zerolog.Logger) {
    err := env.Parse(&cfg)
    if err != nil {
        logg.Print("fail to parse cfg:  ", err)
    }
	urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    reportIntervalRealPtr := flag.Duration("r", 10000000000, "interval of perort")
    pollIntervalRealPtr := flag.Duration("p", 2000000000, "interval of poll")
    keyPtr := flag.String("k", "", "secret key")
    batchPtr := flag.Bool("b", false, "batch style")
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
	reportIntervalReal = *reportIntervalRealPtr
	pollIntervalReal = *pollIntervalRealPtr
	keyReal = *keyPtr
	batchReal = *batchPtr
	// debugReal = *debug

	if cfg.Batch {
        batchReal = cfg.Batch
    }
	if cfg.Address != "" {
        urlReal	= cfg.Address
    }
	if cfg.Key != "" {
        keyReal	= cfg.Key
    }
	if cfg.ReportInterval != 0 {
        reportIntervalReal = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
        pollIntervalReal = cfg.PollInterval
	}
    logg.Print("--------------------------agent-------------------------------")
    logg.Print("--------------------------const-------------------------------")
	logg.Print("URL:     ", URL)
	logg.Print("ReportInterval:     ", ReportInterval)
	logg.Print("PollInterval:     ", PollInterval)
	logg.Print("Key:     ", Key)
	logg.Print("Batch:     ", Batch)
	logg.Print("--------------------------flag-------------------------------")
	logg.Print("urlRealPtr:     ", *urlRealPtr)
	logg.Print("reportIntervalRealPtr:     ", *reportIntervalRealPtr)
	logg.Print("pollIntervalRealPtr:     ", *pollIntervalRealPtr)
	logg.Print("keyPtr:     ", *keyPtr)
	logg.Print("batchPtr:     ", *batchPtr)
	logg.Print("debug:     ", *debug)
	logg.Print("--------------------------cfg-------------------------------")
	logg.Print("cfg.Address:     ", cfg.Address)
	logg.Print("cfg.ReportInterval:     ", cfg.ReportInterval)
	logg.Print("cfg.PollInterval:     ", cfg.PollInterval)
	logg.Print("cfg.Key:     ", cfg.Key)
	logg.Print("cfg.Batch:     ", cfg.Batch)
	logg.Print("--------------------------real-------------------------------")
	logg.Print("URLReal:     ", urlReal)
	logg.Print("ReportIntervalReal:     ", reportIntervalReal)
	logg.Print("PollIntervalReal:     ", pollIntervalReal)
	logg.Print("KeyReal:     ", keyReal)
	logg.Print("BatchReal:     ", batchReal)
	logg.Print("--------------------------Ok-------------------------------")
	return pollIntervalReal, reportIntervalReal, urlReal, keyReal, batchReal, log.Logger
}
