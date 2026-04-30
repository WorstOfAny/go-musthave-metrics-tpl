package repository

import(
	"context"
	"time"
	"os"
	"iter"
	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	StoreInterval int `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	RestoreStorage bool `env:"RESTORE"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

type Repository[T allowedObject] interface {
	Set(context.Context, T)
	BulkSet(context.Context, []T)
	Get(context.Context, string) (T, bool)
	Remove(context.Context, string)
	All(context.Context) iter.Seq[T]
	Ping(context.Context) error
	Err() error
}

type hasKey interface {
	Key() string
}

type hasSQL interface {
	InsertSQL() string
	SelectSQL() string
	GetSQL() string
	DeleteSQL() string
	ToPgxNamedArgs() pgx.NamedArgs
}

type allowedObject interface {
	hasKey
	hasSQL
}

func NewRepository[T allowedObject](ctx context.Context, errCh chan error, cfg Config) (Repository[T]) {

	switch {
		case cfg.DatabaseDSN != "":
			db, err := pgxpool.New(ctx, cfg.DatabaseDSN)
			if err != nil { errCh <- err }

			go func() {
				select {
					case <-ctx.Done(): db.Close()
				}
			}()

			return NewDBDecorator[T](db)
		case cfg.FileStoragePath != "":
			file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil { errCh <- err }
			go func() {
				select {
					case <-ctx.Done(): file.Close()
				}
			}()

			repo, err := NewStorage[T](file, cfg.RestoreStorage)
			if err != nil { errCh <- err }
			go repo.WriteToFile(ctx, errCh, time.Duration(cfg.StoreInterval) * time.Second)
			return repo
		default:
			repo, err := NewStorage[T](nil, false)
			if err != nil { errCh <- err }
			return repo
	}

}
