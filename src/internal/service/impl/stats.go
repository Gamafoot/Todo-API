package impl

import (
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"
	"time"
)

type statsService struct {
	config  *config.Config
	storage *storage.Storage
}

func NewStatsService(cfg *config.Config, storage *storage.Storage) *statsService {
	return &statsService{
		config:  cfg,
		storage: storage,
	}
}

func (s *statsService) GetDailyStats(userId uint, date time.Time) (*domain.DailyStats, error) {
	return s.storage.Stats.GetDailyStats(userId, date)
}

func (s *statsService) GetWeeklyStats(userId uint, date time.Time) (*domain.WeeklyStats, error) {
	return s.storage.Stats.GetWeeklyStats(userId, date)
}

func (s *statsService) GetMonthlyStats(userId uint, year, month int) (*domain.MonthlyStats, error) {
	return s.storage.Stats.GetMonthlyStats(userId, year, month)
}

func (s *statsService) GetYearlyStats(userId uint, year int) (*domain.YearlyStats, error) {
	return s.storage.Stats.GetYearlyStats(userId, year)
}
