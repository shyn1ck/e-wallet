package cache

import (
	"fmt"
	"time"
)

const (
	PrefixAPIClient = "api_client"
)

const (
	TTLAPIClient = 15 * time.Minute // API clients are cached for 15 minutes
)

func BuildAPIClientKey(userID string) string {
	return fmt.Sprintf("%s:%s", PrefixAPIClient, userID)
}
