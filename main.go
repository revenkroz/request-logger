package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Log struct {
	Method     string        `json:"method,omitempty"`
	URL        url.URL       `json:"url"`
	FullURL    string        `json:"full_url,omitempty"`
	Status     string        `json:"status,omitempty"`
	StatusCode int           `json:"status_code,omitempty"`
	Req        string        `json:"req,omitempty"`
	Res        string        `json:"res,omitempty"`
	Elapsed    time.Duration `json:"elapsed,omitempty"`
	StartedAt  time.Time     `json:"started_at"`
	DoneAt     time.Time     `json:"done_at"`
}

func main() {
	proxyURLString := flag.String("url", "", "Set the proxy URL")
	proxyAddr := flag.String("addr", "0.0.0.0:21000", "Address of the proxy server to listen on")
	frontendAddr := flag.String("faddr", "0.0.0.0:21001", "Address of the frontend server to listen on")
	maxLogs := flag.Int("maxlogs", 20, "Set the maximum number of logs to keep in memory")
	stdout := flag.Bool("stdout", false, "Set to true to log to stdout")
	flag.Parse()

	if *proxyURLString == "" {
		fmt.Println("You must provide a URL. Use -url")
		os.Exit(1)
	}

	URL, err := url.Parse(*proxyURLString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logChan := make(chan Log, 10)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		fmt.Println("Frontend server listening on", *frontendAddr)
		err := http.ListenAndServe(*frontendAddr, NewFrontendServer(logChan, *maxLogs))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()

		fmt.Println("Proxy server listening on", *proxyAddr)
		err := http.ListenAndServe(*proxyAddr, NewProxyHandler(URL, logChan, *stdout))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	wg.Wait()
}
