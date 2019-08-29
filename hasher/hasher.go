package hasher

import (
	"crypto/md5"
	"encoding/hex"
)

// GetHash calculates an MD5 hash.
func GetHash(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
