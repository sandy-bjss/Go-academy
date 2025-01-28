package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

const PORT = ":8080"

func Api() {

	// basic logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// start new servemux
	mux := http.NewServeMux()

	// define handlers for routing
	mux.HandleFunc("/create", testHandler)
	mux.HandleFunc("/get", testHandler)
	mux.HandleFunc("/update", testHandler)
	mux.HandleFunc("/delete", testHandler)

	// start server
	if err := http.ListenAndServe(PORT, mux); err != nil {
		logger.Error("Server couldn't start")
		return
	}
	slog.Info("Server Started, listening...", "PORT", PORT)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test endpoint")
}
