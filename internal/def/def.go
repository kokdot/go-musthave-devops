package def
import(
	"time"
	"github.com/rs/zerolog"
)

type Gauge float64
type Counter int64
type GaugeMap map[string] Gauge

type Conf struct {
	PollInterval time.Duration
	ReportInterval time.Duration
	URL string
	Key string
	Batch bool
	Logg zerolog.Logger
}