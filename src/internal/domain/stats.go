package domain

import "time"

type DailyStats struct {
	Data []struct {
		Hour  int `json:"hour"`
		Count int `json:"count"`
	} `json:"data"`
}

type WeeklyStats struct {
	Data []struct {
		Day   time.Time `json:"day"`
		Count int       `json:"count"`
	} `json:"data"`
}

type MonthlyStats struct {
	Data []struct {
		Day   time.Time `json:"day"`
		Count int       `json:"count"`
	} `json:"data"`
}

type YearlyStats struct {
	Data []struct {
		Month int `json:"month"`
		Count int `json:"count"`
	} `json:"data"`
}
