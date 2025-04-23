package impl

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

type authService struct {
	config       *config.Config
	storage      *storage.Storage
	tokenManager jwt.TokenManager
}

func NewAuthService(cfg *config.Config, storage *storage.Storage, tokenManager jwt.TokenManager) *authService {
	return &authService{
		config:       cfg,
		storage:      storage,
		tokenManager: tokenManager,
	}
}

func (s *authService) Login(input *domain.LoginInput) (*domain.Tokens, error) {
	hasher := hash.NewSHA1Hasher(s.config.Hash.Salt)

	user, err := s.storage.User.GetByUsername(input.Username)
	if err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrInvalidLoginOrPassword
		}
		return nil, err
	}

	ok, err := hasher.CheckHash(user.Password, input.Password)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, domain.ErrInvalidLoginOrPassword
	}

	return s.createSession(user.Id)
}

func (s *authService) Register(input *domain.RegisterInput) error {
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

	_, err = s.storage.User.Create(&domain.User{
		Username: input.Username,
		Password: hash,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) RefreshToken(userId uint, refreshToken string) (*domain.Tokens, error) {
	_, err := s.storage.Session.Get(userId, refreshToken)
	if err != nil {
		if pkgErrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrReshreshTokenNotFound
		}
		return nil, err
	}

	return s.createSession(userId)
}

func (s *authService) createSession(userId uint) (*domain.Tokens, error) {
	var (
		tokens domain.Tokens
		err    error
	)

	tokens.AccessToken, err = s.tokenManager.NewJWT(userId, s.config.Jwt.AccessTokenTtl)
	if err != nil {
		return nil, err
	}

	tokens.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	session := &domain.Session{
		UserId:       userId,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(s.config.Jwt.RefreshTokenTtl),
	}

	err = s.storage.Session.Set(session)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}
