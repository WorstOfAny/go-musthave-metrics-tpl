package repository

import(
	"context"
	"iter"
	"errors"
	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type dbDecorator[T hasSQL] struct {
	db *pgxpool.Pool
	err error
}

func NewDBDecorator[T hasSQL](db *pgxpool.Pool) *dbDecorator[T] {
	return &dbDecorator[T]{db: db}
}

func (dbd *dbDecorator[T]) Set(ctx context.Context, obj T) {
	_, err := dbd.db.Exec(ctx, obj.InsertSQL(), obj.ToPgxNamedArgs())
	if err != nil { dbd.err = err }
}

func (dbd *dbDecorator[T]) BulkSet(ctx context.Context, objs []T) {
	tx, err := dbd.db.Begin(ctx)
	defer tx.Rollback(ctx)

	for _, obj := range objs {
		dbd.db.Exec(ctx, obj.InsertSQL(), obj.ToPgxNamedArgs())
		if err != nil {
			dbd.err = err
			tx.Rollback(ctx)
			return
		}
	}
	tx.Commit(ctx)
}

func (dbd *dbDecorator[T]) Get(ctx context.Context, key string) (T, bool) {
	var obj T
	rows, err := dbd.db.Query(ctx, obj.GetSQL(), key)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		dbd.err = err
		return obj, false
	}

	obj, err = pgx.CollectExactlyOneRow[T](rows, pgx.RowToStructByName)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) { dbd.err = err }
		return obj, false
	}

	return obj, true
}

func (dbd *dbDecorator[T]) Remove(ctx context.Context, key string) {
	var obj T
	_, err := dbd.db.Exec(ctx, obj.DeleteSQL(), key)

	if err != nil { dbd.err = err }
}

func (dbd *dbDecorator[T]) All(ctx context.Context) iter.Seq[T] {
	var obj T
	rows, err := dbd.db.Query(ctx, obj.SelectSQL())
	if err != nil { dbd.err = err }

	objs, err := pgx.CollectRows[T](rows, pgx.RowToStructByName)
	if err != nil { dbd.err = err }

	return func(yield func(T) bool) {
		for _, v := range objs {
			if !yield(v) { return }
		}
	}
}

func (dbd *dbDecorator[T]) Ping(ctx context.Context) error {
	return dbd.db.Ping(ctx)
}

func (dbd *dbDecorator[T]) Err() error {
	return dbd.err
}
