package handler

import (
	"encoding/json"
	"fmt"
	"io"

	// "github.com/caarlos0/env/v6"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kokdot/go-musthave-devops/internal/store"
	// "flag"
	// "strconv"
)

type key int

const (
	nameDataKey key = iota
	valueDataKey
)

// const (
//     url = "127.0.0.1:8080"
//     StoreInterval time.Duration = 200
//     StoreFile = "/tmp/devops-metrics-db.json"
//     Restore = false
// )
// type Config struct {
//     Address  string 		`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
//     StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"30s"`
//     StoreFile  string 		`env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
//     Restore  bool 		`env:"RESTORE" envDefault:"true"`
// }
var (
	// ms = new(store.MemStorage)
	M  store.Repo
	//     UrlReal = url
	StoreInterval time.Duration// = StoreInterval

// 	storeFile = StoreFile
// 	restore = Restore
//     cfg Config
)

func init() {
	// fmt.Println("---------init-------------------")/>;./
	// onboarding()
}
func InterfaceInit(storeInterval time.Duration, storeFile string, restore bool) {
    StoreInterval = storeInterval
	if storeInterval > 0 {
		//---------save to MemStorage-------------------
		// fmt.Println("---------save to MemStorage-------------------")
		if storeFile == "" {
			// fmt.Println("---------storeFile == \"\"-------------------")

			ms, err := store.NewMemStorage()
            if err != nil {
				log.Fatalf("failed to create MemStorage, err: %s", err)
			}
			M = ms
		} else {
			// fmt.Println("---------storeFile == \"\" else------------------")

			ms, err := store.NewMemStorageWithFile(storeFile)
			// fmt.Println("---------ms-------------------", ms)

			if err != nil {
				log.Fatalf("failed to create MemStorage, err: %s", err)
			}
			M = ms
			// fmt.Println("---------m-------------------", M, "restore:  :", restore)

			if restore {
				// fmt.Println("---------restore-------------------", M)

				_, err := M.ReadStorage()
				// fmt.Println("---------M.ReadStorage()-------------------", M)

				if err != nil {
					log.Printf("Can't to read data froM file, err: %s", err)
				}
			}
			DownloadingToFile()
			M = ms
		}
	} else {
		//---------save to FileStorage------------------
		// fmt.Println("---------save to FileStorage-------------------")

		if storeFile == "" {
			ms, err := store.NewFileStorage()
			if err != nil {
				log.Fatalf("failed to create FileStorage, err: %s", err)
			}
			M = ms
		} else {
			ms, err := store.NewFileStorageWithFile(storeFile)
			if err != nil {
				log.Fatalf("failed to create FileStorage, err: %s", err)
			}
			M = ms
			if restore {
				_, err := M.ReadStorage()
				if err != nil {
					log.Printf("Can't to read data from file, err: %s", err)
				}
			}
			// DownloadingToFile()
		}
	}
	// fmt.Printf("---------init exit %+v: -------------------", M)
}

// func onboarding() {
// 	// fmt.Println("---------onboarding-------------------")
// 	// fmt.Println("------1---storeInterval-------------------", storeInterval)
// 
//     err := env.Parse(&cfg)
//     if err != nil {
//         fmt.Println(err)
//     }
//     // fmt.Printf("main:  %+v\n", cfg)
// 
//     urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
//     restorePtr := flag.Bool("r", true, "restore Metrics(Bool)")
//     storeFilePtr := flag.String("f", "/tmp/devops-metrics-db.json", "file name")
//     storeIntervalPtr := flag.Duration("i", 300, "interval of download")
//
//     flag.Parse()
// 	// fmt.Println("-----2----storeInterval-------------------", storeInterval)
//
//     UrlReal = *urlRealPtr
//     storeInterval = *storeIntervalPtr
//     storeFile = *storeFilePtr
//     restore = *restorePtr
// 	// fmt.Println("----3-----storeInterval-------------------", storeInterval)
//
//     UrlReal	= cfg.Address
//     storeInterval = cfg.StoreInterval
//     storeFile = cfg.StoreFile
//     restore = cfg.Restore
// 	// fmt.Println("---4------storeInterval-------------------", storeInterval)
//
//     // fmt.Println("---------onboarding m: -------------------", M)
//     // fmt.Println("---------onboarding m: -------------------", ms)
// 	// fmt.Println("---5------cfg.StoreInterval-------------------", cfg.StoreInterval)
//     fmt.Println("UrlReal:", UrlReal)
//     fmt.Println("restore:", restore)
//     fmt.Println("storeFile:", storeFile)
//     fmt.Println("storeInterval:", storeInterval)
//
// }

