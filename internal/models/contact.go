package models

import "time"

type ContactMessage struct {
	ID        int
	Name      string
	Phone     string
	Email     string
	Topic     string
	Message   string
	CreatedAt time.Time
}
