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

func Push(exp ...interface{}) bson.M {
	if len(exp) == 1 {
		return bson.M{"$push": exp[0]}
	}

	if len(exp)%2 != 0 {
		panic("Push expects an even number of arguments")
	}

	push := bson.M{}
	for i := 0; i < len(exp); i += 2 {
		push[exp[i].(string)] = exp[i+1]
	}
	return bson.M{"$push": push}
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
