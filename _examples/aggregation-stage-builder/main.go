package main

import (
	"context"
	"fmt"

	"github.com/pmatteo/friendlymongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	friendlymongo.SetInstance("mongodb://username:password@localhost:27017")
}

func main() {
	defer func() {
		i := friendlymongo.GetInstance()

		if err := i.Database("test").Drop(context.Background()); err != nil {
			fmt.Println("Error dropping database:", err)
		}

		if err := i.Client().Disconnect(context.Background()); err != nil {
			fmt.Println("Error disconnecting client:", err)
		}
	}()

	prodRepo := newProductRepository()
	orderRepo := newOrderRepository()

	products := []*product{
		{
			Name:     "Barialla Pennette Rigate",
			Category: "Rice&Pasta",
			Ean:      "5425380924770",
			Price:    122,
		},
		{
			Name:     "Strawberry",
			Category: "Fresh Products",
			Ean:      "5424382921726",
			Price:    150,
		},
		{
			Name:     "Apple",
			Category: "Fresh Products",
			Ean:      "54243372821722",
			Price:    95,
		},
		{
			Name:     "Rice",
			Category: "Rice&Pasta",
			Ean:      "9244382101364",
			Price:    99,
		},
	}
	prodRepo.InsertMany(context.Background(), products)

	orders := []*order{
		{
			Products: []string{
				"5425380924770",
				"54243372821722",
				"5424382921726",
			},
			Status: "paid",
		},
		{
			Products: []string{
				"5425380924770",
				"9244382101364",
				"54243372821722",
			},
			Status: "paid",
		},
		{
			Products: []string{
				"5425380924770",
				"9244382101364",
			},
			Status: "delivered",
		},
	}
	orderRepo.InsertMany(context.Background(), orders)

	stages := friendlymongo.NewStageBuilder().
		Match("status_filter", bson.M{"status": "paid"}).
		Lookup("product_lookup", "product", "products", "ean", "products").
		Unwind("unwind_products", "$products").
		Group("category_id_group", bson.M{
			"_id": bson.M{
				"orderId":  "$_id",
				"category": "$products.category",
			},
			"status":   bson.M{"$first": "$status"},
			"products": bson.M{"$push": "$products"},
		}).
		Group("productsByCategory", bson.M{
			"_id":    "$_id.orderId",
			"status": bson.M{"$first": "$status"},
			"productsByCategory": bson.M{"$push": bson.M{
				"category": "$_id.category",
				"products": "$products",
			}},
		}).
		Project("final_project", bson.M{
			"_id":    1,
			"status": 1,
			"grouped_products": bson.M{
				"$arrayToObject": bson.M{
					"$map": bson.M{
						"input": "$productsByCategory",
						"as":    "cat",
						"in": bson.M{
							"k": "$$cat.category",
							"v": "$$cat.products",
						},
					},
				},
			},
		})

	results := []struct {
		ProductsByCategory map[string]*product `bson:"grouped_products"`
		Status             string              `bson:"status"`
		Orderid            primitive.ObjectID  `bson:"_id"`
	}{}
	orderRepo.Aggregate(context.Background(), stages.Build(), results)

	fmt.Println("Aggregation result:", results)

	for _, r := range results {
		fmt.Printf("Status: %s - Products: %v\n", r.Status, r.ProductsByCategory)
	}
}

type product struct {
	friendlymongo.BaseModel `bson:",inline"`

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

type order struct {
	friendlymongo.BaseModel `bson:",inline"`

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
