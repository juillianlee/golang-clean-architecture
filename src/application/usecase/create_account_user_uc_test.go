package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/juillianlee/helley-server/src/domain/entities"
	domainerror "github.com/juillianlee/helley-server/src/domain/errors"
	"github.com/juillianlee/helley-server/src/domain/models"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	MockGetUserByEmail func(email string) (*entities.User, error)
	MockStore          func(user *entities.User) (*entities.User, error)
}

func (mock *mockRepository) GetUserByEmail(email string) (*entities.User, error) {
	return mock.MockGetUserByEmail(email)
}

func (mock *mockRepository) Store(user *entities.User) (*entities.User, error) {
	return mock.MockStore(user)
}

type mockCriptHasher struct {
	MockCriptoHasher func(value string) (string, error)
}

func (cript *mockCriptHasher) Hash(plaintext string) (string, error) {
	return cript.MockCriptoHasher(plaintext)
}

var mockCriptHasherSuccessfully = &mockCriptHasher{
	MockCriptoHasher: func(plaintext string) (string, error) {
		return "passwordhash", nil
	},
}

var mockCriptHasherError = &mockCriptHasher{
	MockCriptoHasher: func(plaintext string) (string, error) {
		return "", errors.New("fail generate hash password")
	},
}

func TestCreateAccountUserUCSuccessfully(t *testing.T) {
	var m = &mockRepository{
		MockGetUserByEmail: func(email string) (*entities.User, error) {
			return &entities.User{}, nil
		},
		MockStore: func(user *entities.User) (*entities.User, error) {
			user.ID = "1"
			return user, nil
		},
	}

	createAccountUserUC := NewCreateAccountUserUC(m, m, mockCriptHasherSuccessfully)
	user, err := createAccountUserUC.Handle(models.CreateAccountUserModel{
		Name:     "Juillian Lee",
		Email:    "juillian.lee@gmail.com",
		Password: "abc123",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, user.Name, "Juillian Lee")
	assert.Equal(t, user.Email, "juillian.lee@gmail.com")
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)
}

func TestCreateAccountUserUCEmailAlreadyExists(t *testing.T) {
	var m = &mockRepository{
		MockGetUserByEmail: func(email string) (*entities.User, error) {
			return entities.NewUser("1", "Juillian Lee", "juillian.lee@gmail.com", "abc123", time.Now(), time.Now())
		},
	}
	createAccountUserUC := NewCreateAccountUserUC(m, m, mockCriptHasherSuccessfully)
	_, err := createAccountUserUC.Handle(models.CreateAccountUserModel{Name: "Juillian Lee",
		Email:    "juillian.lee@gmail.com",
		Password: "abc123",
	})

	appError := err.(domainerror.DomainError)
	assert.Error(t, err)
	assert.ErrorIs(t, err, appError)
	assert.Equal(t, appError.Code, domainerror.ErrUserAlreadyExists)
}

func TestCreateAccountUserUCGetUserEmailError(t *testing.T) {
	var m = &mockRepository{
		MockGetUserByEmail: func(email string) (*entities.User, error) {
			return &entities.User{}, errors.New("Fail on get user by e-mail mock")
		},
	}
	createAccountUserUC := NewCreateAccountUserUC(m, m, mockCriptHasherSuccessfully)
	_, err := createAccountUserUC.Handle(models.CreateAccountUserModel{Name: "Juillian Lee",
		Email:    "juillian.lee@gmail.com",
		Password: "abc123",
	})

	assert.Error(t, err)
	appError := err.(domainerror.DomainError)
	assert.Error(t, err)
	assert.ErrorIs(t, err, appError)
	assert.Equal(t, appError.Code, domainerror.ErrGetDataRepository)
}

