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
)

type Config struct {
    Address  string 		`env:"ADDRESS"`
    ReportInterval time.Duration	 `env:"REPORT_INTERVAL"`
    PollInterval time.Duration	 `env:"POLL_INTERVAL"`
	Key string 			`env:"KEY"`
}

var( 
	PollIntervalReal = PollInterval
	ReportIntervalReal = ReportInterval
	URLReal = URL
	KeyReal string
	cfg Config
)

func OnboardingAgent() {
    err := env.Parse(&cfg)
    if err != nil {
        fmt.Println("fail to parse cfg:  ", err)
    }
	urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    reportIntervalRealPtr := flag.Duration("r", 10000000000, "interval of perort")
    pollIntervalRealPtr := flag.Duration("p", 2000000000, "interval of poll")
    keyPtr := flag.String("k", "", "secret key")
	
	flag.Parse()
	URLReal = *urlRealPtr
	ReportIntervalReal = *reportIntervalRealPtr
	PollIntervalReal = *pollIntervalRealPtr
	KeyReal = *keyPtr

	if cfg.Address != "" {
        URLReal	= cfg.Address
    }
	if cfg.Key != "" {
        KeyReal	= cfg.Key
    }
	if cfg.ReportInterval != 0 {
        ReportIntervalReal = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
        PollIntervalReal = cfg.PollInterval
	}
    fmt.Println("--------------------------2023-------------------------------")
    fmt.Println("--------------------------const-------------------------------")
	fmt.Println("URL:     ", URL)
	fmt.Println("ReportInterval:     ", ReportInterval)
	fmt.Println("PollInterval:     ", PollInterval)
	fmt.Println("Key:     ", Key)
	fmt.Println("--------------------------flag-------------------------------")
	fmt.Println("urlRealPtr:     ", urlRealPtr)
	fmt.Println("reportIntervalRealPtr:     ", reportIntervalRealPtr)
	fmt.Println("pollIntervalRealPtr:     ", pollIntervalRealPtr)
	fmt.Println("keyPtr:     ", keyPtr)
	fmt.Println("--------------------------cfg-------------------------------")
	fmt.Println("cfg.Address:     ", cfg.Address)
	fmt.Println("cfg.ReportInterval:     ", cfg.ReportInterval)
	fmt.Println("cfg.PollInterval:     ", cfg.PollInterval)
	fmt.Println("cfg.Key:     ", cfg.Key)
	fmt.Println("--------------------------real-------------------------------")
	fmt.Println("URLReal:     ", URLReal)
	fmt.Println("ReportIntervalReal:     ", ReportIntervalReal)
	fmt.Println("PollIntervalReal:     ", PollIntervalReal)
	fmt.Println("KeyReal:     ", KeyReal)
	fmt.Println("--------------------------Ok-------------------------------")
}