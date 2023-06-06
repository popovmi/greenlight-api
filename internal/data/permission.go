package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for _, p := range p {
		if p == code {
			return true
		}
	}

	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
	SELECT p.code
	FROM permissions p
	INNER JOIN users_permissions up
	ON up.permission_id = p.id
	INNER JOIN users u
	ON u.id = up.user_id
	WHERE u.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return permissions, nil
}

func (m PermissionModel) GrantForUser(userID int64, permissions ...string) error {
	query := `
	INSERT INTO users_permissions
	SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(permissions))

	return err
}

type MockPermissionModel struct{}

func (m MockPermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	return nil, nil
}

func (m MockPermissionModel) GrantForUser(userID int64, permissions ...string) error {
	return nil
}
