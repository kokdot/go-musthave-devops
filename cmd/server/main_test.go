package main

import (
	// "encoding/json"
	// "fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	// "time"

	// "bytes"

	"github.com/stretchr/testify/assert"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kokdot/go-musthave-devops/internal/handler"
	"github.com/kokdot/go-musthave-devops/internal/interfaceinit"
	"github.com/kokdot/go-musthave-devops/internal/onboardingserver"
	"github.com/kokdot/go-musthave-devops/internal/repo"
	"github.com/kokdot/go-musthave-devops/internal/store"
	"github.com/kokdot/go-musthave-devops/internal/metricsserver"
	// "github.com/kokdot/go-musthave-devops/internal/downloadingtofile"
)

//test git git test what
func TestHandler(t *testing.T) {
url, storeFile, key, restore, storeInterval, dataBaseDSNReal, logg  := onboardingserver.OnboardingServer()
    logg.Print("--------------------main-------------------------------------------")
    logg.Print("url:  ", url)
    logg.Print("storeInterval:  ", storeInterval)
    logg.Print("storeFile:  ", storeFile)
    logg.Print("restore:  ", restore)
    logg.Print("key:  ", key)
    logg.Print("dataBaseDSNReal:  ", dataBaseDSNReal)

    m, err := interfaceinit.InterfaceInit(storeInterval, storeFile, restore, url, key, dataBaseDSNReal, logg)
    if err != nil {
        logg.Printf("\nthere in error in starting interface and restore data: %s", err)
    }
    handler.PutM(m)
    handler.GetLogg(logg)
    store.GetLogg(logg)
    metricsserver.GetLogg(logg)
    logg.Printf("interface m:   %#v", m)
    logg.Print("--------------------main--started-----------------------------------------")
    
    // if m.GetDataBaseDSN() != "" {
    //     downloading_to_file.DownloadingToFile(m)
    // }

    // определяем роутер chi
    r := chi.NewRouter()
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)    
    r.Use(middleware.Recoverer)
    // r.Use(middleware.Compress(5, "gzip"))
    r.Use(middleware.Compress(5))
    r.Get("/", handler.GetAll)
    r.Get("/ping", handler.GetPing)

    r.Post("/updates/", handler.PostUpdateByBatch)
    r.Post("/updates1/", handler.PostUpdateByBatch1)
    
    r.Route("/update", func(r chi.Router) {
        r.Post("/", handler.PostUpdate)
        r.Route("/counter", func(r chi.Router) {
            r.Route("/{nameData}/{valueData}", func(r chi.Router) {
                r.Use(handler.PostCounterCtx)
                r.Post("/", handler.PostUpdateCounter)
            })
        })
        r.Route("/gauge", func(r chi.Router) {
            r.Route("/{nameData}/{valueData}", func(r chi.Router) {
                r.Use(handler.PostGaugeCtx)
                r.Post("/", handler.PostUpdateGauge)
            })
        })
        r.Route("/",func(r chi.Router) {
            r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
		        w.Header().Set("content-type", "text/plain; charset=utf-8")
                w.WriteHeader(http.StatusNotImplemented)
                // fmt.Fprint(w, "line: 52; http.StatusNotImplemented")
	        })
        })
    })

    r.Route("/value", func(r chi.Router) {
        r.Post("/", handler.GetValue)
		r.Route("/counter", func(r chi.Router){
            r.Route("/{nameData}", func(r chi.Router) {
                r.Use(handler.GetCtx)
                r.Get("/", handler.GetCounter)
            })
        })
       	r.Route("/gauge", func(r chi.Router){
            r.Route("/{nameData}", func(r chi.Router) {
                r.Use(handler.GetCtx)
                r.Get("/", handler.GetGauge)
            })
        })
	})

	type want struct {
		StatusCode  int
		contentType string
		result      string
		// mtxNew store.Metrics
	}
	// var valGauge store.Gauge = 4.6 
	// var valCounter store.Counter = 64
	tests := []struct {
		name   string
		want   want
		url    string
		method string
		mtxOld repo.Metrics
	}{
		{
			name: "counter norm",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				result:      "5",
			},
			url:    "/update/counter/PollCount/5",
			method: http.MethodPost,
		},
		{
			name: "counter update norm",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				result:      "10",
			},
			url:    "/update/counter/PollCount/5",
			method: http.MethodPost,
		},
		{
			name: "counter error",
			want: want{
				StatusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			url:    "/update/counter/PollCount/none",
			method: http.MethodPost,
		},
		{
			name: "gauge norm",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				result:      "3.6",
			},
			url:    "/update/gauge/Alloc/3.6",
			method: http.MethodPost,
		},
		{
			name: "gauge error",
			want: want{
				StatusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			url:    "/update/gauge/Alloc/none",
			method: http.MethodPost,
		},
		{
			name: "no counter no gauge",
			want: want{
				StatusCode:  http.StatusNotImplemented,
				contentType: "text/plain; charset=utf-8",
			},
			url:    "/update/error/PollCount/100",
			method: http.MethodPost,
		},
		{
			name: "default",
			want: want{
				StatusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
			},
			url:    "/",
			method: http.MethodPost,
		},
		{
			name: "get counter",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				result:      "10",
			},
			method: http.MethodGet,
			url:    "/value/counter/PollCount",
		},
		{
			name: "get gauge",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				result:      "3.6",
			},
			url:    "/value/gauge/Alloc",
			method: http.MethodGet,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// tt.Skip()
			request := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			result := w.Result()
			body, err := io.ReadAll(result.Body)
			result.Body.Close()

			assert.NoError(t, err)
			assert.Equal(t, tt.want.StatusCode, result.StatusCode)
			logg.Print(tt.name)
			if tt.name != "default" {
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			}
			// if tt.method == http.MethodGet {
				assert.Equal(t, tt.want.result, string(body))
			// }
		})
	}
	// for _, tt := range testsJson {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		logg.Print(tt.name)
	// 		logg.Print("tt.mtxOld:   ", int(*tt.mtxOld.Delta))
	// 		bodyBytes, _ := json.Marshal(tt.mtxOld)
	// 		buf := bytes.NewReader(bodyBytes)
	// 		request := httptest.NewRequest(tt.method, tt.url, buf)
	// 		w := httptest.NewRecorder()
	// 		r.ServeHTTP(w, request)
	// 		// r.ServeHTTP(w, request)
	// 		// r.ServeHTTP(w, request)
	// 		result := w.Result()
	// 		body, _ := io.ReadAll(result.Body)
    // 		var mtxNew store.Metrics
	// 		_ = json.Unmarshal(bodyBytes, &mtxNew)
	// 		result.Body.Close()
	// 		// assert.NoError(t, err)
	// 		logg.Print("mtxNew:   ", int(*mtxNew.Delta))
	// 		assert.Equal(t, tt.want.StatusCode, result.StatusCode)
	// 		assert.Equal(t, tt.want.mtxNew, mtxNew)
	// 		if tt.name != "default" {
	// 			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
	// 		}
	// 		if tt.method == http.MethodGet {
	// 			assert.Equal(t, tt.want.result, string(body))
	// 		}
	// 	})
	// }
}

	//json-----------------------------------------------------
	// testsJson := []struct {
	// 	name   string
	// 	want   want
	// 	url    string
	// 	method string
	// 	mtxOld store.Metrics
	// }{
	// 	{
	// 		name: "counter norm update json",
	// 		want: want{
	// 			StatusCode:  http.StatusOK,
	// 			contentType: "application/json",
	// 			mtxNew:  store.Metrics{
	// 				ID: "PollCount",
	// 				MType: "Counter",
	// 				Delta: &valCounter,
	// 			},
	// 		},
	// 		url:    "/update/",
	// 		method: http.MethodPost,
	// 		mtxOld:  store.Metrics{
	// 			ID: "PollCount",
	// 			MType: "Counter",
	// 			Delta: &valCounter,
	// 		},
	// 	},
	// 	// {
	// 	// 	name: "counter error update json",
	// 	// 	want: want{
	// 	// 		StatusCode:  http.StatusOK,
	// 	// 		contentType: "application/json",
	// 	// 		mtxNew:  store.Metrics{
	// 	// 			ID: "PollCount",
	// 	// 			MType: "Counter",
	// 	// 			// Delta: &valCounter,
	// 	// 		},
	// 	// 	},
	// 	// 	url:    "/update/",
	// 	// 	method: http.MethodPost,
	// 	// 	mtxOld:  store.Metrics{
	// 	// 		ID: "PollCount",
	// 	// 		MType: "Counter",
	// 	// 		Delta: &valCounter,
	// 	// 	},
	// 	// },
	// 	// {
	// 	// 	name: "gauge norm update json",
	// 	// 	want: want{
	// 	// 		StatusCode:  http.StatusOK,
	// 	// 		contentType: "application/json",
	// 	// 		mtxNew:  store.Metrics{
	// 	// 			ID: "PollCount",
	// 	// 			MType: "Counter",
	// 	// 			Value: &valGauge,
	// 	// 		},
	// 	// 	},
	// 	// 	url:    "/update/",
	// 	// 	method: http.MethodPost,
	// 	// 	mtxOld:  store.Metrics{
	// 	// 		ID: "PollCount",
	// 	// 		MType: "Counter",
	// 	// 		Value: &valGauge,
	// 	// 	},
	// 	// },
	// }
	// ts := httptest.NewServer(r)
	// defer ts.Close()
