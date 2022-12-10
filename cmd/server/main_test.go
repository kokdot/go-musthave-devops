package main

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"io"
	// "github.com/kokdot/go-musthave-devops/internal/store"
)

func TesPostUpdateCounter(t *testing.T) {
	// var ms = new(store.MemStorage)
	// ms.GaugeMap = make(store.GaugeMap)
	// ms.CounterMap = make(store.CounterMap)
	// var m store.Repo = ms
	type want struct {
		StatusCode  int
		contentType string
	}
	tests := []struct {
		name string
		want want
		url  string
	}{
		{
			name: "gauge norm",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/gauge/Alloc/3.6",
		},
		{
			name: "gauge error",
			want: want{
				StatusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/gauge/Alloc/none",
		},
		{
			name: "counter norm",
			want: want{
				StatusCode:  http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/PollCount/5",
		},
		{
			name: "counter error",
			want: want{
				StatusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/counter/PollCount/none",
		},
		{
			name: "no counter no gauge",
			want: want{
				StatusCode:  http.StatusNotImplemented,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/update/error/PollCount/100",
		},
		{
			name: "default",
			want: want{
				StatusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
			url: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(PostUpdateGauge)
			// запускаем сервер
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			_, err := io.ReadAll(result.Body) 
			assert.NoError(t, err)
			assert.Equal(t, tt.want.StatusCode, result.StatusCode)
			// fmt.Println("result.Header.Get(\"Content-Type\"): ", result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

		})
	}
}
