package main

import (
	"context"

	"github.com/pmatteo/friendlymongo"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	friendlymongo.SetInstance("mongodb://username:password@localhost:27017")
}

func main() {
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
			Name:     "Rice",
			Category: "Rice&Pasta",
			Ean:      "9244382101364",
			Price:    99,
		},
	}

	prodRepo.InsertMany(context.Background(), products)

	orderRepo.InsertOne(context.Background(), &order{
		Products: []string{
			"5425380924770",
			"9244382101364",
		},
		Status: "paid",
	})

	stages := friendlymongo.NewStageBuilder().
		Match("status_filter", bson.M{"status": "paid"}).
		Lookup("product_lookup", bson.M{
			"from":         "product",
			"localField":   "products",
			"foreignField": "ean",
			"as":           "products",
		}).
		Group("category_group", bson.M{}).
		Project("project", bson.M{})

	result := []struct {
		Products []*product
		Status   string
	}{}
	orderRepo.Aggregate(context.Background(), stages.Build(), result)

}
