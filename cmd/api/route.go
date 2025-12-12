package main

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (app *Application) route(mux *http.ServeMux) *http.Server {

	srv := http.Server{
		Addr:    app.config.Addr,
		Handler: mux,
	}

	// Swagger UI
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// API routes
	mux.HandleFunc("GET /notes", app.GetNotes)
	mux.HandleFunc("GET /note", app.GetNoteById)
	mux.HandleFunc("POST /create", app.InsetNote)
	mux.HandleFunc("PUT /update", app.UpdateNote)
	mux.HandleFunc("DELETE /delete", app.DeleteNote)

	return &srv

}
