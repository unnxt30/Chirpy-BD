package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type parameters struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}



func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {

	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.WriteHeader(code)

	w.Write(response)

	return nil

}

func validateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithJSON(w, 404, map[string]string{"error": "Something went wrong"})
	}

	if len(params.Body) > 140 {
		respondWithJSON(w, 400, map[string]string{"error": "Chirp is too long"})
	} else {
		// found := false
		// var cleaned_body_words []string;
		// bad_words := []string{"kerfuffle", "sharbert", "fornax"};
		// //respondWithJSON(w, 200, map[string]bool{"valid" : true})
		// input_str := params.Body;
		// for _,word := range strings.Split(input_str, " "){
		// 	found = false;
		// 	for _, bad_word := range bad_words {
		// 		if strings.ToLower(word) == bad_word{
		// 			found = true;
		// 			break;
		// 		}
		// 	}
		// 	if found{
		// 			cleaned_body_words = append(cleaned_body_words, "****");
		// 	}else{
		// 		cleaned_body_words = append(cleaned_body_words, word);
		// 	}
		// }
		//cleaned_body := strings.Join(cleaned_body_words, " ");

		params.ID = idCount
		idCount += 1
		write_data, err := json.Marshal(params)
		if err != nil {
			fmt.Println(err)
		}
		os.WriteFile("database.json", write_data, 0644)
		respondWithJSON(w, 200, map[string]string{"body": params.Body})
	}

}
