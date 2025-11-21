package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/devonLoen/leave-request-service/internal/app/rest_api/database"
	entity "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
)

type User struct {
	database.BaseSQLRepository[entity.User]
}

func NewUserRepository(db *sql.DB) *User {
	return &User{
		BaseSQLRepository: database.BaseSQLRepository[entity.User]{DB: db},
	}
}

func mapUser(rows *sql.Row, u *entity.User) error {
	return rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Role)
}

func mapUserWithPassword(rows *sql.Row, u *entity.User) error {
	return rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.Password)
}

func mapUsers(rows *sql.Rows, u *entity.User) error {
	return rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Role)
}

func (r *User) FindByEmail(email string) (*entity.User, error) {
	return r.SelectSingle(
		mapUser,
		"SELECT u.id, u.full_name, u.email, u.role FROM users u WHERE u.email = $1",
		email,
	)
}

func (r *User) FindByEmailWithPassword(email string) (*entity.User, error) {
	return r.SelectSingle(
		mapUserWithPassword,
		"SELECT u.id, u.full_name, u.email, u.role, u.password FROM users u WHERE u.email = $1",
		email,
	)
}

func (r *User) FindById(id int) (*entity.User, error) {
	return r.SelectSingle(
		mapUser,
		"SELECT u.id, u.full_name, u.email, u.role FROM users u WHERE u.id = $1",
		id,
	)
}

func (r *User) GetAllUsers(limit, offset int, sortBy, orderBy, search string, filter entity.UserFilter) ([]*entity.User, error) {
	baseQuery := "SELECT u.id, u.full_name, u.email, u.role FROM users u"
	var conditions []string
	var args []interface{}

	argId := 1

	if filter.Role != "" {
		conditions = append(conditions, fmt.Sprintf("(u.role = $%d)", argId))
		args = append(args, filter.Role)
		argId++
	}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(u.full_name ILIKE $%d OR u.role::text ILIKE $%d OR u.email ILIKE $%d)",
			argId, argId, argId,
		))
		args = append(args, "%"+search+"%")
		argId++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += fmt.Sprintf(" ORDER BY u.%s %s LIMIT $%d OFFSET $%d", sortBy, orderBy, argId, argId+1)
	fmt.Println((baseQuery))
	args = append(args, limit, offset)

	return r.SelectMultiple(
		mapUsers,
		baseQuery,
		args...,
	)
}

func (r *User) Create(user *entity.User) error {
	_, err := r.Insert(
		"INSERT INTO users (full_name, email, password, role) VALUES ($1, $2, $3, $4)",
		user.FullName, user.Email, user.Password, user.Role,
	)

	return err
}
