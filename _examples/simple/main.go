package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pmatteo/friendlymongo"
	"go.mongodb.org/mongo-driver/bson"
)

type userProfile struct {
	*friendlymongo.BaseModel

	Name      string
	Surname   string
	Email     string
	BirthDate time.Time
}

func (u *userProfile) String() string {
	return fmt.Sprintf(
		"ID: %s, Name: %s, Surname: %s, Email: %s, BirthDate: %s",
		u.ID.Hex(), u.Name, u.Surname, u.Email, u.BirthDate.Local().Format("2006-01-02"),
	)
}

func init() {
	friendlymongo.SetInstance("mongodb://username:password@localhost:27017")
}

func main() {
	db := friendlymongo.GetInstance().Database("user")
	repo := friendlymongo.NewBaseRepository(db, "userProfile", &userProfile{})

	// Create an instance of the UserProfile
	birthday := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)
	user := &userProfile{
		BaseModel: friendlymongo.NewBaseModel(),
		Name:      "John",
		Surname:   "Doe",
		Email:     "john.doe@test.com",
		BirthDate: birthday,
	}

	// Insert the user into the database
	err := repo.InsertOne(context.Background(), user)
	if err != nil {
		panic(err)
	}

	// Find the user by email
	filter := bson.M{"email": "john.doe@test.com"}
	foundUser, err := repo.FindOne(context.Background(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("found user %v\n", foundUser)

	// Update the user
	update := bson.M{"$set": bson.M{"name": "updateOne bson"}}
	updatedUser, err := repo.UpdateOne(context.Background(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Printf("updated user with bson %v\n", updatedUser)

	foundUser.Name = "updateOne model"
	updatedUser, err = repo.UpdateOne(context.Background(), filter, foundUser)
	if err != nil {
		panic(err)
	}

	fmt.Printf("updated user with model %v\n", updatedUser)

	// Replace the user
	foundUser.Name = "replace"
	err = repo.ReplaceOne(context.Background(), filter, foundUser)
	if err != nil {
		panic(err)
	}

	fmt.Printf("replaced user with model %v\n", foundUser)

	// Delete the user
	deleted, err := repo.Delete(context.Background(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("deleted count %d\n", deleted)

}
