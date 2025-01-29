package main

import (
	"log/slog"
	"net/http"
	"os"
)

type Todo struct {
	Item   string
	Status bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func FileServer() {
	// basic logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Check if the directory exists
	directoryPath := "."
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		slog.Error("Directory not found.\n", "path", directoryPath)
		return
	}

	//=====================================================================================
	// start new servemux
	mux := http.NewServeMux()

	// static
	// Create a file server handler to serve the directory's contents
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/", fileServer) // http.StripPrefix("/", fileServer)

	// dynamic
	//mux.HandleFunc("/list", dynamicHandler)

	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		logger.Error("File Server couldn't start")
	}
	//=====================================================================================
}

// func dynamicHandler(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseFS(tmplFS, "index.html")
// 	if err != nil {
// 		slog.Error("Could not parse html template")
// 		return
// 	}
// 	data := TodoPageData{
// 		PageTitle: "My TODO list",
// 		Todos: []Todo{
// 			{Item: "Task 1", Status: true},
// 			{Item: "Task 2", Status: false},
// 			{Item: "Task 3", Status: false},
// 		},
// 	}

// 	tmpl.Execute(w, data)
// }
