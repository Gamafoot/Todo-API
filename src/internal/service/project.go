package service

import (
	"root/internal/config"
	"root/internal/domain"
	"root/internal/storage"
	"root/pkg/hash"
	"root/pkg/jwt"
	"time"

	pkgErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type ProjectService interface {
	Create(input LoginInput) (Tokens, error)
	Delete(input RegisterInput) error
	Update(userId uint, refreshToken string) (Tokens, error)
}

type projectService struct {
	config       *config.Config
	storage      *storage.Storage
	tokenManager jwt.TokenManager
}

func newProjectService(cfg *config.Config, storage *storage.Storage, tokenManager jwt.TokenManager) *authService {
	return &authService{
		config:       cfg,
		storage:      storage,
		tokenManager: tokenManager,
	}
}

type (
	CreateProjectInput struct {
		UserId uint
		Name   string
	}

	DeleteProjectInput struct {
		Id uint
	}
)

func (s *authService) CreateProject(input CreateProjectInput) (Tokens, error) {
	var (
		user   *domain.User
		tokens Tokens
	)

	hasher := hash.NewSHA1Hasher(s.config.Hash.Salt)

	user, err := s.storage.User.GetByUsername(input.Username)
	if err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return tokens, domain.ErrInvalidLoginOrPassword
		}
		return tokens, err
	}

	ok, err := hasher.CheckHash(user.Password, input.Password)
	if err != nil {
		return tokens, err
	}

	if !ok {
		return tokens, domain.ErrInvalidLoginOrPassword
	}

	return s.createSession(user.Id)
}

func (s *authService) Register(input RegisterInput) error {
	if input.Password != input.RePassword {
		return domain.ErrPasswordsDontMatch
	}

	hasher := hash.NewSHA1Hasher(s.config.Hash.Salt)

	hash, err := hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	_, err = s.storage.User.GetByUsername(input.Username)
	if err != nil {
		if !pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		return domain.ErrUsernameIsOccupied
	}

	err = s.storage.User.Create(&domain.User{
		Username: input.Username,
		Password: hash,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) RefreshToken(userId uint, refreshToken string) (Tokens, error) {
	tokens := Tokens{}

	_, err := s.storage.Session.Get(userId, refreshToken)
	if err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return tokens, domain.ErrReshreshTokenNotFound
		}
		return tokens, err
	}

	return s.createSession(userId)
}

func (s *authService) createSession(userId uint) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId, s.config.Jwt.AccessTokenTtl)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := &domain.Session{
		UserId:       userId,
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.config.Jwt.RefreshTokenTtl),
	}

	err = s.storage.Session.Set(session)
	if err != nil {
		return res, err
	}

	return res, nil
}
