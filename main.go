package main

import (
	"log"
	"net/http"
	"sfladmin/app"
	"sfladmin/controllers"

	u "sfladmin/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	resp := u.Message(true, "Welcome to sfl!")
	resp["user"] = "Susannah Parsons"
	u.Respond(w, resp)
}

//Curl test
/*
 curl -H "Origin: http://localhost:3001" \
 -H "Access-Control-Request-Method: POST" \
 -H "Access-Control-Request-Headers: X-Requested-With, Authorization" \
 -X OPTIONS --verbose http://localhost:8080/login
*/
func main() {
	//Update security group for DB with current ip if running on dev
	//https://console.aws.amazon.com/ec2/v2/home?region=us-east-1#SecurityGroups:search=sg-1cdc4442;sort=groupId

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) //attach JWT auth middleware
	router.HandleFunc("/", homeLink).Methods("GET")
	router.HandleFunc("/listAdmin", controllers.ListAdmin).Methods("GET")
	router.HandleFunc("/runShell", controllers.RunShell).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Accept", "Origin", "Referer", "Sec-Fetch-Mode", "User-Agent"})
	originsOk := handlers.AllowedOrigins([]string{"*", "http://localhost:3001", "https://sfl.formyer.com", "http://86.184.53.15:3001"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// start server listen
	// with error handling
	log.Fatal(http.ListenAndServeTLS(":6060", "apiserver.crt", "apiserver.key", handlers.CORS(originsOk, headersOk, methodsOk)(router))) //Launch the app, visit https://localhost:6060/api
}
