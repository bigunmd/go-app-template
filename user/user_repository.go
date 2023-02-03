package user

import (
	"app/domain/model"
	"app/domain/repository"
	"app/pkg/logger"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type r struct {
	log logger.Logger
	db  *sqlx.DB
}

// SelectUser implements repository.UserRepository
func (r *r) SelectUser(opts *model.UserOptions) (*model.Users, error) {
	us := new(model.Users)
	us.Limit = opts.Limit
	us.Offset = opts.Offset
	err := r.db.Get(&us.Total, selectUserCount())
	if err != nil {
		return nil, err
	}
	if us.Total == 0 {
		return us, nil
	}
	q := selectUserWithOptions(opts)
	r.log.AddField("query", q).Debug()
	err = r.db.Select(&us.Data, q)
	if err != nil {
		return nil, err
	}
	return us, nil
}

// UpsertUser implements UserRepository
func (r *r) UpsertUser(u *model.User) (*model.User, error) {
	stmt, err := r.db.PrepareNamed(upsertUserQuery())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.Get(u, u)
	if err != nil {
		return nil, err
	}
	u, err = r.SelectUserById(u.Id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// SelectUserById implements UserRepository
func (r *r) SelectUserById(id uuid.UUID) (*model.User, error) {
	u := new(model.User)
	err := r.db.Get(u, selectUserByIdQuery(id))
	if err != nil {
		return nil, err
	}
	return u, nil
}

// DeleteUser implements UserRepository
func (r *r) DeleteUserById(id uuid.UUID) error {
	res, err := r.db.Exec(deleteUserByIdQuery(id))
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// InsertUser implements UserRepository
func (r *r) InsertUser(u *model.User) (*model.User, error) {
	stmt, err := r.db.PrepareNamed(insertUserQuery())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.Get(u, u)
	if err != nil {
		return nil, err
	}
	u, err = r.SelectUserById(u.Id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func NewUserRepository(db *sqlx.DB, log logger.Logger) repository.UserRepository {
	return &r{
		log: log,
		db:  db,
	}
}
