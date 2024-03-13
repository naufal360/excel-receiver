package entity

import "time"

type Request struct {
	// ID        int       `db:"id" json:"id,omitempty"`
	RequestID string    `db:"request_id" json:"request_id"`
	Status    string    `db:"status" json:"status"`
	Filepath  string    `db:"file_path" json:"file_path"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
