package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"showmaster/src/config"
)

func GetPSQLDatabase(cfg *config.Config) *gorm.DB {
	var (
		dbURI = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s TimeZone=%s",
			cfg.Database.PSQL.Host,
			cfg.Database.PSQL.User,
			cfg.Database.PSQL.DBName,
			cfg.Database.PSQL.Password,
			cfg.Database.PSQL.TimeZone,
		)

		err error
	)

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v\n", err)
	}

	return db
}
