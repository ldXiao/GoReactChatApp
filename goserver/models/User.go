package models

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/ldXiao/GoReactChatApp/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User is a struct in model
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Email    string             `json:"email,omitempty"`
	Password string             `json:"password,omitempty"`
	LastName string             `json:"lastname, omitempty"`
	Role     int                `json:"role",omitempty` //0 for ordinary user 1 for amin
	Image    string             `json:"image", omotiempty`
	Token    string             `json:"token, omitempty"`
	TokenExp string             `json:"tokenExp, omitempty"`
}

//Save is exported
func (u *User) Save() bool {
	if u.Password != "" {
		password := []byte(u.Password)
		fmt.Println("input ", u.Password)
		hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		// the salt is automatically generated
		if err != nil {
			log.Println(err)
		}
		u.Password = string(hash)
		u.GenerateToken()

		fmt.Println(u.Password)
		_, err = middleware.UsersCollection.InsertOne(context.Background(), u)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	return false
}

func (u *User) UpdateToken() error {
	// find the unique user with the email and only change the token
	// TODO maybe add a method to update multiple file at once
	fmt.Println("calledUpdateToken", u.Token)
	singres := middleware.UsersCollection.FindOneAndUpdate(
		context.TODO(),
		bson.D{{"email", u.Email}},
		bson.D{{"$set", bson.D{{"token", u.Token}}}},
	)
	return singres.Err()
}

// ComparePassword is exported
func (u *User) ComparePassword(plainPassword string) bool {
	fmt.Println(u.Password)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		return false
	}
	return true

}

func (u *User) GenerateToken() {
	os.Setenv("ACCESS_SECRET", "secret") //this should be in an env file
	atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
	atClaims["user_id"] = u.ID
	// atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		log.Fatal(err)
	}
	u.Token = token
}
