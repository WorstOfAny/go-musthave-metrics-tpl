package repository

import(
	"context"
	"iter"
	"errors"
	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgconn"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	"time"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
)

type dbDecorator struct {
	db *pgxpool.Pool
}

func NewDBDecorator(db *pgxpool.Pool) *dbDecorator {
	return &dbDecorator{db: db}
}

func (dbd *dbDecorator) Set(ctx context.Context, obj models.Metrics) error {
	_, err := retry(
		func() (any, error) {
			_, err := dbd.db.Exec(
				ctx,
				"INSERT INTO metrics (id, mtype, value, delta, hash) VALUES (@id, @mtype, @value, @delta, @hash) ON CONFLICT ON CONSTRAINT metrics_pkey DO UPDATE SET value = EXCLUDED.value, delta = EXCLUDED.delta, hash = EXCLUDED.hash",
				pgx.NamedArgs{
					"id": obj.ID,
					"mtype": obj.MType,
					"value": obj.Value,
					"delta": obj.Delta,
					"hash": obj.Hash,
				},
			)

			if err != nil {
				return nil, fmt.Errorf("failed to insert object to db: %w", err)
			}

			return nil, nil
		},
		3,
	)

	if err != nil {
		return fmt.Errorf("failed to save object to repository: %w", err)
	}
	return nil
}

func (dbd *dbDecorator) BulkSet(ctx context.Context, objs []models.Metrics) error {
	_, err := retry(
		func() (any, error) {
			tx, err := dbd.db.Begin(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to start transaction: %w", err)
			}
			defer tx.Rollback(ctx)

			for _, obj := range objs {
				_, err := dbd.db.Exec(
					ctx,
					"INSERT INTO metrics (id, mtype, value, delta, hash) VALUES (@id, @mtype, @value, @delta, @hash) ON CONFLICT ON CONSTRAINT metrics_pkey DO UPDATE SET value = EXCLUDED.value, delta = EXCLUDED.delta, hash = EXCLUDED.hash",
					pgx.NamedArgs{
						"id": obj.ID,
						"mtype": obj.MType,
						"value": obj.Value,
						"delta": obj.Delta,
						"hash": obj.Hash,
					},
				)
				if err != nil { return nil, fmt.Errorf("failed to insert object to db: %w", err) }
			}
			err = tx.Commit(ctx)
			if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				return nil, fmt.Errorf("failed to commit transaction: %w", err)
			}
			return nil, nil
		},
		3,
	)

	if err != nil {
		return fmt.Errorf("failed to save objects to repository: %w", err)
	}

	return nil
}

func (dbd *dbDecorator) Get(ctx context.Context, key string) (models.Metrics, error) {
	rows, err := retry(
		func() (any, error) {
			return dbd.db.Query(
				ctx,
				"SELECT * FROM metrics WHERE CONCAT(mtype, id) = $1 LIMIT 1",
				key,
			)
		},
		3,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.Metrics{}, fmt.Errorf("failed to fetch obj from db: %w", err)
	}

	obj, err := pgx.CollectExactlyOneRow[models.Metrics](rows.(pgx.Rows), pgx.RowToStructByName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Metrics{}, ErrNotFound
		} else {
			return models.Metrics{}, fmt.Errorf("failed to parse obj from db to model: %w", err)
		}
	}

	return obj, nil
}

func (dbd *dbDecorator) Remove(ctx context.Context, key string) error {
	_, err := retry(
		func() (any, error) {
			return dbd.db.Exec(
				ctx,
				"DELETE FROM metrics WHERE CONCAT(mtype, id) = $1",
				key,
			)
		},
		3,
	)

	if err != nil {
		return fmt.Errorf("failed to delete object from db: %w", err)
	}

	return nil
}

func (dbd *dbDecorator) All(ctx context.Context) (iter.Seq[models.Metrics], error) {
	rows, err := retry(
		func() (any, error) {
			return dbd.db.Query(ctx, "SELECT * FROM metrics")
		},
		3,
	)
	if err != nil {
		return func(yield func(models.Metrics) bool) {}, fmt.Errorf("failed to fetch metrics from db")
	}

	objs, err := pgx.CollectRows[models.Metrics](rows.(pgx.Rows), pgx.RowToStructByName)

	if err != nil {
		return func(yield func(models.Metrics) bool) {}, fmt.Errorf("failed to parse metrics from db")
	}

	return func(yield func(models.Metrics) bool) {
		for _, v := range objs { if !yield(v) { return } }
	}, nil
}

func (dbd *dbDecorator) Ping(ctx context.Context) error {
	_, err := retry(
		func() (any, error) { return nil, dbd.db.Ping(ctx) },
		3,
	)
	if err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}
	return nil
}

func retry(request func() (any, error), maxRetries int) (any, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		log.Debug().Str("db attempt", strconv.Itoa(attempt)).Msg("db retry")
		res, err := request()
		if err == nil { return res, nil }

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code[:2] == "08" {
			<-time.After(time.Duration(2 * (attempt + 1) - 1) * time.Second)
			continue
		}

		return nil, fmt.Errorf("failed request: %w", err)
	}

	return nil, errors.New("db unreachable")
}
