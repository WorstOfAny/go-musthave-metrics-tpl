package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/migrations"
	models "github.com/WorstOfAny/go-musthave-metrics-tpl/internal/model"
)

// Config настройки для репозитория
type Config struct {
	StoreInterval   time.Duration `env:"STORE_INTERVAL" json:"store_interval"` // время обновления файлового хранилища
	FileStoragePath string        `env:"FILE_STORAGE_PATH" json:"store_file"`  // путь к файловому хранилищу
	RestoreStorage  bool          `env:"RESTORE" json:"restore"`               // нужно ли восстанавливать хранилище из файла при запуске приложения
	DatabaseDSN     string        `env:"DATABASE_DSN" json:"database_dsn"`     // адрес БД
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	al := &struct {
		StoreInterval string `json:"store_interval"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, al); err != nil {
		return err
	}

	duration, err := time.ParseDuration(al.StoreInterval)
	if err != nil {
		return fmt.Errorf("failed parse duration: %w", err)
	}

	c.StoreInterval = duration
	return nil
}

// Repository интерфейс для файлового хранилища и обертки БД
type Repository interface {
	Set(context.Context, models.Metrics) error
	BulkSet(context.Context, []models.Metrics) error
	Get(context.Context, string) (models.Metrics, error)
	Remove(context.Context, string) error
	All(context.Context) (iter.Seq[models.Metrics], error)
	Ping(context.Context) error
}

type repositoryError string

func (re repositoryError) Error() string {
	return string(re)
}

// ErrNotFound ошибка, возвращаемая, если метрика не найдена в репозитории
const ErrNotFound = repositoryError("metric not found")

// NewRepository конструктор, возвращающий интерфейс для работы с хранилищем
// в зависимости от настроек, будет возвращаться либо repository.storage, либо repository.dbDecorator
func NewRepository(ctx context.Context, cfg Config) (Repository, error) {

	switch {
	case cfg.DatabaseDSN != "":
		sourceDriver, err := iofs.New(migrations.MigrationsFS, "migrations")
		if err != nil {
			return nil, fmt.Errorf("failed initialize source driver for migrations: %w", err)
		}
		m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, cfg.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("failed initialize migrations: %w", err)
		}
		err = m.Up()

		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}

		db, err := pgxpool.New(ctx, cfg.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize db connections pool: %w", err)
		}

		go func() {
			<-ctx.Done()
			db.Close()
		}()

		return NewDBDecorator(db), nil
	case cfg.FileStoragePath != "":
		file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		go func() {
			<-ctx.Done()
			file.Close()
		}()

		repo, err := NewStorage(file, cfg.RestoreStorage)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize mem storage: %w", err)
		}
		go repo.WriteToFile(ctx, time.Duration(cfg.StoreInterval)*time.Second)
		return repo, nil
	default:
		repo, err := NewStorage(nil, false)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize mem storage: %w", err)
		}
		return repo, nil
	}

}
