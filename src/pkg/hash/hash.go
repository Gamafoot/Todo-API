package hash

import (
	"crypto/sha1"
	"fmt"

	pkgErrors "github.com/pkg/errors"
)

type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

func (h SHA1Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", pkgErrors.WithStack(err)
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}

func (h SHA1Hasher) CheckHash(originHash, password string) (bool, error) {
	hash, err := h.Hash(password)
	if err != nil {
		return false, err
	}

	if originHash == hash {
		return true, nil
	}

	return false, nil 
}
