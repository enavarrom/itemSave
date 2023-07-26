package model

import "time"

type Item struct {
	Id          string    `json:"id"`
	Site        string    `json:"site"`
	Price       float64   `json:"price"`
	StartTime   time.Time `json:"start_time"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Nickname    string    `json:"nickname"`
	Error       string
	EventId     string
}
