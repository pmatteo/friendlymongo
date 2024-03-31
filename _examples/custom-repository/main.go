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

// userProfileRepository is a custom repository for the userProfile model
// It embeds the BaseRepository and it's used to add custom methods to the repository
// It's not mandatory to create a custom repository, but it's a good practice to keep the code organized
// and to avoid code duplication.
// Alternatively, you can alias the BaseRepository and add custom methods to the alias.
type userProfileRepository struct {
	*friendlymongo.BaseRepository[*userProfile]
}

func newUserProfileRepo() userProfileRepository {
	i := friendlymongo.GetInstance()

	return userProfileRepository{
		BaseRepository: friendlymongo.NewBaseRepository(i.Database("test"), "userProfile", new(userProfile)),
	}
}

func init() {
	friendlymongo.SetInstance("mongodb://username:password@localhost:27017")
}

func main() {
	repo := newUserProfileRepo()

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
	update := bson.M{"$set": bson.M{"name": "Jane"}}
	updatedUser, err := repo.UpdateOne(context.Background(), filter, update)
	if err != nil {
		panic(err)
	}

	fmt.Printf("updated user %v\n", updatedUser)

	// Delete the user
	deleted, err := repo.Delete(context.Background(), filter)
	if err != nil {
		panic(err)
	}

	fmt.Printf("deleted count %d\n", deleted)

}
