package main

import (
	"github.com/pmatteo/friendlymongo"
)

type order struct {
	*friendlymongo.BaseModel

	Products []string `bson:"products"`
	Status   string   `bson:"status"`
}

// orderRepository is a custom repository for the order model
// It embeds the BaseRepository and it's used to add custom methods to the repository
// It's not mandatory to create a custom repository, but it's a good practice to keep the code organized
// and to avoid code duplication.
// Alternatively, you can alias the BaseRepository and add custom methods to the alias.
type orderRepository struct {
	*friendlymongo.BaseRepository[*order]
}

func newOrderRepository() orderRepository {
	i := friendlymongo.GetInstance()

	return orderRepository{
		BaseRepository: friendlymongo.NewBaseRepository(i.Database("test"), "order", new(order)),
	}
}
