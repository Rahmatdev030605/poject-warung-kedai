package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config adalah struktur data untuk konfigurasi database.
type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

// NewConnection membuat dan mengembalikan koneksi ke database.
func NewConnection(config *Config) (*gorm.DB, error) {
	// Membuat string koneksi (DSN) untuk PostgreSQL.
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Membuka koneksi ke database PostgreSQL.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	return db, nil
}
