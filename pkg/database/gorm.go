package database

import (
	"database/sql"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/opentelemetry/tracing"
)

var WithUpdate = clause.Locking{Strength: clause.LockingStrengthUpdate}

func Gorm(conn *sql.DB) (*gorm.DB, error) {
	if dbLogger == nil {
		dbLogger = zap.L().WithOptions(zap.WithCaller(false)).Named("database")
	}

	dialector := postgres.New(postgres.Config{
		Conn: conn,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		dbLogger.Error("failed to open gorm connection", zap.Error(err))
		return nil, err
	}

	dbLogger.Debug("connected to gorm", zap.String("dialect", "pgx"))

	if err := db.Use(tracing.NewPlugin()); err != nil {
		dbLogger.Error("failed to use with update", zap.Error(err))
		return nil, err
	}

	lock.Lock()
	orm = db
	lock.Unlock()

	if err := AutoMigrate(db); err != nil {
		dbLogger.Error("failed to run auto migrations", zap.Error(err))
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255),
			username VARCHAR(255) UNIQUE,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			phone VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)
	`).Error
	if err != nil {
		return err
	}

	db.Exec(`ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key`)

	db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'users_username_unique') THEN
				ALTER TABLE users ADD CONSTRAINT users_username_unique
				UNIQUE (username) DEFERRABLE INITIALLY DEFERRED;
			END IF;
		END $$
	`)

	err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`).Error
	if err != nil {
		return err
	}

	err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`).Error
	if err != nil {
		return err
	}

	err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at)`).Error
	if err != nil {
		return err
	}

	return nil
}
