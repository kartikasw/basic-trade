package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New(" invalid")
)

type Payload struct {
	TokenID   uuid.UUID `json:"id"`
	UUID      uuid.UUID `json:"uid"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(uuidUsr uuid.UUID, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		TokenID:   uuid.New(),
		UUID:      uuidUsr,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
