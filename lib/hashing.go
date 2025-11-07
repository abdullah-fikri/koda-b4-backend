package lib

import "github.com/matthewhartstonge/argon2"

func HashPassword(password string) []byte {
	argon := argon2.DefaultConfig()
	bytePassword, _ := argon.HashEncoded([]byte(password))
	return bytePassword
}