func DownloadingToFile() {
	// fmt.Println("---------DownloadingToFile m: -------------------", M)
	// fmt.Println("---------DownloadingToFile-------------------", storeInterval)

	go func() {
		// var interval = StoreInterval
		// var interval = time.Duration(storeInterval) * time.Second
		for {
			<-time.After(StoreInterval)
			fmt.Println("main; line: 67; DownloadToFile", ";  file:  ")
			err := M.WriteStorage()
			if err != nil {
				log.Printf("StoreMap did not been saved in file, err: %s", err)
			}
		}
	}()
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("------PostUpdate---m-------------------", M)
    fmt.Println("r.Header: ", r.Header)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "http.StatusBadRequest")
		return
	}
	var mtxNew store.Metrics
	err = json.Unmarshal(bodyBytes, &mtxNew)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "http.StatusBadRequest")
		return
	}
        // fmt.Println("-197-handler--PostUpdate-----m/txNew-------------------", mtxNew)
    // if mtxNew.Delta != nil {
        // fmt.Println("--199--handler---PostUpdate-----mtxNew-------------------", mtxNew, "delta:  ", *mtxNew.Delta)
    // }
    // if mtxNew.Value != nil {
    //     fmt.Println("--202---handler----PostUpdate-----mtxNew-------------------", mtxNew, "Value:  ", *mtxNew.Value)
    // }
    // sm, _ := M.GetAllValues()
    // fmt.Println("----205---handler-------------/-------------------------------------------")
	// fmt.Println("--206--handler-----PostUpdate----sm-!!!!!!!!-----------------\n:  ", sm)
    // fmt.Println("---207---handler---------------------------------------------------------")
	if mtxNew.ID == "Alloc"{
		fmt.Println("-------Post-----------Alloc start----------------------")
		 sm := ""
		if M != nil {
		fmt.Println("-----PostUpdate----M-------------------", M)

			sm, _ = M.GetAllValues()

		}
		fmt.Println("-----PostUpdate----sm-------------------", sm)
		fmt.Println("----------------------------------------")
	}
	mtxOld, err := M.Save(&mtxNew)
	if mtxNew.ID == "Alloc"{
		 sm := ""
		if M != nil {
		fmt.Println("-----PostUpdate----M-------------------", M)

			sm, _ = M.GetAllValues()

		}
		fmt.Println("-----PostUpdate----sm-------------------", sm)
		fmt.Println("-----Post--Alloc-----finish----------------------------")
	}
    // fmt.Println("---210----handler--------M.Save----------------")
        // fmt.Println("--211---handler-----PostUpdate-----mtxOld-------------------", mtxNew)
    // if mtxNew.Delta != nil {
        // fmt.Println("--213----handler------PostUpdate-----mtxOld-------------------", mtxNew, "delta:  ", *mtxNew.Delta)
    // }
    // if mtxNew.Value != nil {
        // fmt.Println("--216----handler------PostUpdate-----mtxOld-------------------", mtxNew, "Value:  ", *mtxNew.Value)
    // }
    // sm, _ = M.GetAllValues()
	// fmt.Println("--219-------handler----PostUpdate----sm-------------------", sm)
    // fmt.Println("---220-------handler---------------------------------------------------------------------------")
    // fmt.Println("----221--------handler-------------------------------------------------------------------------")

	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprint(w, "http.StatusBadRequest")
		return
	}
	bodyBytes, err = json.Marshal(mtxOld)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "http.StatusBadRequest")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "%v", bodyBytes)
	w.Write(bodyBytes)
	// if syncDownload {
	//     DownloadMemStorageToFile(StoreFile)
	// }
}

