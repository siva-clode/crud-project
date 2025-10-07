package main

import "net/http"

func (app *Application) route(mux *http.ServeMux) *http.Server {

	srv := http.Server{
		Addr:    app.config.Addr,
		Handler: mux,
	}

	mux.HandleFunc("GET /notes", app.GetNotes)
	mux.HandleFunc("GET /note", app.GetNoteById)
	mux.HandleFunc("POST /create", app.InsetNote)
	mux.HandleFunc("PUT /update", app.UpdateNote)
	mux.HandleFunc("DELETE /delete", app.DeleteNote)

	return &srv

}
