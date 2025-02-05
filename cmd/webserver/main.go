package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	// create a channel for waiting for signal from OS
	sigs := make(chan os.Signal, 1)
	// notify the channel of a signal from the OS
	signal.Notify(sigs, os.Interrupt)
	fmt.Println("Starting server...\nCTRL-C to shutdown")

	// start api
	go FileServer()

	// capture signal
	<-sigs
	fmt.Println("\nKeyboard interupt...\nShutting down server")
}

const PORT = ":8080"

var tmpl *template.Template

type Todo struct {
	Item   string
	Status bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func FileServer() {
	// logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// start new servemux
	mux := http.NewServeMux()

	// static
	fileServer := http.FileServer(http.Dir("../../web/static/"))
	mux.Handle("/about", http.StripPrefix("/about", fileServer))

	// dynamic
	mux.HandleFunc("/list", dynamicHandler)

	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		logger.Error("File Server couldn't start")
	}
}

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../web/templates/index.html")
	if err != nil {
		slog.Error("Could not parse html template")
		return
	}
	data := TodoPageData{
		PageTitle: "My TODO list",
		Todos: []Todo{
			{Item: "Task 1", Status: true},
			{Item: "Task 2", Status: false},
			{Item: "Task 3", Status: false},
		},
	}

	tmpl.Execute(w, data)
}
