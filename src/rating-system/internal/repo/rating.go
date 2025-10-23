package repo

import (
	"context"
	"rating-system/internal/models"
	"rating-system/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type RatingRepository interface {
	GetRatingRepo(ctx context.Context, username string) (*models.Rating, error)
	UpdateRatingRepo(ctx context.Context, username string, stars int) error
}

type ratingRepo struct {
	conn postgres.Connection
}

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func NewRatingRepo(client postgres.Client) RatingRepository {
	return &ratingRepo{conn: client.Conn()}
}

func (r *ratingRepo) GetRatingRepo(ctx context.Context, username string) (*models.Rating, error) {
	query := qb.Select("id, username, stars").
		From("rating").
		Where("username = ?", username)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rate, err := pgx.CollectOneRow[models.Rating](rows, pgx.RowToStructByName)
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *ratingRepo) UpdateRatingRepo(ctx context.Context, username string, stars int) error {
	query := qb.Update("rating").
		Set("stars", stars).
		Where("username = ?", username)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
