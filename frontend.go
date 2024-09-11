package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed all:frontend/dist/*
var frontend embed.FS

type LogStore struct {
	limit int
	Logs  []Log
}

func (l *LogStore) Add(log Log) {
	if len(l.Logs) >= l.limit {
		l.Logs = l.Logs[1:]
	}

	l.Logs = append(l.Logs, log)
}

func NewFrontendServer(
	logChan chan Log,
	logLimit int,
) http.Handler {
	mux := http.NewServeMux()

	logStore := LogStore{
		limit: logLimit,
	}

	indexHTML, err := frontend.ReadFile("frontend/dist/index.html")
	if err != nil {
		log.Println("You need to build the frontend first")
		panic(err)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write(indexHTML)
	})

	mux.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		fsSub, err := fs.Sub(frontend, "frontend/dist")
		if err != nil {
			log.Println("You need to build the frontend first")
			panic(err)
		}

		http.FileServer(http.FS(fsSub)).ServeHTTP(w, r)
	})

	mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			select {
			case c := <-logChan:
				logStore.Add(c)
				jsn, _ := json.Marshal(c)
				fmt.Fprintf(w, "event: log\ndata: %s\n\n", jsn)
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				return
			}
		}
	})

	mux.HandleFunc("/logs/all", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsn, _ := json.Marshal(logStore.Logs)
		w.Write(jsn)
	})

	return mux
}
