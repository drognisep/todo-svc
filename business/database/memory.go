package database

import (
	"fmt"
	"sync"
	"todo-svc/foundation/model"
)

var _ Persistence = (*MemoryPersistence)(nil)

type MemoryPersistence struct {
	mux    sync.RWMutex
	data   map[uint64]*model.TodoItem
	nextID uint64
}

// NewMemoryPersistence creates a new Persistence implementation that stores data in memory.
func NewMemoryPersistence() (*MemoryPersistence, error) {
	return &MemoryPersistence{
		data:   map[uint64]*model.TodoItem{},
		nextID: 1,
	}, nil
}

func (m *MemoryPersistence) CreateTodo(item *model.TodoItem) (*model.TodoItem, error) {
	if item == nil {
		return nil, fmt.Errorf("nil item passed to CreateTodo: %w", ErrBadInput)
	}

	m.mux.Lock()
	defer m.mux.Unlock()
	newID := m.nextID
	m.nextID += 1

	item.Id = newID
	m.data[newID] = item
	return item, nil
}

func (m *MemoryPersistence) GetAllTodos() ([]*model.TodoItem, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	if len(m.data) == 0 {
		return nil, nil
	}

	buf := make([]*model.TodoItem, len(m.data))
	idx := 0
	for _, v := range m.data {
		buf[idx] = v
		idx++
	}
	return buf, nil
}

func (m *MemoryPersistence) GetTodo(id uint64) (*model.TodoItem, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	item, ok := m.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return item, nil
}

func (m *MemoryPersistence) UpdateTodo(id uint64, newState *model.TodoItem) error {
	if newState == nil {
		return fmt.Errorf("new state cannot be nil: %w", ErrBadInput)
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	newState.Id = id
	_, ok := m.data[id]
	if !ok {
		return ErrNotFound
	}
	m.data[id] = newState
	return nil
}

func (m *MemoryPersistence) DeleteTodo(id uint64) error {
	if id == 0 {
		return fmt.Errorf("0 is not a valid ID: %w", ErrBadInput)
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	_, ok := m.data[id]
	if !ok {
		return ErrNotFound
	}
	delete(m.data, id)
	return nil
}