func TestCreateAccountUserUCStoreUserError(t *testing.T) {
	var m = &mockRepository{
		MockGetUserByEmail: func(email string) (*entities.User, error) {
			return &entities.User{}, nil
		},
		MockStore: func(user *entities.User) (*entities.User, error) {
			return user, errors.New("Fail on store user")
		},
	}
	createAccountUserUC := NewCreateAccountUserUC(m, m, mockCriptHasherSuccessfully)
	_, err := createAccountUserUC.Handle(models.CreateAccountUserModel{Name: "Juillian Lee",
		Email:    "juillian.lee@gmail.com",
		Password: "abc123",
	})

	assert.Error(t, err)
	appError := err.(domainerror.DomainError)
	assert.Error(t, err)
	assert.ErrorIs(t, err, appError)
	assert.Equal(t, appError.Code, domainerror.ErrStoreRepository)
}

func TestCreateAccountUserUCUserNameRequiredError(t *testing.T) {
	var m = &mockRepository{}
	createAccountUserUC := NewCreateAccountUserUC(m, m, mockCriptHasherSuccessfully)
	_, err := createAccountUserUC.Handle(models.CreateAccountUserModel{Name: ""})

	assert.Error(t, err)
	appError := err.(domainerror.DomainError)
	assert.Error(t, err)
	assert.ErrorIs(t, err, appError)
	assert.Equal(t, appError.Code, domainerror.ErrRequiredField)
	assert.Equal(t, appError.Message, "name is required")
}

func TestCreateAccountUserUCUserEmailRequiredError(t *testing.T) {
	var m = &mockRepository{}
	createAccountUserUC := NewCreateAccountUserUC(m, m, mockCriptHasherSuccessfully)
	_, err := createAccountUserUC.Handle(models.CreateAccountUserModel{Name: "Juillian Lee"})

	assert.Error(t, err)
	appError := err.(domainerror.DomainError)
	assert.Error(t, err)
	assert.ErrorIs(t, err, appError)
	assert.Equal(t, appError.Code, domainerror.ErrRequiredField)
	assert.Equal(t, appError.Message, "e-mail is required")
}

func TestCreateAccountUserUCUserInvalidEmailError(t *testing.T) {
	var m = &mockRepository{}
	createAccountUserUC := NewCreateAccountUserUC(m, m, mockCriptHasherSuccessfully)
	_, err := createAccountUserUC.Handle(models.CreateAccountUserModel{Name: "Juillian Lee", Email: "juillian"})

	assert.Error(t, err)
	appError := err.(domainerror.DomainError)
	assert.Error(t, err)
	assert.ErrorIs(t, err, appError)
	assert.Equal(t, appError.Code, domainerror.ErrRequiredField)
	assert.Equal(t, appError.Message, "e-mail is not valid")
}

func TestCreateAccountUserUCUserPasswordRequiredError(t *testing.T) {
	var m = &mockRepository{}
	createAccountUserUC := NewCreateAccountUserUC(m, m, mockCriptHasherSuccessfully)
	_, err := createAccountUserUC.Handle(models.CreateAccountUserModel{Name: "Juillian Lee", Email: "juillianlee@gmail.com"})

	assert.Error(t, err)
	appError := err.(domainerror.DomainError)
	assert.Error(t, err)
	assert.ErrorIs(t, err, appError)
	assert.Equal(t, appError.Code, domainerror.ErrRequiredField)
	assert.Equal(t, appError.Message, "password is required")
}

func TestCreateAccountUserUCUserPasswordEncriptHashError(t *testing.T) {
	var m = &mockRepository{
		MockGetUserByEmail: func(email string) (*entities.User, error) {
			return &entities.User{}, nil
		},
		MockStore: func(user *entities.User) (*entities.User, error) {
			user.ID = "1"
			return user, nil
		},
	}

	createAccountUserUC := NewCreateAccountUserUC(m, m, mockCriptHasherError)
	_, err := createAccountUserUC.Handle(models.CreateAccountUserModel{
		Name:     "Juillian Lee",
		Email:    "juillian.lee@gmail.com",
		Password: "abc123",
	})
	assert.Error(t, err)
	appError := err.(domainerror.DomainError)
	assert.Error(t, err)
	assert.ErrorIs(t, err, appError)
	assert.Equal(t, appError.Code, domainerror.ErrGenerateHashPassword)
	assert.Equal(t, appError.Message, "fail on generate hash to password user: fail generate hash password")
}
