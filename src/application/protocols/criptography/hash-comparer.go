package criptography

type HashComparer interface {
	Compare(value string, hash string) bool
}
