package user

import (
	"app/domain/model"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func insertUserQuery() string {
	return fmt.Sprintf(
		`INSERT INTO %s(email, first_name, last_name)
		VALUES (:email, :first_name, :last_name)
		RETURNING id;`,
		model.UserTable,
	)
}
func deleteUserByIdQuery(id uuid.UUID) string {
	return fmt.Sprintf(
		`DELETE FROM %s
		WHERE id = '%s'`,
		model.UserTable,
		id.String(),
	)
}

func selectUserByIdQuery(id uuid.UUID) string {
	return fmt.Sprintf(
		`SELECT * FROM %s
		WHERE id = '%s' LIMIT 1`,
		model.UserTable,
		id.String(),
	)
}
func upsertUserQuery() string {
	return fmt.Sprintf(
		`INSERT INTO %s(id, email, first_name, last_name)
		VALUES (:id, :email, :first_name, :last_name)
		ON CONFLICT (id)
		DO UPDATE SET email=:email, first_name=:first_name, last_name=:last_name
		RETURNING id;`,
		model.UserTable,
	)
}

func selectUserWithOptions(opts *model.UserOptions) string {
	s := fmt.Sprintf(
		`SELECT * FROM %s`,
		model.UserTable,
	)
	if !opts.GtDatetime.IsZero() {
		s = fmt.Sprintf(
			`%s WHERE created_at >= '%v'`,
			s,
			opts.GtDatetime.Format(time.RFC3339),
		)
	}
	if !opts.LtDatetime.IsZero() {
		if !opts.GtDatetime.IsZero() {
			s = fmt.Sprintf(`%s AND`, s)
		} else {
			s = fmt.Sprintf(`%s WHERE`, s)
		}
		s = fmt.Sprintf(
			`%s created_at <= '%v'`,
			s,
			opts.LtDatetime.Format(time.RFC3339),
		)
	}
	if opts.Limit == 0 {
		opts.Limit = 25
	}
	return fmt.Sprintf(
		`%s LIMIT %v OFFSET %v;`,
		s,
		opts.Limit,
		opts.Offset,
	)
}

func selectUserCount() string {
	return fmt.Sprintf(
		`SELECT COUNT(*) from %s`,
		model.UserTable,
	)
}
