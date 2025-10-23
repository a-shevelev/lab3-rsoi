package models

import "github.com/google/uuid"

type Library struct {
	ID         uint64    `db:"id"`
	LibraryUID uuid.UUID `db:"library_uid"`
	Name       string    `db:"name"`
	City       string    `db:"city"`
	Address    string    `db:"address"`
}
