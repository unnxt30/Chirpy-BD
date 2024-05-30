package src

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "os"
	"strconv"
	"strings"

	// "github.com/golang-jwt/jwt/v5"
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

	header := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(header, "Bearer ")

	decoder := json.NewDecoder(r.Body)
	params := database.Parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		RespondWithJSON(w, 404, map[string]string{"error": "Something went wrong"})
	} else if database.CheckToken(jwtToken) != -1 {
		auth_id := database.CheckToken(jwtToken)
		// fmt.Println("recieved Chirp Succesfully.")
		chirpBody := GiveValidBody(w, params.Body, params)
		responseChirp, _ := MyDatabase.CreateChirp(chirpBody, auth_id)
		fmt.Println("Created chirp")
		RespondWithJSON(w, 201, responseChirp)
	}
}

func ChirpsGET(w http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")
	order := r.URL.Query().Get("sort")
	MyDatabase, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	if author_id != "" {
		id, _ := strconv.Atoi(author_id)
		chirpsByAuth := MyDatabase.GetChirpByAuthor(id, order)
		RespondWithJSON(w, 200, chirpsByAuth)
	}

	chirpArray, err := MyDatabase.GetChirp(order)
	if err != nil {
		RespondWithJSON(w, 404, map[string]string{"error": err.Error()})
	}

	RespondWithJSON(w, 200, chirpArray)
}

func ChirpGETbyID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	MyDatabase, _ := database.NewDB("database.json")

	chirpArray, _ := MyDatabase.GetChirp("")

	if id > len(chirpArray) {
		RespondWithJSON(w, 404, map[string]string{"error": "ya bish"})
		return
	}

	RespondWithJSON(w, 200, chirpArray[id-1])
}

func DeleteChirp(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	MyDatabse, _ := database.NewDB("database.json")

	header := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(header, "Bearer ")
	chirpArray, _ := MyDatabse.GetChirp("")

	if database.CheckToken(jwtToken) != chirpArray[id-1].Author_id {
		RespondWithJSON(w, 403, map[string]string{"error": "unauthorized"})
		return
	}

	code := MyDatabse.DeleteChirpByID(id)
	RespondWithJSON(w, code, "")

}
