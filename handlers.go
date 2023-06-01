package mse6

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

func connect(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		w.Header().Set("Server", "mse6 "+Version)
		w.WriteHeader(200)
		log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
		if len(r.URL.Query()["body"]) > 0 {
			w.Write([]byte(`{"mse6":"Hello from the connect endpoint"}`))
		}
	} else {
		send405(w, r)
	}
}

func options(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		code := 0
		if len(r.URL.Query()["code"]) > 0 {
			code, _ = strconv.Atoi(r.URL.Query()["code"][0])
		} else {
			code = 200
		}
		w.Header().Add("Allow", "OPTIONS")
		w.Header().Add("Allow", "GET")
		w.Header().Set("Server", "mse6 "+Version)
		w.WriteHeader(code)
		if len(r.URL.Query()["body"]) > 0 {
			w.Write([]byte(`{"mse6":"Hello from the options endpoint"}`))
		}

		log.Info().Msgf("served %v OPTIONS request with X-Request-Id %s code %d", r.URL.Path, getXRequestId(r), code)
	} else {
		send405(w, r)
	}
}

func trace(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if r.Method == "TRACE" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Type", "message/http")
		w.WriteHeader(200)
		w.Write([]byte(`{"mse6":"Hello from the trace endpoint"}`))
		log.Info().Msgf("served %v trace with X-Request-Id %s,%s reading %d bytes from inbound", r.URL.Path, getXRequestId(r), expectContinue(r), len(body))
	} else {
		send405(w, r)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(`{"mse6":"Hello from the get endpoint"}`))
		log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
	} else {
		send405(w, r)
	}
}

func nocontentenc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.WriteHeader(200)
	w.Write([]byte(`{"mse6":"Hello from the nocontentenc endpoint"}`))
	log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func unknowncontentenc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "unknown")
	w.WriteHeader(200)
	w.Write([]byte(`{"mse6":"Hello from the unknowncontentenc endpoint"}`))
	log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func echoquery(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Encode()
	if r.Method == "GET" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf(`{"mse6":"Hello from the echo query endpoint. Your query string was %v"}`, q)))
		log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
	} else {
		send405(w, r)
	}
}

func echoport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf(`{"mse6":"Hello from the echo port endpoint. My port is %v"}`, Port)))
		log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
	} else {
		send405(w, r)
	}
}

func echoheader(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf(`{"mse6":"Hello from the echo header endpoint. %v"}`, r.Header)))
		log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
	} else {
		send405(w, r)
	}
}

