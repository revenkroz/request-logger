package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
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

const defaultProxyListenAddr = "0.0.0.0:21001"

var (
	// Server configuration
	fromToAddresses arrayFlags = getFromEnvStringSlice("PROXY_ADDR", []string{})
	frontendAddress            = getFromEnvString("FRONTEND_ADDR", "0.0.0.0:21000")

	// Other
	maxLogs       = getFromEnvInt("MAX_LOGS", 30)
	printToStdout = getFromEnvBool("USE_STDOUT", false)
)

func init() {
	flag.Var(&fromToAddresses, "proxy", "Multiple values, proxy listen address and target address (if empty, env:PROXY_ADDR will be used).")
	flag.StringVar(&frontendAddress, "front", frontendAddress, "Frontend listen address (if empty, env:FRONTEND_ADDR will be used).")
	flag.IntVar(&maxLogs, "maxlogs", maxLogs, "Maximum number of logs to keep in memory (if empty, env:MAX_LOGS will be used).")
	flag.BoolVar(&printToStdout, "stdout", printToStdout, "Print logs to stdout (if empty, env:USE_STDOUT will be used).")
	flag.Parse()
}

func main() {
	logChan := make(chan Log, 10)

	wg := sync.WaitGroup{}
	wg.Add(1 + len(fromToAddresses))

	go func() {
		defer wg.Done()

		fmt.Printf("Frontend server listening on http://%s\n", frontendAddress)
		err := http.ListenAndServe(frontendAddress, NewFrontendServer(logChan, maxLogs))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	addresses := prepareFromToAddresses(fromToAddresses)
	for _, fromToAddr := range addresses {
		from, to := fromToAddr[0], fromToAddr[1]

		go func() {
			defer wg.Done()

			fmt.Printf("Proxy server listening on http://%s, forwarding to %s\n", from, to)
			err := http.ListenAndServe(from, NewProxyHandler(parseUrl(to), logChan, printToStdout))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()
	}

	wg.Wait()
}

func prepareFromToAddresses(fromToAddresses []string) [][]string {
	fromTos := make([][]string, 0)
	usedAddrs := make([]string, 0)

	for _, fromToAddr := range fromToAddresses {
		split := strings.Split(fromToAddr, "::")
		if len(split) > 2 || len(split) == 0 {
			fmt.Println("Invalid proxy address:", fromToAddr)
			os.Exit(1)
		}

		if len(split) == 1 {
			split = []string{defaultProxyListenAddr, split[0]}
		}

		if slices.Contains(usedAddrs, split[0]) {
			fmt.Println("Duplicate proxy listen address:", split[0])
			os.Exit(1)
		}

		if split[1] == "" {
			fmt.Println("Target address is empty")
			os.Exit(1)
		}

		fromTos = append(fromTos, split)
		usedAddrs = append(usedAddrs, split[0])
	}

	return fromTos
}
