package utils

import (
	"encoding/json"
	"net/http"
)

type key string

//Userid references the postgres userid in context
const (
	UserID key = "userid"
)

//Userkey references the user data in context
const (
	Userkey key = "user"
)

//Message given a status boolean and message string, returns a map of the status and message
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

//Respond json encodes the data and adds a json header
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//GetContext returns the value from the given request and key
func GetContext(w http.ResponseWriter, r *http.Request, akey key) interface{} {
	if value := r.Context().Value(akey); value != nil {
		return value
	} else {
		Respond(w, Message(false, "Context key not found"))
		return nil
	}
}
