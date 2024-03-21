package entity

import "time"

type Token struct {
	ID        int       `db:"id"`
	Token     string    `db:"token"`
	ExpiredAt time.Time `db:"expired_at"`
	CreatedAt time.Time `db:"created_at"`
}

func (t *Token) IsExpired() bool {
	return t.ExpiredAt.Before(time.Now().Add(7 * time.Hour))
}
