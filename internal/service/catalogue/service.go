package catalogue

import (
	"product-service/internal/domain/tasks"
)

// Конфигурация — это псевдоним функции, которая принимает указатель на службу и изменяет ее.
type Configuration func(s *Service) error

// Реализация сервиса
type Service struct {
	tasksRepository tasks.Repository
	tasksCache      tasks.Cache
}

func New(configs ...Configuration) (s *Service, err error) {
	// создаем сервис
	s = &Service{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}
	return
}

// WithCategoryRepository применяет данный репозиторий категорий к Сервису
func WithTasksRepository(tasksRepository tasks.Repository) Configuration {
	// вернуть функцию, соответствующую псевдониму конфигурации,
	// Вам нужно вернуть это, чтобы родительская функция могла принимать все необходимые параметры
	return func(s *Service) error {
		s.tasksRepository = tasksRepository
		return nil
	}
}

// WithCategoryCache применяет заданный кеш категории к сервису
func WithTasksCache(tasksCache tasks.Cache) Configuration {
	return func(s *Service) error {
		s.tasksCache = tasksCache
		return nil
	}
}
