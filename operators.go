package friendlymongo

import "go.mongodb.org/mongo-driver/bson"

func Map(input, as, k, v string) bson.M {
	return bson.M{
		"$map": bson.M{
			"input": input,
			"as":    as,
			"in": bson.M{
				"k": k,
				"v": v,
			},
		},
	}
}

func ArrayToObject(exp interface{}) bson.M {
	return bson.M{"$arrayToObject": exp}
}

func Push(exp interface{}) bson.M {
	return bson.M{"$push": exp}
}

func First(exp interface{}) bson.M {
	return bson.M{"$first": exp}
}

func Each(values interface{}) bson.M {
	return bson.M{"$each": values}
}

func AddToSet(exp map[string]interface{}) bson.M {
	addToSet := bson.M{}
	for k, v := range exp {
		addToSet[k] = v
	}

	return bson.M{"$addToSet": addToSet}
}
