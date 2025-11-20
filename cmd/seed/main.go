package main

import (
	"database/sql"
	"log"

	"github.com/devonLoen/leave-request-service/config"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/seeder"

	_ "github.com/lib/pq"
)

func main() {
	conf := config.NewConfig()

	db, err := sql.Open(conf.Database.DatabaseDriver, conf.Database.DatabaseSource)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	seeder.SeedSuperAdmin(db)
}
