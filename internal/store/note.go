package store

import (
	"database/sql"
	"errors"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type NoteModel struct {
	Db *sql.DB
}

func NewNoteModel(db *sql.DB) *NoteModel {

	return &NoteModel{
		Db: db,
	}
}

func (n *NoteModel) GetById(id int) (*Note, error) {
	note := &Note{}
	query := "SELECT id, title, content, createdat, updatedat FROM note WHERE id = $1"
	row := n.Db.QueryRow(query, id)

	err := row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return note, nil
}

func (n *NoteModel) Update(update *Note, id int) (*Note, error) {
	query := `
		UPDATE note
		SET title = $1,
		    content = $2,
		    updatedat = $3
		WHERE id = $4;`

	_, err := n.Db.Exec(query, update.Title, update.Content, update.UpdatedAt, id)
	if err != nil {
		return nil, err
	}

	// 🔥 Re-fetch the updated note from DB
	updatedNote, err := n.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return updatedNote, nil
}

func (n *NoteModel) Latest() ([]*Note, error) {
	notes := []*Note{}
	query := "SELECT Id,Title,Content,CreatedAt,UpdatedAt FROM note LIMIT 10"
	rows, err := n.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		n := &Note{}
		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (n *NoteModel) Create(note *Note) (*int, error) {
	query := `INSERT INTO note(title, content, createdat, updatedat)
          VALUES ($1,$2,$3,$4)
          RETURNING id;`

	err := n.Db.QueryRow(query, note.Title, note.Content, note.CreatedAt, note.UpdatedAt).Scan(&note.ID)
	if err != nil {
		return nil, err
	}

	return &note.ID, nil
}

func (n *NoteModel) Delete(id int) (int, error) {
	query := "Delete FROM note WHERE id = $1 RETURNING id"
	var deletedID int
	err := n.Db.QueryRow(query, id).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrNoRecord
		}
		return -1, err
	}

	return deletedID, nil
}
