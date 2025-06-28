package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sync"
	"time"

	"pokedex_backend_go/pkg/logger"

	"github.com/XSAM/otelsql"
	"github.com/pressly/goose/v3"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var fsys embed.FS

var (
	lock     sync.Mutex
	once     sync.Once
	dbLogger *zap.Logger
	orm      *gorm.DB
)

const driverName = "pgx"

func Invoke(_ *sql.DB, _ *gorm.DB) {
	once.Do(func() {
		if err := goose.SetDialect("pgx"); err != nil {
			dbLogger.Error("failed to set goose dialect", zap.Error(err))
			panic(err)
		}
		goose.SetBaseFS(fsys)
	})
}

func Connection() (*sql.DB, error) {
	dbLogger = logger.NewLogger("database")

	const url = "host=localhost port=5432 user=pokedex_backend_go password=pokedex_backend_go dbname=pokedex_backend_go sslmode=disable search_path=public timezone=UTC"

	opts := []otelsql.Option{
		otelsql.WithSQLCommenter(true),
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
			semconv.ServiceNameKey.String(driverName),
			semconv.TelemetrySDKLanguageGo,
		),
	}

	db, err := otelsql.Open(driverName, url, opts...)
	if err != nil {
		dbLogger.Error("failed to open database", zap.Error(err))
		return nil, err
	}

	err = otelsql.RegisterDBStatsMetrics(db, opts...)
	if err != nil {
		dbLogger.Error("failed to register database stats metrics", zap.Error(err))
		return nil, err
	}

	if err := db.Ping(); err != nil {
		dbLogger.Error("failed to ping database", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var version string
	err = db.QueryRowContext(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		dbLogger.Error("failed to get database version", zap.Error(err))
		return nil, err
	}

	dbLogger.Info("connected to database", zap.String("version", version))

	return db, nil
}

type gooseLogger struct {
	*zap.Logger
}

// Fatalf implements goose.Logger.
func (g *gooseLogger) Fatalf(format string, v ...interface{}) {
	g.Logger.Fatal(fmt.Sprintf(format, v...))
}

// Printf implements goose.Logger.
func (g *gooseLogger) Printf(format string, v ...interface{}) {
	g.Logger.Info(fmt.Sprintf(format, v...))
}

func MigrationUp(ctx context.Context, db *sql.DB) error {
	dbLogger.Info("running database migrations...")
	goose.SetLogger(&gooseLogger{dbLogger})
	return goose.UpContext(ctx, db, "migrations")
}

func MigrationDown(ctx context.Context, db *sql.DB) error {
	dbLogger.Info("rolling back database migrations...")
	goose.SetLogger(&gooseLogger{dbLogger})
	return goose.DownContext(ctx, db, "migrations")
}
