package argonhasher

import "strings"

func ValidHash(h string) bool {

	/*
		This is highly effective at catching potential bad hash strings
		but is not really checking a great deal, so is used as a quick win
		highly effective at weeding out missing salts and key values from the string

		The 'Magic Number' comes from:
		23 from the basic string
		+ 2 from version
		+ 1 from time
		+ 4 from memory
		+ 1 from threads
		+ 86 from the salt
		+ 22 from the hash
	*/
	if len(h) < 139 {
		return false
	}

	/*
		Each hash string section should be prefixed by a `$`
		check there are the expected 5 sections
	*/
	if strings.Count(h, "$") != 5 {
		return false
	}

	/*
		Check for component parts that the
		argon key should contain
	*/
	if !strings.Contains(h, "$argon2id") ||
		!strings.Contains(h, "$v=") ||
		!strings.Contains(h, "$t=") ||
		!strings.Contains(h, ",m=") ||
		!strings.Contains(h, ",p=") {
		return false
	}

	/*
		Check for component parts that the
		argon key should NOT contain
		these would suggest missing data/arguments
	*/
	if strings.Contains(h, "=$") ||
		strings.Contains(h, "=,") ||
		strings.Contains(h, "$$") ||
		strings.HasSuffix(h, "$") {
		return false
	}

	return true
}
