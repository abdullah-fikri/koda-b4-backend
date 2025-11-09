package lib

import (
	"fmt"

	"github.com/matthewhartstonge/argon2"
)

func HashPassword(password string) []byte {
	argon := argon2.DefaultConfig()
	bytePassword, _ := argon.HashEncoded([]byte(password))
	return bytePassword
}

func VerifyPassword(password string, hashedPassword string) bool {
	ok, _ := argon2.VerifyEncoded([]byte(password), []byte(hashedPassword))
	fmt.Println(hashedPassword)
	fmt.Println(password)
	return ok
}
