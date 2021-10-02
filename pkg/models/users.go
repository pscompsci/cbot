package models

import "time"

type User struct {
	ID      int
	Name    string
	Email   string
	Hash    []byte
	Created time.Time
	Active  bool
}
