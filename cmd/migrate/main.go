package main

import (
	"flag"
	"log"

	config "github.com/devonLoen/leave-request-service/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(conf *config.Config, isDown bool) {
	m, err := migrate.New("file://./migration", conf.Database.DatabaseSource)
	if err != nil {
		log.Fatal("Migration init error:", err)
	}

	if isDown {
		log.Println("Rolling back migration...")
		err = m.Down()
	} else {
		log.Println("Running migration...")
		err = m.Up()
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No changes applied (Database is already up to date)")
		} else {
			log.Fatal("Migration failed:", err)
		}
	} else {
		log.Println("Migration finished successfully!")
	}
}

func main() {
	downFlag := flag.Bool("down", false, "Set to true to rollback migrations")
	flag.Parse()

	conf := config.NewConfig()

	RunMigration(conf, *downFlag)
}
