package utils

import (
	"crypto/sha256"
	"strconv"
	"time"
)

//GenerateSHA256 generates a SHA-256 hash
func GenerateSHA256() []byte {
	// get unix time  convert to string
	unixNow := strconv.FormatInt(time.Now().Unix(), 10)

	// create sha256 from unix time
	h := sha256.New()
	h.Write([]byte(unixNow))
	hash := h.Sum(nil)

	return hash
}
