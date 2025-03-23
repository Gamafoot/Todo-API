package jwt

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	pkgErrors "github.com/pkg/errors"
)

type TokenManager interface {
	NewJWT(id uint, ttl time.Duration) (string, error)
	Parse(accessToken string) (uint, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId uint `json:"user_id"`
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, pkgErrors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(id uint, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: id,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, pkgErrors.Errorf(`unexpected signing method: #{token.Header["alg"]}`)
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return 0, pkgErrors.WithStack(err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, pkgErrors.New("error get user claims from token")
	}

	return claims.UserId, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
