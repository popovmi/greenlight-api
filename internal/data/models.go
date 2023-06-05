package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
		Update(user *User) error
		GetByToken(tokenPlainText, scope string) (*User, error)
	}
	Permissions interface {
		GetAllForUser(userID int64) (Permissions, error)
		GrantForUser(userID int64, permissions ...string) error
	}
	Tokens interface {
		New(userID int64, ttl time.Duration, scope string) (*Token, error)
		Insert(token *Token) error
		DeleteAllForUser(scope string, userID int64) error
	}
	Movies interface {
		GetMany(title string, genres []string, lp ListParams) ([]*Movie, Metadata, error)
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
	}
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Users:       UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Movies:      MovieModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Users:       MockUserModel{},
		Permissions: MockPermissionModel{},
		Tokens:      MockTokenModel{},
		Movies:      MockMovieModel{},
	}
}
