package friendlymongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ensure BaseModel implements the Model interface
var _ Model = &BaseModel{}

type Model interface {
	OnCreate()
	OnUpdate()
	OnReplace()
}

type BaseModel struct {
	// ID must have bson tag `omitempty` to allow the document to either have a database generated ID or
	// a custom one, and being elegible for `ReplaceOne`. See method `ReplaceOne` in `BaseRepository` for more info.
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func NewBaseModel() *BaseModel {
	now := time.Now()

	return &BaseModel{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (b *BaseModel) OnCreate() {
	if b.ID.IsZero() {
		b.setID(primitive.NewObjectID())
	}

	b.setCreatedAt()
	b.setUpdatedAt()
}

func (b *BaseModel) OnReplace() {

	b.setID(primitive.NilObjectID)
	b.setUpdatedAt()
}

func (b *BaseModel) OnUpdate() {

	b.setID(primitive.NilObjectID)
	b.setUpdatedAt()
}

func (b *BaseModel) setID(id primitive.ObjectID) {

	b.ID = id
}

func (b *BaseModel) setUpdatedAt() {

	b.UpdatedAt = time.Now()
}

func (b *BaseModel) setCreatedAt() {

	b.CreatedAt = time.Now()
}
