package metricsserver
import(
	"fmt"
	"github.com/kokdot/go-musthave-devops/internal/repo"
	"crypto/hmac"
	"github.com/rs/zerolog"
	"crypto/sha256"

)
type Counter = repo.Counter
type Gauge = repo.Gauge
type Metrics = repo.Metrics
var zeroG Gauge = 0
var zeroC Counter = 0
var logg zerolog.Logger

func GetLogg(loggReal zerolog.Logger)  {
	logg = loggReal
}
func NewMetrics(id string, mType string) Metrics {
	if mType == "gauge" {
		return Metrics{
		ID: id,
		MType: "gauge",
		Value: &zeroG,
		}
	} else {
		return Metrics{
			ID: id,
			MType: "counter",
			Delta: &zeroC,
		}
	}
}
func NewCounterMetrics(id string, counter Counter, key string) *Metrics {
	var varMetrics = Metrics{
			ID: id,
			MType: "counter",
 			Delta: &counter,
	}
	if key == "" {
		return &varMetrics
	} 
	src := []byte(fmt.Sprintf("%s:counter:%d", id, counter))
	h := hmac.New(sha256.New, []byte(key))
    h.Write(src)
    dst := h.Sum(nil)
	varMetrics.Hash = fmt.Sprintf("%x", dst)
	return &varMetrics
}

func NewGaugeMetrics(id string, gauge Gauge, key string) *Metrics {
	var varMetrics = Metrics{
			ID: id,
			MType: "gauge",
			Value: &gauge,
		}
	if key == "" {
		return &varMetrics
	} 
	src := []byte(fmt.Sprintf("%s:gauge:%f", id, float64(gauge)))
	h := hmac.New(sha256.New, []byte(key))
	h.Write(src)
	dst := h.Sum(nil)
	varMetrics.Hash = fmt.Sprintf("%x", dst)
	return &varMetrics
}
func MtxValid(mtx *Metrics, key string) bool {
	if key == "" {
		return true
	}
	if mtx.Hash == "" {
		return false
	}
	// var src []byte
	// if mtx.MType == "gauge"{
	// 	src = []byte((fmt.Sprintf("%s:gauge:%f", mtx.ID, *mtx.Value)))
	// } else if mtx.MType == "counter" {
	// 	src = []byte((fmt.Sprintf("%s:counter:%v", mtx.ID, *mtx.Delta)))
	// } else {
		
	// 	return false
	// }
	logg.Print("-------------------------------MtxValid---------------------------start-----------------key not nil--------")
	logg.Print("mtx.ID =  ", mtx.ID)
	logg.Print("mtx.MType =  ", mtx.MType)
	logg.Print("mtx.MType =  ", mtx.MType)
	logg.Print("key =  ", key)
	if mtx.Value != nil {
		logg.Print("mtx.Value =  ", *mtx.Value)
	}
	// h := hmac.New(sha256.New, []byte(key))
    // h.Write(src)
    // dst := h.Sum(nil)
	hash := Hash(mtx, key)
	logg.Print("hash: ", hash)
	// logg.Print("hash old: ", fmt.Sprintf("%x", dst))
	logg.Print("hash is come: ", mtx.Hash)
	logg.Print("hmac.Equal(dst, mtx.Hash): ", (hash == mtx.Hash))
	// logg.Print("hmac.Equal(dst, mtx.Hash): ", (fmt.Sprintf("%x", dst) == mtx.Hash))
	// logg.Print("hmac.Equal(dst, mtx.Hash): ", hmac.Equal(dst, mtx.Hash))
	return (hash == mtx.Hash)
	// return (fmt.Sprintf("%x", dst) == mtx.Hash)
	// return hmac.Equal(dst, mtx.Hash)
}
func Hash(m *Metrics, key string) string {
	var data string
	switch m.MType {
	case "counter":
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case "gauge":
		data = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}
