package mse6

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
	"time"
)

var zipPool = sync.Pool{
	New: func() interface{} {
		var buf bytes.Buffer
		return gzip.NewWriter(&buf)
	},
}



func Bootstrap(port int, waitSeconds float64) {
	waitDuration := time.Second * time.Duration(waitSeconds)
	log.Info().Msgf("wait duration for slow requests seconds %v", waitDuration.Seconds())
	log.Info().Msgf("mse6 starting http server on port %d", port)

	http.HandleFunc("/mse6/getting", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(`{"MSE6":"Hello from the getting endpoint"}`))
	})

	http.HandleFunc("/mse6/posting", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(`{"MSE6":"Hello from the posting endpoint"}`))
	})

	http.HandleFunc("/mse6/slowbody", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "identity")
		//we must have this, else golang sets it to 'chunked' after 2nd write
		w.Header().Set("Transfer-Encoding", "identity")
		w.WriteHeader(200)
		hj, _ := w.(http.Hijacker)
		conn, bufrw, _ := hj.Hijack()
		defer conn.Close()
		time.Sleep(waitDuration/2)
		bufrw.WriteString(`[{"mse6":"Hello from the slowbody endpoint"}`)
		bufrw.Flush()
		time.Sleep(waitDuration/2)
		bufrw.WriteString(`,{"mse6":"and some more data from the slowbody endpoint"}]`)
		bufrw.Flush()
	})

	http.HandleFunc("/mse6/slowheader", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(waitDuration)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(`{"MSE6":"Hello from the slowheader endpoint"}`))
	})

	http.HandleFunc("/mse6/gzip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(200)
		w.Write(gzipenc([]byte(`{"MSE6":"Hello from the gzip endpoint"}`)))
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err.Error())
	}
}
