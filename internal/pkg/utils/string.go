package utils

import "math/rand"

var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// RandString generates a random string with the length provided
func RandString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// TruncateString returns a truncated version of a string by maxLength with ellipsis
func TruncateString(str string, maxLength int) string {
	truncated := str
	if len(str) > maxLength {
		if maxLength > 3 {
			maxLength -= 3
		}
		truncated = str[0:maxLength] + "..."
	}
	return truncated
}
