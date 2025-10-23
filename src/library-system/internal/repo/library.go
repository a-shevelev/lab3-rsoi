package repo

import (
	"context"
	"fmt"

	//"fmt"
	"lab2-rsoi/library-system/internal/models"
	"lab2-rsoi/library-system/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LibraryRepository interface {
	FetchLibrariesByCity(ctx context.Context, city string, page, size int) ([]models.Library, error)
	FetchBooksByLibrary(ctx context.Context, libraryUID uuid.UUID, showAll bool, page, size int) ([]BookWithCount, error)
	UpdateCount(ctx context.Context, bookID, libraryID uuid.UUID, count int) error
	UpdateCondition(ctx context.Context, bookUID uuid.UUID, condition string) error
	GetLibraryByUID(ctx context.Context, uid uuid.UUID) (*models.Library, error)
	GetBookByUID(ctx context.Context, uid uuid.UUID) (*BookWithCount, error)
	CountLibrariesByCity(ctx context.Context, city string) (int, error)
	CountBooksByLibrary(ctx context.Context, libraryUID uuid.UUID, showAll bool) (int, error)
	//IncreaseCount(ctx context.Context, i int, i2 int) interface{}
}

type libraryRepo struct {
	conn postgres.Connection
}

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func NewLibraryRepo(client postgres.Client) LibraryRepository {
	return &libraryRepo{conn: client.Conn()}
}

type BookWithCount struct {
	models.Book
	AvailableCount int `db:"available_count"`
}

func (r *libraryRepo) CountLibrariesByCity(ctx context.Context, city string) (int, error) {
	query := qb.Select("COUNT(*)").
		From("library").
		Where("city = ?", city)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return pgx.CollectOneRow(rows, pgx.RowTo[int])
}

func (r *libraryRepo) CountBooksByLibrary(ctx context.Context, libraryUID uuid.UUID, showAll bool) (int, error) {
	query := qb.Select("COUNT(*)").From("books b").
		Join("library l ON l.library_uid = ?", libraryUID).
		Join("library_books lb ON lb.book_id = b.id AND lb.library_id = l.id")

	if !showAll {
		query = query.Where("lb.available_count > 0")
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return pgx.CollectOneRow(rows, pgx.RowTo[int])
}

func (r *libraryRepo) FetchLibrariesByCity(ctx context.Context, city string, page, size int) ([]models.Library, error) {
	offset := (page - 1) * size
	query := qb.Select("id", "library_uid", "name", "city", "address").
		From("library").
		Where(squirrel.Eq{"city": city})
	query = query.OrderBy("name ASC").Limit(uint64(size)).Offset(uint64(offset))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	libraries, err := pgx.CollectRows[models.Library](rows, pgx.RowToStructByNameLax)
	if err != nil {
		return nil, err
	}
	return libraries, nil
}

func (r *libraryRepo) FetchBooksByLibrary(ctx context.Context, libraryUID uuid.UUID, showAll bool, page, size int) ([]BookWithCount, error) {
	offset := (page - 1) * size

	query := qb.Select(
		"b.id", "b.book_uid", "b.name", "b.author", "b.genre", "b.condition",
		"lb.available_count",
	).
		From("books b").
		Join("library l ON l.library_uid = ?", libraryUID).
		Join("library_books lb ON lb.book_id = b.id AND lb.library_id = l.id")

	if !showAll {
		query = query.Where("lb.available_count > 0")
	}

	query = query.OrderBy("b.name ASC").Limit(uint64(size)).Offset(uint64(offset))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books, err := pgx.CollectRows[BookWithCount](rows, pgx.RowToStructByNameLax)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (r *libraryRepo) UpdateCount(ctx context.Context, bookUID, libraryUID uuid.UUID, count int) error {
	sql := `
        UPDATE library_books lb
        SET available_count = $1
        FROM books b, library l
        WHERE b.book_uid = $2
          AND l.library_uid = $3
          AND lb.book_id = b.id
          AND lb.library_id = l.id;
    `
	result, err := r.conn.Exec(ctx, sql, count, bookUID, libraryUID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no record found for book_uid %s and library_uid %s", bookUID, libraryUID)
	}

	return nil
}

func (r *libraryRepo) UpdateCondition(ctx context.Context, bookUID uuid.UUID, condition string) error {
	query := qb.Update("books").
		Set("condition", condition).
		Where("book_uid = ?", bookUID)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.conn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no record found for book_uid %d", bookUID)
	}

	return nil
}

func (r *libraryRepo) GetBookByUID(ctx context.Context, uid uuid.UUID) (*BookWithCount, error) {
	query := qb.Select("b.id, b.book_uid, b.name, b.author, b.genre, b.condition, lb.available_count").
		From("books b").
		Join("library_books lb ON lb.book_id = b.id").
		Join("library l ON l.id = lb.library_id").
		Where("b.book_uid = ?", uid)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	book, err := pgx.CollectOneRow[BookWithCount](rows, pgx.RowToStructByName)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *libraryRepo) GetLibraryByUID(ctx context.Context, uid uuid.UUID) (*models.Library, error) {
	query := qb.Select("id, library_uid, name, city, address").
		From("library").
		Where("library_uid = ?", uid)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	library, err := pgx.CollectOneRow[models.Library](rows, pgx.RowToStructByName)
	if err != nil {
		return nil, err
	}
	return &library, nil
}
