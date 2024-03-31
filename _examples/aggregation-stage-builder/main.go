package main

import (
	"context"
	"fmt"

	fm "github.com/pmatteo/friendlymongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	fm.SetInstance("mongodb://username:password@localhost:27017")
}

func main() {

	// Clean up after the example
	defer func() {
		i := fm.GetInstance()

		if err := i.Database("test").Drop(context.Background()); err != nil {
			fmt.Println("Error dropping database:", err)
		}

		if err := i.Client().Disconnect(context.Background()); err != nil {
			fmt.Println("Error disconnecting client:", err)
		}
	}()

	prodRepo := newProductRepository()
	orderRepo := newOrderRepository()

	setupOrderAndProducts(prodRepo, orderRepo)

	// The aggregation pipeline returns the orders with the status "paid" with their products
	// grouped the by category
	stages := fm.NewStageBuilder().
		Match("status_filter", bson.M{"status": "paid"}).
		Lookup("product_lookup", "product", "products", "ean", "products").
		Unwind("unwind_products", "$products").
		Group("category_id_group", bson.M{
			"_id": bson.M{
				"orderId":  "$_id",
				"category": "$products.category",
			},
			"status":   fm.First("$status"),
			"products": fm.Push("$products"),
		}).
		Group("productsByCategory", bson.M{
			"_id":                "$_id.orderId",
			"status":             fm.First("$status"),
			"productsByCategory": fm.Push("category", "$_id.category", "products", "$products"),
		}).
		Project("final_project", bson.M{
			"_id":    1,
			"status": 1,
			"grouped_products": fm.ArrayToObject(
				fm.Map("$productsByCategory", "cat", "$$cat.category", "$$cat.products"),
			),
		})

	var results []result
	err := orderRepo.Aggregate(context.Background(), stages.Build(), &results)
	if err != nil {
		fmt.Println("Error aggregating orders:", err)
		return
	}

	printResults(results)
}

type result struct {
	GroupedProducts map[string][]*product `bson:"grouped_products"`
	Status          string                `bson:"status"`
	OrderId         primitive.ObjectID    `bson:"_id"`
}

func printResults(results []result) {
	for _, r := range results {
		fmt.Printf("Order ID: %s\n", r.OrderId.String())
		fmt.Printf("Order status: %s\n", r.Status)
		fmt.Println("Products:")
		for cat, prods := range r.GroupedProducts {
			fmt.Printf("\tCategory: %s\n", cat)
			for _, p := range prods {
				fmt.Printf("\t\t * %s: %.2f\n", p.Name, float64(p.Price)/100)
			}
		}
		fmt.Println()
	}
}

func setupOrderAndProducts(prodRepo productRepository, orderRepo orderRepository) {
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
}

type product struct {
	fm.BaseModel `bson:",inline"`

	Name     string `bson:"name"`
	Category string `bson:"category"`
	Ean      string `bson:"ean"`
	Price    int    `bson:"price"`
}

type productRepository struct {
	*fm.BaseRepository[*product]
}

func newProductRepository() productRepository {
	i := fm.GetInstance()

	return productRepository{
		BaseRepository: fm.NewBaseRepository(i.Database("test"), "product", new(product)),
	}
}

type order struct {
	fm.BaseModel `bson:",inline"`

	Products []string `bson:"products"`
	Status   string   `bson:"status"`
}

type orderRepository struct {
	*fm.BaseRepository[*order]
}

func newOrderRepository() orderRepository {
	i := fm.GetInstance()

	return orderRepository{
		BaseRepository: fm.NewBaseRepository(i.Database("test"), "order", new(order)),
	}
}
