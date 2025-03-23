package v1

import (
	"errors"
	"net/http"
	"root/internal/domain"
	"root/internal/service"
	customErrors "root/pkg/errors"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *handler) initAuthRoutes(api *echo.Group) {
	api.POST("/auth/login", h.Login)
	api.POST("/auth/register", h.Register)
	api.GET("/auth/refresh", h.RefreshToken, h.requiredAuth)
}

type (
	loginInput struct {
		Username string `json:"username" binding:"required,min=3,max=25"`
		Password string `json:"password" binding:"required,min=8,max=64"`
	}

	loginResponse struct {
		AccessToken string `json:"access_token"`
	}
)

func (h *handler) Login(c echo.Context) error {
	input := new(loginInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	tokens, err := h.service.Auth.Login(service.LoginInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidLoginOrPassword) {
			return newResponse(c, http.StatusUnauthorized, err.Error())
		}
		return err
	}

	return setTokensToResponse(c, tokens, h.config.Jwt.RefreshTokenTtl)
}

type registerInput struct {
	Username   string `json:"username" binding:"required,min=3,max=25"`
	Password   string `json:"password" binding:"required,min=8,max=64"`
	RePassword string `json:"re_password" binding:"required,min=8,max=64"`
}

func (h *handler) Register(c echo.Context) error {
	input := new(registerInput)

	if err := c.Bind(input); err != nil {
		return newResponse(c, http.StatusBadRequest, "invalid request body")
	}

	err := h.service.Auth.Register(service.RegisterInput{
		Username:   input.Username,
		Password:   input.Password,
		RePassword: input.RePassword,
	})
	if err != nil {
		if customErrors.MatchIn(err, domain.ErrPasswordsDontMatch, domain.ErrUsernameIsOccupied) {
			return newResponse(c, http.StatusBadRequest, err.Error())
		}
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *handler) RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.NoContent(http.StatusForbidden)
	}

	userId, err := getUserIdFromContext(c)
	if err != nil {
		return err
	}

	tokens, err := h.service.Auth.RefreshToken(userId, cookie.Value)
	if err != nil {
		if errors.Is(err, domain.ErrReshreshTokenNotFound) {
			return c.NoContent(http.StatusForbidden)
		}

		return err
	}

	return setTokensToResponse(c, tokens, h.config.Jwt.RefreshTokenTtl)
}

func setTokensToResponse(c echo.Context, tokens service.Tokens, refreshTokenTtl time.Duration) error {
	refreshTokenMaxAge := time.Now().Add(refreshTokenTtl).Second()

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   refreshTokenMaxAge,
	})

	return c.JSON(http.StatusOK, loginResponse{
		AccessToken: tokens.AccessToken,
	})
}
