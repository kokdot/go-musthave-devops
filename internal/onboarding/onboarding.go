package onboarding
import (
	"github.com/caarlos0/env/v6"
	"time"
	"flag"
	// "github.com/kokdot/go-musthave-devops/internal/def"
	"fmt"
)

const (
	Url           				 	= "127.0.0.1:8080"
	PollInterval time.Duration  	= time.Second * 2
	ReportInterval time.Duration 	= time.Second * 10
)

type Config struct {
    Address  string 		`env:"ADDRESS"`
    ReportInterval time.Duration	 `env:"REPORT_INTERVAL"`
    PollInterval time.Duration	 `env:"POLL_INTERVAL"`
}

var( 
	PollIntervalReal = PollInterval
	ReportIntervalReal = ReportInterval
	UrlReal = Url
	cfg Config
)

func Onboarding() {
    err := env.Parse(&cfg)
    if err != nil {
        fmt.Println("fail to parse cfg:  ", err)
    }
	urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    reportIntervalRealPtr := flag.Duration("r", 10000000000, "interval of perort")
    pollIntervalRealPtr := flag.Duration("p", 2000000000, "interval of poll")
    flag.Parse()
	UrlReal = *urlRealPtr
	ReportIntervalReal = *reportIntervalRealPtr
	PollIntervalReal = *pollIntervalRealPtr

	if cfg.Address != "" {
        UrlReal	= cfg.Address
    }
	if cfg.ReportInterval != 0 {
        ReportIntervalReal = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
        PollIntervalReal = cfg.PollInterval
	}
	// UrlReal = "http://" + UrlReal
    fmt.Println("--------------------------const-------------------------------")
	fmt.Println("Url:     ", Url)
	fmt.Println("ReportIntervalReal:     ", ReportInterval)
	fmt.Println("pollIntervalReal:     ", PollInterval)
	fmt.Println("--------------------------flag-------------------------------")
	fmt.Println("urlRealPtr:     ", urlRealPtr)
	fmt.Println("reportIntervalRealPtr:     ", reportIntervalRealPtr)
	fmt.Println("pollIntervalRealPtr:     ", pollIntervalRealPtr)
	fmt.Println("--------------------------cfg-------------------------------")
	fmt.Println("cfg.Address:     ", cfg.Address)
	fmt.Println("cfg.ReportInterval:     ", cfg.ReportInterval)
	fmt.Println("cfg.PollInterval:     ", cfg.PollInterval)
	fmt.Println("--------------------------real-------------------------------")
	fmt.Println("urlReal:     ", UrlReal)
	fmt.Println("reportIntervalReal:     ", ReportIntervalReal)
	fmt.Println("pollIntervalReal:     ", PollIntervalReal)
	fmt.Println("--------------------------Ok-------------------------------")
}