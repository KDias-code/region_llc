package tasks

import (
	"errors"
	"net/http"
)

// Запросы, ответы , парсинги
type Request struct {
	Title    string `json:"title"`
	ActiveAt string `json:"active_at"`
	Status   bool   `json:"status"`
}

// Проверка на пустоту
func (s *Request) Bind(r *http.Request) error {
	if s.Title == "" {
		return errors.New("title: cannot be blank")
	}

	if s.ActiveAt == "" {
		return errors.New("active_at: cannot be blank")
	}

	return nil
}

type Response struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ActiveAt string `json:"active_at"`
	Status   bool   `json:"status"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:       data.ID,
		Title:    *data.Title,
		ActiveAt: *data.ActiveAt,
	}

	if data.Status != nil {
		res.Status = *data.Status
	}
	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
