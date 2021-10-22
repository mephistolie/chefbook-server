package hash

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type HashManager interface {
	Hash(data string) (string, error)
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
	return fmt.Sprintf("%x", hashedData), nil
}