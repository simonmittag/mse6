package mse6

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

var waitDuration time.Duration
var version = "v0.1.0"

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{"mse6":"Hello from the getting endpoint"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{"mse6":"Hello from the posting endpoint"}`))
}

func slowbody(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "identity")
	//we must have this, else golang sets it to 'chunked' after 2nd write
	w.Header().Set("Transfer-Encoding", "identity")
	w.WriteHeader(200)

	hj, _ := w.(http.Hijacker)
	conn, bufrw, _ := hj.Hijack()
	defer conn.Close()

	//sleep half the wait duration and send a few bytes
	time.Sleep(waitDuration / 2)
	bufrw.WriteString(`[{"mse6":"Hello from the slowbody endpoint"}`)
	bufrw.Flush()

	//sleep some more and send the rest
	time.Sleep(waitDuration / 2)
	bufrw.WriteString(`,{"mse6":"and some more data from the slowbody endpoint"}]`)
	bufrw.Flush()
}

func slowheader(w http.ResponseWriter, r *http.Request) {
	time.Sleep(waitDuration)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{"mse6":"Hello from the slowheader endpoint"}`))
}

func gzipf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	w.Write(gzipenc([]byte(`{"mse6":"Hello from the gzip endpoint"}`)))
}

func Bootstrap(port int, waitSeconds float64) {
	waitDuration = time.Second * time.Duration(waitSeconds)
	log.Info().Msgf("wait duration for slow requests seconds %v", waitDuration.Seconds())
	log.Info().Msgf("mse6 starting http server %s on port %d", version, port)

	http.HandleFunc("/mse6/get", get)
	http.HandleFunc("/mse6/post", post)
	http.HandleFunc("/mse6/slowbody", slowbody)
	http.HandleFunc("/mse6/slowheader", slowheader)
	http.HandleFunc("/mse6/gzip", gzipf)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err.Error())
	}
}
