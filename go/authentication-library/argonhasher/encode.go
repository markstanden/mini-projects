package argonhasher

import (
	"encoding/base64"
	"fmt"

	"github.com/markstanden/securerandom"
	"golang.org/x/crypto/argon2"
)

/*
	Encode creates an argon2 hash from a plaintext password.
	It encodes using hardcoded defaults, with the cost providing
	a multiplier to the resources required.
	A cost set to zero will provide a strong default option, and is recommended
	The function returns a standard format argon2 hash string if the hash completes without error,
	otherwise an empty string is returned.
*/
func Encode(pw string, cost uint) (hashWithConfig string) {

	if pw == "" {
		return ""
	}

	/*
		Set the default option cost to be a
		sensible encoding time on modern hardware.
	*/
	if cost == 0 {
		cost = 3
	}

	newArgon := KDFconfig{
		SaltLength: 64,
		Time:       uint32(2 * cost),
		Memory:     uint32((1 + cost) * 32 * 1024),
		Threads:    uint8(1 + (cost * 3 / 2)),
		KeyLen:     16,
	}

	/*
		call our salt generator function to produce a cryptographically secure salt
		that is the specified length.
	*/
	newArgon.Salt = securerandom.ByteSlice(newArgon.SaltLength)
	if newArgon.Salt == nil {
		return ""
	}

	/*
		from "golang.org/x/crypto/argon2"
		func IDKey(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte
	*/
	key := argon2.IDKey([]byte(pw), newArgon.Salt, newArgon.Time, newArgon.Memory, newArgon.Threads, newArgon.KeyLen)

	hash := fmt.Sprintf("$argon2id$v=%v$t=%v,m=%v,p=%v$%s$%s",
		argon2.Version,
		newArgon.Time,
		newArgon.Memory,
		newArgon.Threads,
		base64.RawStdEncoding.EncodeToString(newArgon.Salt),
		base64.RawStdEncoding.EncodeToString(key))

	if !ValidHash(hash) {
		return ""
	}
	return hash
}
