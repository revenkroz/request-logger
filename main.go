package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	proxyAddr    = getFromEnvString("PROXY_ADDR", "0.0.0.0:21000")
	frontendAddr = getFromEnvString("FRONTEND_ADDR", "0.0.0.0:21001")

	// Proxy configuration
	targetUrlRaw = getFromEnvString("TARGET_URL", "")

	// Other
	maxLogs       = getFromEnvInt("MAX_LOGS", 20)
	printToStdout = getFromEnvBool("USE_STDOUT", false)

	// Internal variables
	targetUrl *url.URL
)

func init() {
	flag.StringVar(&proxyAddr, "addr", proxyAddr, "Proxy listen address (if empty, env:LISTEN_ADDR will be used).")
	flag.StringVar(&frontendAddr, "faddr", frontendAddr, "Frontend listen address (if empty, env:FRONTEND_ADDR will be used).")
	flag.StringVar(&targetUrlRaw, "target", targetUrlRaw, "Target URL (if empty, env:TARGET_URL will be used).")
	flag.IntVar(&maxLogs, "maxlogs", maxLogs, "Maximum number of logs to keep in memory (if empty, env:MAX_LOGS will be used).")
	flag.BoolVar(&printToStdout, "stdout", printToStdout, "Print logs to stdout (if empty, env:USE_STDOUT will be used).")
	flag.Parse()

	targetUrl = parseUrl(targetUrlRaw)
}

func main() {
	logChan := make(chan Log, 10)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		fmt.Println("Frontend server listening on", frontendAddr)
		err := http.ListenAndServe(frontendAddr, NewFrontendServer(logChan, maxLogs))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()

		fmt.Println("Proxy server listening on", proxyAddr)
		err := http.ListenAndServe(proxyAddr, NewProxyHandler(targetUrl, logChan, printToStdout))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	wg.Wait()
}

func getFromEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getFromEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalln("Invalid value for " + key + ": " + value)
	}

	return i
}

func getFromEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalln("Invalid value for " + key + ": " + value)
	}

	return b
}

func parseUrl(rawUrl string) *url.URL {
	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Fatalln("Invalid URL: " + rawUrl)
	}

	return u
}
