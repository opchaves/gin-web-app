package fixture

import (
	"fmt"
	"math/rand"
	"strings"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("1234567890")

func randStringLowerRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes)/2)]
	}

	return string(b)
}

// RandInt returns a random int within the given range
func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// Email returns a random email ending with @example.com
func Email() string {
	email := fmt.Sprintf("%s@example.com", randStringLowerRunes(RandInt(5, 10)))
	return strings.ToLower(email)
}
