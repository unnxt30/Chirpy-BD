package src

import (
	"net/http"
	"strings"

	database "github.com/unnxt30/Chirpy-BD/internal"

)


func CheckRefToken(w http.ResponseWriter, r *http.Request){
	tokenString := r.Header.Get("Authorization");
	trimmedTokenString := strings.TrimPrefix(tokenString, "Bearer ");
	code, returnVal := database.ValidateRefToken(trimmedTokenString);
	RespondWithJSON(w, code, map[string]string{"token":returnVal});
}

func RevokeToken(w http.ResponseWriter, r *http.Request){
	tokenString:= r.Header.Get("authorization");
	trimmedTokenString := strings.TrimPrefix(tokenString, "Bearer ");
	code := database.DeleteRefToken(trimmedTokenString);
	RespondWithJSON(w, code, "");
}