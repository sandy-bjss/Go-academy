package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

const PORT = ":8080"

// const for traceID key of TraceIDType
const TraceIDString = TraceIDType("traceID")

type TraceIDType string

func Api() {
	// basic logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// start new servemux
	mux := http.NewServeMux()

	// define handlers for routing
	gh := http.HandlerFunc(genericHandler)
	mux.Handle("/", middlewareOne(middlewareTwo(middlewareTraceID(gh))))
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

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("executing middleware 1")
		next.ServeHTTP(w, r)
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("executing middleware 2")
		next.ServeHTTP(w, r)
	})
}

func middlewareTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("middlewareTRaceID: adding context traceID")
		r.WithContext(context.WithValue(r.Context(), TraceIDString, uuid.New()))
		next.ServeHTTP(w, r)
	})
}

func genericHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("executing generic handler", string(TraceIDString), r.Context())
	w.Header().Set("Content-Type", "application/json")
	endpoint := strings.TrimPrefix(r.URL.Path, "")
	fmt.Fprintf(w, "Api endpoint: %s\nEndpoint type: %v", endpoint, w.Header().Get("Content-Type"))
}
