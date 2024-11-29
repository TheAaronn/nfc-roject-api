package utils

import "time"

type JSON map[string]string

// Data type to simplify the Log entity to append to the array of logs, this is the equivalent of []string
type Log struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}
