package p

import (
	"encoding/json"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {

	type Bird struct {
		Species     string
		Description string
	}

	birdJson := `{"species": "pigeon","description": "likes to perch on rocks"}`
	var bird Bird
	json.Unmarshal([]byte(birdJson), &bird)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bird)

}
