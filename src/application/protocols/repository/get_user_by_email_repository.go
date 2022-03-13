package repository

import (
	"github.com/juillianlee/helley-server/src/domain/entities"
)

type GetUserByEmailRepository interface {
	GetUserByEmail(email string) (*entities.User, error)
}
