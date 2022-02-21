package database

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"todo-svc/foundation/model"
)

func TestNewMemoryPersistence(t *testing.T) {
	mem := getPersistence(t)
	require.NotNil(t, mem.data, "Data map should have been initialized")
	require.Equal(t, uint64(1), mem.nextID, "First ID should be 1")
}

func TestMemoryPersistence_CreateTodo(t *testing.T) {
	tests := map[string]struct {
		newItem     *model.TodoItem
		expectedErr error
	}{
		"Nil item": {
			newItem:     nil,
			expectedErr: ErrBadInput,
		},
		"Happy path": {
			newItem: &model.TodoItem{
				Summary: "Some summary",
			},
			expectedErr: nil,
		},
		"ID replaced": {
			newItem: &model.TodoItem{
				Id:      5,
				Summary: "Some summary",
			},
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mem := getPersistence(t)
			created, err := mem.CreateTodo(tc.newItem)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, created)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, created)
				assert.Equal(t, uint64(1), created.Id)
			}
		})
	}
}

func TestMemoryPersistence_GetAllTodos(t *testing.T) {
	mem := getPersistence(t)

	first, err := mem.GetAllTodos()
	assert.NoError(t, err)
	assert.Len(t, first, 0)

	created, err := mem.CreateTodo(&model.TodoItem{
		Summary: "Some summary",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	second, err := mem.GetAllTodos()
	assert.NoError(t, err)
	assert.Len(t, second, 1)
	assert.Equal(t, created, second[0])
}

func TestMemoryPersistence_GetTodo(t *testing.T) {
	mem := getPersistence(t)

	first, err := mem.GetTodo(1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, first)

	created, err := mem.CreateTodo(&model.TodoItem{
		Summary: "Some summary",
	})
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, uint64(1), created.Id)

	second, err := mem.GetTodo(1)
	assert.NoError(t, err)
	assert.NotNil(t, second)
	assert.Equal(t, created, second)
}

func TestMemoryPersistence_UpdateTodo(t *testing.T) {
	mem := getPersistence(t)

	newState := &model.TodoItem{
		Id:      5,
		Summary: "Another Summary",
		Done:    true,
	}

	err := mem.UpdateTodo(1, nil)
	assert.ErrorIs(t, err, ErrBadInput)

	err = mem.UpdateTodo(1, newState)
	assert.ErrorIs(t, err, ErrNotFound)

	created, err := mem.CreateTodo(&model.TodoItem{
		Summary: "Some summary",
	})
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, uint64(1), created.Id)

	err = mem.UpdateTodo(1, newState)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), newState.Id)
}

func TestMemoryPersistence_DeleteTodo(t *testing.T) {
	mem := getPersistence(t)

	err := mem.DeleteTodo(1)
	assert.ErrorIs(t, err, ErrNotFound)

	err = mem.DeleteTodo(0)
	assert.ErrorIs(t, err, ErrBadInput)

	created, err := mem.CreateTodo(&model.TodoItem{
		Summary: "Some summary",
	})
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, uint64(1), created.Id)

	err = mem.DeleteTodo(1)
	assert.NoError(t, err)

	retrieved, err := mem.GetTodo(1)
	assert.Nil(t, retrieved)
	assert.ErrorIs(t, err, ErrNotFound)
}

func getPersistence(t *testing.T) *MemoryPersistence {
	mem, err := NewMemoryPersistence()
	require.NoError(t, err)
	return mem
}
