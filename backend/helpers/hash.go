package helpers

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func CompareHash(userPassword string, comparePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(comparePassword))
	return err == nil
}

func ToMd5(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
