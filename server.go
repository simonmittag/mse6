package mse6

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var waitDuration time.Duration
var Version = "v0.1.8"
var Port int
var Prefix string

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{"mse6":"Hello from the get endpoint"}`))
	log.Info().Msgf("served %v request", r.URL.Path)
}

func redirected(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{"mse6":"Hello from the redirected endpoint. if you're reading this and you didn't load this URL, chances are you've been redirected.'"}`))
	log.Info().Msgf("served %v request", r.URL.Path)
}

func die(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("served %v request, process exiting with -1", r.URL.Path)
	os.Exit(-1)
}

func post(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(`{"mse6":"Hello from the post endpoint"}`))
		log.Info().Msgf("served %v post request", r.URL.Path)
	} else {
		send404(w, r)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(`{"mse6":"Hello from the put endpoint"}`))
		log.Info().Msgf("served %v put request", r.URL.Path)
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

	log.Info().Msgf("served %v request in %d seconds", r.URL.Path, int(wd.Seconds()))
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

	log.Info().Msgf("served %v request", r.URL.Path)
}

func slowheader(w http.ResponseWriter, r *http.Request) {
	wd := parseWaitDuration(r)
	time.Sleep(wd)
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"mse6":"Hello from the slowheader endpoint", "waitSeconds":"%d"}`, int(wd.Seconds()))))

	log.Info().Msgf("served %v request in %v seconds", r.URL.Path, int(wd.Seconds()))
}

func gzipf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	w.Write(gzipenc([]byte(`{"mse6":"Hello from the gzip endpoint"}`)))

	log.Info().Msgf("served %v request", r.URL.Path)
}

func badgzipf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	gzipBytes := gzipenc([]byte(`{"mse6":"Hello from the gzip endpoint"}`))
	badBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0}
	copy(gzipBytes, badBytes)
	w.Write(gzipBytes)

	log.Info().Msgf("served %v request", r.URL.Path)
}

func send404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(404)
	w.Write([]byte(`{"mse6":"404"}`))

	log.Info().Msgf("served %v request", r.URL.Path)
}

func send(w http.ResponseWriter, r *http.Request) {
	code := 0
	location := ""
	if len(r.URL.Query()["code"]) > 0 {
		code, _ = strconv.Atoi(r.URL.Query()["code"][0])
		if !(code > 99 && code < 1000) {
			code = 200
		}
	} else {
		code = 200
	}

	host := ""
	if strings.Contains(r.Host, ":") {
		host = strings.Split(r.Host, ":")[0]
	} else {
		host = r.Host
	}

	if len(r.URL.Query()["url"]) > 0 {
		location = r.URL.Query()["url"][0]
	} else {
		location = fmt.Sprintf("http://%s:%d%sredirected", host, Port, Prefix)
	}

	redirect := ""
	if code >= 300 && code <= 303 {
		w.Header().Set("Location", location)
		redirect = fmt.Sprintf("redirect to %s ", location)
	}
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(code)
	if code >= 200 {
		w.Write([]byte(fmt.Sprintf(`{"mse6":"%d"}`, code)))
	} else {
		//the headers are flushed at this point, it doesn't send more. you need to hijack the connection.
		log.Info().Msgf("sending additional 200 header after %d for %v", code, r.URL.Path)
		w.WriteHeader(200)
		w.Write([]byte("\r\n\r\n"))
	}

	log.Info().Msgf("served %v %vrequest with code %d", r.URL.Path, redirect, code)
}

func Bootstrap(port int, waitSeconds float64, prefix string) {
	waitDuration = time.Second * time.Duration(waitSeconds)
	log.Info().Msgf("wait duration for slow requests seconds %v", waitDuration.Seconds())

	Port = port
	Prefix = prefix
	log.Info().Msgf("mse6 %s starting http server on port %d with prefix '%s'", Version, Port, Prefix)

	http.HandleFunc(prefix+"die", die)
	http.HandleFunc(prefix+"get", get)
	http.HandleFunc(prefix+"redirected", redirected)
	http.HandleFunc(prefix+"post", post)
	http.HandleFunc(prefix+"put", put)
	http.HandleFunc(prefix+"slowbody", slowbody)
	http.HandleFunc(prefix+"slowheader", slowheader)
	http.HandleFunc(prefix+"badcontentlength", badcontentlength)
	http.HandleFunc(prefix+"send", send)
	http.HandleFunc(prefix+"gzip", gzipf)
	http.HandleFunc(prefix+"badgzip", badgzipf)
	http.HandleFunc("/", send404)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err.Error())
	}
}
