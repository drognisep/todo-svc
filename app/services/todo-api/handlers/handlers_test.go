package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
	"todo-svc/app/services/todo-api/config"
	"todo-svc/business/database"
	"todo-svc/foundation/model"
)

func TestGetAllTodos(t *testing.T) {
	app, pers := setupHandlers(t)
	req, err := http.NewRequest("GET", "/api/v1/todo", nil)
	assert.NoError(t, err)
	req.SetBasicAuth("bob", "bob")

	{
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, resp.Header.Get(fiber.HeaderContentType))
		buf := []*model.TodoItem{}
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&buf))
		assert.Len(t, buf, 0)
	}

	_, _ = pers.CreateTodo(&model.TodoItem{
		Summary: "Some summary",
	})

	{
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, resp.Header.Get(fiber.HeaderContentType))
		buf := []*model.TodoItem{}
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&buf))
		assert.Len(t, buf, 1)
	}
}

func TestGetSingleTodo(t *testing.T) {
	app, pers := setupHandlers(t)
	req, err := http.NewRequest("GET", "/api/v1/todo/1", nil)
	assert.NoError(t, err)
	req.SetBasicAuth("bob", "bob")

	{
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	}

	created, err := pers.CreateTodo(&model.TodoItem{
		Summary: "New item",
	})
	assert.NoError(t, err)

	{
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		body := new(model.TodoItem)
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(body))
		assert.Equal(t, created, body)
	}
}

func TestCreateTodo(t *testing.T) {
	app, pers := setupHandlers(t)
	toCreate := &model.TodoItem{
		Id:      1,
		Summary: "Created",
	}
	buf, err := json.Marshal(toCreate)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", "/api/v1/todo", bytes.NewReader(buf))
	assert.NoError(t, err)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.SetBasicAuth("bob", "bob")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	body := new(model.TodoItem)
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(body))
	assert.Equal(t, toCreate, body)
	found, err := pers.GetTodo(1)
	assert.NoError(t, err)
	assert.Equal(t, toCreate, found)
}

func TestUpdateTodo(t *testing.T) {
	app, pers := setupHandlers(t)
	created := &model.TodoItem{
		Summary: "Created",
	}
	update := &model.TodoItem{
		Summary: "Updated",
		Done:    true,
	}
	updateBody, err := json.Marshal(update)
	assert.NoError(t, err)
	req, err := http.NewRequest("PUT", "/api/v1/todo/1", bytes.NewReader(updateBody))
	assert.NoError(t, err)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.SetBasicAuth("bob", "bob")

	{
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	}

	_, _ = pers.CreateTodo(created)

	{
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}
}

func TestDeleteTodo(t *testing.T) {
	app, pers := setupHandlers(t)
	created := &model.TodoItem{
		Summary: "Created",
	}

	req, err := http.NewRequest("DELETE", "/api/v1/todo/1", nil)
	assert.NoError(t, err)
	req.SetBasicAuth("bob", "bob")

	{
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}

	_, err = pers.CreateTodo(created)
	assert.NoError(t, err)

	{
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}
}

func setupHandlers(t *testing.T) (*fiber.App, database.Persistence) {
	logg := log.New(os.Stderr, "Testing!!! ", log.LstdFlags)
	cfg := config.Config{
		Web: config.WebConfig{
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			IdleTimeout:     10 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		Auth: config.AuthConfig{
			AuthMode: "dev",
		},
	}
	pers, err := database.NewMemoryPersistence()
	require.NoError(t, err)

	app, err := SetupHandlers(logg, cfg, pers)
	require.NoError(t, err)

	return app, pers
}