func redirected(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{"mse6":"Hello from the redirected endpoint. if you're reading this and you didn't load this URL, chances are you've been redirected.'"}`))
	log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func post(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if r.Method == "POST" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(201)
		w.Write([]byte(`{"mse6":"Hello from the post endpoint"}`))
		log.Info().Msgf("served %v post request with X-Request-Id %s,%s reading %d bytes from inbound", r.URL.Path, getXRequestId(r), expectContinue(r), len(body))
	} else {
		send405(w, r)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if r.Method == "PUT" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(`{"mse6":"Hello from the put endpoint"}`))
		log.Info().Msgf("served %v put request with X-Request-Id %s,%s reading %d bytes from inbound", r.URL.Path, getXRequestId(r), expectContinue(r), len(body))
	} else {
		send405(w, r)
	}
}

func patch(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if r.Method == "PATCH" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(200)
		w.Write([]byte(`{"mse6":"Hello from the patch endpoint"}`))
		log.Info().Msgf("served %v patch request with X-Request-Id %s,%s reading %d bytes from inbound", r.URL.Path, getXRequestId(r), expectContinue(r), len(body))
	} else {
		send405(w, r)
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if r.Method == "DELETE" {
		w.Header().Set("Server", "mse6 "+Version)
		w.WriteHeader(204)
		log.Info().Msgf("served %v delete request with X-Request-Id %s,%s reading %d bytes from inbound", r.URL.Path, getXRequestId(r), expectContinue(r), len(body))
	} else {
		send405(w, r)
	}
}

func expectContinue(r *http.Request) string {
	c100 := strings.ToLower(r.Header.Get("Expect"))
	if c100 == "100-continue" {
		return " incoming request with Expect: 100-continue "
	} else {
		return " "
	}
}

func chunked(w http.ResponseWriter, r *http.Request) {
	wd := parseWaitDuration(r)

	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Encoding", "identity")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(200)

	f, _ := w.(http.Flusher)

	w.Write([]byte(`[{"mse6":"Hello from the chunked endpoint"}`))
	f.Flush()
	//sleep some
	time.Sleep(wd)
	w.Write([]byte(fmt.Sprintf(`,{"mse6":"and some more data from the chunked endpoint", "waitSeconds":"%d"}]`, int(wd.Seconds()))))
	f.Flush()

	log.Info().Msgf("served %v chunked request with X-Request-Id %s in %d seconds", r.URL.Path, getXRequestId(r), int(wd.Seconds()))
}

func slowbody(w http.ResponseWriter, r *http.Request) {
	wd := parseWaitDuration(r)

	hj, _ := w.(http.Hijacker)
	conn, bufrw, _ := hj.Hijack()
	defer conn.Close()

	bufrw.WriteString("HTTP/1.1 200 OK")
	bufrw.WriteString(fmt.Sprintf("\nServer: mse6 %s", Version))
	bufrw.WriteString("\nContent-Encoding: identity")
	bufrw.WriteString("\nConnection: close")
	bufrw.WriteString("\n")
	bufrw.WriteString("\n")
	bufrw.Flush()

	//sleep half the wait duration and send a few bytes
	time.Sleep(wd / 2)
	bufrw.WriteString(`[{"mse6":"Hello from the slowbody endpoint"}`)
	bufrw.Flush()

	//sleep some more and send the rest
	time.Sleep(wd / 2)
	bufrw.WriteString(fmt.Sprintf(`,{"mse6":"and some more data from the slowbody endpoint", "waitSeconds":"%d"}]`, int(wd.Seconds())))
	bufrw.Flush()

	log.Info().Msgf("served %v request with X-Request-Id %s in %d seconds", r.URL.Path, getXRequestId(r), int(wd.Seconds()))
}

func hangupConnDuringHeadersSend(w http.ResponseWriter, r *http.Request) {
	wd := time.Duration(time.Second * 2)

	hj, _ := w.(http.Hijacker)
	conn, bufrw, _ := hj.Hijack()

	bufrw.WriteString("HTTP/1.1 200 OK")
	bufrw.WriteString(fmt.Sprintf("\nServer: mse6 %s", Version))
	bufrw.WriteString("\nContent-Encoding: identity")
	bufrw.WriteString("\nContent-Length: 1024")
	bufrw.Flush()

	//sleep half the wait duration and send a few bytes
	time.Sleep(wd)
	bufrw.Flush()
	conn.Close()

	log.Info().Msgf("served %v incomplete request during headers send, initiated hard conn close X-Request-Id %s in %d seconds", r.URL.Path, getXRequestId(r), int(wd.Seconds()))
}

func hangupConnAfterHeadersSent(w http.ResponseWriter, r *http.Request) {
	wd := time.Duration(time.Second * 2)

	hj, _ := w.(http.Hijacker)
	conn, bufrw, _ := hj.Hijack()

	bufrw.WriteString("HTTP/1.1 200 OK")
	bufrw.WriteString(fmt.Sprintf("\nServer: mse6 %s", Version))
	bufrw.WriteString("\nContent-Encoding: identity")
	bufrw.WriteString("\nContent-Length: 1024")
	bufrw.WriteString("\n")
	bufrw.WriteString("\n")
	bufrw.Flush()

	//sleep half the wait duration and send a few bytes
	time.Sleep(wd)
	bufrw.Flush()
	conn.Close()

	log.Info().Msgf("served %v incomplete request with headers sent, but no body, initiated hard conn close X-Request-Id %s in %d seconds", r.URL.Path, getXRequestId(r), int(wd.Seconds()))
}

func hangupConnDuringBodySend(w http.ResponseWriter, r *http.Request) {
	wd := time.Duration(time.Second * 2)

	hj, _ := w.(http.Hijacker)
	conn, bufrw, _ := hj.Hijack()

	bufrw.WriteString("HTTP/1.1 200 OK")
	bufrw.WriteString(fmt.Sprintf("\nServer: mse6 %s", Version))
	bufrw.WriteString("\nContent-Encoding: identity")
	bufrw.WriteString("\nContent-Length: 1024")
	bufrw.WriteString("\n")
	bufrw.WriteString("\n")
	bufrw.WriteString(`[{"mse6":"Hello from the /hangupduringbody endpoint"}`)
	bufrw.Flush()

	//sleep half the wait duration and send a few bytes
	time.Sleep(wd)
	bufrw.Flush()
	conn.Close()

	log.Info().Msgf("served %v incomplete request with partial body sent, initiated hard conn close, X-Request-Id %s in %d seconds", r.URL.Path, getXRequestId(r), int(wd.Seconds()))
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
	hj, _ := w.(http.Hijacker)
	conn, bufrw, _ := hj.Hijack()
	defer conn.Close()

	bufrw.WriteString("HTTP/1.1 200 OK")
	bufrw.WriteString(fmt.Sprintf("\nServer: mse6 %s", Version))
	bufrw.WriteString("\nContent-Encoding: identity")
	bufrw.WriteString("\nContent-Length: 2048")
	bufrw.WriteString("\nConnection: close")
	bufrw.WriteString("\n")
	bufrw.WriteString("\n")
	bufrw.Flush()

	//sleep half the wait duration and send a few bytes
	time.Sleep(waitDuration / 2)
	bufrw.WriteString(`[{"mse6":"Hello from the badcontentlength endpoint"}`)
	bufrw.Flush()

	//sleep some more and send the rest
	time.Sleep(waitDuration / 2)
	bufrw.WriteString(`,{"mse6":"and some more data from the badcontentlength endpoint"}]`)
	bufrw.Flush()

	log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func slowheader(w http.ResponseWriter, r *http.Request) {
	wd := parseWaitDuration(r)
	time.Sleep(wd)
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"mse6":"Hello from the slowheader endpoint", "waitSeconds":"%d"}`, int(wd.Seconds()))))

	log.Info().Msgf("served %v request with X-Request-Id %s in %v seconds", r.URL.Path, getXRequestId(r), int(wd.Seconds()))
}

