package friendlymongo_test

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/pmatteo/friendlymongo"
)

var uri string
var testDB string
var testCollection string

var repo *friendlymongo.BaseRepository[*customModel]
var artworksRepo *friendlymongo.BaseRepository[*artwork]

func newCustomModelRepo() *friendlymongo.BaseRepository[*customModel] {
	i := friendlymongo.GetInstance()

	return friendlymongo.NewBaseRepository(i.Database(testDB), testCollection, new(customModel))
}

func newOtherModelRepo() *friendlymongo.BaseRepository[*otherModel] {
	i := friendlymongo.GetInstance()

	return friendlymongo.NewBaseRepository(i.Database(testDB), "otherCollection", new(otherModel))
}

func newArtworkRepo() *friendlymongo.BaseRepository[*artwork] {
	i := friendlymongo.GetInstance()

	return friendlymongo.NewBaseRepository(i.Database(testDB), "artwork", new(artwork))
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
	artworksRepo = newArtworkRepo()

	initArtwork()

	m.Run()
}

func cleanUp() {
	i := friendlymongo.GetInstance()

	// CleanUp database
	err := i.Database(testDB).Drop(context.Background())
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

func initArtwork() {
	var artworks = []*artwork{
		{
			Title: "The Pillars of Society", Artist: "Grosz", Year: 1926, Price: 199.99,
			Tags: []string{"painting", "satire", "Expressionism", "caricature"},
		},
		{
			Title: "Melancholy III", Artist: "Munch", Year: 1902, Price: 280.00,
			Tags: []string{"woodcut", "Expressionism"},
		},
		{
			Title: "Dancer", Artist: "Miro", Year: 1925, Price: 76.04,
			Tags: []string{"oil", "Surrealism", "painting"},
		},
		{
			Title: "The Great Wave off Kanagawa", Artist: "Hokusai", Price: 167.30,
			Tags: []string{"woodblock", "ukiyo-e"},
		},
		{
			Title: "The Persistence of Memory", Artist: "Dali", Year: 1931, Price: 483.00,
			Tags: []string{"Surrealism", "painting", "oil"},
		},
		{
			Title: "Composition VII", Artist: "Kandinsky", Year: 1913, Price: 385.00,
			Tags: []string{"oil", "painting", "abstract"},
		},
	}

	artworksRepo.InsertMany(context.Background(), artworks)
}
