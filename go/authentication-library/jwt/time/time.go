package time

import "time"

ProjectStart = 

func GetUnix() int64 {
	fmt.Println(time.Now().UTC().Unix())
	return time.Now().UTC().Unix()
}

// withinRange checks the number lies between the low and high
// withinRange returns false if the provided number is too low or too high
func WithinRange(num, low, high int64) bool {

	if num < low {
		return false
	}

	if num > high {
		return false
	}

	return true
}