func tinyidentityf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{}`))

	log.Info().Msgf("served %v tiny identity request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func tinygzipf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	w.Write(gzipenc([]byte(`{}`)))

	log.Info().Msgf("served %v tiny gzip request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func gzipf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	w.Write(gzipenc([]byte(`{"mse6":"Hello from the gzip endpoint"}`)))

	log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func brotlif(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "br")
	w.WriteHeader(200)
	w.Write(*BrotliEncode([]byte(`{"mse6":"Hello from the brotli endpoint"}`)))

	log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func chooseaef(w http.ResponseWriter, r *http.Request) {
	ae := r.Header.Get("Accept-Encoding")
	log.Info().Msgf("incoming %v request with Accept-Encoding %s and X-Request-Id %s", r.URL.Path, ae, getXRequestId(r))
	if strings.Contains(ae, "br") {
		brotlif(w, r)
	} else if strings.Contains(ae, "gzip") {
		gzipf(w, r)
	} else if strings.Contains(ae, "deflate") {
		deflatef(w, r)
	} else {
		get(w, r)
	}
}

func deflatef(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "deflate")
	w.WriteHeader(200)
	w.Write(*Deflate([]byte(`{"mse6":"Hello from the deflate endpoint"}`)))

	log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func badgzipf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	gzipBytes := gzipenc([]byte(`{"mse6":"Hello from the gzip endpoint"}`))
	badBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0}
	copy(gzipBytes, badBytes)
	w.Write(gzipBytes)

	log.Info().Msgf("served %v request with X-Request-Id %s", r.URL.Path, getXRequestId(r))
}

func send404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(404)
	w.Write([]byte(`{"mse6":"404"}`))

	log.Info().Msgf("served %v request with X-Request-Id %s response code 404", r.URL.Path, getXRequestId(r))
}

func send405(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(405)
	w.Write([]byte(`{"mse6":"405"}`))

	log.Info().Msgf("served %v request with X-Request-Id %s response code 405", r.URL.Path, getXRequestId(r))
}

func getorhead(w http.ResponseWriter, r *http.Request) {
	cl := false
	if len(r.URL.Query()["cl"]) > 0 {
		cl = true
	}

	code := 200
	b := []byte(`{"mse6":"Hello from the getorhead endpoint"}`)
	cls := "0"

	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	if r.Method == "HEAD" {
		w.Header().Add("ETag", "W/0815")
		if cl {
			cls = fmt.Sprintf("%d", len(b))
		} else {
			cls = "0"
		}
		w.Header().Set("Content-Length", cls)
		w.WriteHeader(code)
		log.Info().Msgf("served %v request with X-Request-ID %s method %s content-length %s code 200", r.URL.Path, getXRequestId(r), r.Method, cls)
	} else if r.Method == "GET" {
		cls = fmt.Sprintf("%d", len(b))
		w.Header().Set("Content-Length", cls)
		w.WriteHeader(code)
		w.Write(b)
		log.Info().Msgf("served %v request with X-Request-ID %s method %s content-length %s code 200", r.URL.Path, getXRequestId(r), r.Method, cls)
	} else {
		send405(w, r)
	}

}

func getXRequestId(r *http.Request) string {
	xrid := r.Header.Get("X-Request-Id")
	if len(xrid) == 0 {
		xrid = "none"
	}
	return xrid
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

	locHeader := map[int]int{300: 300, 301: 301, 302: 302, 303: 303, 305: 305, 307: 307, 308: 308}
	if _, ok := locHeader[code]; ok {
		w.Header().Set("Location", location)
		redirect = fmt.Sprintf("redirect to %s ", location)
	}

	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(code)
	if code >= 200 {
		w.Write([]byte(fmt.Sprintf(`{"mse6":"%d"}`, code)))
	}

	log.Info().Msgf("served %v %v request with X-Request-Id %s code %d", r.URL.Path, redirect, getXRequestId(r), code)
}

func jwks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{
  "keys": [
    {
      "alg":"RS256",
      "kty":"RSA",
      "n":"w6hKIwXAuI5VqmQjvQmgZdNbV80GMC3UkPmm-OQDjzOjeLRA6yLPYLZHaGhONx37DWMA-a3D_Zg_-oueYuZlrhusbTDC-bt1JSctAJV3ollQaalmJQHhLfyL54Y6Cgt3H_68u4Q3kLrFOmdFJwRswHR-1m-Oh_-uphL9IYR5U0zYcPH05Qwg2YYP4LiIV8inYQEeCjWXIAc3L3cqHAawLSDfcGs3ZnClZrJQ9lmMZgUzB6pGoKohOi_QVA_uN_86PSeA04rXwHFRmU5B6UEhT81kDo5VTnPAbK1eUtn13UQlqie5KMPQ7uBV3O7iASqVDzxIj4ov1YxHMIvIVSUPCw",
      "e":"AQAB",
      "kid": "k1"
    },
    {
      "alg":"RS256",
      "kty":"RSA",
      "n":"uvtFgDnIcdB_jqSLICnsz7FXU_uiFSdJGVpGc5Dy-xm8wZwgiy6lJdL9_TtYjnmJefkPVyYdazabvGvOcns73rshkt0g6Ackqa72yiUEsv1kzCvBObPYNXgr1dNda8_F_ZiO3V9BtcTgQs9Y6rdOWJq7zNpees8pfuhEamk3sQp8AmKImFNfuZceNeglMHLLt0NcmSQp4VmhDCladFa1EdLirtFM9BtEIOlX20SRcN1LjeRsos8JywpQRxe6M3bnGFXcDQHqrsvwkkzu-vBtnPFa2e-jkBSDWkf6ZwvdJnEEUiJkHYTgJuXD1sbGeUkQL1Jb5NaQHhQ1mt3xn1z0tw",
      "e":"AQAB",
      "kid": "k2"
    }
  ]
}
`))
	log.Info().Msgf("served %v jwks request with X-Request-Id %s code %d", r.URL.Path, getXRequestId(r), 200)
}

