package seeder

import (
	"database/sql"
	"log"

	config "github.com/devonLoen/leave-request-service/config"
	entities "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
	"golang.org/x/crypto/bcrypt"
)

func SeedSuperAdmin(db *sql.DB) {
	conf := config.NewConfig()
	email := conf.SuperAdmin.Email
	password := conf.SuperAdmin.Password

	if email == "" || password == "" {
		log.Println("Skipping SuperAdmin seeding: ENV variables not set.")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	user := entities.User{
		FullName: "Super Admin",
		Email:    email,
		Password: string(hashedPassword),
		Role:     "superadmin",
	}

	query := `
		INSERT INTO users (full_name, email, password, role) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO NOTHING;
	`

	_, err = db.Exec(query, user.FullName, user.Email, user.Password, user.Role)
	if err != nil {
		log.Println("Failed to insert SuperAdmin:", err)
	} else {
		log.Println("SuperAdmin seeder finished successfully.")
	}
}
