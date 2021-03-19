package mse6

import (
	"crypto/tls"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var waitDuration time.Duration
var Version = "v0.2.16"
var Port int
var Prefix string
var rc = 0

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

func die(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("served %v request with X-Request-Id %s, process exiting with -1", r.URL.Path, getXRequestId(r))
	os.Exit(-1)
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

func trace(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if r.Method == "TRACE" {
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Type", "message/http")
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(200)
		log.Info().Msgf("served %v delete trace with X-Request-Id %s,%s reading %d bytes from inbound", r.URL.Path, getXRequestId(r), expectContinue(r), len(body))
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

func gzipf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "mse6 "+Version)
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	w.Write(gzipenc([]byte(`{"mse6":"Hello from the gzip endpoint"}`)))

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
		w.Header().Set("Content-Length", "0")
		w.Header().Set("Server", "mse6 "+Version)
		w.Header().Set("Content-Encoding", "identity")
		w.WriteHeader(code)

		log.Info().Msgf("served %v OPTIONS request with X-Request-Id %s code %d", r.URL.Path, getXRequestId(r), code)
	} else {
		send405(w, r)
	}
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
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Error().Msgf("error during websocket upgrade request: %s", err)
	}
	//go func() {
	defer conn.Close()

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			// handle error
		}
		err = wsutil.WriteServerMessage(conn, op, []byte(fmt.Sprintf("mse6 websocket echo: %s", msg)))
		if err != nil {
			// handle error
		}
	}
	//}()

	//no http response
}

func Bootstrap(port int, waitSeconds float64, prefix string, tlsMode bool) {
	waitDuration = time.Second * time.Duration(waitSeconds)
	log.Info().Msgf("wait duration for slow requests seconds %v", waitDuration.Seconds())

	Port = port
	Prefix = prefix
	mode := "http"
	if tlsMode {
		mode = "tls"
	}
	log.Info().Msgf("mse6 %s starting %s server on port %d with prefix '%s'", Version, mode, Port, Prefix)

	http.HandleFunc(prefix+"badcontentlength", badcontentlength)
	http.HandleFunc(prefix+"badgzip", badgzipf)
	http.HandleFunc(prefix+"chunked", chunked)
	http.HandleFunc(prefix+"delete", delete)
	http.HandleFunc(prefix+"die", die)
	http.HandleFunc(prefix+"echoheader", echoheader)
	http.HandleFunc(prefix+"jwks", jwks)
	http.HandleFunc(prefix+"jwkses256", jwkses256)
	http.HandleFunc(prefix+"jwksbad", jwksbad)
	http.HandleFunc(prefix+"jwksmix", jwksmix)
	http.HandleFunc(prefix+"jwksrotate", jwksrotate)
	http.HandleFunc(prefix+"jwksbadrotate", jwksbadrotate)
	http.HandleFunc(prefix+"get", get)
	http.HandleFunc(prefix+"gzip", gzipf)
	http.HandleFunc(prefix+"getorhead", getorhead)
	http.HandleFunc(prefix+"redirected", redirected)
	http.HandleFunc(prefix+"options", options)
	http.HandleFunc(prefix+"patch", patch)
	http.HandleFunc(prefix+"post", post)
	http.HandleFunc(prefix+"put", put)
	http.HandleFunc(prefix+"send", send)
	http.HandleFunc(prefix+"slowbody", slowbody)
	http.HandleFunc(prefix+"slowheader", slowheader)
	http.HandleFunc(prefix+"trace", trace)
	http.HandleFunc(prefix+"websocket", websocket)

	//catchall
	http.HandleFunc("/", send404)

	var err error
	if tlsMode == false {
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	} else {
		chain, _ := getCert()
		server := &http.Server{
			Addr: fmt.Sprintf(":%d", port),
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{chain},
			},
		}
		err = server.ListenAndServeTLS("", "")
	}

	if err != nil {
		panic(err.Error())
	}
}

