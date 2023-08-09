package tasks

// Структура с телом задач
type Entity struct {
	ID       string  `db:"id"`
	Title    *string `db:"title"`
	ActiveAt *string `db:"active_at"`
	Status   *bool   `db:"status"`
}
