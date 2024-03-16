package friendlymongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ MongoRepository[Model] = &BaseRepository[Model]{}

// MongoRepository is an interface that defines the methods for interacting with a MongoDB collection.
type MongoRepository[T Model] interface {
	// InsertOne inserts a single document into the collection.
	InsertOne(ctx context.Context, document Model) error

	// InsertMany inserts multiple documents into the collection.
	InsertMany(ctx context.Context, documents []Model) error

	// FindOne finds a single document in the collection.
	FindOne(ctx context.Context, filter interface{}) (Model, error)

	// Find finds multiple documents in the collection.
	Find(ctx context.Context, filter interface{}) ([]Model, error)

	// UpdateOne finds a single document and updates it.
	UpdateOne(ctx context.Context, filters interface{}, update interface{}) (Model, error)

	// Delete deletes multiple documents from the collection.
	Delete(ctx context.Context, filter interface{}) (int64, error)

	// ReplaceOne replaces a single document in the collection.
	ReplaceOne(ctx context.Context, filter interface{}, replacement T) error

	// Aggregate runs an aggregation framework pipeline on the collection.
	Aggregate(ctx context.Context, pipeline mongo.Pipeline, result interface{}) error
}

// BaseRepository is a base implementation of the MongoRepository interface.
type BaseRepository[T Model] struct {
	collection *mongo.Collection
}

// NewBaseRepository creates a new instance of BaseRepository.
func NewBaseRepository[T Model](db *mongo.Database, collectionName string, t T) *BaseRepository[T] {

	return &BaseRepository[T]{
		collection: db.Collection(collectionName),
	}
}

// InsertOne inserts a single document into the collection.
//
// The document parameter must be a pointer to a struct that implements the Model interface.
func (r *BaseRepository[T]) InsertOne(ctx context.Context, document T) error {
	document.Init()

	_, err := r.collection.InsertOne(ctx, document)
	return err
}

// InsertMany inserts multiple documents into the collection.
func (r *BaseRepository[T]) InsertMany(ctx context.Context, documents []T) error {

	var interfaceSlice = make([]interface{}, len(documents))
	for i, d := range documents {
		if d.GetID().IsZero() {
			d.Init()
		}

		interfaceSlice[i] = d
	}

	_, err := r.collection.InsertMany(ctx, interfaceSlice)
	return err
}

// FindOne finds a single document in the collection.
func (r *BaseRepository[T]) FindOne(ctx context.Context, filter interface{}) (T, error) {

	var document T

	err := r.collection.FindOne(ctx, filter).Decode(&document)

	return document, err
}

// Find finds multiple documents in the collection.
func (r *BaseRepository[T]) Find(ctx context.Context, filter interface{}) ([]T, error) {

	var documents []T

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var document T
		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}

	return documents, nil
}

// FindOneAndUpdate finds a single document and updates it.
// The update parameter must be a bson.M or a struct that implements the Model interface.
func (r *BaseRepository[T]) UpdateOne(ctx context.Context, filters interface{}, update interface{}) (T, error) {
	var document T
	var updateQuery bson.M

	switch u := update.(type) {
	case T:
		u.SetUpdatedAt()
		u.UnsetID()
		updateQuery = bson.M{"$set": u}
	case bson.M:
		u["$currentDate"] = bson.M{"updatedAt": true}
		updateQuery = u
	default:
		return document, fmt.Errorf("update parameter must be a bson.M or a Model")
	}

	singleRes := r.collection.FindOneAndUpdate(ctx, filters, updateQuery)
	if singleRes.Err() != nil {
		return document, singleRes.Err()
	}

	err := singleRes.Decode(&document)

	return document, err
}

// Delete deletes multiple documents from the collection.
func (r *BaseRepository[T]) Delete(ctx context.Context, filter interface{}) (int64, error) {

	deleteRes, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return deleteRes.DeletedCount, err
}

// Aggregate runs an aggregation framework pipeline on the collection.
func (r *BaseRepository[T]) Aggregate(ctx context.Context, pipeline mongo.Pipeline, result interface{}) error {

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

// FindOneReplaceOne replaces a single document in the collection.
// Replaced document must have the same ID as the one being replaced or not have it serializible at all. It is
// strongly suggested to have the ID field with the `omitempty` bson tag in case of structs.
func (r *BaseRepository[Model]) ReplaceOne(ctx context.Context, filter interface{}, replacement Model) error {

	replacement.SetUpdatedAt()
	replacement.UnsetID()

	singleRes := r.collection.FindOneAndReplace(ctx, filter, replacement)
	if singleRes.Err() != nil {
		return singleRes.Err()
	}

	return nil
}
