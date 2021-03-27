package token

import "time"

func getUnixTime() int64 {
	return time.Now().UTC().Unix()
}

// withinRange checks the number lies between the low and high
// withinRange returns false if the provided number is too low or too high
func withinRange(num, low, high int64) bool {

	if num < low {
		return false
	}

	if num > high {
		return false
	}

	return true
}
