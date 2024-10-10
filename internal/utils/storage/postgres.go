package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testEffectiveMobile/internal/utils/config"
)

func MustLoadPostgres(cfg config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %s", err.Error()))
	}
	return db
}