func GetValue(w http.ResponseWriter, r *http.Request) {
    // sm := ""
    // if M != nil {
	// fmt.Println("-----PostUpdate----M-------------------", M)

        // sm, _ = M.GetAllValues()

    // }
	// fmt.Println("-----PostUpdate----sm-------------------", sm)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
        // fmt.Println("-----------------------------------err line 247")

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprint(w, "http.StatusBadRequest")
		return
	}
	var mtxNew store.Metrics
	err = json.Unmarshal(bodyBytes, &mtxNew)
	if err != nil {
        // fmt.Println("-----------------------------------err line 254")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// fmt.Fprint(w, "http.StatusBadRequest")
		return
	}
	if mtxNew.ID == "Alloc"{
		fmt.Println("-------Get-----------Alloc start----------------------")
		 sm := ""
		if M != nil {
		fmt.Println("-----GetValue----M-------------------", M)

			sm, _ = M.GetAllValues()

		}
		fmt.Println("-----GetValue----sm-------------------", sm)
		fmt.Println("----------------------------------------")
	}
	mtxOLd, err := M.Get(mtxNew.ID) 
	if mtxNew.ID == "Alloc"{
		 sm := ""
		if M != nil {
		fmt.Println("-----GetValue----M-------------------", M)

			sm, _ = M.GetAllValues()

		}
		fmt.Println("-----GetValue----sm-------------------", sm)
		fmt.Println("----Get---------Alloc finish---------------------------")
	}

	fmt.Println("----------GetValue------mtxNew.----:   ", mtxNew, "--id---:   ", mtxNew.ID)
	// // fmt.Println("----------GetValue------mtxOld:----:  ", mtxOLd)
    // if mtxNew.Delta != nil {
    //     fmt.Println("----GetValue-----mtxOld-------------------", mtxOLd, "delta:  ", *mtxOLd.Delta)
    // }
    // if mtxNew.Value != nil {
    //     fmt.Println("----GetValue-----mtxOld-------------------", mtxOLd, "delta:  ", *mtxOLd.Value)
    // }
    // fmt.Println("-------------------------------------------------------------------------------------")
    // fmt.Println("-------------------------------------------------------------------------------------")

	if err != nil {
        fmt.Println("-----------------------------------err line 274, err:  ", err)

		w.Header().Set("content-type", "application/json")
		// w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "http.StatusBadRequest")
		return
	}
	bodyBytes, err = json.Marshal(mtxOLd)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		// w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "http.StatusBadRequest")
		return
	}
	// w.Header().Set("content-type", "text/html")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "%v", bodyBytes)
	w.Write(bodyBytes)
}

