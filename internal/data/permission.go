package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (self Permissions) Include(code string) bool {
	for _, p := range self {
		if p == code {
			return true
		}
	}

	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (self PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
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

	rows, err := self.DB.QueryContext(ctx, query, userID)
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

func (self PermissionModel) GrantForUser(userID int64, permissions ...string) error {
	query := `
	INSERT INTO users_permissions
	SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := self.DB.ExecContext(ctx, query, userID, pq.Array(permissions))

	return err
}

type MockPermissionModel struct{}

func (self MockPermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	return nil, nil
}

func (self MockPermissionModel) GrantForUser(userID int64, permissions ...string) error {
	return nil
}
