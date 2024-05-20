package database

import (
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]chirps.Parameters `json:"chirps"`
}
