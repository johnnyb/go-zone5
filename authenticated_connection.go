package zone5

import (
	"time"
)

type AuthenticatedConnection struct {
	Connection *Connection
	AuthToken string
	RefreshToken string
	Expiration *time.Time
	Details map[string]interface{}
}

