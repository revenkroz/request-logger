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

func NewFrontendServer(
	logChan chan Log,
) http.Handler {
	mux := http.NewServeMux()

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
				jsn, _ := json.Marshal(c)
				fmt.Fprintf(w, "event: log\ndata: %s\n\n", jsn)
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				return
			}
		}
	})

	return mux
}
