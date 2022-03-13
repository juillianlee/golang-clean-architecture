package repository

import (
	"github.com/juillianlee/helley-server/src/domain/entities"
)

type StoreUserRepository interface {
	Store(user *entities.User) (*entities.User, error)
}
