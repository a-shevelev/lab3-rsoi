package repo

import (
	"context"
	"fmt"
	"reservation-system/internal/models"
	"reservation-system/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type ReservationRepo interface {
	CreateReservation(ctx context.Context, res models.Reservation) (*models.Reservation, error)
	GetReservationByUID(ctx context.Context, uid string) (*models.Reservation, error)
	GetCurrentReservationsAmount(ctx context.Context, username string) (uint64, error)
	GetReservations(ctx context.Context, username string) ([]models.Reservation, error)
	UpdateReservationStatus(ctx context.Context, reservationUID uuid.UUID, status string) error
	Delete(ctx context.Context, reservationUID string) error
}

type reservationRepo struct {
	conn postgres.Connection
}

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func NewReservationRepo(conn postgres.Client) ReservationRepo {
	return &reservationRepo{conn: conn.Conn()}
}

func (r *reservationRepo) CreateReservation(ctx context.Context, res models.Reservation) (*models.Reservation, error) {
	query := qb.Insert("reservation").
		Columns("username", "library_uid", "book_uid", "start_date", "till_date", "status", "reservation_uid").
		Values(
			res.Username,
			res.LibraryUID,
			res.BookUID,
			res.StartDate,
			res.TillDate,
			res.Status,
			res.ReservationUID,
		).
		Suffix("RETURNING *")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	model, err := pgx.CollectOneRow[models.Reservation](rows, pgx.RowToStructByName)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *reservationRepo) GetReservationByUID(ctx context.Context, uid string) (*models.Reservation, error) {
	query := qb.Select("id", "reservation_uid", "username",
		"book_uid", "library_uid", "start_date", "till_date", "status").
		From("reservation").
		Where(squirrel.Eq{"reservation_uid": uid})
	sql, args, err := query.ToSql()
	fmt.Println(sql)
	if err != nil {
		return nil, err
	}
	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	model, err := pgx.CollectOneRow[models.Reservation](rows, pgx.RowToStructByNameLax)
	if err != nil {
		log.WithError(err).Errorf("get reservation by uid: %s", uid)
		return nil, err
	}

	return &model, nil
}

func (r *reservationRepo) GetReservations(ctx context.Context, username string) ([]models.Reservation, error) {
	query := qb.Select("id", "reservation_uid", "username",
		"book_uid", "library_uid", "start_date", "till_date", "status").
		From("reservation").
		Where(squirrel.Eq{"username": username})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows[models.Reservation](rows, pgx.RowToStructByName)
}

func (r *reservationRepo) GetCurrentReservationsAmount(ctx context.Context, username string) (uint64, error) {
	query := qb.Select("count(*)").
		From("reservation").
		Where(squirrel.Eq{
			"username": username,
			"status":   "RENTED",
		})
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}
	rows, err := r.conn.Query(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowTo[uint64])
}

func (r *reservationRepo) UpdateReservationStatus(ctx context.Context, reservationUID uuid.UUID, status string) error {
	query := qb.Update("reservation").
		Set("status", status).
		Where(squirrel.Eq{"reservation_uid": reservationUID})
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

func (r *reservationRepo) Delete(ctx context.Context, reservationUID string) error {
	sql, args, err := qb.Delete("reservation").Where(squirrel.Eq{"reservation_uid": reservationUID}).ToSql()
	if err != nil {
		return err
	}
	_, err = r.conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil

}
