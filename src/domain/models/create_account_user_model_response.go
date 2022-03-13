package models

import (
	"time"

	"github.com/juillianlee/helley-server/src/domain/entities"
)

type CreateAccountUserModelResponse struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func MakeCreateAccountUserModelResponseFromEntity(u entities.User) CreateAccountUserModelResponse {
	return CreateAccountUserModelResponse{
		ID:        u.ID,
		Name:      u.Name.Value,
		Email:     u.Email.Value,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
