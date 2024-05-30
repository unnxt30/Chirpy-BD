package src

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "os"
	// "strconv"
	"strings"

	// "github.com/golang-jwt/jwt/v5"
	database "github.com/unnxt30/Chirpy-BD/internal"
)

func WriteUser(w http.ResponseWriter, r *http.Request) {
	MyDatabase, err := database.NewDB("database_user.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err = decoder.Decode(&user)
	if err != nil {
		RespondWithJSON(w, 404, map[string]string{"error": "couldn't create user"})
	}
	userData, _ := MyDatabase.CreateUser(user.Email, user.Password)

	RespondWithJSON(w, 201, map[string]any{"id": userData.ID, "email": userData.Email, "is_chirpy_red": userData.Chirpy_red});
	// RespondWithJSON(w, 201, userData);
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	MyDatabase, err := database.NewDB("database_user.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	decoder.Decode(&user)

	code, returnUser := MyDatabase.VerifyUser(user.Email, user.Password, user.ExpirationTime)

	RespondWithJSON(w, code, map[string]any{"id": returnUser.ID, "email": returnUser.Email, "token": returnUser.GenerateToken(), "refresh_token": returnUser.CreateRefreshToken(), "is_chirpy_red": returnUser.Chirpy_red})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var Uid int

	givenToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	// fmt.Printf("header token: %v", givenToken)

	Uid = database.CheckToken(givenToken)

	if Uid == -1  {
		RespondWithJSON(w, 401, map[string]string{"error": "Token doesn't match database"})
		return
	}

	MyDatabase, _ := database.NewDB("database_user.json")
	user := database.User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&user)

	updatedUser := MyDatabase.UpdateUserInDB(Uid, user.Email, user.Password)
	RespondWithJSON(w, 200, map[string]any{"id": updatedUser.ID, "email": updatedUser.Email})
}
