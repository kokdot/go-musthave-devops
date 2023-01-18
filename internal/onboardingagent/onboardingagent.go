package onboardingagent
import (
	"github.com/caarlos0/env/v6"
	"time"
	"flag"
	"fmt"
)

const (
	URL           				 	= "127.0.0.1:8080"
	PollInterval time.Duration  	= time.Second * 2
	ReportInterval time.Duration 	= time.Second * 10
	Key = ""
	Batch 							= false
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

func OnboardingAgent() (time.Duration, time.Duration, string, string, bool) {
    err := env.Parse(&cfg)
    if err != nil {
        fmt.Println("fail to parse cfg:  ", err)
    }
	urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    reportIntervalRealPtr := flag.Duration("r", 10000000000, "interval of perort")
    pollIntervalRealPtr := flag.Duration("p", 2000000000, "interval of poll")
    keyPtr := flag.String("k", "", "secret key")
    batchPtr := flag.Bool("b", false, "batch style")
	
	flag.Parse()
	urlReal = *urlRealPtr
	reportIntervalReal = *reportIntervalRealPtr
	pollIntervalReal = *pollIntervalRealPtr
	keyReal = *keyPtr
	batchReal = *batchPtr

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
    fmt.Println("--------------------------agent-------------------------------")
    fmt.Println("--------------------------const-------------------------------")
	fmt.Println("URL:     ", URL)
	fmt.Println("ReportInterval:     ", ReportInterval)
	fmt.Println("PollInterval:     ", PollInterval)
	fmt.Println("Key:     ", Key)
	fmt.Println("Batch:     ", Batch)
	fmt.Println("--------------------------flag-------------------------------")
	fmt.Println("urlRealPtr:     ", *urlRealPtr)
	fmt.Println("reportIntervalRealPtr:     ", *reportIntervalRealPtr)
	fmt.Println("pollIntervalRealPtr:     ", *pollIntervalRealPtr)
	fmt.Println("keyPtr:     ", *keyPtr)
	fmt.Println("batchPtr:     ", *batchPtr)
	fmt.Println("--------------------------cfg-------------------------------")
	fmt.Println("cfg.Address:     ", cfg.Address)
	fmt.Println("cfg.ReportInterval:     ", cfg.ReportInterval)
	fmt.Println("cfg.PollInterval:     ", cfg.PollInterval)
	fmt.Println("cfg.Key:     ", cfg.Key)
	fmt.Println("cfg.Batch:     ", cfg.Batch)
	fmt.Println("--------------------------real-------------------------------")
	fmt.Println("URLReal:     ", urlReal)
	fmt.Println("ReportIntervalReal:     ", reportIntervalReal)
	fmt.Println("PollIntervalReal:     ", pollIntervalReal)
	fmt.Println("KeyReal:     ", keyReal)
	fmt.Println("BatchReal:     ", batchReal)
	fmt.Println("--------------------------Ok-------------------------------")
	return pollIntervalReal, reportIntervalReal, urlReal, keyReal, batchReal
}
