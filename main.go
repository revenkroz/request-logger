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

var (
	// Server configuration
	frontendAddress            = getFromEnvString("FRONTEND_ADDR", "0.0.0.0:21000")
	proxyAddresses  arrayFlags = getFromEnvStringSlice("PROXY_ADDR", []string{
		"0.0.0.0:21001",
	})

	// Proxy configuration
	targetUrlRaw = getFromEnvString("TARGET_URL", "")

	// Other
	maxLogs       = getFromEnvInt("MAX_LOGS", 30)
	printToStdout = getFromEnvBool("USE_STDOUT", false)

	// Internal variables
	targetUrl *url.URL
)

func init() {
	flag.Var(&proxyAddresses, "addr", "Frontend listen address (if empty, env:FRONTEND_ADDR will be used).")
	flag.StringVar(&frontendAddress, "faddr", frontendAddress, "Proxy listen address (if empty, env:LISTEN_ADDR will be used).")
	flag.StringVar(&targetUrlRaw, "target", targetUrlRaw, "Target URL (if empty, env:TARGET_URL will be used).")
	flag.IntVar(&maxLogs, "maxlogs", maxLogs, "Maximum number of logs to keep in memory (if empty, env:MAX_LOGS will be used).")
	flag.BoolVar(&printToStdout, "stdout", printToStdout, "Print logs to stdout (if empty, env:USE_STDOUT will be used).")
	flag.Parse()

	targetUrl = parseUrl(targetUrlRaw)
}

func main() {
	logChan := make(chan Log, 10)

	wg := sync.WaitGroup{}
	wg.Add(1 + len(proxyAddresses))

	go func() {
		defer wg.Done()

		fmt.Println("Frontend server listening on", frontendAddress)
		err := http.ListenAndServe(frontendAddress, NewFrontendServer(logChan, maxLogs))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	for _, proxyAddr := range proxyAddresses {
		go func() {
			defer wg.Done()

			fmt.Println("Proxy server listening on", proxyAddr)
			err := http.ListenAndServe(proxyAddr, NewProxyHandler(targetUrl, logChan, printToStdout))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()
	}

	wg.Wait()
}
