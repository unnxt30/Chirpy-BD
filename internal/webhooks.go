package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type WebhookStruct struct {
	Event string         `json:"event"`
	Data  map[string]int `json:"data"`
}

var upgradeQuery string = "user.upgraded"

func (db *DB) UpgradeUser(data WebhookStruct) int {
	if data.Event != upgradeQuery {
		return 204
	}

	readData, _ := os.ReadFile(db.path)
	params := &DBStructUsers{}
	json.Unmarshal(readData, &params)

	id := data.Data["user_id"]
	// fmt.Println(id)
	found := false
	for i := 0; i < len(params.Users); i++ {
		currentUser := params.Users[i]
		if params.Users[i].ID == id {
			var newUser User
			newUser.ID = currentUser.ID
			newUser.Chirpy_red = true
			newUser.Email = currentUser.Email
			newUser.Password = currentUser.Password
			newUser.ExpirationTime = currentUser.ExpirationTime

			params.Users[i] = newUser
			found = true;
		}
	}

	fmt.Println(params);

	fmt.Printf("%v", found)
	if !found {
		return 403
	}

	os.Remove(db.path)
	openOrCreateFile(db.path)
	writeData, _ := json.Marshal(params)
	os.WriteFile(db.path, writeData, 0666)
	return 204

}
