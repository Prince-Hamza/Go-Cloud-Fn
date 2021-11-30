package Cors

import (
	//"encoding/json"
	// "github.com/gorilla/handlers"
	// "github.com/gorilla/mux"
	 "net/http"
)
type CorsAccess struct{}


func (cors CorsAccess) Cors(respWriter http.ResponseWriter, req *http.Request) {
	enableCors(&respWriter)

	if req.Method == http.MethodOptions {
		respWriter.Header().Set("Access-Control-Allow-Credentials", "true")
		respWriter.Header().Set("Access-Control-Allow-Headers", "Authorization")
		respWriter.Header().Set("Access-Control-Allow-Methods", "POST")
		respWriter.Header().Set("Access-Control-Allow-Origin", "*")
		respWriter.Header().Set("Access-Control-Max-Age", "3600")
		respWriter.Header().Set("Content-Type", "application/json")
		respWriter.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	respWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	respWriter.Header().Set("Access-Control-Allow-Origin", "*")
	respWriter.Header().Set("Content-Type", "Application/Json")

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}


