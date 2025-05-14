package postgres

import (
	"root/internal/domain"
	"time"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type statsStorage struct {
	db *gorm.DB
}

func NewStatsStorage(db *gorm.DB) *statsStorage {
	return &statsStorage{db: db}
}

func (s *statsStorage) GetDailyStats(userId uint, date time.Time) (*domain.DailyStats, error) {
	stats := new(domain.DailyStats)

	err := s.db.Raw("SELECT * FROM daily_stats(?, ?)", userId, date).Scan(&stats.Data).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return stats, nil
}

func (s *statsStorage) GetWeeklyStats(userId uint, date time.Time) (*domain.WeeklyStats, error) {
	stats := new(domain.WeeklyStats)

	err := s.db.Raw("SELECT * FROM weekly_stats(?, ?)", userId, date).Scan(&stats.Data).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return stats, nil
}

func (s *statsStorage) GetMonthlyStats(userId uint, year, month int) (*domain.MonthlyStats, error) {
	stats := new(domain.MonthlyStats)

	err := s.db.Raw("SELECT * FROM monthly_stats(?, ?, ?)", userId, year, month).Scan(&stats.Data).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return stats, nil
}

func (s *statsStorage) GetYearlyStats(userId uint, year int) (*domain.YearlyStats, error) {
	stats := new(domain.YearlyStats)

	err := s.db.Raw("SELECT * FROM yearly_stats(?, ?)", userId, year).Scan(&stats.Data).Error
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return stats, nil
}
