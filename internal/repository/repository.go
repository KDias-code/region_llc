package repository

import (
	"product-service/internal/domain/tasks"
	"product-service/internal/repository/memory"
	"product-service/internal/repository/postgres"
	"product-service/pkg/store"
)

// Конфигурация — это псевдоним для функции, которая принимает указатель на репозиторий и изменяет его.
type Configuration func(r *Repository) error

// Репозиторий — это реализация репозитория.
type Repository struct {
	postgres *store.Database

	Tasks tasks.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	// создаем репозиторий
	s = &Repository{}

	// Применить все переданные конфигурации
	for _, cfg := range configs {
		// Передайте репозиторий в функцию конфигурации
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

// Close закрывает репозиторий и предотвращает запуск новых запросов.
// Close затем ожидает завершения всех запросов, которые начали обрабатываться на сервере.
func (r *Repository) Close() {
	if r.postgres != nil {
		r.postgres.Client.Close()
	}
}

// WithMemoryStore применяет хранилище памяти к репозиторию
func WithMemoryStore() Configuration {
	return func(s *Repository) (err error) {
		// Create the memory store, if we needed parameters, such as connection strings they could be inputted here
		s.Tasks = memory.NewTasksRepository()

		return
	}
}

// WithPostgresStore применяет хранилище postgres к репозиторию
func WithPostgresStore(schema, dataSourceName string) Configuration {
	return func(s *Repository) (err error) {
		// Create the postgres store, if we needed parameters, such as connection strings they could be inputted here
		s.postgres, err = store.NewDatabase(schema, dataSourceName)
		if err != nil {
			return
		}

		err = s.postgres.Migrate()
		if err != nil {
			return
		}

		s.Tasks = postgres.NewTasksRepository(s.postgres.Client)
		return
	}
}
