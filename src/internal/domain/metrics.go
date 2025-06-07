package domain

import "time"

type PreMetrics struct {
	TotalTasks  int `json:"total_tasks"`
	DoneTasks   int `json:"done_tasks"`
	RemTasks    int `json:"rem_tasks"`
	DaysElapsed int `json:"days_elapsed"`
	DaysLeft    int `json:"days_left"`
}

type Metrics struct {
	TotalTasks          int       `json:"total_tasks"`
	DoneTasks           int       `json:"done_tasks"`
	RemTasks            int       `json:"rem_tasks"`
	DaysElapsed         int       `json:"days_elapsed"`
	DaysLeft            int       `json:"days_left"`
	VReal               float64   `json:"v_real"`
	VReq                float64   `json:"v_req"`
	PerceptionDone      int       `json:"perception_done"`
	ProjectedFinishDate time.Time `json:"projected_finish_date"`
	Status              string    `json:"status"`
}
