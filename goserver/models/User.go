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

	password := []byte(u.Password)
	fmt.Println("register called2")
	fmt.Println("input ", u.Password)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	fmt.Println("register called1")
	// the salt is automatically generated
	if err != nil {
		log.Println(err)
	}
	u.Password = string(hash)
	u.GenerateToken()
	fmt.Println("register called")
	fmt.Println(u.Password)
	fmt.Println("register called")
	_, err = middleware.UsersCollection.InsertOne(context.Background(), u)
	fmt.Println("register called")
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := "secret"
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func (u *User) LoadByToken(tok string) bool {
	claim, _ := extractClaims(tok)
	fmt.Println("acfaew")
	id := claim["user_id"].(string)
	fmt.Println(id)
	fmt.Println("acfaew")
	objid, _ := primitive.ObjectIDFromHex(id)
	singres := middleware.UsersCollection.FindOne(
		context.TODO(),
		bson.D{{"_id", objid}, {"token", tok}},
	)

	if singres.Err() == nil {
		singres.Decode(&u)

		fmt.Println("load by token", u.Email)
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
	id := u.ID.Hex()

	// atClaims["authorized"] = true
	atClaims["user_id"] = id
	// atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		log.Fatal(err)
	}
	u.Token = token
}
