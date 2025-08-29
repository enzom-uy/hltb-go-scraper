package helpers

import (
	"encoding/json"
	"net/http"
)

func SendJSONSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// func SendJSONError(w http.ResponseWriter, errorMessage string) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusBadRequest)
// 	json.NewEncoder(w)
// }
