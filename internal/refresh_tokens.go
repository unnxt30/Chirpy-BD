package database

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	// "golang.org/x/tools/go/analysis/passes/defers"
)

const RefTokenPath string = "database_refToken.json"

var index int = 0
var mut sync.RWMutex

type dbStruct struct {
	Tokens map[int]storeTokenStruct
}

type storeTokenStruct struct {
	Id          int       `json:"id"`
	RefToken    string    `json:"token"`
	CreatedAt   time.Time `json:"created_at"`
	ExpDuration int       `json:"exp_duration"`
}

func generateRefToken() string {
	length := 32
	Token := make([]byte, length)

	_, err := rand.Read(Token)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	returnToken := hex.EncodeToString(Token)
	return returnToken
}

func (user *User) CreateRefreshToken() string {
	openOrCreateFile(RefTokenPath)

	length := 32
	Token := make([]byte, length)

	_, err := rand.Read(Token)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	returnToken := hex.EncodeToString(Token)
	newToken := storeTokenStruct{
		Id:          user.ID,
		RefToken:    returnToken,
		CreatedAt:   time.Now(),
		ExpDuration: 60 * 24 * 60,
	}

	fileStatus, _ := os.Stat(RefTokenPath)

	if fileStatus.Size() == 0 {
		writeData := &dbStruct{
			Tokens: make(map[int]storeTokenStruct),
		}
		writeData.Tokens[index] = newToken
		write_data, _ := json.Marshal(writeData)
		os.WriteFile(RefTokenPath, write_data, 0666)
		index++
		return returnToken
	}

	file, _ := os.ReadFile(RefTokenPath)
	params := &dbStruct{}
	json.Unmarshal(file, &params)
	params.Tokens[index] = newToken

	writData, _ := json.Marshal(params)
	os.WriteFile(RefTokenPath, writData, 0666)
	index++
	return returnToken
}

func ValidateRefToken(Token string) (int, string) {
	mut.Lock()
	defer mut.Unlock()
	file, _ := os.ReadFile(RefTokenPath)
	params := &dbStruct{}
	json.Unmarshal(file, &params)
	fmt.Printf("%v", params.Tokens)
	myDB, _ := NewDB("database_user.json")
	var user User
	for i := range params.Tokens {
		currentToken := params.Tokens[i]
		if currentToken.RefToken == Token {
			if time.Now().After(currentToken.CreatedAt.Add(time.Duration(currentToken.ExpDuration) * time.Minute)) {
				return 401, "Token Expired"
			} else {
				userArray := myDB.getUserArray()
				for i := range userArray {
					if userArray[i].ID == currentToken.Id {
						user.ID = currentToken.Id
						user.Email = userArray[i].Email
						user.Password = userArray[i].Password
						user.ExpirationTime = 3600
						myDB.UpdateUserInDB(currentToken.Id, userArray[i].Email, userArray[i].Password)
						break
					}
				}
			}
			break
		} else {
			return 401, "token doesn't exist"
		}
	}

	return 200, user.GenerateToken()

}

func DeleteRefToken(Token string) int {
	mut.Lock()
	defer mut.Unlock()
	file, _ := os.ReadFile(RefTokenPath)
	params := &dbStruct{}
	json.Unmarshal(file, &params)

	newToken := &storeTokenStruct{}

	for i := range params.Tokens {
		currentToken := params.Tokens[i]
		if currentToken.RefToken == Token {
			params.Tokens[i] = *newToken
		}
		break
	}

	os.Remove(RefTokenPath)
	openOrCreateFile(RefTokenPath)
	writeData, _ := json.Marshal(params)
	os.WriteFile(RefTokenPath, writeData, 0666)

	return 204
}
