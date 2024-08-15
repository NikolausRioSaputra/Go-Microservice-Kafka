package domain

import "time"

type TTransaction struct {
	IncomingMessage
	ID         int64     `json:"id"`
	Topic      string    `json:"topic"`
	StepStatus string    `json:"stepStatus"`
	CreatedAt  time.Time `json:"createdAt"`
}
