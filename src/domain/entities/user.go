package entities

import (
	"net/mail"
	"time"

	domainerror "github.com/juillianlee/helley-server/src/domain/errors"
)

type User struct {
	ID        string
	Name      UserName
	Email     UserEmail
	Password  UserPassword
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id string, name string, email string, password string, createdAt time.Time, updatedAt time.Time) (*User, error) {
	fieldName, err := NewUserName(name)

	if err != nil {
		return &User{}, err
	}

	fieldEmail, err := NewUserEmail(email)
	if err != nil {
		return &User{}, err
	}

	fieldPassword, err := NewUserPassword(password)
	if err != nil {
		return &User{}, err
	}

	return &User{
		ID:        id,
		Name:      fieldName,
		Email:     fieldEmail,
		Password:  fieldPassword,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

type UserEmail struct {
	Value string
}

type UserName struct {
	Value string
}

type UserPassword struct {
	Value string
}

func NewUserName(name string) (UserName, error) {
	if name == "" {
		return UserName{}, domainerror.NewDomainError(domainerror.ErrRequiredField, "name is required")
	}
	return UserName{
		Value: name,
	}, nil
}

func NewUserEmail(email string) (UserEmail, error) {
	if email == "" {
		return UserEmail{}, domainerror.NewDomainError(domainerror.ErrRequiredField, "e-mail is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return UserEmail{}, domainerror.NewDomainError(domainerror.ErrRequiredField, "e-mail is not valid")
	}

	return UserEmail{
		Value: email,
	}, nil
}

func NewUserPassword(password string) (UserPassword, error) {
	if password == "" {
		return UserPassword{}, domainerror.NewDomainError(domainerror.ErrRequiredField, "password is required")
	}
	return UserPassword{
		Value: password,
	}, nil
}
