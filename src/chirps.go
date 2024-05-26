package src

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	database "github.com/unnxt30/Chirpy-BD/internal"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {

	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.WriteHeader(code)

	w.Write(response)

	return nil

}

func GiveValidBody(w http.ResponseWriter, str string, params database.Parameters) string {
	if len(params.Body) > 140 {
		RespondWithJSON(w, 400, map[string]string{"error": "Chirp is too long"})
		return ""
	}
	found := false
	var cleaned_body_words []string
	bad_words := []string{"kerfuffle", "sharbert", "fornax"}
	//respondWithJSON(w, 200, map[string]bool{"valid" : true})
	input_str := params.Body
	for _, word := range strings.Split(input_str, " ") {
		found = false
		for _, bad_word := range bad_words {
			if strings.ToLower(word) == bad_word {
				found = true
				break
			}
		}
		if found {
			cleaned_body_words = append(cleaned_body_words, "****")
		} else {
			cleaned_body_words = append(cleaned_body_words, word)
		}
	}
	cleaned_body := strings.Join(cleaned_body_words, " ")
	return cleaned_body
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	MyDatabase, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println(err)
	}

	decoder := json.NewDecoder(r.Body)
	params := database.Parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		RespondWithJSON(w, 404, map[string]string{"error": "Something went wrong"})
	} else {
		// fmt.Println("recieved Chirp Succesfully.")
		chirpBody := GiveValidBody(w, params.Body, params)
		responseChirp, _ := MyDatabase.CreateChirp(chirpBody)
		fmt.Println("Created chirp")
		RespondWithJSON(w, 201, responseChirp)
	}
}

func ChirpsGET(w http.ResponseWriter, r *http.Request) {

	MyDatabase, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	chirpArray, err := MyDatabase.GetChirp()
	if err != nil {
		RespondWithJSON(w, 404, map[string]string{"error": err.Error()})
	}

	RespondWithJSON(w, 200, chirpArray)
}

func ChirpGETbyID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	MyDatabase, _ := database.NewDB("database.json")

	chirpArray, _ := MyDatabase.GetChirp()

	if id > len(chirpArray) {
		RespondWithJSON(w, 404, map[string]string{"error": "ya bish"})
		return
	}

	RespondWithJSON(w, 200, chirpArray[id-1])
}

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

	RespondWithJSON(w, 201, map[string]any{"id": userData.ID, "email": userData.Email})
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

	code, returnUser := MyDatabase.VerifyUser(user.Email, user.Password)

	RespondWithJSON(w, code, map[string]any{"id": returnUser.ID, "email": returnUser.Email})
}
