package migrations

import "embed"

//go:embed *

// MigrationsFS содержит миграции для БД
var MigrationsFS embed.FS
