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
	// keyCrypto := key
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
			Hash: dst,
		}
	return &varMetrics
}

func NewGaugeMetrics(id string, gauge Gauge, key string) *Metrics {
	src := []byte(fmt.Sprintf("%s:gauge:%f", id, float64(gauge)))
	// keyCrypto := key
	h := hmac.New(sha256.New, []byte(key))
    h.Write(src)
    dst := h.Sum(nil)
	var varMetrics Metrics = Metrics{
			ID: id,
			MType: "gauge",
			Value: &gauge,
			Hash: dst,
		}
	return &varMetrics
}
func MtxValid(mtx *Metrics, key string) bool {
	// fmt.Println("-----------------------------------MtxValid-----start----")

	if key == "" {
	fmt.Println("-----------------------------------MtxValid-------if key == nil--")

		return true
	}
	if mtx.Hash == nil {
		fmt.Println("--------------------------------------------------------------------------------------------------mtx.Hash is ampty----")
		return false
	}
	var src []byte
	if mtx.MType == "gauge"{
	// fmt.Println("-----------------------------------MtxValid-------if mtx.MType == gauge--")

		src = []byte((fmt.Sprintf("%s:gauge:%f", mtx.ID, *mtx.Value)))
		// src = []byte((fmt.Sprintf("%s:gauge:%f", mtx.ID, float64(*mtx.Value))))
	} else if mtx.MType == "counter" {
	// fmt.Println("-----------------------------------MtxValid-------else if mtx.MType == counter--")

		src = []byte((fmt.Sprintf("%s:counter:%v", mtx.ID, *mtx.Delta)))
	} else {
		fmt.Println("-----------------------------------MtxValid-------else --false--")
		fmt.Printf("not counter not gauge: %#v:    \n", mtx)
		
		return false
	}
	
	h := hmac.New(sha256.New, []byte(key))
    h.Write(src)
    dst := h.Sum(nil)
	// fmt.Println("-----------------------------------MtxValid-----return----")
	fmt.Println("hash old: ", dst)
	fmt.Println("hash new: ", mtx.Hash)
	fmt.Println("hmac.Equal(dst, mtx.Hash): ", hmac.Equal(dst, mtx.Hash))
	// fmt.Println("-----------------------MtxValid-----finish---------------------")
	return hmac.Equal(dst, mtx.Hash)
}