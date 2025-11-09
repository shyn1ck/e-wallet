package entity

import "time"

type APIClient struct {
	ID        int64
	UserID    string
	SecretKey string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Authenticate checks if the given digest matches the expected digest
func (c *APIClient) Authenticate(userID, digest, body string, hmacFunc func(secret, data string) string) bool {
	if !c.IsActive {
		return false
	}

	if c.UserID != userID {
		return false
	}

	expectedDigest := hmacFunc(c.SecretKey, body)
	return expectedDigest == digest
}
