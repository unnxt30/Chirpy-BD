package src

import (
	"encoding/json"
	"fmt"
	"net/http"
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