func jwksrotate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)

	now := time.Now().Unix()
	if now%2 == 0 {
		w.Write([]byte(`{
  "keys": [
    {
      "alg":"RS256",
      "kty":"RSA",
      "n":"w6hKIwXAuI5VqmQjvQmgZdNbV80GMC3UkPmm-OQDjzOjeLRA6yLPYLZHaGhONx37DWMA-a3D_Zg_-oueYuZlrhusbTDC-bt1JSctAJV3ollQaalmJQHhLfyL54Y6Cgt3H_68u4Q3kLrFOmdFJwRswHR-1m-Oh_-uphL9IYR5U0zYcPH05Qwg2YYP4LiIV8inYQEeCjWXIAc3L3cqHAawLSDfcGs3ZnClZrJQ9lmMZgUzB6pGoKohOi_QVA_uN_86PSeA04rXwHFRmU5B6UEhT81kDo5VTnPAbK1eUtn13UQlqie5KMPQ7uBV3O7iASqVDzxIj4ov1YxHMIvIVSUPCw",
      "e":"AQAB",
      "kid": "k1"
    }
  ]
}` + "\n"))
	} else {
		w.Write([]byte(`{
  "keys": [
    {
      "alg":"RS256",
      "kty":"RSA",
      "n":"uvtFgDnIcdB_jqSLICnsz7FXU_uiFSdJGVpGc5Dy-xm8wZwgiy6lJdL9_TtYjnmJefkPVyYdazabvGvOcns73rshkt0g6Ackqa72yiUEsv1kzCvBObPYNXgr1dNda8_F_ZiO3V9BtcTgQs9Y6rdOWJq7zNpees8pfuhEamk3sQp8AmKImFNfuZceNeglMHLLt0NcmSQp4VmhDCladFa1EdLirtFM9BtEIOlX20SRcN1LjeRsos8JywpQRxe6M3bnGFXcDQHqrsvwkkzu-vBtnPFa2e-jkBSDWkf6ZwvdJnEEUiJkHYTgJuXD1sbGeUkQL1Jb5NaQHhQ1mt3xn1z0tw",
      "e":"AQAB",
      "kid": "k2"
    }
  ]
}` + "\n"))
	}

	log.Info().Msgf("served %v rotating jwks request with X-Request-Id %s code %d", r.URL.Path, getXRequestId(r), 200)
}

