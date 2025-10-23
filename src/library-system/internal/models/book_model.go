package models

import "github.com/google/uuid"

type Book struct {
	ID        uint64    `db:"id"`
	BookUID   uuid.UUID `db:"book_uid"`
	Name      string    `db:"name"`
	Author    *string   `db:"author"`
	Genre     *string   `db:"genre"`
	Condition string    `db:"condition"  validate:"oneof=EXCELLENT GOOD BAD"`
}
