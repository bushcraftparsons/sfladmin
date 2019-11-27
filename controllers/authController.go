package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sfladmin/models"
	u "sfladmin/utils"

	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
)

//FileName holds the filename
type FileName struct {
	Name string `json:"name"`
}

//RunShell checks if the request is from an admin, then runs the shell script detailed in the request params
var RunShell = func(w http.ResponseWriter, r *http.Request) {

	var claimSet *googleAuthIDTokenVerifier.ClaimSet
	claimSet, ok := u.GetContext(w, r, u.Userkey).(*googleAuthIDTokenVerifier.ClaimSet)
	if !ok {
		fmt.Printf("%T\n", u.GetContext(w, r, u.Userkey))
		u.Respond(w, u.Message(false, "Failed to get user context"))
	}
	email := claimSet.Email
	if models.VerifyAdmin(email) {
		//User is admin, get the name of the shell script and run it
		fileName := &FileName{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&fileName)
		if err != nil {
			fmt.Println("Error", err)
			u.Respond(w, u.Message(false, "Failed decoding to filename struct"))
			return
		} else {
			out, err := exec.Command(fileName.Name).Output()
			if err != nil {
				log.Fatal(err)
			}
			u.Respond(w, u.Message(true, fmt.Printf("output is %s\n", out)))
		}
	} else {
		u.Respond(w, u.Message(false, fmt.Printf("User is not admin %s", email)))
	}
	u.Respond(w, resp)
}

//ListAdmin lists all the administrators
var ListAdmin = func(w http.ResponseWriter, r *http.Request) {
	resp := models.ListAdmin()
	u.Respond(w, resp)
}