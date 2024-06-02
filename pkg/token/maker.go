package token

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(uuid uuid.UUID, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
