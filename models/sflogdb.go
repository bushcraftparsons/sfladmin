package models

import (
	"fmt"
	"os"

	//See example at https://gorm.io/docs/connecting_to_the_database.html
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB //database

func init() {
	//Loads the .env file with db connection info
	//This method gets called before the main.go main() function runs
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print("Failed to load environment variables: ", e)
	} else {
		fmt.Println(fmt.Sprintf("Loaded environment variables %s", os.Getenv("user")))
	}

	username := os.Getenv("user")
	password := os.Getenv("password")
	dbName := os.Getenv("dbname")
	dbHost := os.Getenv("host")
	port := os.Getenv("port")
	sslrootcert := os.Getenv("sslrootcert")
	sslmode := os.Getenv("sslmode")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslrootcert=%s sslmode=%s",
		dbHost, port, username, dbName, password, sslrootcert, sslmode) //Build connection string
	fmt.Println("Connecting to db")
	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Print("Failed to connect to db: ", err)
		fmt.Println(dbURI)
	}

	db = conn
	//db.Debug().AutoMigrate(&Account{}, &Contact{}) //Database migration
}

//GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
