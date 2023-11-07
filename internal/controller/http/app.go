package http

import (
	"context"
	"goapptemplate/internal/domain"
	"goapptemplate/internal/usecase"
	"time"

	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type AppHTTPController interface {
	CreateBook() func(*fiber.Ctx) error
	GetBooks() func(*fiber.Ctx) error
	GetBook() func(*fiber.Ctx) error
	UpdateBook() func(*fiber.Ctx) error
	DeleteBook() func(*fiber.Ctx) error
}

type AppHTTPControllerConfig struct {
	Timeout time.Duration
}

type appHTTPController struct {
	f      *fiber.App
	books  usecase.Books
	config *AppHTTPControllerConfig
	log    *logrus.Entry
}

// CreateBook implements AppHTTPController.
func (hc *appHTTPController) CreateBook() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		book := &domain.Book{}
		err := c.BodyParser(book)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		ctx, cancel := context.WithTimeout(context.Background(), hc.config.Timeout)
		defer cancel()
		b, err := hc.books.New(ctx, book)
		if err != nil {
			if errors.Is(err, domain.ErrValidation) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
			}
			hc.log.WithFields(structs.Map(book)).Error("cannot add new book")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
		}
		c.Append(fiber.HeaderLocation, c.Path()+"/"+b.ID.String())
		return c.SendStatus(fiber.StatusCreated)
	}
}

// DeleteBook implements AppHTTPController.
func (hc *appHTTPController) DeleteBook() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bookID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		ctx, cancel := context.WithTimeout(context.Background(), hc.config.Timeout)
		defer cancel()
		err = hc.books.Remove(ctx, bookID)
		if err != nil {
			hc.log.WithField("book_id", bookID).Error("cannot remove book")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// GetBook implements AppHTTPController.
func (hc *appHTTPController) GetBook() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bookID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		ctx, cancel := context.WithTimeout(context.Background(), hc.config.Timeout)
		defer cancel()
		b, err := hc.books.View(ctx, bookID)
		if err != nil {
			if errors.Is(err, domain.ErrBookNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
			}
			hc.log.WithField("book_id", bookID).Error("cannot view book")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
		}
		return c.Status(fiber.StatusOK).JSON(b)
	}
}

// GetBooks implements AppHTTPController.
func (hc *appHTTPController) GetBooks() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), hc.config.Timeout)
		defer cancel()
		filters := new(domain.BookFilters)
		err := c.QueryParser(filters)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		p, err := hc.books.List(ctx, filters)
		if err != nil {
			if errors.Is(err, domain.ErrValidation) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
			}
			hc.log.WithFields(structs.Map(filters)).Error("cannot list books")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
		}
		return c.Status(fiber.StatusOK).JSON(p)
	}
}

// UpdateBook implements AppHTTPController.
func (hc *appHTTPController) UpdateBook() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bookID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		book := &domain.Book{ID: bookID}
		err = c.BodyParser(book)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
		}
		ctx, cancel := context.WithTimeout(context.Background(), hc.config.Timeout)
		defer cancel()
		_, err = hc.books.Modify(ctx, book)
		if err != nil {
			if errors.Is(err, domain.ErrValidation) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.NewError(fiber.StatusBadRequest, err.Error()))
			}
			if errors.Is(err, domain.ErrBookNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
			}
			hc.log.WithFields(structs.Map(book)).Error("cannot modify book")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
		}
		c.Append(fiber.HeaderLocation, c.Path())
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func NewAppHTTPController(
	f *fiber.App,
	bu usecase.Books,
	config *AppHTTPControllerConfig,
	logger *logrus.Logger,
) AppHTTPController {
	hc := &appHTTPController{
		f:      f,
		books:  bu,
		config: config,
		log:    logger.WithField("layer", "internal.controller.http.appHTTPController"),
	}
	books := hc.f.Group("/books")
	books.Post("", hc.CreateBook())
	books.Get("/:id", hc.GetBook())
	books.Get("", hc.GetBooks())
	books.Put("/:id", hc.UpdateBook())
	books.Delete("/:id", hc.DeleteBook())

	return hc
}
