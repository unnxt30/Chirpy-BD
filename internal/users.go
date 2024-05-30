package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var userCount int = 0

const hashValue int = 8

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ExpirationTime int    `json:"expires_in_secods"`
}

type DBStructUsers struct {
	Users map[int]User `json:"user"`
}

func (db *DB) CreateUser(args ...string) (User, error) {
	userCount++

	if len(args) < 2 {
		return User{}, errors.New("too few arguments")
	}

	var newUser User
	newUser.ID = userCount
	newUser.Email = args[0]
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(args[1]), hashValue)
	newUser.Password = string(hashedPassword)
	newUser.ExpirationTime = 3600;

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

func (db *DB) VerifyUser(userEmail, userPassword string, expirationTime int) (int, User) {

	userDataArray := db.getUserArray()
	// fmt.Println(userDataArray)

	for i := range userDataArray {
		if userDataArray[i].Email == userEmail {
			hashedPass := userDataArray[i].Password
			err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(userPassword))
			if err == nil {
				db.mux.Lock()
				defer db.mux.Unlock()
				file, _ := os.ReadFile(db.path)
				readData := &DBStructUsers{}
				json.Unmarshal(file, &readData)
				newUser := User{
					Email:          userDataArray[i].Email,
					ID:             userDataArray[i].ID,
					Password:       userDataArray[i].Password,
					ExpirationTime: expirationTime,
				}
				readData.Users[i] = newUser
				writeData, _ := json.Marshal(readData)
				os.Remove(db.path)
				openOrCreateFile(db.path)
				os.WriteFile(db.path, writeData, 0666)
				return 200, userDataArray[i]

			} else {
				fmt.Println(err.Error())
				return 401, User{}
			}
		}
		continue
	}

	return 401, User{}
}

func (db *DB) UpdateUserInDB(id int, email, password string) User {

	file, _ := os.ReadFile(db.path)
	readData := &DBStructUsers{}
	json.Unmarshal(file, &readData)
	
	var userToReturn User = User{}
	writeData := &DBStructUsers{
		Users: make(map[int]User),
	}
	for i := 0; i < len(readData.Users); i++ {
		if readData.Users[i].ID == id {
			var newUser User
			newUser, _ = db.CreateUser(email, password)
			newUser.ID = id
			newUser.ExpirationTime = 3600;
			writeData.Users[i] = newUser

			userToReturn = newUser
		} else {
			var newUser User
			newUser.Email = readData.Users[i].Email
			newUser.Password = readData.Users[i].Password
			newUser.ExpirationTime = 3600;
			writeData.Users[i] = newUser
		}
	}

	os.Remove(db.path)
	openOrCreateFile(db.path)
	newDbData, _ := json.Marshal(writeData)

	os.WriteFile(db.path, newDbData, 0666)
	return userToReturn

}
