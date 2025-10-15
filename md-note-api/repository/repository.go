package repository

import (
	"context"
	"database/sql"
	"fmt"
	"md-note-api/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Note struct {
	db *pgxpool.Pool
}

func NewNote(db *pgxpool.Pool) *Note {
	return &Note{db}
}

func (r *Note) FindById(ctx context.Context, id uuid.UUID) (*model.Note, error) {
	var note model.Note
	if err := r.db.QueryRow(ctx,
		"SELECT * FROM notes WHERE id = $1",
		id).Scan(&note.ID, &note.FileName); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note %d: not found", id)
		}
		return nil, fmt.Errorf("note %d: %v", id, err)
	}
	return &note, nil
}

func (r *Note) Create(ctx context.Context, note model.Note) error {
	if _, err := r.db.Exec(
		ctx,
		"INSERT INTO notes (id, filename) VALUES ($1, $2)", note.ID, note.FileName,
	); err != nil {
		return fmt.Errorf("error inserting row: %w", err)
	}

	return nil
}
