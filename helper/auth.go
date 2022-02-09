package helper

import (
	"api/config"
	"api/db"
	"context"
	"crypto/md5"
	"encoding/hex"
)

func ValidateStr(str1 string, str2 string) (equal bool) {
	var len1 = len(str1)
	var len2 = len(str2)

	var max = len2
	if len1 > len2 {
		max = len1
	}

	equal = true
	for i := 0; i < max; i++ {
		if i >= len1 || i >= len2 || str1[i] != str2[i] {
			equal = false
		}
	}

	return
}

func RegenerateToken(token string) {
	// Create New Token
	hasher := md5.New()
	hasher.Write([]byte(token + config.ENV.Pepper))
	result := hex.EncodeToString(hasher.Sum(nil))

	// Hash it one more time
	hasher = md5.New()
	hasher.Write([]byte(result))

	db.Redis.Set(context.Background(), "TOKEN:"+hex.EncodeToString(hasher.Sum(nil)), "OK", 0)
}
