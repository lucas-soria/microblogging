package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBClient interface {
	AutoMigrate(value any) error
	Create(ctx context.Context, value any) error
	First(ctx context.Context, dest any, conds ...any) error
	Save(ctx context.Context, value any) error
	Delete(ctx context.Context, value any, conds ...any) error
	WithContext(ctx context.Context) *gorm.DB
}

type PostgresClient struct {
	db *gorm.DB
}

func NewPostgresClient(dsn string) (*PostgresClient, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// For UUID generation
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";").Error; err != nil {
		log.Fatalf("Failed to create pgcrypto extension: %v", err)
	}

	return &PostgresClient{db: db}, nil
}

func (p *PostgresClient) AutoMigrate(value any) error {
	return p.db.AutoMigrate(value)
}

func (p *PostgresClient) Create(ctx context.Context, value any) error {
	return p.db.WithContext(ctx).Create(value).Error
}

func (p *PostgresClient) First(ctx context.Context, dest any, conds ...any) error {
	return p.db.WithContext(ctx).First(dest, conds...).Error
}

func (p *PostgresClient) Save(ctx context.Context, value any) error {
	return p.db.WithContext(ctx).Save(value).Error
}

func (p *PostgresClient) Delete(ctx context.Context, value any, conds ...any) error {
	return p.db.WithContext(ctx).Delete(value, conds...).Error
}

func (p *PostgresClient) WithContext(ctx context.Context) *gorm.DB {
	return p.db.WithContext(ctx)
}
