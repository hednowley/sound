package hasher

import (
	"crypto/md5"
	"encoding/hex"
)

/*
func GetHash(path string) string {
	file, err := os.Open(path)
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "oops"
	}

	return hex.EncodeToString(hash.Sum(nil))
}
*/

func GetHash(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
