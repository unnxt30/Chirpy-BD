package database

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var userCount int = 0

const hashValue int = 8

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DBStructUsers struct {
	Users map[int]User `json:"user"`
}

func (db *DB) CreateUser(userEmail string, userPassword string) (User, error) {
	userCount++

	var newUser User
	newUser.ID = userCount
	newUser.Email = userEmail
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userPassword), hashValue)
	newUser.Password = string(hashedPassword)

	db.mux.Lock()
	defer db.mux.Unlock()

	DB_file, err := os.ReadFile(db.path)

	if err != nil {
		return User{}, err
	}

	fileStatus, _ := os.Stat(db.path)

	if fileStatus.Size() == 0 {
		dbData := &DBStructUsers{
			Users: make(map[int]User), // Initialize the chirps map
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

func (db *DB) getUserArray() []User {
	db.mux.Lock()
	defer db.mux.Unlock()

	file, err := os.ReadFile(db.path)
	if err != nil {
		return nil
	}
	readData := &DBStructUsers{}
	json.Unmarshal(file, &readData)

	var userDataArray []User
	for i := 0; i < len(readData.Users); i++ {
		userDataArray = append(userDataArray, readData.Users[i])
	}

	return userDataArray
}

func (db *DB) VerifyUser(userEmail, userPassword string) (int, User) {

	userDataArray := db.getUserArray()
	// fmt.Println(userDataArray)

	for i := range userDataArray {
		if userDataArray[i].Email == userEmail {
			hashedPass := userDataArray[i].Password
			err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(userPassword))
			if err == nil {
				return 200, userDataArray[i]

			} else {
				return 401, User{}
			}
		}
		continue
	}

	return 401, User{}
}
