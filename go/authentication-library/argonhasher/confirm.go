package argonhasher

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInputInvalid = fmt.Errorf("input arguments invalid")
)

// We will hash the provided strings using the arguments stored in the
// provided hash to compare against.
// this will allow us change the defaults and still compare against existing
// passwords.:
func Confirm(password string, hash string) (valid bool) {

	if password == "" || hash == "" {
		return false
	}

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
		return false
	}
	newArgon.Time = uint32(t)

	m, err := strconv.ParseUint(tmp[1][2:], 10, 32)
	if err != nil {
		return false
	}
	newArgon.Memory = uint32(m)

	p, err := strconv.ParseUint(tmp[2][2:], 10, 8)
	if err != nil {
		return false
	}
	newArgon.Threads = uint8(p)

	// Add the salt to the config
	newArgon.Salt, err = base64.RawStdEncoding.DecodeString(config[4])
	if err != nil {
		return false
	}

	// Match the key length in the current hash rather than using our defaults.
	key, err := base64.RawStdEncoding.DecodeString(config[5])
	if err != nil {
		return false
	}
	newArgon.KeyLen = uint32(len(key))

	// Create a new hash.  This is the expensive part of the comparison,
	// so we should only check if necessary
	// from "golang.org/x/crypto/argon2"
	// func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
	var newHash string
	if worthHashing {
		newKey := argon2.IDKey([]byte(password), newArgon.Salt, newArgon.Time, newArgon.Memory, newArgon.Threads, newArgon.KeyLen)

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
		return true
	}

	return false
}
