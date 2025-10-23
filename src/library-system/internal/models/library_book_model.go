package models

type LibraryBook struct {
	BookID         uint64 `db:"book_id"`
	LibraryID      uint64 `db:"library_id"`
	AvailableCount int    `db:"available_count"`
}
