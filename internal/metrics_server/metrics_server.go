package metrics_server
import(
	"fmt"
	"github.com/kokdot/go-musthave-devops/internal/repo"
	"crypto/hmac"
	"crypto/sha256"

)
type Counter = repo.Counter
type Gauge = repo.Gauge
type Metrics = repo.Metrics
var zeroG Gauge = 0
var zeroC Counter = 0

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
	src := []byte(fmt.Sprintf("%s:counter:%d", id, counter))
	if key == "" {
		panic("bad")
	}
	h := hmac.New(sha256.New, []byte(key))
    h.Write(src)
    dst := h.Sum(nil)
	var varMetrics Metrics = Metrics{
			ID: id,
			MType: "counter",
			Delta: &counter,
			Hash: fmt.Sprintf("%x", dst),
		}
	return &varMetrics
}

func NewGaugeMetrics(id string, gauge Gauge, key string) *Metrics {
	src := []byte(fmt.Sprintf("%s:gauge:%f", id, float64(gauge)))
	h := hmac.New(sha256.New, []byte(key))
    h.Write(src)
    dst := h.Sum(nil)
	var varMetrics Metrics = Metrics{
			ID: id,
			MType: "gauge",
			Value: &gauge,
			Hash: fmt.Sprintf("%x", dst),
		}
	return &varMetrics
}
func MtxValid(mtx *Metrics, key string) bool {
	if key == "" {
		return true
	}
	if mtx.Hash == "" {
		return false
	}
	var src []byte
	if mtx.MType == "gauge"{
		src = []byte((fmt.Sprintf("%s:gauge:%f", mtx.ID, *mtx.Value)))
	} else if mtx.MType == "counter" {
		src = []byte((fmt.Sprintf("%s:counter:%v", mtx.ID, *mtx.Delta)))
	} else {
		
		return false
	}
	
	h := hmac.New(sha256.New, []byte(key))
    h.Write(src)
    dst := h.Sum(nil)
	fmt.Println("hash old: ", dst)
	fmt.Println("hash new: ", mtx.Hash)
	fmt.Println("hmac.Equal(dst, mtx.Hash): ", hmac.Equal(dst, mtx.Hash))
	return (fmt.Sprintf("%x", dst) == mtx.Hash)
	// return hmac.Equal(dst, mtx.Hash)
}