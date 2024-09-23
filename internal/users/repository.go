package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go-web-api-starter/internal/database"
	"strings"
	"time"
)

type UserPsqlRepo struct {
	DB *database.DB
}

func (m UserPsqlRepo) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var roleID int64
	err := m.DB.QueryRowContext(ctx, "SELECT id from roles WHERE name = $1", user.Role.Name).Scan(&roleID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRoleNotFound
		default:
			return fmt.Errorf("error getting role: %w", err)
		}
	}

	query := `INSERT INTO users(id, email, role_id)
              VALUES ($1, $2, $3)
              RETURNING id, created_at, updated_at`

	args := []any{user.ID, user.Email, roleID}

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m UserPsqlRepo) GetById(id uuid.UUID) (*User, error) {
	query := `SELECT users.id, users.email, users.is_deleted, users.created_at, users.updated_at,
                  roles.name,
                  STRING_AGG(DISTINCT permissions.code, ',') as permissions
              FROM users
              LEFT JOIN roles ON users.role_id = roles.id
              LEFT JOIN roles_permissions ON roles.id = roles_permissions.role_id
              LEFT JOIN permissions ON roles_permissions.permission_id = permissions.id
              WHERE users.id = $1
              GROUP BY users.id, users.email, users.is_deleted, users.created_at, users.updated_at, roles.name`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	var roleName, permissions sql.NullString

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.IsDeleted,
		&user.CreatedAt,
		&user.UpdatedAt,
		&roleName,
		&permissions,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrRecordNotFound
		}
		return nil, err
	}

	user.Role.Name = roleName.String
	if permissions.Valid {
		user.Role.Permissions = strings.Split(permissions.String, ",")
	} else {
		user.Role.Permissions = []string{}
	}

	return &user, nil
}

func (m UserPsqlRepo) Update(user *User) error {
	query := `UPDATE users
              SET email = $1, updated_at = CURRENT_TIMESTAMP
              WHERE id = $2 AND updated_at <= CURRENT_TIMESTAMP
              RETURNING updated_at`

	args := []any{
		user.Email,
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return database.ErrEditConflict
		default:
			var pqErr *pq.Error
			ok := errors.As(err, &pqErr)
			if ok {
				if pqErr.Code == "23505" {
					return ErrDuplicateEmail
				}
			}
			return err
		}
	}

	return nil
}

func (m UserPsqlRepo) UpdateEmail(user *User) error {
	query := `UPDATE users 
              SET email = $1, updated_at = CURRENT_TIMESTAMP
              WHERE id = $2
              RETURNING updated_at`

	args := []any{user.Email, user.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return database.ErrRecordNotFound
		default:
			var pqErr *pq.Error
			ok := errors.As(err, &pqErr)
			if ok {
				if pqErr.Code == "23505" {
					return ErrDuplicateEmail
				}
			}
			return err
		}
	}

	return nil
}

func (m UserPsqlRepo) Delete(id uuid.UUID, delEmail string) error {
	query := `UPDATE users
              SET is_deleted = TRUE, email = $1, updated_at = CURRENT_TIMESTAMP
              WHERE id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{delEmail, id}

	res, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return database.ErrRecordNotFound
	}

	return nil
}

func (m RolePsqlRepo) GetRoleForUser(userID uuid.UUID) (*Role, error) {
	query := `SELECT roles.name, permissions.code
              FROM roles
              INNER JOIN roles_permissions ON roles_permissions.role_id = roles.id
              INNER JOIN permissions ON permissions.id = roles_permissions.permission_id
              INNER JOIN users on roles.id = users.role_id
              WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	role := &Role{}

	for rows.Next() {
		var permission string

		err = rows.Scan(&role.Name, &permission)
		if err != nil {
			return nil, err
		}

		role.Permissions = append(role.Permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return role, nil
}

func (m RolePsqlRepo) UpdateRoleForUser(userID uuid.UUID, roleName string) error {
	query := `UPDATE users SET role_id = (SELECT id FROM roles WHERE name = $2) WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, roleName)
	return err
}
