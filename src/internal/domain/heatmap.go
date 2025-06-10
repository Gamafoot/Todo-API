package domain

import "time"

type Heatmap struct {
	Day   time.Time `json:"day"`
	Count int       `json:"count"`
}
