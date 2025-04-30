package domain

type ProjectStats struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Overdue int `json:"overdue"`
}
