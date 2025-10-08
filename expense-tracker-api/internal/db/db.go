package db

import (
	"spendime/internal/transaction"
	"spendime/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&user.User{}, &transaction.Transaction{}); err != nil {
		return nil, err
	}

	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
}
