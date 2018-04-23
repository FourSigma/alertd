package util

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var randSrc rand.Source

const (
	saltLEN  = 30
	hashCOST = 14
)

func init() {
	randSrc = rand.NewSource(time.Now().Unix())
}

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//TODO: Placeholder bycrypt
func EncryptPassword(password string) (salt string, hash string, err error) {
	cs := []rune(charSet)
	rs := make([]rune, saltLEN)
	for i := 0; i < saltLEN; i++ {
		rs[i] = cs[rand.Intn(len(cs))]
	}
	salt = string(rs)
	password = fmt.Sprintf("%s:%s", salt, password)
	hash, err = bcrypt.GenerateFromPassword([]byte(password), hashCOST)
	return
}
