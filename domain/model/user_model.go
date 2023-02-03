package model

import (
	"time"

	"github.com/google/uuid"
)

const UserTable string = "service.user"

type UserOptions struct {
	Limit      int       `json:"limit,omitempty" example:"100"`
	Offset     int       `json:"offset,omitempty" example:"100"`
	GtDatetime time.Time `json:"gt_datetime,omitempty" example:"1674582282"`
	LtDatetime time.Time `json:"lt_datetime,omitempty" example:"1674582282"`
}

type UserBase struct {
	Email     string `json:"email,omitempty" db:"email" validate:"required,email" example:"example@email.com"`
	FirstName string `json:"first_name,omitempty" db:"first_name" validate:"required,gte=3,lte=255" example:"John"`
	LastName  string `json:"last_name,omitempty" db:"last_name" validate:"gte=3,lte=255" example:"Doe"`
}

type User struct {
	Id uuid.UUID `json:"id,omitempty" db:"id" example:"31afbd8d-0312-4f18-87ee-24d5881a619e"`
	*UserBase
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" example:"2023-01-07T18:09:15.237672Z"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" example:"2023-01-07T18:09:15.237672Z"`
}

type Users struct {
	Total  int     `json:"total" example:"1000"`
	Limit  int     `json:"limit" example:"100"`
	Offset int     `json:"offset" example:"100"`
	Data   []*User `json:"data,omitempty"`
}

func (User) TableName() string {
	return UserTable
}
