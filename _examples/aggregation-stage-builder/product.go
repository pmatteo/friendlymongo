package main

import (
	"github.com/pmatteo/friendlymongo"
)

type product struct {
	*friendlymongo.BaseModel

	Name     string `bson:"name"`
	Category string `bson:"category"`
	Ean      string `bson:"ean"`
	Price    int    `bson:"price"`
}

// productRepository is a custom repository for the product model
// It embeds the BaseRepository and it's used to add custom methods to the repository
// It's not mandatory to create a custom repository, but it's a good practice to keep the code organized
// and to avoid code duplication.
// Alternatively, you can alias the BaseRepository and add custom methods to the alias.
type productRepository struct {
	*friendlymongo.BaseRepository[*product]
}

func newProductRepository() productRepository {
	i := friendlymongo.GetInstance()

	return productRepository{
		BaseRepository: friendlymongo.NewBaseRepository(i.Database("test"), "product", new(product)),
	}
}
