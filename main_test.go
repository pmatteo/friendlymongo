package friendlymongo_test

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/pmatteo/friendlymongo"
	"go.mongodb.org/mongo-driver/bson"
)

var uri string
var testDB string
var testCollection string

var repo *friendlymongo.BaseRepository[*customModel]

func newCustomModelRepo() *friendlymongo.BaseRepository[*customModel] {
	i := friendlymongo.GetInstance()

	return friendlymongo.NewBaseRepository(i.Database(testDB), testCollection, new(customModel))
}

func newOtherModelRepo() *friendlymongo.BaseRepository[*otherModel] {
	i := friendlymongo.GetInstance()

	return friendlymongo.NewBaseRepository(i.Database(testDB), "otherCollection", new(otherModel))
}

func TestMain(m *testing.M) {
	flag.StringVar(&uri, "uri", "mongodb://root:toor@localhost:27017", "MongoDB URI")
	flag.StringVar(&testDB, "db", "testDatabase", "database name")
	flag.StringVar(&testCollection, "collection", "testCollection", "collection name")

	flag.Parse()

	if friendlymongo.GetInstance() != nil {
		panic("Instance already set")
	}

	i := friendlymongo.SetInstance(uri)
	err := i.Connect()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %v", err))
	}
	defer disconnect()
	defer cleanUp()

	repo = newCustomModelRepo()

	m.Run()
}

func cleanUp() {
	i := friendlymongo.GetInstance()

	// CleanUp database
	_, err := i.Database(testDB).Collection(testCollection).DeleteMany(context.Background(), bson.D{})
	if err != nil {
		fmt.Println("Error deleting documents:", err)
		return
	}
}

func disconnect() {
	i := friendlymongo.GetInstance()
	err := i.Disconnect()
	if err != nil {
		fmt.Println("Error disconnecting from database:", err)
	}
}
