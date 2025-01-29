package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"test2/config"
	"testing"
)

func TestGenerateHmacSHA256(t *testing.T) {
	data := "hello world"
	expectedH := hmac.New(sha256.New, []byte(config.Conf.Auth.HashSecret))
	expectedH.Write([]byte(data))
	expectedHash := hex.EncodeToString(expectedH.Sum(nil))

	hash := GenerateHmacSHA256(data)
	fmt.Println(hash)
	if hash != expectedHash {
		t.Errorf("Expected %s, got %s", expectedHash, hash)
	}
}
