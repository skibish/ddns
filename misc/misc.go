package misc

// Success return true if status code is in success (2xx) range,
// false otherwise
func Success(status int) bool {
	return status > 199 && status < 300
}