func jwksbadrotate(w http.ResponseWriter, r *http.Request) {
	k1 := `{
  "keys": [
    {
      "alg":"RS256",
      "kty":"RSA",
      "n":"w6hKIwXAuI5VqmQjvQmgZdNbV80GMC3UkPmm-OQDjzOjeLRA6yLPYLZHaGhONx37DWMA-a3D_Zg_-oueYuZlrhusbTDC-bt1JSctAJV3ollQaalmJQHhLfyL54Y6Cgt3H_68u4Q3kLrFOmdFJwRswHR-1m-Oh_-uphL9IYR5U0zYcPH05Qwg2YYP4LiIV8inYQEeCjWXIAc3L3cqHAawLSDfcGs3ZnClZrJQ9lmMZgUzB6pGoKohOi_QVA_uN_86PSeA04rXwHFRmU5B6UEhT81kDo5VTnPAbK1eUtn13UQlqie5KMPQ7uBV3O7iASqVDzxIj4ov1YxHMIvIVSUPCw",
      "e":"AQAB",
      "kid": "k1"
    }
  ]
}` + "\n"
	k2 := `{
  "keys": [
    {
      "alg":"RS256",
      "kty":"RSA",
      "n":"uvtFgDnIcdB_jqSLICnsz7FXU_uiFSdJGVpGc5Dy-xm8wZwgiy6lJdL9_TtYjnmJefkPVyYdazabvGvOcns73rshkt0g6Ackqa72yiUEsv1kzCvBObPYNXgr1dNda8_F_ZiO3V9BtcTgQs9Y6rdOWJq7zNpees8pfuhEamk3sQp8AmKImFNfuZceNeglMHLLt0NcmSQp4VmhDCladFa1EdLirtFM9BtEIOlX20SRcN1LjeRsos8JywpQRxe6M3bnGFXcDQHqrsvwkkzu-vBtnPFa2e-jkBSDWkf6ZwvdJnEEUiJkHYTgJuXD1sbGeUkQL1Jb5NaQHhQ1mt3xn1z0tw",
      "e":"AQAB",
      "kid": "k2"
    }
  ]
}` + "\n"
	k2b := `{
  "keys": [
    {
      "alg":"HS256",
      "kty":"RSA",
      "n":"xxuvtFgDnIcdB_jqSLICnsz7FXU_uiFSdJGVpGc5Dy-xm8wZwgiy6lJdL9_TtYjnmJefkPVyYdazabvGvOcns73rshkt0g6Ackqa72yiUEsv1kzCvBObPYNXgr1dNda8_F_ZiO3V9BtcTgQs9Y6rdOWJq7zNpees8pfuhEamk3sQp8AmKImFNfuZceNeglMHLLt0NcmSQp4VmhDCladFa1EdLirtFM9BtEIOlX20SRcN1LjeRsos8JywpQRxe6M3bnGFXcDQHqrsvwkkzu-vBtnPFa2e-jkBSDWkf6ZwvdJnEEUiJkHYTgJuXD1sbGeUkQL1Jb5NaQHhQ1mt3xn1z0tw",
      "e":"AQAB",
      "kid": "k2b"
    }
  ]
}` + "\n"
	k3k4 := `{
  "keys": [
    {
      "alg":"RS256",
      "kty":"RSA",
      "n":"qr_urqMJcKsSPs8VsIbzfM7WMFffvPnno3vUpvxR-ONiS83y46Wz-cYVmrA5so9MvCCmMaRamhTyTrT75cPV5ec2v1wu5UXQKPsyGaK8F_dYcuFd4G-fnZFRGT_EBHhjXU_d2eEdgxjtcHRWMJZoE4tR5ojQFOn6kbCT8HeoZriBnCRTunL7eBjjPQk7_sxjomJanJqh7JkprakA3SnSzUacetjXeFrTw5T4XCvzSAY8T75CUlho0019itG4BPw_zkqlEZZKC0J-vz6Q1Au0Di6kyoRnZJCDHzPa7gStZ1-MqIZNyYr2-eHLECYU9L04pIBtKHufaNjF4jR3XXz_nQ",
      "e":"AQAB",
      "kid": "k3"
    },
	{
      "alg":"RS256",
      "kty":"RSA",
      "n":"3-qU4mZcorx_oZWwQ_N-_tJgxk2Fz8EStfCM5xuumdhp26cF9yARrvd4mPcYxLN-d0R98gYg8Z_QwvqIzicp22ziGWLo2URH8lSxUcsEco2qCCyetqh1BXYId8fQiUOTq3yxdQ_LpEetrAuTfoiKRP_B3njV8VLmtNEc11WJx5bMg31QV4Y132aiaTlLiwgLqHLoDI2JLGE9NP4kH_MElPYKRTCTXh6Jj7aBp7PCTYahCt07qzOGVZby7L87gDOkJKPUu7GM4qPFcB0T7OL_KlimSIJclAfxk1UxrOkfdYw09K6BhR8LIjrHEm1JtI5pLyEzGlhjkhBp_wBta_bChw",
      "e":"AQAB",
      "kid": "k4"
    }
  ]
}` + "\n"
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)

	if len(r.URL.Query()["rc"]) > 0 {
		c, _ := strconv.Atoi(r.URL.Query()["rc"][0])
		if c == 0 {
			rc = 0
		}
	}
	rc++

	if rc == 1 {
		w.Write([]byte(k1))
	} else if rc == 2 {
		w.Write([]byte(k2))
	} else if rc == 3 {
		w.Write([]byte(k2b))
	} else {
		w.Write([]byte(k3k4))
	}

	log.Info().Msgf("served %v rotating jwks request count %d with X-Request-Id %s code %d", r.URL.Path, rc, getXRequestId(r), 200)
}

