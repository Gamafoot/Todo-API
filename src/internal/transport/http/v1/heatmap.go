package v1

import (
	"github.com/labstack/echo/v4"
)

func (h *handler) initHeatmapRoutes(api *echo.Group) {
	api.GET("/heatmap", h.HeatmapActivity)
}

// @Summary Heatmap активность
// @Tags heatmap
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Heatmap
// @Failure 401
// @Router /api/v1/heatmap [get]
func (h *handler) HeatmapActivity(c echo.Context) error {
	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	heatmaps, err := h.service.Heatmap.GetActivity(userId)
	if err != nil {
		return err
	}

	return c.JSON(200, heatmaps)
}