func getCert() (tls.Certificate, error) {
	certPem := "-----BEGIN CERTIFICATE-----\nMIIFtDCCA5ygAwIBAgICEAIwDQYJKoZIhvcNAQELBQAwdjELMAkGA1UEBhMCQVUx\nDDAKBgNVBAgMA05TVzENMAsGA1UECgwEbXljYTEaMBgGA1UECwwRbXljYSBpbnRl\ncm1lZGlhdGUxLjAsBgNVBAMMJW15IGNlcnRpZmljYXRlIGF1dGhvcml0eSBpbnRl\ncm1lZGlhdGUwIBcNMjAwOTE5MDQxNjAwWhgPMjEwMjExMDkwNDE2MDBaMGIxCzAJ\nBgNVBAYTAkFVMQwwCgYDVQQIDANOU1cxDzANBgNVBAcMBlN5ZG5leTEOMAwGA1UE\nCgwFY2VydDIxETAPBgNVBAsMCGNlcnQyIG91MREwDwYDVQQDDAhjZXJ0MiBjbjCC\nASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN3FFHDc3fWIyxukMDRriEbY\ntVA41EeiQiwf7RLdDxh+N2VAazUbbxUJ06nKAslX2+6ZmJrMlS+ionX1BvPhPy3s\nnuZI1movXcvH6ZV5yUGZyJDocjOTHHqNwPSDOAQX87tLjQbCa8Rw//B488GoPbaZ\nlWYDvZQ0Mw5rasiu0B+OI6PL8+Vnc2jXdPlc3tiNoIVXRZ14TNei7bUDA3O1y593\nift2tQ/TZxlY7fylZWhTV4sUm/9yk/zob+dyzro795Jy8vThlePAN//tZGLWFzG7\na8o9Mx36BPncSZ0v+EfEvP24ZffIDFRtysBewu2+33IVpISlbaHgj6nsuv8GFM0C\nAwEAAaOCAVwwggFYMAkGA1UdEwQCMAAwEQYJYIZIAYb4QgEBBAQDAgZAMDMGCWCG\nSAGG+EIBDQQmFiRPcGVuU1NMIEdlbmVyYXRlZCBTZXJ2ZXIgQ2VydGlmaWNhdGUw\nHQYDVR0OBBYEFEEFnOmrROOjNNQrLRoXPXsJLPkeMIGdBgNVHSMEgZUwgZKAFOlV\np+B1WShwNQuSAVuQAvCe/teZoXakdDByMQswCQYDVQQGEwJBVTEMMAoGA1UECAwD\nTlNXMQ8wDQYDVQQHDAZTeWRuZXkxDTALBgNVBAoMBG15Y2ExEjAQBgNVBAsMCW15\nY2Egcm9vdDEhMB8GA1UEAwwYbXkgY2VydGlmaWNhdGUgYXV0aG9yaXR5ggIQADAO\nBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwHwYDVR0RBBgwFoII\nY2VydC5jb22CCiouY2VydC5jb20wDQYJKoZIhvcNAQELBQADggIBADSL8AegMDhJ\nUgRfP6CQeAcLgbHAb9cS7vo0ju9E38pSVDBKA1VachXgwf6630XJ4/YrHzCNgbGO\neX3GcwwcD8oWopnPX4bnGdwZaQ52qd4yUNgErNFpsZU02+ohgJew1Wx+caGNQ5F3\nMqsIy8X86a5FOFCGa0CUx4Iv4JieD6kKFWzJwvXwbWS6tFUxUOlpxYZRpZj4ZPb/\nyz65PBHeH9K+A0q+upwvVdK3Gp0qbcl7ZEE3rVR1GB5VSGnyG4YG0Y59Ys0JlsgR\n2jY0zdC2DTAGQdPL6u1HsNgCDz2nzUDaYGOMb1NVRTsRZ/25irkAsOJFHP/CkuSy\nW/xogRCbX5WhHwxIzucpj+tnB7Hi9TBJLcsl7MNHuhUz5vtkl3d3dUEipIESwKmC\nn/avv2+6/8tm3UV2ji2N10246nPHZX8IddAAMwdfNriwPsz5XfXaF7czgaWYvBsu\nxkd5b2mbGH3BVJiEwRDeQRo8WGBfs0vqAF3abqjIrTiJikZpcI7GbqzhYzLcJnlX\nbOc2Xo8PXj9mE7dQ0Tfkd2wAovQ2xnuKBQgu14adFJjLFhSk5xRuu3274Kn3CmrK\nlo4FaSoSIw1vHs3J9VH6z/VRpu6dwC1U3XMSSlNSLcVO7UI/FJdacypAz6NDgpfP\nlv2z+ne8UwT6KGgPPTvCs0T5kbIXjOPY\n-----END CERTIFICATE-----\n"
	keyPem := "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA3cUUcNzd9YjLG6QwNGuIRti1UDjUR6JCLB/tEt0PGH43ZUBr\nNRtvFQnTqcoCyVfb7pmYmsyVL6KidfUG8+E/Leye5kjWai9dy8fplXnJQZnIkOhy\nM5Mceo3A9IM4BBfzu0uNBsJrxHD/8Hjzwag9tpmVZgO9lDQzDmtqyK7QH44jo8vz\n5WdzaNd0+Vze2I2ghVdFnXhM16LttQMDc7XLn3eJ+3a1D9NnGVjt/KVlaFNXixSb\n/3KT/Ohv53LOujv3knLy9OGV48A3/+1kYtYXMbtryj0zHfoE+dxJnS/4R8S8/bhl\n98gMVG3KwF7C7b7fchWkhKVtoeCPqey6/wYUzQIDAQABAoIBACPJuhK8kdUdzikX\nxe+vqr5EGn5nrVoiBSu5uzhgFB+PvsDINITNeI+cllvADdMQKp3Gi6nveePGCxGe\nCREyOE/g74OaHX/lRO2txTQqAyBjAMrhuAw6oU3lsk3DHzcJ5ntDJe8BUQLSeXsF\nCdEmpU7iWgmscNuJ0PNywjjAfTWaHXNgXbcragVT/El53/fAnO36aDmd5SP4BiiQ\n984Hig9Z+B9AuqYzKour8o96+IC8eD6EzSVbyvE7WnUZiVV2Opf4mJ8qUEw1NQlg\nGScrcF5RSCJTmB1lt9/mLE1PFS2SZpt2u3iCyKPAqWLa3oAzWMqD9X45+UV2UFlV\nnrfkrsECgYEA/rEE64qKiR5dgjvZVJls6dVu2WYy+EXCSqY2mYFbzHP+rw/xs7oZ\nk39/c0QghZJXDzzxXFUgKa5oeKrkYefPBWFquUfZx/OltbWfjdk8L/z8kfpYJetB\nySELnZiq9mb0JcDPGT5TJVR/udTlCtz89VPeYVt7dOypsAF0uvSrrUcCgYEA3ujC\nvvlughdm7oqhIgaRsIZKQedXLQVb8B1X1HnrbDgnuvBXEKioxIZT6Aw73scl5IFU\n7VBA+tasm9MdwtM18wJ62XCKuN3EgAA0/XpiuageWxSMfwm4Gy2t6FnV5CM+3in/\nmEPDG4NiUqyhk8eDuuuPLWtnXpRN+HQKM5xHp0sCgYEA3dZb/bkXP6WGNxhgDRLx\nzZ6MxakBvkQsng62QfBtf+CMtfjCQxRWkKWd4k01soIreGdRp2Wx9PwnnOrkr+5T\n4FDgv2843rN245XF2qybgwTtDU0rmCOYklJJJsTCLIqyH2wYNtmVXE+ETN2FfnfL\nkPezG8Ot/cLhbh9miCzyl6MCgYEAgYU9oznLvEtcw75JYjvu62McQq7pOH+krCBg\nqFUvNfJrI3QDIurdJVPn7S0unIOawOtlLX80Qov6P5Cr+kg/ULRgLXf3IvO4+acl\nIyO5uaa1/LYz7Jz5HNGt+xQ39BeGsBA3M4IsHBB7UQ591CBZqoK07u85YPtLUtIa\nG2LzP4ECgYBHOPg3ndFMe5EBql/92nSH+RILE6ADUCa+oQUOKa5p/cdWMt6ClT0m\n6cMOJN8lMmtVzwRG/aLPhN2L/vCbtBFDBDIm8PM5gg0340uFv5Mo4p1Sf8iRZG4B\nmzl86a1/OBk4MrtJqoqKrR9yg5/BXlvwuXBJRHaLjGERxhzyhk/WaQ==\n-----END RSA PRIVATE KEY-----\n"
	interCAPem := "-----BEGIN CERTIFICATE-----\nMIIFzDCCA7SgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwcjELMAkGA1UEBhMCQVUx\nDDAKBgNVBAgMA05TVzEPMA0GA1UEBwwGU3lkbmV5MQ0wCwYDVQQKDARteWNhMRIw\nEAYDVQQLDAlteWNhIHJvb3QxITAfBgNVBAMMGG15IGNlcnRpZmljYXRlIGF1dGhv\ncml0eTAgFw0yMDA5MTcyMTE3MjdaGA8yMjM5MDkzMDIxMTcyN1owdjELMAkGA1UE\nBhMCQVUxDDAKBgNVBAgMA05TVzENMAsGA1UECgwEbXljYTEaMBgGA1UECwwRbXlj\nYSBpbnRlcm1lZGlhdGUxLjAsBgNVBAMMJW15IGNlcnRpZmljYXRlIGF1dGhvcml0\neSBpbnRlcm1lZGlhdGUwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCm\n1qRLG0twxwbxBdt/nDeUr0Ia8LtLDtvPjdVVUDTpCp4gnlkEEzHu8JXPPsIami2C\nt5vo35JW1AI223FD5eef54wZG2rXlJbzwlB+yyE5+/V/6WSKe42rePvZDCD+Ym/Q\nyYeqzObViXGnmIvta2aEYZzLeTJPzppvQws/bM+d5IhRa43JuJOVYmjPdp1cjaOm\ntmW3zQSj/00a3i/97SHoyqaJX+y2bPQIJ+yScdBSn9W+Ke3o7/WnuP0HO/ST1fZM\nyzorGbso6aGnTswFbOdWMDUpauE97SL1M6ztoaI4a0HHD8Z8dPhtAmXWbs/5hmQr\njZqBj5W4oUik9iIjUhC2l1aYUf934Om62JjMn9if/mIIA5UTorddj/wKtIsd0n4X\nq5nhJ+X4yVXi3YjqW8iegenaq6UGuvNsm6m/JRAf+5n3FuspHH4WrCgAaIrYg0ZY\nDDu5ro6zHxTcHF6j01CXlJTDEJlStoZ6N9cIKVT94pUPM+EZBq3DGlhBDKipZWk7\n+sEu7sZoQ51WoV4haMY+4Wd7ea8o4sE50eoW+DN2o9lIPHMyxY5uFD7CluUt/b37\ntCcOYAV86JWBN5htTPYAH3wXsDBU/KFSJPLRPF96cuHL6Dq++Gvlqw0rKDKQ/gKh\nDma8lZ9SjVTskqk3l5wzHyNjy7nYFSIRItGIhVbp0wIDAQABo2YwZDAdBgNVHQ4E\nFgQU6VWn4HVZKHA1C5IBW5AC8J7+15kwHwYDVR0jBBgwFoAUD+ANepMk0O9Poxxx\nMpCnxxVyHNMwEgYDVR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAYYwDQYJ\nKoZIhvcNAQELBQADggIBAAQKj6FLBiy23kqHB7iUrl3dSXjJEsPm03zApRhWhr3e\nuxGVYO/YM6RlcJlc7RiKrQAO7XMuOfGbV/TedKPYz+SAeoHCdAVmT21o9HqgwRJ2\nkbJulqIF7oRmmqFOUDIUNg+ZC68QvR9cfuhzcLsEdmfEhXvI5j6CvrhOUN3UHw8A\nO7b4kiymBVT88uXUC0i3bGeEI3h6Fz/RZLbShcvTz2BwcuqoWdInyKi+8mKNfc1O\n+HGBMjnPahNAiovaEuUGErloETdjhmSOkbPBG8h9KpkndCwclEhsBN1+skKiDzKa\nMk53cXXKjqPvPEG9dfQQu0NEnOeY3ZtyVpMqnbo+G0MtyzkozvAB5WjWlpaWZYV2\nnw/wnyCi57ruYI7UjUp+NvFDiIRlOysLC7K6xia+8m7mP8MaFJibQh0tA2UDmdXs\nwy/Z87c6KUCyDB8Hl//rLWbWg6JpHTcH+81yDkVeq2TvJkB6P8jThv51Pz1z4b6U\ndHWAMK5kLmHv+P6sw0JkE5fwszoFOaqSxABq02Pkt5+Hv2EvwxpJZvySkdp7s+Xn\nGUwXhduMscVL/Yd62ES5dYSQ+vbmZIEK3PIttcIyleif6DLFZijJnywf5etYxvrK\nY9wgX6D9PwShl32sf3nzHXh3npLdbio3XwJQUcO6c/lm49rKD7L9L5RM6FNShl8R\n-----END CERTIFICATE-----\n"
	rootCAPem := "-----BEGIN CERTIFICATE-----\nMIIFzDCCA7SgAwIBAgIJAKdYQFPloO6RMA0GCSqGSIb3DQEBCwUAMHIxCzAJBgNV\nBAYTAkFVMQwwCgYDVQQIDANOU1cxDzANBgNVBAcMBlN5ZG5leTENMAsGA1UECgwE\nbXljYTESMBAGA1UECwwJbXljYSByb290MSEwHwYDVQQDDBhteSBjZXJ0aWZpY2F0\nZSBhdXRob3JpdHkwIBcNMjAwOTE3MjA1MTM1WhgPMjI5NDA3MDMyMDUxMzVaMHIx\nCzAJBgNVBAYTAkFVMQwwCgYDVQQIDANOU1cxDzANBgNVBAcMBlN5ZG5leTENMAsG\nA1UECgwEbXljYTESMBAGA1UECwwJbXljYSByb290MSEwHwYDVQQDDBhteSBjZXJ0\naWZpY2F0ZSBhdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoIC\nAQDgvRdI24Rv9XBnirlB1LwS32MYyVM2mksTF52E0qrg1OKcMs1D2737BrgaUD6C\nB1I2lMAKR25Q3+x9fSutyww8KZ7yQkFcX2lhwsyYll0j1rvkjek0M1K4787ZFrXS\ncRihE6BSvP5886O+v7a30HxtKbI9oFHdbzgpLpTzvVAn53tokRgAJNtQZWpyJ5Qq\nIG7c96dG9zsXE5+tYT0E0p3ec1z/Ucdx6SKOFjCR8bVLX+Y97mxypOMaPEhGJ4D3\nBlxlCvwDo5sF46e/ntie3Fqghk3jRZTUXedB0IjN8iJCKODPMO1j1cESqVg21xGZ\nyZxIn/ra1iqx9VDCP8egfUOmmMF8flGV08qOGDLGEc/dpVe/yHvG3lmld3MBsW+3\nu6O2l7GIKdLHKibe3uGHhmuPbHq2vlc6IIlRtpsZtK3IXt+bpvlKdI3rxbl4MbT7\n8Z09IUpTsT5jDPEVRnX0zV78Gs4TyKqJKxJJaINx9n0AuXJ8b3jmth/Bb6OkoPgv\nsbFS2QER2Yp8whE1W2PMwtJ06u20YX0RSwuKD+CsnTVmtQwWLBXescCNRH372HwS\nLHO8dvyFWfekLaB2LfciJWYBd8thO5Y4O65FnKLGDvEUh6Ew2OOnhOpy4flWAng6\n39r5uuDQqmWrPFjDNR5HvQjQu1Bv0j81cFY4qZqSIskR9wIDAQABo2MwYTAdBgNV\nHQ4EFgQUD+ANepMk0O9PoxxxMpCnxxVyHNMwHwYDVR0jBBgwFoAUD+ANepMk0O9P\noxxxMpCnxxVyHNMwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAYYwDQYJ\nKoZIhvcNAQELBQADggIBAFdMzocjv6RojMXft1TnwYKb8H0ce6qcsBHZmd/M7IXf\nhyedRkcm7RuN7ayNjFA+44pwAr4jMMNklBQDpGD5yYt1jsltiYoYX5bwZdn2I/rr\nwNQ/FfNSp8rJWqtBhaEt0VI+snHuy0Gdx1eQGf4bJNzvsDLjjJuQ32VUjaCOzsd1\n7d8jR/yjR3Sq20oFEu3HqFSC9OCH2QORTqf6i2IkaUeJbkVTa8+uVceDDbRs3CwY\nVgk/4WcOzcrz0F2BJPpFQ4knrSuHgUbElPHPVuZcn3XZ0n1KBXZdNVCIyLVRowdr\nI+gNEgWE3670Osx55QWg7depP7hU30nQlC1cm2ej2MxM48ddbAL4Zqs8/W1gm+Xb\nDkTsfh81QZQaw6qFVGHJNRIyfMT68ekFB8AgqntulIFR2RJTr/3QJBMhGHKQkmcT\nsa0z0ZrmS/ieurRUjaCsud10Y5VbY5Y8ll5kPsuRWuyijftjcPFqHBzLSSdLacO9\nlVIGkTA3ARCGgym3v5+ZZJ4DeLOJRz9c9OCIASlCkNFFEm1aJ8oagynh2tYqe5TK\nCva1MX8QW5OjHbrm1xvQ8uZOSj55yuBQWKH47GF4QxiojzKikLv4Cpv2Tk5SR9qv\nq3C4t8B26KurNb4z99eo5XhW5XXvQdKZTQC9BqZDN7xhQlwm5lbRSuhZMBJJaQOS\n-----END CERTIFICATE-----\n"
	cB := []byte(certPem + interCAPem + rootCAPem)
	kB := []byte(keyPem)
	chain, err1 := tls.X509KeyPair(cB, kB)
	return chain, err1
}
