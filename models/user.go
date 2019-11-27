package models

import (
	"fmt"
	u "sfladmin/utils"

	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
	"github.com/jinzhu/gorm"
)

//User is a struct to rep user account
type User struct {
	ID      int    `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	Email   string `json:"email,omitempty"`
	IsAdmin bool   `gorm:"is_admin" json:"isAdmin,omitempty"`
}

//GetUserID returns the user id for the given email
func GetUserID(email string) int {
	var user User
	err := GetDB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0
		}
		return 0
	}
	return user.ID
}

//Login authenticate
func Login(email string) map[string]interface{} {
	var user User
	fmt.Println(fmt.Sprintf("Logging in user %s", email))
	err := GetDB().Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

//VerifyUser returns an error if the user is not registered on our system
func VerifyUser(email string) (int, error) {
	user := User{}
	GetDB().Table("users").Where("email = ?", email).First(&user)
	if user.Email == "" { //User not found!
		return 0, fmt.Errorf("User not found: %s", email)
	}
	return user.ID, nil
}

func VerifyAdmin(email string) (int, error) {
	user := User{}
	GetDB().Table("users").Where("email = ? and is_admin = true", email).First(&user)
	if user.Email == "" { //User not found!
		return 0, fmt.Errorf("User not found: %s", email)
	}
	return user.ID, nil
}

//TestDB is for testing the database
func TestDB(email string) {
	user := User{}
	err := GetDB().Table("users").Where("email = ?", email).First(&user).Error
	if err != nil {
		fmt.Println("Error", err)
	}
}

//TestToken takes a google authorisation token and verifies it
func TestToken(tok string) (googleAuthIDTokenVerifier.ClaimSet, error) {
	v := googleAuthIDTokenVerifier.Verifier{}
	aud := "948082053040-r1tead48gksuq902m1g4fo4rsk5qj1tu.apps.googleusercontent.com"
	err := v.VerifyIDToken(tok, []string{
		aud,
	})

	if err == nil {
		var claimSet *googleAuthIDTokenVerifier.ClaimSet
		claimSet, err := googleAuthIDTokenVerifier.Decode(tok)
		if err != nil {
			var cs googleAuthIDTokenVerifier.ClaimSet
			fmt.Println("Token not verified against google requirements", err)
			return cs, err
		}
		return *claimSet, nil
	}
	fmt.Println("Token not verified against client account ", err)
	var cs googleAuthIDTokenVerifier.ClaimSet
	return cs, err

}

//ListAdmin lists all the users who are administrators
func ListAdmin() map[string]interface{} {
	users := []User{}
	err := GetDB().Where("is_admin = true").Find(&users).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "No admin users found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	resp := u.Message(true, "Found admin users")
	resp["users"] = users
	return resp
}