func jwksbad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{
  "keys": [
    {
      "kty":"RSA",
      "n":"uFrwC7xqek3lA7TkRMBr7koamTCE5DF0UxVPd0FbmloGTkkLLXW3R6fOxubi8O2PXk_tN-TfJZiOYswUE_-ngR7gEXLebosLtVdmbGraTGwtoGmpSe3FRr9ZmQu74pZsAzwqZVMqz6CINc7uvxTIDjd98ORUrnuxqgHE9Yz_uo2qvnaOgWIXKhkDkMqA8O0Fk_kaCfeeZQMN70OnCwIS-LPFE8uYGIdbaEIkjZfMxm_iNRENOV849vwOiOuWruCyp-YMqTVtcW49Q1mcZfyGT7B5GHWe7MtxqQNhf1m2Nvo1m_LvaLap_EM3684xOa6RexB1XdB8oegpMRygPx7orw",
      "e":"AQAB",
      "kid": "k1"
    },
    {
      "kty":"RSA",
      "n":"tXhyIjACJ9I_1RLe6ewuBIzZ1275BUssbeUdE87qSNpkJHsn6lNKPUQVix_Hk8MDME6Et1zmyK7a2XoTovMELgaHFSpH3i-Eqdl1jG9c0_vkHlwC6Ba-MLxvSCn6HVrcSMMGpOdVHUU4cuqDRpVO4owby8e1ZSS1hdhaqs5t464BID7e907oe7hE8deqD9MXmGEimcXXEJTF84wH2xcBqUO35dcc5SBJfPAibZ6U2AaNIEZJouUYMJOqwVttTBvKYwhuEwcxsPrYfkufbmGb9dnTfKMJamujAwFf-YUwifYfpY763cQ4Ex7eHWVp4LlBB9zYYBBGp2ueLuhJSMWhk0yP4KBk8ZDcIgLZKsTzYDdnvbecii7qAxRYMaSEkdjSj2JTmV_GtDBLmkejVNqo9s_BvgEIDiPipTWesPKsaNigyhs6p6POJvOHkAAc3-88cfShLuDpobWmNEO6eOAGGvACbWs-EOepMrvWuL53QWgJzJaKsxgGejQ1jVCIRZeaVsWiPrJFSUk87lWwxGpRcSdvOATlGgjz28jL_CqtuAySGTb4S0LsBFgdpykrGChjbajxeMMjnV3khI4c_KXlSmOsxHfJ5vzfbicw1Inn_4RoVxw72p4t1NN3va1W6jZt_FZ5R8xgV5T5zgeAEkSmHJa_PXCQoBYwK7cuMJhjRaM",
      "e":"AQAB",
      "kid": "k2"
    }
  ]
}
`))
	log.Info().Msgf("served %v bad jwks request with X-Request-Id %s code %d", r.URL.Path, getXRequestId(r), 200)
}

func jwksmix(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{
  "keys": [
    {
	  "alg":"RS256",
      "kty":"RSA",
      "n":"uFrwC7xqek3lA7TkRMBr7koamTCE5DF0UxVPd0FbmloGTkkLLXW3R6fOxubi8O2PXk_tN-TfJZiOYswUE_-ngR7gEXLebosLtVdmbGraTGwtoGmpSe3FRr9ZmQu74pZsAzwqZVMqz6CINc7uvxTIDjd98ORUrnuxqgHE9Yz_uo2qvnaOgWIXKhkDkMqA8O0Fk_kaCfeeZQMN70OnCwIS-LPFE8uYGIdbaEIkjZfMxm_iNRENOV849vwOiOuWruCyp-YMqTVtcW49Q1mcZfyGT7B5GHWe7MtxqQNhf1m2Nvo1m_LvaLap_EM3684xOa6RexB1XdB8oegpMRygPx7orw",
      "e":"AQAB",
      "kid": "k1"
    },
    {
      "alg": "RS384",
	  "kty":"RSA",
	  "n":"tXhyIjACJ9I_1RLe6ewuBIzZ1275BUssbeUdE87qSNpkJHsn6lNKPUQVix_Hk8MDME6Et1zmyK7a2XoTovMELgaHFSpH3i-Eqdl1jG9c0_vkHlwC6Ba-MLxvSCn6HVrcSMMGpOdVHUU4cuqDRpVO4owby8e1ZSS1hdhaqs5t464BID7e907oe7hE8deqD9MXmGEimcXXEJTF84wH2xcBqUO35dcc5SBJfPAibZ6U2AaNIEZJouUYMJOqwVttTBvKYwhuEwcxsPrYfkufbmGb9dnTfKMJamujAwFf-YUwifYfpY763cQ4Ex7eHWVp4LlBB9zYYBBGp2ueLuhJSMWhk0yP4KBk8ZDcIgLZKsTzYDdnvbecii7qAxRYMaSEkdjSj2JTmV_GtDBLmkejVNqo9s_BvgEIDiPipTWesPKsaNigyhs6p6POJvOHkAAc3-88cfShLuDpobWmNEO6eOAGGvACbWs-EOepMrvWuL53QWgJzJaKsxgGejQ1jVCIRZeaVsWiPrJFSUk87lWwxGpRcSdvOATlGgjz28jL_CqtuAySGTb4S0LsBFgdpykrGChjbajxeMMjnV3khI4c_KXlSmOsxHfJ5vzfbicw1Inn_4RoVxw72p4t1NN3va1W6jZt_FZ5R8xgV5T5zgeAEkSmHJa_PXCQoBYwK7cuMJhjRaM",
      "e":"AQAB",
      "kid": "k2"
    }
  ]
}
`))
	log.Info().Msgf("served %v mixed jwks request with X-Request-Id %s code %d", r.URL.Path, getXRequestId(r), 200)
}

