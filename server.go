package mse6

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var waitDuration time.Duration
var Version = "v0.1.5"

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{"mse6":"Hello from the get endpoint"}`))
	log.Info().Msg("served /get request")
}

func die(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("serving /die request, process exiting with -1")
	os.Exit(-1)
}

func post(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(`{"mse6":"Hello from the post endpoint"}`))
		log.Info().Msg("served /post request")
	} else {
		send404(w, r)
	}
}

func slowbody(w http.ResponseWriter, r *http.Request) {
	wd := parseWaitDuration(r)

	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	//we must have this, else golang sets it to 'chunked' after 2nd write
	w.Header().Set("Transfer-Encoding", "identity")
	w.WriteHeader(200)

	hj, _ := w.(http.Hijacker)
	conn, bufrw, _ := hj.Hijack()
	defer conn.Close()

	//sleep half the wait duration and send a few bytes
	time.Sleep(wd / 2)
	bufrw.WriteString(`[{"mse6":"Hello from the slowbody endpoint"}`)
	bufrw.Flush()

	//sleep some more and send the rest
	time.Sleep(wd / 2)
	bufrw.WriteString(fmt.Sprintf(`,{"mse6":"and some more data from the slowbody endpoint", "waitSeconds":"%d"}]`, int(wd.Seconds())))
	bufrw.Flush()

	log.Info().Msgf("served /slowbody request in %d seconds", int(wd.Seconds()))
}

func parseWaitDuration(r *http.Request) time.Duration {
	var wd time.Duration = waitDuration
	if len(r.URL.Query()["wait"]) > 0 {
		ws, _ := strconv.Atoi(r.URL.Query()["wait"][0])
		if ws > 0 {
			wd = time.Duration(time.Duration(ws) * time.Second)
		} else {
			log.Warn().Msgf("unable to parse wait parameter, using default %d seconds", int(waitDuration.Seconds()))
		}
	}
	return wd
}

func badcontentlength(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.Header().Set("Content-Length", "2048")
	//we must have this, else golang sets it to 'chunked' after 2nd write
	w.Header().Set("Transfer-Encoding", "identity")
	w.WriteHeader(200)

	hj, _ := w.(http.Hijacker)
	conn, bufrw, _ := hj.Hijack()
	defer conn.Close()

	//sleep half the wait duration and send a few bytes
	time.Sleep(waitDuration / 2)
	bufrw.WriteString(`[{"mse6":"Hello from the badcontentlength endpoint"}`)
	bufrw.Flush()

	//sleep some more and send the rest
	time.Sleep(waitDuration / 2)
	bufrw.WriteString(`,{"mse6":"and some more data from the badcontentlength endpoint"}]`)
	bufrw.Flush()

	log.Info().Msg("served /badcontentlength request")
}

func slowheader(w http.ResponseWriter, r *http.Request) {
	wd := parseWaitDuration(r)
	time.Sleep(wd)
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"mse6":"Hello from the slowheader endpoint", "waitSeconds":"%d"}`, int(wd.Seconds()))))

	log.Info().Msgf("served /slowheader request in %v seconds", int(wd.Seconds()))
}

func gzipf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	w.Write(gzipenc([]byte(`{"mse6":"Hello from the gzip endpoint"}`)))

	log.Info().Msg("served /gzip request")
}

func badgzipf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	gzipBytes := gzipenc([]byte(`{"mse6":"Hello from the gzip endpoint"}`))
	badBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0}
	copy(gzipBytes, badBytes)
	w.Write(gzipBytes)

	log.Info().Msg("served /badgzip request")
}

func send404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(404)
	w.Write([]byte(`{"mse6":"404"}`))

	log.Info().Msg("served /send404 request")
}

func send(w http.ResponseWriter, r *http.Request) {
	code := 0
	if len(r.URL.Query()["code"]) > 0 {
		code, _ = strconv.Atoi(r.URL.Query()["code"][0])
		if !(code > 99 && code < 1000) {
			code = 200
		}
	} else {
		code = 200
	}
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`{"mse6":"%d"}`, code)))

	log.Info().Msgf("served /send request with code %d", code)
}

func Bootstrap(port int, waitSeconds float64) {
	waitDuration = time.Second * time.Duration(waitSeconds)
	log.Info().Msgf("wait duration for slow requests seconds %v", waitDuration.Seconds())
	log.Info().Msgf("mse6 %s starting http server on port %d", Version, port)

	http.HandleFunc("/mse6/die", die)
	http.HandleFunc("/mse6/get", get)
	http.HandleFunc("/mse6/post", post)
	http.HandleFunc("/mse6/slowbody", slowbody)
	http.HandleFunc("/mse6/slowheader", slowheader)
	http.HandleFunc("/mse6/badcontentlength", badcontentlength)
	http.HandleFunc("/mse6/send", send)
	http.HandleFunc("/mse6/gzip", gzipf)
	http.HandleFunc("/mse6/badgzip", badgzipf)
	http.HandleFunc("/", send404)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err.Error())
	}
}
