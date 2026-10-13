package course

import (
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db,err := gorm.open((
		postgres.Open(dsn),
		&grm.config{},
	))

	if err != nil {
		return nil ,err
	}

	ftm.println("Database connected");

	return db,nil

}