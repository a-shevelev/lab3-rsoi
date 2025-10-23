package models

type Rating struct {
	ID       uint64 `db:"id"`
	Username string `db:"username"`
	Stars    int    `db:"stars"`
}
