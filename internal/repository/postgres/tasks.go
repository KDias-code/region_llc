package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"

	"product-service/internal/domain/tasks"
)

// реализация всех методов с помощью постгрес
type TasksRepository struct {
	db *sqlx.DB
}

func NewTasksRepository(db *sqlx.DB) tasks.Repository {
	return &TasksRepository{
		db: db,
	}
}

func (s *TasksRepository) Select(ctx context.Context) (dest []tasks.Entity, err error) {
	query := `
		SELECT id, title, status, active_at
		FROM tasks
		WHERE DATE(active_at) = DATE($1)
		ORDER BY id`

	currentDate := time.Now().Format("2006-01-02")

	err = s.db.SelectContext(ctx, &dest, query, currentDate)

	return
}

func (s *TasksRepository) Create(ctx context.Context, data tasks.Entity) (id string, err error) {
	query := `
		INSERT INTO tasks (title, active_at)
		VALUES ($1, $2)
		RETURNING id`

	args := []any{data.Title, data.ActiveAt}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *TasksRepository) Status(ctx context.Context, id string, data tasks.Entity) (err error) {
	query := `
		UPDATE tasks
		SET status = true
		WHERE id = $1
	`

	_, err = s.db.ExecContext(ctx, query, id)
	trueValue := true
	data.Status = &trueValue
	return err
}

func (s *TasksRepository) Update(ctx context.Context, id string, data tasks.Entity) (err error) {
	sets, args := s.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE tasks SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))
		_, err = s.db.ExecContext(ctx, query, args...)
	}

	return
}

func (s *TasksRepository) prepareArgs(data tasks.Entity) (sets []string, args []any) {
	if data.Title != nil {
		args = append(args, data.Title)
		sets = append(sets, fmt.Sprintf("title=$%d", len(args)))
	}

	if data.ActiveAt != nil {
		args = append(args, data.ActiveAt)
		sets = append(sets, fmt.Sprintf("active_at=$%d", len(args)))
	}

	return
}

func (s *TasksRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE 
		FROM tasks
		WHERE id=$1`

	args := []any{id}

	_, err = s.db.ExecContext(ctx, query, args...)

	return
}
