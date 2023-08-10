package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger/v2"
	"net/url"
	"product-service/docs"
	_ "product-service/docs"
	"product-service/internal/config"
	"product-service/internal/handler/http"
	"product-service/internal/service/catalogue"
	"product-service/pkg/server/router"
)

type Dependencies struct {
	Configs          config.Configs
	CatalogueService *catalogue.Service
}

// Конфигурация — это псевдоним функции, которая принимает указатель на обработчик и изменяет его.
type Configuration func(h *Handler) error

// Хендлер является реализацией хендлера
type Handler struct {
	dependencies Dependencies

	HTTP *chi.Mux
}

// New принимает переменное количество функций конфигурации и возвращает новый обработчик.
// Каждая конфигурация будет вызываться в порядке их передачи.
func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	// создаем хендлер
	h = &Handler{
		dependencies: d,
	}

	// Применить все переданные конфигурации
	for _, cfg := range configs {
		// Передать сервис в функцию конфигурации
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

// WithHTTPHandler применяет обработчик http к обработчику
func WithHTTPHandler() Configuration {
	return func(h *Handler) (err error) {
		// Создайте обработчик http, если нам нужны параметры, такие как строки подключения, их можно ввести здесь.
		h.HTTP = router.New()

		docs.SwaggerInfo.BasePath = "/api/todo-list"
		docs.SwaggerInfo.Host = h.dependencies.Configs.HTTP.Host
		docs.SwaggerInfo.Schemes = []string{h.dependencies.Configs.HTTP.Schema}

		swaggerURL := url.URL{
			Scheme: h.dependencies.Configs.HTTP.Schema,
			Host:   h.dependencies.Configs.HTTP.Host,
			Path:   "swagger/doc.json",
		}

		h.HTTP.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(swaggerURL.String()),
		))

		tasksHandler := http.NewTasksHandler(h.dependencies.CatalogueService)

		h.HTTP.Route("/api/todo-list", func(r chi.Router) {
			r.Mount("/tasks", tasksHandler.Routes())
		})

		return
	}
}
