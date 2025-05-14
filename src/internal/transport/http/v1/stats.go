package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *handler) initStatsRoutes(api *echo.Group) {
	api.GET("/stats/daily", h.DailyStats)
	api.GET("/stats/weekly", h.WeeklyStats)
	api.GET("/stats/monthly", h.MonthlyStats)
	api.GET("/stats/yearly", h.YearlyStats)
}

// @Summary Статистика за день
// @Description Поле "date" будет брать текущую дату, если его не заполнить
// @Tags stats
// @Produce json
// @Security BearerAuth
// @Param date query string false "Дата в формат: year-month-day"
// @Success 200 {array} domain.DailyStats
// @Failure 400
// @Failure 401
// @Router /api/v1/stats/daily [get]
func (h *handler) DailyStats(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	var date time.Time

	input := c.QueryParam("date")
	if len(input) > 0 {
		date, err = stringToTime(input)
		if err != nil {
			return c.NoContent(400)
		}
	} else {
		date = time.Now().UTC()
	}

	stats, err := h.service.Stats.GetDailyStats(userId, date)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

type weeklyStatsOutput struct {
	Data []struct {
		Day   string `json:"day"`
		Count int    `json:"count"`
	} `json:"data"`
}

// @Summary Статистика за неделю
// @Description Поле "date" будет брать текущую дату, если его не заполнить
// @Tags stats
// @Produce json
// @Security BearerAuth
// @Param date query string false "Дата в формат: year-month-day"
// @Success 200 {array} weeklyStatsOutput
// @Failure 400
// @Failure 401
// @Router /api/v1/stats/weekly [get]
func (h *handler) WeeklyStats(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	var date time.Time

	input := c.QueryParam("date")
	if len(input) > 0 {
		date, err = stringToTime(input)
		if err != nil {
			return c.NoContent(400)
		}
	} else {
		date = time.Now().UTC()
	}

	stats, err := h.service.Stats.GetWeeklyStats(userId, date)
	if err != nil {
		return err
	}

	output := weeklyStatsOutput{
		Data: make([]struct {
			Day   string `json:"day"`
			Count int    `json:"count"`
		}, len(stats.Data)),
	}

	for i, data := range stats.Data {
		output.Data[i].Day = data.Day.Format("2006-01-02")
		output.Data[i].Count = data.Count
	}

	return c.JSON(http.StatusOK, output)
}

type monthlyStatsOutput struct {
	Data []struct {
		Day   string `json:"day"`
		Count int    `json:"count"`
	} `json:"data"`
}

// @Summary Статистика за месяц
// @Description Поля "month" и "year" будут брать текущую дату, если они не заполнены
// @Tags stats
// @Produce json
// @Security BearerAuth
// @Param month query int false "Месяц"
// @Param year query int false "Год"
// @Success 200 {array} monthlyStatsOutput
// @Failure 400
// @Failure 401
// @Router /api/v1/stats/monthly [get]
func (h *handler) MonthlyStats(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	var (
		year  = 0
		month = 0
	)

	yearInput := c.QueryParam("year")
	monthInput := c.QueryParam("month")

	if len(yearInput) > 0 || len(monthInput) > 0 {
		year, err = getIntFromQuery(c, "year")
		if err != nil {
			return c.NoContent(400)
		}

		month, err = getIntFromQuery(c, "month")
		if err != nil {
			return c.NoContent(400)
		}
	} else {
		now := time.Now().UTC()
		year = now.Year()
		month = int(now.Month())
	}

	stats, err := h.service.Stats.GetMonthlyStats(userId, year, month)
	if err != nil {
		return err
	}

	output := monthlyStatsOutput{
		Data: make([]struct {
			Day   string `json:"day"`
			Count int    `json:"count"`
		}, len(stats.Data)),
	}

	for i, data := range stats.Data {
		output.Data[i].Day = data.Day.Format("2006-01-02")
		output.Data[i].Count = data.Count
	}

	return c.JSON(http.StatusOK, output)
}

// @Summary Статистика за год
// @Description Поле "year" будет брать текущую дату, если его не заполнить
// @Tags stats
// @Produce json
// @Security BearerAuth
// @Param year query int false "Год"
// @Success 200 {array} domain.YearlyStats
// @Failure 400
// @Failure 401
// @Router /api/v1/stats/yearly [get]
func (h *handler) YearlyStats(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	year := 0

	yearInput := c.QueryParam("year")

	if len(yearInput) > 0 {
		year, err = getIntFromQuery(c, "year")
		if err != nil {
			return c.NoContent(400)
		}
	} else {
		now := time.Now().UTC()
		year = now.Year()
	}

	stats, err := h.service.Stats.GetYearlyStats(userId, year)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

func stringToTime(date string) (time.Time, error) {
	result := time.Time{}

	result, err := time.Parse("2006-01-02", date)
	if err != nil {
		return result, errors.New("fail parse date")
	}

	return result, nil
}
