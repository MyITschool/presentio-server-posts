package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	v0 "presentio-server-posts/src/v0"
	"time"
)

func createConn(dbUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbUrl,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDb, _ := db.DB()

	sqlDb.SetConnMaxIdleTime(time.Minute * 10)

	sqlDb.SetMaxIdleConns(10)

	sqlDb.SetMaxOpenConns(100)

	sqlDb.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func main() {
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := createConn(dbUrl)

	if err != nil {
		log.Fatalf("Unable to connect to postgres database at %s\n", dbUrl)
	}

	router := gin.Default()

	v0.SetupRouter(router.Group("/v0"), &v0.Config{Db: db})

	err = router.Run()

	if err != nil {
		log.Fatalf("Failed to start server on port %s\n", os.Getenv("PORT"))
	}
}
