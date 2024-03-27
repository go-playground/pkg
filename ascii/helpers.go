package asciiext

// IsAlphanumeric returns true if the byte is an ASCII letter or digit.
func IsAlphanumeric(c byte) bool {
	return IsLower(c) || IsUpper(c) || IsDigit(c)
}

// IsUpper returns true if the byte is an ASCII uppercase letter.
func IsUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

// IsLower returns true if the byte is an ASCII lowercase letter.
func IsLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

// IsDigit returns true if the byte is an ASCII digit.
func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
