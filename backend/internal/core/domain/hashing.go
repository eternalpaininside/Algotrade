package domain

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	PassSaltSize = 16
	PassTime     = 2
	PassMemory   = 19 * 1024
	PassThreads  = 1
	PassKeySize  = 32
)

func PasswordHashing(
	password string,
) (salt, hash []byte, err error) {
	salt = make([]byte, PassSaltSize)

	if _, err = rand.Read(salt); err != nil {
		return nil, nil, fmt.Errorf(
			"generate salt: %w",
			err,
		)
	}

	hash = argon2.IDKey(
		[]byte(password),
		salt,
		PassTime,
		PassMemory,
		PassThreads,
		PassKeySize,
	)

	return
}

func ConcatHash(
	salt, hash []byte,
) string {
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		PassMemory,
		PassTime,
		PassThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}
