package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dmc0001/crud-project/internal/store"
)

func (app *Application) GetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := app.config.newNoteModel.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}

func (app *Application) GetNoteById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.badRequest(w, err)
		return
	}
	note, err := app.config.newNoteModel.GetById(id)
	if err != nil {
		if errors.Is(err, store.ErrNoRecord) {
			app.notFound(w, err)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)

}

func (app *Application) InsetNote(w http.ResponseWriter, r *http.Request) {
	note := &store.Note{}
	if err := json.NewDecoder(r.Body).Decode(note); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now
	id, err := app.config.newNoteModel.Create(note)

	if err != nil {
		app.serverError(w, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := map[string]any{
		"message": fmt.Sprintf("Note with id %d has been created", id),
		"id":      id,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		app.serverError(w, err)
		return
	}

}

func (app *Application) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}
	note := &store.Note{}

	if err := json.NewDecoder(r.Body).Decode(note); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	now := time.Now()
	note.UpdatedAt = now
	updatedNote, err := app.config.newNoteModel.Update(note, id)
	if err != nil {
		if errors.Is(err, store.ErrNoRecord) {
			app.notFound(w, err)
		} else {
			app.serverError(w, err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updatedNote); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

}

func (app *Application) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}
	deletedId, err := app.config.newNoteModel.Delete(id)
	if err != nil {
		if errors.Is(err, store.ErrNoRecord) {
			app.notFound(w, err)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deletedId)

}

func (app *Application) notFound(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusNotFound)
}

func (app *Application) badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
