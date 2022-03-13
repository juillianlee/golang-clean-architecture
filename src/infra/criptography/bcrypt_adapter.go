package criptography

import (
	app_criptography "github.com/juillianlee/helley-server/src/application/protocols/criptography"
	"golang.org/x/crypto/bcrypt"
)

type BcryptAdapter interface {
	app_criptography.Hasher
	app_criptography.HashComparer
}

type bcryptAdapter struct {
	salt int
}

func (c *bcryptAdapter) Hash(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), c.salt)
	return string(bytes), err
}

func (c *bcryptAdapter) Compare(value string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}

func NewBcryptAdapter(salt int) BcryptAdapter {
	return &bcryptAdapter{
		salt: salt,
	}
}
