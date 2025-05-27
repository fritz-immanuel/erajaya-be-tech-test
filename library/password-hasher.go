package library

import (
	"crypto/sha512"
	"encoding/base64"
)

const passwordSalt = "a99VVoWzmd1C9ujcitK0fIVNE0I5I61AC47C852RoLTsHDyLCltvP+ZHEkIl/2hkzTOW90c3ZEjtYRkdfTWJ1Q=="

func PasswordHasher(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(passwordSalt))
	hasher.Write([]byte(password))
	pwd := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return pwd
}
