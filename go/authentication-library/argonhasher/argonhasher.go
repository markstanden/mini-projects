package argonhasher

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/markstanden/securerandom"
	"golang.org/x/crypto/argon2"
)

// KDFconfig is the base struct for my Argon2id wrapper
// we will use the standard library's argon2 IDKey function
// func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
// $argon2id$v=19$t=10,m=65536,p=8$SALT$HASH

type KDFconfig struct {

	// Salt is the base64 string used to salt our derived keys.
	Salt []byte

	// SaltLength
	// length of random-generated salt
	// (min 16 bytes recommended for password hashing)
	SaltLength int

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

func Encode(pw string) (hashWithConfig string, err error) {

	newArgon := KDFconfig{
		SaltLength: 64,
		Time:       10,
		Memory:     64 * 1024,
		Threads:    8,
		KeyLen:     16,
	}

	// call our salt generator function to produce a salt the required length.
	newArgon.Salt, err = securerandom.ByteSlice(newArgon.SaltLength)
	if err != nil {
		return "", err
	}

	// from "golang.org/x/crypto/argon2"
	// func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
	hash := argon2.IDKey([]byte(pw), newArgon.Salt, newArgon.Time, newArgon.Memory, newArgon.Threads, newArgon.KeyLen)

	return fmt.Sprintf("$argon2id$v=%v$t=%v,m=%v,p=%v$%s$%s",
		argon2.Version,
		newArgon.Time,
		newArgon.Memory,
		newArgon.Threads,
		base64.RawStdEncoding.EncodeToString(newArgon.Salt),
		base64.RawStdEncoding.EncodeToString(hash)), nil
}

// We will hash the provided strings using the arguments stored in the
// provided hash to compare against.
// this will allow us change the defaults and still compare against existing
// passwords.:
func Compare(ptPassword string, hash string) (err error) {
	// If we fail at any point prior to hashing,
	// it won't be worth the expense of hashing.
	worthHashing := true

	// Create a new conig struct
	newArgon := &KDFconfig{}

	// get the configuration from the hash
	config := strings.Split(hash, "$")

	// split the arguments
	// t=time, m=memory, p=threads
	tmp := strings.Split(config[3], ",")

	t, err := strconv.ParseUint(tmp[0][2:], 10, 32)
	if err != nil {
		return errors.New("failed to retrieve config.Time")
	}
	newArgon.Time = uint32(t)

	m, err := strconv.ParseUint(tmp[1][2:], 10, 32)
	if err != nil {
		return errors.New("failed to retrieve config.Memory")
	}
	newArgon.Memory = uint32(m)

	p, err := strconv.ParseUint(tmp[2][2:], 10, 8)
	if err != nil {
		return errors.New("failed to retrieve config.Threads")
	}
	newArgon.Threads = uint8(p)

	// Add the salt to the config
	newArgon.Salt, err = base64.RawStdEncoding.DecodeString(config[4])
	if err != nil {
		log.Println("Failed to decode Salt")
	}

	// Match the key length in the current hash rather than using our defaults.
	key, err := base64.RawStdEncoding.DecodeString(config[5])
	if err != nil {
		return fmt.Errorf("failed to decode hash")
	}
	newArgon.KeyLen = uint32(len(key))

	// Create a new hash.  This is the expensive part of the comparison,
	// so we should only check if necessary
	// from "golang.org/x/crypto/argon2"
	// func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
	var newHash string
	if worthHashing {
		newKey := argon2.IDKey([]byte(ptPassword), newArgon.Salt, newArgon.Time, newArgon.Memory, newArgon.Threads, newArgon.KeyLen)

		// format the key as a full argon hash string containing the config arguments
		newHash = fmt.Sprintf("$argon2id$v=%v$t=%v,m=%v,p=%v$%s$%s",
			argon2.Version,
			newArgon.Time,
			newArgon.Memory,
			newArgon.Threads,
			base64.RawStdEncoding.EncodeToString(newArgon.Salt),
			base64.RawStdEncoding.EncodeToString(newKey))
	}

	// Shamefully stolen from the x/crypto/bcrypt source code, we want all comparisons to take equal time,
	// whether it fails on the first bit or the last
	if subtle.ConstantTimeCompare([]byte(hash), []byte(newHash)) == 1 {
		// ptPassword hash equals hashedpassword
		return nil
	}

	return fmt.Errorf("mismatched password")
}