// {
		// 	name: "update counter json norm",
		// 	want: want{
		// 		StatusCode:  http.StatusOK,
		// 		contentType: "application/json",
		// 		mtxNew:  store.Metrics{
		// 			ID: "PollCount",
		// 			MType: "Counter",
		// 			Delta: &valCounter,
		// 		},
		// 	},
		// 	url:    "/update/",
		// 	method: http.MethodPost,
		// 	mtxOld:  store.Metrics{
		// 		ID: "PollCount",
		// 		MType: "Counter",
		// 		Delta: &valCounter,
		// 	},
		// },
		// {
		// 	name: "counter error update json",
		// 	want: want{
		// 		StatusCode:  http.StatusOK,
		// 		contentType: "application/json",
		// 		mtxNew:  store.Metrics{
		// 			ID: "PollCount",
		// 			MType: "Counter",
		// 			// Delta: &valCounter,
		// 		},
		// 	},
		// 	url:    "/update/",
		// 	method: http.MethodPost,
		// 	mtxOld:  store.Metrics{
		// 		ID: "PollCount",
		// 		MType: "Counter",
		// 		Delta: &valCounter,
		// 	},
		// },
			// 		{
		// 	name: "gauge norm update json",
		// 	want: want{
		// 		StatusCode:  http.StatusOK,
		// 		contentType: "application/json",
		// 		mtxNew:  store.Metrics{
		// 			ID: "PollCount",
		// 			MType: "Counter",
		// 			Value: &valGauge,
		// 		},
		// 	},
		// 	url:    "/update/",
		// 	method: http.MethodPost,
		// 	mtxOld:  store.Metrics{
		// 		ID: "PollCount",
		// 		MType: "Counter",
		// 		Value: &valGauge,
		// 	},
		// },