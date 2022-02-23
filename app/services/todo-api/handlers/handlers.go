// Package handlers provides HTTP handlers to implement the REST service.
package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovery "github.com/gofiber/fiber/v2/middleware/recover"
	"golang.org/x/crypto/bcrypt"
	"log"
	"todo-svc/app/services/todo-api/config"
	"todo-svc/business/database"
	"todo-svc/foundation/model"
	"todo-svc/zarf/keys"
)

const (
	badRequestBody         = "Unrecognized request body or content type"
	badRequestBodyCreation = "Invalid request body for resource creation"
	badRequestId           = "Invalid ID format"
	resourceNotFound       = "TodoItem not found"
)

func SetupHandlers(logg *log.Logger, cfg config.Config, pers database.Persistence) (*fiber.App, error) {
	a := fiber.New(fiber.Config{
		ServerHeader: cfg.Version.Desc,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	})

	api := a.Group("/api/v1", recovery.New(), logger.New(logger.Config{
		TimeZone: "UTC",
		Output:   logg.Writer(),
	}), basicauth.New(basicauth.Config{
		Realm: "todo-api",
		Authorizer: func() func(string, string) bool {
			devCreds, err := keys.GetDevAuth()
			if err != nil {
				panic("failed to get development credential store")
			}
			return func(user string, pass string) bool {
				var hash []byte
				var ok bool
				if cfg.Auth.Mode == "dev" {
					logg.Println("Checking against DEV credentials")
					hash, ok = devCreds[user]
					if !ok {
						return false
					}
				} else {
					//TODO: Provide real auth source
					panic("unimplemented")
				}
				if err := bcrypt.CompareHashAndPassword(hash, []byte(pass)); err != nil {
					return false
				}
				return true
			}
		}(),
	}))

	todos := api.Group("/todo")
	todos.Post("/", func(ctx *fiber.Ctx) error {
		newTodo := new(model.TodoItem)
		if err := ctx.BodyParser(newTodo); err != nil {
			return ctx.
				Status(fiber.StatusBadRequest).
				SendString(badRequestBody)
		}
		created, err := pers.CreateTodo(newTodo)
		if err != nil {
			if errors.Is(err, database.ErrBadInput) {
				return ctx.
					Status(fiber.StatusBadRequest).
					SendString(badRequestBodyCreation)
			}
			return err
		}
		return ctx.Status(fiber.StatusCreated).JSON(created)
	})
	todos.Get("/", func(ctx *fiber.Ctx) error {
		allTodos, err := pers.GetAllTodos()
		if err != nil {
			return err
		}
		return ctx.Status(fiber.StatusOK).JSON(allTodos)
	})
	todos.Get("/:id", func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id", 0)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(badRequestId)
		}
		todo, err := pers.GetTodo(uint64(id))
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ctx.Status(404).SendString(resourceNotFound)
			}
			return err
		}
		return ctx.Status(fiber.StatusOK).JSON(todo)
	})
	todos.Put("/:id", func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id", 0)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(badRequestId)
		}
		body := new(model.TodoItem)
		if err := ctx.BodyParser(body); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(badRequestBody)
		}
		if err := pers.UpdateTodo(uint64(id), body); err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(resourceNotFound)
			}
			return err
		}
		ctx.Status(fiber.StatusOK)
		return nil
	})
	todos.Delete("/:id", func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id", 0)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(badRequestId)
		}
		if err := pers.DeleteTodo(uint64(id)); err != nil {
			if !errors.Is(err, database.ErrNotFound) {
				return err
			}
		}
		ctx.Status(fiber.StatusOK)
		return nil
	})

	return a, nil
}
