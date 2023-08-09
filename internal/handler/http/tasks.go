package http

import (
	"net/http"
	"product-service/internal/domain/tasks"
	"product-service/internal/service/catalogue"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"product-service/pkg/server/response"
	"product-service/pkg/store"
)

type TasksHandler struct {
	tasksService *catalogue.Service
}

func NewTasksHandler(s *catalogue.Service) *TasksHandler {
	return &TasksHandler{tasksService: s}
}

func (h *TasksHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)

	r.Route("/{id}", func(r chi.Router) {
		r.Put("/done", h.status)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return r
}

// List of products from the database
//
//	@Summary	List of tasks from the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Success	200		{array}		response.Object
//	@Failure	500		{object}	response.Object
//	@Router		/tasks 	[get]
func (h *TasksHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.tasksService.ListTasks(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// Add a new product to the database
//
//	@Summary	Add a new tasks to the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		request	body		tasks.Request	true	"body param"
//	@Success	200		{object}	response.Object
//	@Failure	400		{object}	response.Object
//	@Failure	500		{object}	response.Object
//	@Router		/tasks [post]
func (h *TasksHandler) add(w http.ResponseWriter, r *http.Request) {
	req := tasks.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	res, err := h.tasksService.AddTasks(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// Read the product from the database
//
//	@Summary	Read the tasks from the database
//	@Tags		tasks
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"path param"
//	@Success	200	{object}	response.Object
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/tasks/{id}/done [put]
func (h *TasksHandler) status(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := tasks.Request{
		Status: true, // Здесь укажите новое значение статуса
	}

	err := h.tasksService.GetStatus(r.Context(), id, req)
	if err != nil && err != store.ErrorNotFound {
		response.InternalServerError(w, r, err)
		return
	}

	if err == store.ErrorNotFound {
		response.NotFound(w, r, err)
		return
	}
}

// Update the product in the database
//
//	@Summary	Update the tasks in the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string			true	"path param"
//	@Param		request	body	tasks.Request	true	"body param"
//	@Success	200
//	@Failure	400	{object}	response.Object
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/tasks/{id} [put]
func (h *TasksHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := tasks.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	err := h.tasksService.UpdateTasks(r.Context(), id, req)
	if err != nil && err != store.ErrorNotFound {
		response.InternalServerError(w, r, err)
		return
	}

	if err == store.ErrorNotFound {
		response.NotFound(w, r, err)
		return
	}
}

// Delete the product from the database
//
//	@Summary	Delete the tasks from the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"path param"
//	@Success	200
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/tasks/{id} [delete]
func (h *TasksHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.tasksService.DeleteTasks(r.Context(), id)
	if err != nil && err != store.ErrorNotFound {
		response.InternalServerError(w, r, err)
		return
	}

	if err == store.ErrorNotFound {
		response.NotFound(w, r, err)
		return
	}
}
