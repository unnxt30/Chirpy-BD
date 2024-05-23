package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var chirpCount int = 0

type Parameters struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Parameters `json:"chirps"`
}


func openOrCreateFile(filePath string) (*os.File, error) {
	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// If the file does not exist, create it
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return nil, err
		}
		fmt.Println("File created successfully.")
		return file, nil
	} else if err != nil {
		// If there's an error other than file not existing, return it
		fmt.Println("Error checking file:", err)
		return nil, err
	}

	// If the file exists, open it
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	fmt.Println("File opened successfully.")
	return file, nil
}

func NewDB(file_path string) (*DB, error) {

	file, err := openOrCreateFile(file_path);
	if err != nil{
		return nil, err
	}
	file.Close()
	
	myDB := DB{
		path: file_path,
		mux:  &sync.RWMutex{}, // Create a new RWMutex instance
	}

	return &myDB, nil

}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Parameters, error) {
	chirpCount++
	var newChirp Parameters
	newChirp.Body = body
	newChirp.ID = chirpCount

	DB_file, err := os.ReadFile(db.path)

	if err != nil {
		return Parameters{}, err
	}
	db_data := &DBStructure{}
	err = json.Unmarshal(DB_file, &db_data)
	if err != nil {
		return Parameters{}, err
	}

	db_data.Chirps[chirpCount] = newChirp

	json.Marshal(db_data)

	return newChirp, nil

}

// // GetChirps returns all chirps in the database
// func (db *DB) GetChirps() ([]Parameters, error)

// // ensureDB creates a new database file if it doesn't exist
// func (db *DB) ensureDB() error

// // loadDB reads the database file into memory
// func (db *DB) loadDB() (DBStructure, error)

// // writeDB writes the database file to disk
// func (db *DB) writeDB(dbStructure DBStructure) error
