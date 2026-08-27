package database

import (
	"context"
	"errors"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver  string
	DSN     string
	DataDir string
}

func Open(config Config) (*gorm.DB, error) {
	driver := strings.ToLower(strings.TrimSpace(config.Driver))
	if driver == "" {
		driver = "sqlite"
	}
	switch driver {
	case "sqlite":
		dsn := strings.TrimSpace(config.DSN)
		if dsn == "" {
			if err := os.MkdirAll(config.DataDir, 0o755); err != nil {
				return nil, err
			}
			dsn = config.DataDir + "/open_ai_canvas.db?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on&_synchronous=NORMAL"
		}
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: newGORMLogger(os.Stdout)})
	case "postgres", "postgresql":
		dsn := strings.TrimSpace(config.DSN)
		if dsn == "" {
			return nil, errors.New("PostgreSQL 模式必须配置 DATABASE_URL")
		}
		return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newGORMLogger(os.Stdout)})
	default:
		return nil, fmt.Errorf("不支持的数据库驱动：%s", driver)
	}
}

func newGORMLogger(writer io.Writer) logger.Interface {
	return safeGORMLogger{delegate: logger.New(stdlog.New(writer, "", stdlog.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
		Colorful:                  false,
	})}
}

// safeGORMLogger keeps parameter filtering intact even when GORM changes the
// active log mode during initialization. Query values can contain prompts,
// credentials, email addresses, or signed URLs and must never reach SQL logs.
type safeGORMLogger struct {
	delegate logger.Interface
}

func (item safeGORMLogger) LogMode(level logger.LogLevel) logger.Interface {
	return safeGORMLogger{delegate: item.delegate.LogMode(level)}
}

func (item safeGORMLogger) Info(ctx context.Context, message string, values ...interface{}) {
	item.delegate.Info(ctx, message, values...)
}

func (item safeGORMLogger) Warn(ctx context.Context, message string, values ...interface{}) {
	item.delegate.Warn(ctx, message, values...)
}

func (item safeGORMLogger) Error(ctx context.Context, message string, values ...interface{}) {
	item.delegate.Error(ctx, message, values...)
}

func (item safeGORMLogger) Trace(ctx context.Context, begin time.Time, sql func() (string, int64), err error) {
	item.delegate.Trace(ctx, begin, func() (string, int64) {
		_, rows := sql()
		return "[SQL statement redacted]", rows
	}, err)
}

func (safeGORMLogger) ParamsFilter(_ context.Context, sql string, _ ...interface{}) (string, []interface{}) {
	return sql, nil
}

func ConfigurePool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if db.Dialector.Name() == "postgres" {
		sqlDB.SetMaxOpenConns(30)
		sqlDB.SetMaxIdleConns(10)
		return nil
	}
	sqlDB.SetMaxOpenConns(8)
	sqlDB.SetMaxIdleConns(4)
	return nil
}
