package data

import (
	"database/sql"
	"errors"
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
		Movies: MovieModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Users:  MockUserModel{},
		Movies: MockMovieModel{},
	}
}
