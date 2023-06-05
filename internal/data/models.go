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
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
		Movies: MovieModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Users:  MockUserModel{},
		Tokens: MockTokenModel{},
		Movies: MockMovieModel{},
	}
}
