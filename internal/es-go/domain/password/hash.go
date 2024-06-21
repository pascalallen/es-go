package password

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Hash string

func Create(password string) Hash {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return Hash(hash)
}

func (p *Hash) Compare(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*p), []byte(password))
	if err != nil {
		return false
	}

	return true
}
