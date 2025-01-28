package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

const PORT = ":8080"

func Api() {

	// basic logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// start new servemux
	mux := http.NewServeMux()

	// define handlers for routing
	mux.HandleFunc("/create", genericHandler)
	mux.HandleFunc("/get", genericHandler)
	mux.HandleFunc("/update", genericHandler)
	mux.HandleFunc("/delete", genericHandler)

	// start server
	if err := http.ListenAndServe(PORT, mux); err != nil {
		logger.Error("Server couldn't start")
		return
	}
	slog.Info("Server Started, listening...", "PORT", PORT)
}

func genericHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	endpoint := strings.TrimPrefix(r.URL.Path, "")
	fmt.Fprintf(w, "Api endpoint %s", endpoint)
}
