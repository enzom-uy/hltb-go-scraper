package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, func() error, error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=localhost user=enzom_uy password=%s dbname=backloggd port=5432 sslmode=disable timezone=America/Chicago", dbPassword)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, nil, err
	}

	closeFn := func() error {
		return sqlDB.Close()
	}

	log.Println("DB inicializada")
	return gormDB, closeFn, nil
}
