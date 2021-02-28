package argon

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// KDFconfig is the base struct for my Argon2id wrapper
// we will use the standard library's argon2 IDKey function
// func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
// $argon2id$v=19$m=64,t=4,p=8$SSALT$vHASH

type KDFconfig struct {

	// SaltLength
	// length of random-generated salt
	// (16 bytes recommended for password hashing)
	SaltLength uint8

	// Time (i.e. iterations) - t
	// number of iterations or pass throughs to perform
	Time uint32

	// Memory - m
	// amount of memory (in kilobytes) to use
	Memory uint32

	// Threads (parallelism) p: degree of parallelism (i.e. number of threads)
	Threads uint8

	// KeyLen T: desired number of returned bytes
	// 128 bit (16 bytes) sufficient for most applications
	KeyLen uint32
}

func Encode(pw string) (hashWithConfig string) {

	newArgon := KDFconfig{
		SaltLength: 64,
		Time:       10,
		Memory:     64 * 1024,
		Threads:    8,
		KeyLen:     16,
	}

	salt, err := createSalt(newArgon.SaltLength)
	if err != nil {
		return err.Error()
	}

	hash := argon2.IDKey([]byte(pw), salt, newArgon.Time, newArgon.Memory, newArgon.Threads, newArgon.KeyLen)

	hashWithConfig = fmt.Sprintf("$argon2id$v=%v$m=%v,t=%v,p=%v$%s$%s",
		argon2.Version,
		newArgon.Memory,
		newArgon.Time,
		newArgon.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash))

	return hashWithConfig
}

// createSalt creates a random string of bytes of length saltLength
// using the cryptographically secure crypto/rand package.
func createSalt(saltLength uint8) (salt []byte, err error) {

	//create an empty slice of bytes the required size of the salt
	salt = make([]byte, saltLength)

	// fill the slice of bytes with crypto randomness
	_, err = rand.Read(salt)

	// check for errors
	if err != nil {
		return nil, err
	}

	// return the salt if error free
	return salt, nil
}
