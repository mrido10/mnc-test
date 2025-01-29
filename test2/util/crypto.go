package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"test2/config"
)

func GenerateHmacSHA256(data string) string {
	h := hmac.New(sha256.New, []byte(config.Conf.Auth.HashSecret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
