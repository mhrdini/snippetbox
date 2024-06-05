package mocks

import "github.com/mhrdini/snippetbox/internal/models"

var (
	MockEmail    = "astarion@bg3.com"
	MockPassword = "cazadorsucks"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case MockEmail:
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == MockEmail && password == MockPassword {
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
