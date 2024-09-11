package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/andybalholm/brotli"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"time"
)

func NewProxyHandler(
	URL *url.URL,
	logChan chan Log,
	stdout bool,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = URL.Host
		r.URL.Host = URL.Host
		r.URL.Scheme = URL.Scheme

		var reqBody bytes.Buffer
		r.Body = io.NopCloser(io.TeeReader(r.Body, &reqBody))

		res := httptest.NewRecorder()
		start := time.Now()
		httputil.NewSingleHostReverseProxy(URL).ServeHTTP(res, r)
		elapsed := time.Since(start).Truncate(time.Millisecond) / time.Millisecond

		for k, v := range res.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(res.Code)
		w.Write(res.Body.Bytes())

		r.Body = io.NopCloser(bytes.NewReader(reqBody.Bytes()))

		if res.Header().Get("Content-Encoding") == "gzip" {
			reader, _ := gzip.NewReader(res.Body)
			data, _ := io.ReadAll(reader)
			reader.Close()
			res.Body = bytes.NewBuffer(data)
		} else if res.Header().Get("Content-Encoding") == "br" {
			reader := brotli.NewReader(res.Body)
			data, _ := io.ReadAll(reader)
			res.Body = bytes.NewBuffer(data)
		}

		dumpReq, _ := httputil.DumpRequest(r, true)
		dumpRes, _ := httputil.DumpResponse(res.Result(), true)

		if stdout {
			fmt.Println("------------------------")
			fmt.Println(string(dumpReq) + "\n")
			fmt.Println(string(dumpRes) + "\n\n")
		}

		logChan <- Log{
			Method:    r.Method,
			URL:       *r.URL,
			FullURL:   r.URL.String(),
			Status:    res.Result().Status,
			Req:       string(dumpReq),
			Res:       string(dumpRes),
			Elapsed:   elapsed,
			StartedAt: start,
			DoneAt:    time.Now(),
		}
	})
}
