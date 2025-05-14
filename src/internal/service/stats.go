package service

import (
	"root/internal/domain"
	"time"
)

type StatsService interface {
	GetDailyStats(userId uint, date time.Time) (*domain.DailyStats, error)
	GetWeeklyStats(userId uint, date time.Time) (*domain.WeeklyStats, error)
	GetMonthlyStats(userId uint, year, month int) (*domain.MonthlyStats, error)
	GetYearlyStats(userId uint, year int) (*domain.YearlyStats, error)
}
