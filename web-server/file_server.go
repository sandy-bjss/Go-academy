package main

import (
	"log/slog"
	"net/http"
	"os"
	"text/template"
)

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
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/about", http.StripPrefix("/about", fileServer))

	// dynamic
	mux.HandleFunc("/list", dynamicHandler)

	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		logger.Error("File Server couldn't start")
	}
}

func dynamicHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("dynamic/index.html")
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
