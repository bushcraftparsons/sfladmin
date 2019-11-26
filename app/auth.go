package app

import (
	"context"
	"fmt"
	"net/http"
	m "sfladmin/models"
	u "sfladmin/utils"
)

//JwtAuthentication Given an http.Handler, returns the same handler with suitable functions for authentication incorporated
var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/", "/listAdmin"} //List of endpoints that don't require auth
		requestPath := r.URL.Path              //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{}) //initialise
		token := r.Header.Get("Authorization")   //Grab the Google token from the header

		if token == "" { //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		cs, err := m.TestToken(token)
		if err != nil { //Malformed token, returns with http code 403 as usual
			fmt.Println("Malformed authentication token", err)
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		claimSet := &cs //Access memory location of pointer
		userid, err := m.VerifyUser(claimSet.Email)
		if err != nil { //User not registered, returns with http code 403 as usual
			fmt.Println("User not registered", err)
			response = u.Message(false, "User not registered")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), u.Userkey, claimSet)
		ctx = context.WithValue(ctx, u.UserID, userid)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
