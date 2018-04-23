package util

import (
	"fmt"
	"math/rand"
	"time"
)

var randSrc rand.Source

const (
	saltLen = 30
)

func init() {
	randSrc = rand.NewSource(time.Now().Unix())
}

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//TODO: Placeholder bycrypt
func EncryptPassword(password string) (salt string, hash string) {
	cs := []rune(charSet)
	rs := make([]rune, saltLen)
	for i := 0; i < saltLen; i++ {
		rs[i] = cs[rand.Intn(len(cs))]
	}
	salt = string(rs)
	hash = hashPassword(salt, password)
	return
}

func hashPassword(salt string, password string) string {
	return fmt.Sprintf("%s:%s", salt, password)
}
