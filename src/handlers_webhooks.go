package src

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	database "github.com/unnxt30/Chirpy-BD/internal"
)

func CheckUpgradedUser(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("POLKA_KEY")
	myDatabase, _ := database.NewDB("database_user.json")

	header := r.Header.Get("Authorization")
	Key := strings.TrimPrefix(header, "ApiKey ")

	if Key != apiKey {
		RespondWithJSON(w, 401, "")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := database.WebhookStruct{}
	decoder.Decode(&params)

	code := myDatabase.UpgradeUser(params)

	RespondWithJSON(w, code, "")
}
