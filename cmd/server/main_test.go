package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"io"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kokdot/go-musthave-devops/internal/store"

)

func TestHandler(t *testing.T) {
	r := chi.NewRouter()
	var ms = new(store.MemStorage)
	ms.GaugeMap = make(store.GaugeMap)
	ms.CounterMap = make(store.CounterMap)
	m = ms
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Get("/", GetAll)
    r.Route("/update", func(r chi.Router) {
        r.Route("/counter", func(r chi.Router) {
            r.Route("/{nameData}/{valueData}", func(r chi.Router) {
                r.Use(PostCounterCtx)
                r.Post("/", PostUpdateCounter)
            })
        })
        r.Route("/gauge", func(r chi.Router) {
            r.Route("/{nameData}/{valueData}", func(r chi.Router) {
                r.Use(PostGaugeCtx)
                r.Post("/", PostUpdateGauge)
            })
        })
        r.Route("/",func(r chi.Router) {
            r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
		        w.Header().Set("content-type", "text/plain; charset=utf-8")
                w.WriteHeader(http.StatusNotImplemented)
                fmt.Fprint(w, "line: 52; http.StatusNotImplemented")
	        })
        })
    })

    r.Route("/value", func(r chi.Router) {
		r.Route("/counter", func(r chi.Router){
            r.Route("/{nameData}", func(r chi.Router) {
                r.Use(GetCtx)
                r.Get("/", GetCounter)
            })
        })
       	r.Route("/gauge", func(r chi.Router){
            r.Route("/{nameData}", func(r chi.Router) {
                r.Use(GetCtx)
                r.Get("/", GetGauge)
            })
        })
	})

	type want struct {
		StatusCode  int
		contentType string
		result string
	}
	tests := []struct {
		name string
		want want
		url  string
		method string
	}{
		{
			name: "counter norm",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/PollCount/5",
			method: http.MethodPost,
		},
		{
			name: "counter error",
			want: want{
				StatusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/PollCount/none",
			method: http.MethodPost,
		},
		{
			name: "gauge norm",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/gauge/Alloc/3.6",
			method: http.MethodPost,
		},
		{
			name: "gauge error",
			want: want{
				StatusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/gauge/Alloc/none",
			method: http.MethodPost,
		},
		{
			name: "no counter no gauge",
			want: want{
				StatusCode:  http.StatusNotImplemented,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/error/PollCount/100",
			method: http.MethodPost,
		},
		{
			name: "default",
			want: want{
				StatusCode:  http.StatusMethodNotAllowed             ,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/",
			method: http.MethodPost,
		},
		{
			name: "get counter",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				result: "5",
			},
			method: http.MethodGet,
			url: "/value/counter/PollCount",
		},
		{
			name: "get gauge",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				result: "3.6",
			},
			url: "/value/gauge/Alloc",
			method: http.MethodGet,
		},
	}
	ts := httptest.NewServer(r)
	defer ts.Close()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			result := w.Result()

			body, err := io.ReadAll(result.Body) 
			result.Body.Close()

			assert.NoError(t, err)
			assert.Equal(t, tt.want.StatusCode, result.StatusCode)
			fmt.Println(tt.name)
			if tt.name != "default"{
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			}
			if tt.method == http.MethodGet {
				assert.Equal(t, tt.want.result, string(body))
			}
		})
	}
}