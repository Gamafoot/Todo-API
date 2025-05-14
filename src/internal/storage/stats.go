package storage

import (
	"root/internal/domain"
	"time"
)

type StatsStorage interface {
	GetDailyStats(userId uint, date time.Time) (*domain.DailyStats, error)
	GetWeeklyStats(userId uint, date time.Time) (*domain.WeeklyStats, error)
	GetMonthlyStats(userId uint, year, month int) (*domain.MonthlyStats, error)
	GetYearlyStats(userId uint, year int) (*domain.YearlyStats, error)
}
