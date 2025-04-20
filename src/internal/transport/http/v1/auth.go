package v1

import (
	"errors"
	"net/http"
	"root/internal/domain"
	customErrors "root/pkg/errors"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *handler) initAuthRoutes(api *echo.Group) {
	api.POST("/auth/login", h.Login)
	api.POST("/auth/register", h.Register)
	api.GET("/auth/refresh", h.RefreshToken)
}

// @Summary Авторизация
// @Tags auth
// @Accept json
// @Produce json
// @Param body body domain.LoginInput true "Данные для авторизации"
// @Success 200 {object} tokenResponse
// @Header 200 {string} Set-Cookie "Устанавливает refresh_token"
// @Failure 400
// @Router /api/v1/auth/login [post]
func (h *handler) Login(c echo.Context) error {
	input := new(domain.LoginInput)

	if err := c.Bind(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	tokens, err := h.service.Auth.Login(&domain.LoginInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidLoginOrPassword) {
			return c.NoContent(http.StatusBadRequest)
		}
		return err
	}

	return setTokensToResponse(c, tokens, h.config.Jwt.RefreshTokenTtl)
}

// @Summary Регистрация
// @Tags auth
// @Accept json
// @Produce json
// @Param body body domain.RegisterInput true "Данные для регистрации"
// @Success 201
// @Failure 400
// @Router /api/v1/auth/register [post]
func (h *handler) Register(c echo.Context) error {
	input := new(domain.RegisterInput)

	if err := c.Bind(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(input); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err := h.service.Auth.Register(&domain.RegisterInput{
		Username:   input.Username,
		Password:   input.Password,
		RePassword: input.RePassword,
	})
	if err != nil {
		if customErrors.MatchIn(err, domain.ErrPasswordsDontMatch, domain.ErrUsernameIsOccupied) {
			return NewErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return err
	}

	return c.NoContent(http.StatusCreated)
}

// @Summary Обновить токены
// @Description Обновляет Refresh и Access токены
// @Tags auth
// @Produce json
// @Success 200 {object} tokenResponse
// @Header 200 {string} Set-Cookie "Устанавливает refresh_token"
// @Failure 401
// @Router /api/v1/auth/refresh [get]
func (h *handler) RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	tokens, err := h.service.Auth.RefreshToken(userId, cookie.Value)
	if err != nil {
		if errors.Is(err, domain.ErrReshreshTokenNotFound) {
			return c.NoContent(http.StatusUnauthorized)
		}

		return err
	}

	return setTokensToResponse(c, tokens, h.config.Jwt.RefreshTokenTtl)
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func setTokensToResponse(c echo.Context, tokens *domain.Tokens, refreshTokenTtl time.Duration) error {
	refreshTokenMaxAge := time.Now().Add(refreshTokenTtl).Second()

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   refreshTokenMaxAge,
	})

	return c.JSON(http.StatusOK, tokenResponse{
		AccessToken: tokens.AccessToken,
	})
}
