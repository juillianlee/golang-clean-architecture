package usecase

import (
	"fmt"
	"time"

	"github.com/juillianlee/helley-server/src/application/protocols/criptography"
	"github.com/juillianlee/helley-server/src/application/protocols/repository"
	"github.com/juillianlee/helley-server/src/domain/entities"
	domainerror "github.com/juillianlee/helley-server/src/domain/errors"
	"github.com/juillianlee/helley-server/src/domain/models"
)

type CreateAccountUserUC interface {
	Handle(u models.CreateAccountUserModel) (models.CreateAccountUserModelResponse, error)
}

type createAccountUserUC struct {
	getUserByEmailRepository repository.GetUserByEmailRepository
	storeUserRepository      repository.StoreUserRepository
	cryptHash                criptography.Hasher
}

func NewCreateAccountUserUC(getUserByEmailRepository repository.GetUserByEmailRepository, storeUserRepository repository.StoreUserRepository, cryptHash criptography.Hasher) CreateAccountUserUC {
	return &createAccountUserUC{
		getUserByEmailRepository: getUserByEmailRepository,
		storeUserRepository:      storeUserRepository,
		cryptHash:                cryptHash,
	}
}

func (usecase *createAccountUserUC) Handle(model models.CreateAccountUserModel) (models.CreateAccountUserModelResponse, error) {

	u, err := entities.NewUser("", model.Name, model.Email, model.Password, time.Now(), time.Now())
	if err != nil {
		return models.CreateAccountUserModelResponse{}, err
	}

	exists, err := usecase.getUserByEmailRepository.GetUserByEmail(u.Email.Value)
	if err != nil {
		return models.CreateAccountUserModelResponse{}, domainerror.NewDomainError(domainerror.ErrGetDataRepository, fmt.Sprintf("fail to get user by email on repository: %s", err.Error()))
	}

	if exists.ID != "" {
		return models.CreateAccountUserModelResponse{}, domainerror.NewDomainError(domainerror.ErrUserAlreadyExists, "email already exists")
	}

	hashPassword, err := usecase.cryptHash.Hash(u.Password.Value)

	if err != nil {
		return models.CreateAccountUserModelResponse{}, domainerror.NewDomainError(domainerror.ErrGenerateHashPassword, fmt.Sprintf("fail on generate hash to password user: %+v", err))
	}

	u.Password.Value = hashPassword

	userStorage, err := usecase.storeUserRepository.Store(u)
	if err != nil {
		return models.CreateAccountUserModelResponse{}, domainerror.NewDomainError(domainerror.ErrStoreRepository, fmt.Sprintf("fail to store user on repository: %s", err.Error()))
	}

	return models.MakeCreateAccountUserModelResponseFromEntity(*userStorage), nil
}
