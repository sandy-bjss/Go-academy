package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"
)

var static embed.FS

func FileServer() {
	// basic logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Replace "." with the actual path of the directory you want to expose.
	directoryPath := "."

	// Check if the directory exists
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		logger.Error("Directory not found.\n", "path", directoryPath)
		return
	}

	// Create a file server handler to serve the directory's contents
	fileServer := http.FileServer(http.Dir("static"))

	// Create a new HTTP server and handle requests
	http.Handle("/about", http.StripPrefix("/about", fileServer))

	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		logger.Error("File Server couldn't start")
	}
}
