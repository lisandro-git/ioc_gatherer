package db

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/logger"
)

/*
SELECT 'CREATE DATABASE <db_name>'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '<db_name>')
*/

// ConnectToDB creates a database if it doesn't exist
func ConnectToDB() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return db
}

func BuildDatabase(db *gorm.DB) {
	var err error = db.AutoMigrate(
		&SourceIPDescription{},
		&IPDataDescription{},
	)
	if err != nil {
		return
	}
}
