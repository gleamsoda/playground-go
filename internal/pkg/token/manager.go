//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=../../mock/pkg/$GOPACKAGE/$GOFILE
package token

import (
	"time"
)

// Manager is an interface for managing tokens
type Manager interface {
	// Create creates a new token for a specific username and duration
	Create(username string, duration time.Duration) (string, *Payload, error)
	// Verify checks if the token is valid or not
	Verify(t string) (*Payload, error)
}
