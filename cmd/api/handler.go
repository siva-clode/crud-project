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

// GetNotes godoc
// @Summary Get all notes
// @Description Get the latest notes (up to 10)
// @Tags notes
// @Produce json
// @Success 200 {array} store.Note
// @Failure 500 {string} string "Internal Server Error"
// @Router /notes [get]
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

// GetNoteById godoc
// @Summary Get a note by ID
// @Description Get a single note by its ID
// @Tags notes
// @Produce json
// @Param id query int true "Note ID"
// @Success 200 {object} store.Note
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /note [get]
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

// CreateNoteRequest represents the request body for creating a note
type CreateNoteRequest struct {
	Title   string `json:"title" example:"My Note"`
	Content string `json:"content" example:"Note content here"`
}

// CreateNoteResponse represents the response after creating a note
type CreateNoteResponse struct {
	Message string `json:"message" example:"Note with id 1 has been created"`
	ID      int    `json:"id" example:"1"`
}

// InsetNote godoc
// @Summary Create a new note
// @Description Create a new note with title and content
// @Tags notes
// @Accept json
// @Produce json
// @Param note body CreateNoteRequest true "Note to create"
// @Success 201 {object} CreateNoteResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /create [post]
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

// UpdateNote godoc
// @Summary Update an existing note
// @Description Update a note's title and content by ID
// @Tags notes
// @Accept json
// @Produce json
// @Param id query int true "Note ID"
// @Param note body CreateNoteRequest true "Updated note data"
// @Success 200 {object} store.Note
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /update [put]
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

// DeleteNote godoc
// @Summary Delete a note
// @Description Delete a note by ID
// @Tags notes
// @Produce json
// @Param id query int true "Note ID"
// @Success 200 {integer} int "Deleted note ID"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /delete [delete]
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
