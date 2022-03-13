package errors

import "fmt"

type ApplicationError string

const (
	ErrRequiredField        ApplicationError = "REQUIRED_FIELD"
	ErrInvalidEmail         ApplicationError = "INVALID_EMAIL"
	ErrUserAlreadyExists    ApplicationError = "USER_ALREADY_EXISTS"
	ErrGetDataRepository    ApplicationError = "GET_DATA_REPOSITORY"
	ErrStoreRepository      ApplicationError = "STORE_REPOSITORY_ERROR"
	ErrGenerateHashPassword ApplicationError = "GENERATE_HASH_PASSWORD"
)

type DomainError struct {
	Code    ApplicationError
	Message string
}

func (err DomainError) Error() string {
	return fmt.Sprintf("%s: %s", err.Code, err.Message)
}

func NewDomainError(code ApplicationError, message string) DomainError {
	return DomainError{
		Code:    code,
		Message: message,
	}
}
