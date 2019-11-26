package controllers

import (
	"fmt"
	"net/http"
	"sfladmin/models"
	u "sfladmin/utils"

	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
)

//Authenticate sends the request to be authenticated
var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	var claimSet *googleAuthIDTokenVerifier.ClaimSet
	claimSet, ok := u.GetContext(w, r, u.Userkey).(*googleAuthIDTokenVerifier.ClaimSet)
	if !ok {
		fmt.Printf("%T\n", u.GetContext(w, r, u.Userkey))
		u.Respond(w, u.Message(false, "Failed to get user context"))
	}
	email := claimSet.Email
	resp := models.Login(email)
	u.Respond(w, resp)
}

//ListAdmin lists all the administrators
var ListAdmin = func(w http.ResponseWriter, r *http.Request) {
	resp := models.ListAdmin()
	u.Respond(w, resp)
}
