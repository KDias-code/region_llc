package catalogue

import (
	"context"
	"product-service/internal/domain/tasks"
)

func (s *Service) ListTasks(ctx context.Context) (res []tasks.Response, err error) {
	data, err := s.tasksRepository.Select(ctx)
	if err != nil {
		return
	}
	res = tasks.ParseFromEntities(data)

	return
}

func (s *Service) AddTasks(ctx context.Context, req tasks.Request) (res tasks.Response, err error) {
	data := tasks.Entity{
		Title:    &req.Title,
		ActiveAt: &req.ActiveAt,
	}

	data.ID, err = s.tasksRepository.Create(ctx, data)
	if err != nil {
		return
	}
	res = tasks.ParseFromEntity(data)

	return
}

func (s *Service) GetStatus(ctx context.Context, id string, req tasks.Request) (err error) {
	data := tasks.Entity{
		Status: &req.Status,
	}

	err = s.tasksRepository.Status(ctx, id, data)
	if err != nil {
		return
	}

	return
}

func (s *Service) UpdateTasks(ctx context.Context, id string, req tasks.Request) (err error) {
	data := tasks.Entity{
		Title:    &req.Title,
		ActiveAt: &req.ActiveAt,
	}

	err = s.tasksRepository.Update(ctx, id, data)
	if err != nil {
		return
	}

	return
}

func (s *Service) DeleteTasks(ctx context.Context, id string) (err error) {
	err = s.tasksRepository.Delete(ctx, id)
	if err != nil {
		return
	}

	return
}