func jwkses256(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "identity")
	w.WriteHeader(200)
	w.Write([]byte(`{
  "keys": [
    {
    	"alg": "ES256",
    	"created_at": 1560466143,
    	"crv": "P-256",
    	"expired_at": null,
    	"kid": "6c5516e1-92dc-479e-a8ff-5a51992e0001",
    	"kty": "EC",
    	"use": "sig",
    	"x": "35lvC8uz2QrWpQJ3TUH8t9o9DURMp7ydU518RKDl20k",
    	"y": "I8BuXB2bvxelzJAd7OKhd-ZwjCst05Fx47Mb_0ugros"
	}
  ]
}
`))
	log.Info().Msgf("served %v jwkses256 request with X-Request-Id %s code %d", r.URL.Path, getXRequestId(r), 200)
}

func websocket(w http.ResponseWriter, r *http.Request) {
	var n int
	if len(r.URL.Query()["n"]) > 0 {
		ns := r.URL.Query()["n"][0]
		var err error
		n, err = strconv.Atoi(ns)
		if err != nil {
			n = 1
		}
		log.Info().Msgf("websocket echo for connection with remote addr %s set to %d", r.RemoteAddr, n)
	} else {
		n = 1
	}

	var closeProtocol bool
	var closeSocket bool
	if len(r.URL.Query()["c"]) > 0 {
		closeProtocol = true
		closeSocket = true
	}
	if len(r.URL.Query()["c1"]) > 0 {
		closeProtocol = true
	}
	if len(r.URL.Query()["c2"]) > 0 {
		closeSocket = true
	}

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Error().Msgf("error during downstream connection upgrade to websocket with remote addr %s, cause: %s", r.RemoteAddr, err)
	} else {
		log.Info().Msgf("downstream connection with remote addr %s upgraded to websocket", r.RemoteAddr)
	}

	//this will properly close the websocket connection.
	closer := func() {
		if closeProtocol {
			if err := ws.WriteFrame(conn, ws.NewCloseFrame(ws.NewCloseFrameBody(ws.StatusNormalClosure, "close requested"))); err != nil {
				log.Info().
					Msgf("ws protocol connection with remote addr %s mse6 already closed", r.RemoteAddr)
			} else {
				log.Info().Msgf("ws protocol connection with remote addr %s close frame sent by mse6 now", r.RemoteAddr)
			}
		}
		if closeSocket {
			err2 := conn.Close()
			if err2 != nil {
				log.Info().
					Msgf("ws TCP connection with remote addr %s mse6 already closed", r.RemoteAddr)
			} else {
				log.Info().Msgf("ws TCP connection with remote addr %s closed by mse6 now", r.RemoteAddr)
			}
		}

	}
	defer closer()

