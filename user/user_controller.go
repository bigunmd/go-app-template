package user

import (
	"app/domain/model"
	"app/domain/usecase"
	"app/pkg/logger"
	"app/pkg/validator"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/araddon/dateparse"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserController interface {
	RegisterRoutes()
	GetUsers() fiber.Handler
	CreateUser() fiber.Handler
	GetUser() fiber.Handler
	UpdateUser() fiber.Handler
	DeleteUser() fiber.Handler
}

type controller struct {
	log         logger.Logger
	router      fiber.Router
	userUsecase usecase.UserUsecase
}

// CreateUser godoc
//
//	@Summary		Create new user
//	@Description	Create new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.UserBase	true	"create User request"
//	@Success		201
//	@Failure		400	{object}	validator.ErrorResponse
//	@Failure		500	{object}	fiber.Error
//	@Header			201	{string}	Location	"/users/:id"
//	@Router			/users [post]
func (c *controller) CreateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		u := new(model.UserBase)
		err := ctx.BodyParser(u)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		errors := validator.ValidateStruct(u)
		if len(errors.FieldError) != 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(errors)
		}
		usr, err := c.userUsecase.Create(u)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		ctx.Location(fmt.Sprintf("%s/%s", ctx.Path(), usr.Id.String()))
		return ctx.SendStatus(fiber.StatusCreated)
	}
}

// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Delete user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"user id"
//	@Success		204
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Failure		500	{object}	fiber.Error
//	@Router			/users/{id} [delete]
func (c *controller) DeleteUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		err = c.userUsecase.DeleteById(id)
		if err != nil {
			if err == sql.ErrNoRows {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
			} else {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
			}
		}
		return ctx.SendStatus(fiber.StatusNoContent)
	}
}

// GetUser godoc
//
//	@Summary		Get user
//	@Description	Get user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"user id"
//	@Success		200	{object}	model.User
//	@Failure		400	{object}	fiber.Error
//	@Failure		404	{object}	fiber.Error
//	@Failure		500	{object}	fiber.Error
//	@Router			/users/{id} [get]
func (c *controller) GetUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		u, err := c.userUsecase.GetById(id)
		if err != nil {
			if err == sql.ErrNoRows {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		return ctx.Status(fiber.StatusOK).JSON(u)
	}
}

// GetUsers godoc
//
//	@Summary		Get users
//	@Description	Get users with options
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			limit		query		int	false	"limit"
//	@Param			offset		query		int	false	"offset"
//	@Param			lt_datetime	query		int	false	"less than unix timestamp"
//	@Param			gt_datetime	query		int	false	"greater than unix timestamp"
//	@Success		200			{object}	model.Users
//	@Failure		500			{object}	fiber.Error
//	@Router			/users [get]
func (c *controller) GetUsers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		opts := new(model.UserOptions)
		opts.Limit, _ = strconv.Atoi(ctx.Query("limit", "25"))
		if opts.Limit > 100 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, "limit cannot be higher than 100"))
		}
		opts.Offset, _ = strconv.Atoi(ctx.Query("offset"))
		opts.LtDatetime, _ = dateparse.ParseAny(ctx.Query("lt_datetime"))
		opts.GtDatetime, _ = dateparse.ParseAny(ctx.Query("gt_datetime"))

		us, err := c.userUsecase.Get(opts)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		return ctx.Status(fiber.StatusOK).JSON(us)
	}
}

// UpdateUser godoc
//
//	@Summary		Update user
//	@Description	Update existing user or create a new one
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.User	true	"Update User request"\
//	@Param			id		path	string		true	"user id"
//	@Success		201
//	@Failure		400	{object}	validator.ErrorResponse
//	@Failure		500	{object}	fiber.Error
//	@Header			201	{string}	Location	"/users/:id"
//	@Router			/users/{id} [put]
func (c *controller) UpdateUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := uuid.Parse(ctx.Params("id"))
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		u := new(model.User)
		err = ctx.BodyParser(u)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		u.Id = id
		errors := validator.ValidateStruct(u)
		if len(errors.FieldError) != 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(errors)
		}
		_, err = c.userUsecase.Update(u)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
		ctx.Location(ctx.Path())
		return ctx.SendStatus(fiber.StatusCreated)
	}
}

// RegisterRoutes implements UserController
func (c *controller) RegisterRoutes() {
	c.router.Get("/users", c.GetUsers())
	c.router.Post("/users", c.CreateUser())
	c.router.Get("/users/:id", c.GetUser())
	c.router.Put("/users/:id", c.UpdateUser())
	c.router.Delete("/users/:id", c.DeleteUser())
}

func NewUserController(router fiber.Router, db *sqlx.DB, logger logger.Logger) UserController {
	userUsecase := NewUserUsecase(db, logger)
	return &controller{
		log:         logger,
		router:      router,
		userUsecase: userUsecase,
	}
}