func GetAllJson(w http.ResponseWriter, r *http.Request) {
	storeMap, err := M.GetAll()
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bodyBytes, err := json.Marshal(storeMap)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println(string(bodyBytes))
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}
//---------------------------------------------------------------------------GetAll---
func GetAll(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GetAll --------------------------------------enter---------------")
	str, err := M.GetAllValues()
	// storeMap, err := M.GetAll()
	if err != nil {
		w.Header().Set("content-type", "test/html")
		// w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// bodyBytes, err := json.Marshal(storeMap)
	// if err != nil {
	// 	w.Header().Set("content-type", "application/json")
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	fmt.Println("str:   ", str)
	// fmt.Println(string(bodyBytes))
	w.Header().Set("content-type", "text/html")
	// w.Header().Set("content-type", "text/html, text/plain; charset=utf-8, application/json")
	// w.Header().Set("Content-Encoding", "gzip")
	// w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "%s", bodyBytes)
	w.Write([]byte(str))
    fmt.Println("w:   ", w)
}

func PostCounterCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var nameData string
		var valueData int

		nameDataStr := chi.URLParam(r, "nameData")
		valueDataStr := chi.URLParam(r, "valueData")

		if nameDataStr == "" || valueDataStr == "" {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			// fmt.Fprint(w, "http.StatusNotFound")
			return
		}
		nameData = nameDataStr
		valueData, err := strconv.Atoi(valueDataStr)
		if err != nil {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			// fmt.Fprint(w, "http.StatusBadRequest")
			return
		}

		ctx := context.WithValue(r.Context(), nameDataKey, nameData)
		ctx = context.WithValue(ctx, valueDataKey, valueData)
        next.ServeHTTP(w, r.WithContext(ctx))// -------------------------------371
	})
}
func GetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var nameData string
		nameDataStr := chi.URLParam(r, "nameData")
		if nameDataStr == "" {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			// fmt.Fprint(w, "line: 115; http.StatusNotFound")
			return
		}
		nameData = nameDataStr

		ctx := context.WithValue(r.Context(), nameDataKey, nameData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PostGaugeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var nameData string
		var valueData float64

		nameDataStr := chi.URLParam(r, "nameData")
		valueDataStr := chi.URLParam(r, "valueData")

		if nameDataStr == "" || valueDataStr == "" {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			// fmt.Fprint(w, "http.StatusNotFound")
			return
		}
		nameData = nameDataStr
		valueData, err := strconv.ParseFloat(valueDataStr, 64)
		if err != nil {
			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			// fmt.Fprint(w, "http.StatusBadRequest")
			return
		}

		ctx := context.WithValue(r.Context(), nameDataKey, nameData)
		ctx = context.WithValue(ctx, valueDataKey, valueData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func PostUpdateCounter(w http.ResponseWriter, r *http.Request) {
	valueData := r.Context().Value(valueDataKey).(int)
	nameData := r.Context().Value(nameDataKey).(string)
	// fmt.Println("__________________________________\n", M)
    sm, _ := M.GetAllValues()
	fmt.Println("__________________________________\n", sm)
	counter, err := M.SaveCounterValue(nameData, store.Counter(valueData))//------430
    if err != nil {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusBadRequest)
        // fmt.Fprint(w, counter)
        return
    }
    // fmt.Println("__________________________________", M)
    sm, _ = M.GetAllValues()
	fmt.Println("__________________________________\n", sm)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, counter)
}
func PostUpdateGauge(w http.ResponseWriter, r *http.Request) {
	valueData := r.Context().Value(valueDataKey).(float64)
	nameData := r.Context().Value(nameDataKey).(string)
    // fmt.Println("__________________________________\n", M)
    sm, _ := M.GetAllValues()
	fmt.Println("__________________________________\n", sm)
	err := M.SaveGaugeValue(nameData, store.Gauge(valueData))
    if err != nil {
        w.Header().Set("content-type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusBadRequest)
        // fmt.Fprint(w, valueData)``
        return
    }
    // fmt.Println("__________________________________\n", M)
    sm, _ = M.GetAllValues()
	fmt.Println("__________________________________\n", sm)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, valueData)
}
func GetCounter(w http.ResponseWriter, r *http.Request) {
	nameData := r.Context().Value(nameDataKey).(string)
    // fmt.Println("r.Header: ", r.Header)
    sm, _ := M.GetAllValues()
	fmt.Println("__________________________________\n", sm)
	n, err := M.GetCounterValue(nameData)
     sm, _ = M.GetAllValues()
	fmt.Println("__________________________________\n", sm)
	if err != nil {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "line: 175; http.StatusNotFound")
	} else {
	    w.Header().Set("content-type", "text/html")
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", n)
	}
}
func GetGauge(w http.ResponseWriter, r *http.Request) {
	nameData := r.Context().Value(nameDataKey).(string)
     sm, _ := M.GetAllValues()
	fmt.Println("__________________________________\n", sm)
	n, err := M.GetGaugeValue(nameData)
     sm, _ = M.GetAllValues()
	fmt.Println("__________________________________\n", sm)
	if err != nil {
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		// fmt.Fprint(w, "line: 188; http.StatusNotFound, error: ", err)
	} else {
	    w.Header().Set("content-type", "text/html")
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", n)
	}
}
