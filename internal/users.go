package database

import (
	"encoding/json"
	"fmt"
	"os"
)

var userCount int = 0

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
}

type DBStructUsers struct {
	Users map[int]User `json:"user"`
}

func (db *DB) CreateUser(userEmail string) (User, error){
	userCount++

	var newUser User
	newUser.ID = userCount
	newUser.Email = userEmail

	db.mux.Lock()
	defer db.mux.Unlock()


	DB_file, err := os.ReadFile(db.path)

	if err != nil {
		return User{}, err
	}

	fileStatus, _ := os.Stat(db.path)

	if fileStatus.Size() == 0 {
		dbData := &DBStructUsers{
			Users : make(map[int]User), // Initialize the chirps map
		}

		dbData.Users[userCount-1] = newUser
		// fmt.Println(dbData.Chirps);
		writeData, err := json.Marshal(dbData)
		if err != nil {
			fmt.Println(err)
			return User{}, err
		}
		os.WriteFile(db.path, writeData, 0666)
		return newUser, nil
	}

	db_data := &DBStructUsers{}

	err = json.Unmarshal(DB_file, &db_data)
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	db_data.Users[userCount-1] = newUser
	writeData, err := json.Marshal(db_data)
	if err != nil {
		return User{}, err
	}
	os.WriteFile(db.path, writeData, 0666)
	return newUser, nil

}
