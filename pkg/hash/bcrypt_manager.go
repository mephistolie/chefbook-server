package hash

import (
	"github.com/mephistolie/chefbook-server/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type HashManager interface {
	Hash(data string) (string, error)
	ValidateByHash(data string, source string) error
}

type BcryptManager struct {
	saltCost int
}

func NewBcryptManager(saltCost int) *BcryptManager {
	return &BcryptManager{saltCost: saltCost}
}

func (m *BcryptManager) Hash(data string) (string, error)  {
	hashedData, err := bcrypt.GenerateFromPassword([]byte(data), m.saltCost)
	if err != nil {
		return "", err
	}
	return string(hashedData), nil
}

func (m *BcryptManager) ValidateByHash(data string, source string) error  {
	s := []byte(source)
	logger.Error(s)
	err := bcrypt.CompareHashAndPassword(s, []byte(data))
	return err
}