wsloop:
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err == nil {
			log.Info().Msgf("success reading websocket msg from downstream: %s, opcode: %v", msg, op)
		} else if err.Error() == "EOF" {
			log.Info().Msg("downstream hung up, EOF")
			break wsloop
		} else if ce, cet := err.(wsutil.ClosedError); cet {
			log.Info().Msgf("success. downstream requested protocol close: %s", ce)
			closeProtocol = true
			closeSocket = true
			break wsloop
		} else {
			log.Warn().Msgf("error reading websocket msg from downstream, cause: %s", err)
			break wsloop
		}
		for i := 0; i < n; i++ {
			echo := fmt.Sprintf("%s", msg)
			err = wsutil.WriteServerMessage(conn, op, []byte(echo))
			if err == nil {
				log.Info().Msgf("success writing websocket echo %d to downstream: %s, opcode: %v", i+1, echo, op)
			} else {
				log.Warn().Msgf("error writing websocket echo to downstream, cause: %s", err)
				break wsloop
			}
		}
		if closeProtocol || closeSocket {
			break wsloop
		}
	}
}

func formget(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		for k, v := range r.URL.Query() {
			log.Info().Msgf("received GET with URL key/value: %s/%s for X-RequestID: %s", k, v[0], getXRequestId(r))
		}
		raw, _ := httputil.DumpRequest(r, true)
		log.Info().Msgf("raw GET: %v for X-RequestID: %s", string(raw), getXRequestId(r))

		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("received GET: %s", raw)))
	} else {
		send405(w, r)
	}
}

func formpost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			raw, _ := httputil.DumpRequest(r, true)
			err := r.ParseForm()
			if err != nil {
				log.Warn().Msgf("unable to parse form for X-RequestID: %s", getXRequestId(r))
			}
			for k, v := range r.Form {
				log.Info().Msgf("received application/x-www-form-urlencoded POST with form key/value: %s/%s for X-RequestID: %s", k, v[0], getXRequestId(r))
			}

			log.Info().Msgf("raw application/x-www-form-urlencoded POST: %v for X-RequestID: %s", string(raw), getXRequestId(r))

			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf("received post encoded as application/x-www-form-urlencoded: %s", raw)))
		} else if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			raw, _ := httputil.DumpRequest(r, true)
			err := r.ParseMultipartForm(0)
			if err != nil {
				log.Warn().Msgf("unable to parse multipart form for X-RequestID: %s", getXRequestId(r))
			}

			for k, v := range r.Form {
				log.Info().Msgf("received POST with multipart form key/value: %s/%s for X-RequestID: %s", k, v[0], getXRequestId(r))
			}

			log.Info().Msgf("raw multipart/form-data POST: %v for X-RequestID: %s", string(raw), getXRequestId(r))

			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf("received post encoded as multipart/form-data: %s", raw)))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("received post with unsupported encoding as : %s", r.Header.Get("Content-Type"))))
		}

	} else {
		send405(w, r)
	}
}
