package jwt

func GenerateAccessToken(userID string, secret string, ttl int) (string, error) {
	return "", nil
}

func GenerateRefreshToken() (string, error) {
	return "", nil
}

func ValidateAccessToken(token, secret string) (string, error) {
	return "", nil
}

func HashToken(token string) string {
	return ""
}
