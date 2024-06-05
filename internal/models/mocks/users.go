package mocks

import "github.com/mhrdini/snippetbox/internal/models"

const (
	ValidName     = "Astarion Ancunin"
	ValidEmail    = "lilstar@bg3.com"
	ValidPassword = "cazadorsucks"
	DupeEmail     = "dupe@email.com"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case DupeEmail:
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == ValidEmail && password == ValidPassword {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, models.ErrNoRecord
	}
}
