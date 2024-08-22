package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DB_Init() *gorm.DB {
	POOL_SIZE := 100
	postgres_url := os.Getenv("POSTGRES_URL")
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_db := os.Getenv("POSTGRES_DB")

	dbURL := fmt.Sprintf("postgres://%v:%v@%v/%v", postgres_user, postgres_password, postgres_url, postgres_db)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Panic("Error occured while connecting to POSTGRES DB")
	}
	log.Printf("CONNECTED TO POSTGRES...")

	sqlDB, err := db.DB()

	if err != nil {
		log.Panic("Error occured while connecting to POSTGRES DB")
	}

	// set db pool size as POOL_SIZE
	sqlDB.SetMaxIdleConns(POOL_SIZE)

	// auto migrate DB
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Message{})
	db.AutoMigrate(&Hub{})

	log.Printf("Automigrate all tables...")

	return db
}
