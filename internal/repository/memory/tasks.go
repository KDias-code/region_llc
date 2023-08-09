package memory

import (
	"context"
	"database/sql"
	"product-service/pkg/store"
	"sync"

	"github.com/google/uuid"

	"product-service/internal/domain/tasks"
)

type TasksRepository struct {
	db map[string]tasks.Entity
	sync.RWMutex
}

func NewTasksRepository() *TasksRepository {
	return &TasksRepository{
		db: make(map[string]tasks.Entity),
	}
}

func (r *TasksRepository) Select(ctx context.Context) (dest []tasks.Entity, err error) {
	r.RLock()
	defer r.RUnlock()

	dest = make([]tasks.Entity, 0, len(r.db))
	for _, data := range r.db {
		dest = append(dest, data)
	}

	return
}

func (r *TasksRepository) Create(ctx context.Context, data tasks.Entity) (dest string, err error) {
	r.Lock()
	defer r.Unlock()

	id := r.generateID()
	data.ID = id
	r.db[id] = data

	return id, nil
}

func (r *TasksRepository) Status(ctx context.Context, id string, data tasks.Entity) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	src := r.db[id]

	if data.Status != nil {
		src.Status = data.Status
	}

	r.db[id] = src

	return
}

func (r *TasksRepository) Delete(ctx context.Context, id string) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return sql.ErrNoRows
	}
	delete(r.db, id)

	return
}

func (r *TasksRepository) Update(ctx context.Context, id string, data tasks.Entity) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.db[id]; !ok {
		return store.ErrorNotFound
	}
	src := r.db[id]

	if data.Title != nil {
		src.Title = data.Title
	}

	if data.ActiveAt != nil {
		src.ActiveAt = data.ActiveAt
	}

	if data.Status != nil {
		src.Status = data.Status
	}

	r.db[id] = src

	return
}

func (r *TasksRepository) generateID() string {
	return uuid.New().String()
}
