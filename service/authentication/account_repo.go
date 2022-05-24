package authentication

//go:generate mockgen -destination ../../mocks/account_repo/mocks.go -source=./account_repo.go -package=mockAccountRepo

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AccountModel struct {
	gorm.Model
	Password string
}

type AccountRepo interface {
	GetPwdFromDB(accountID string) (string, error)
}

func NewAccountRepo() AccountRepo {
	return &accountRepo{}
}

type accountRepo struct {
}

func (r *accountRepo) GetPwdFromDB(accountID string) (string, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return "", fmt.Errorf("db error: %w", err)
	}

	var ac AccountModel
	db.First(&ac, accountID)

	pwdFromDB := ac.Password
	return pwdFromDB, nil
}
