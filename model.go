package friendlymongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ensure BaseModel implements the Model interface
var _ Model = &BaseModel{}

type getters interface {
	GetID() primitive.ObjectID
	GetCreatedAt() time.Time
	GetUpdateAt() time.Time
}

type setters interface {
	SetID()
	UnsetID()
	SetUpdatedAt()
	SetCreatedAt()
}

type Model interface {
	getters
	setters
	Init()
}

type BaseModel struct {
	// ID must have bson tag `omitempty` to allow the document to either have a database generated ID or
	// a custom one, and being elegible for `ReplaceOne`. See method `ReplaceOne` in `BaseRepository` for more info.
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func NewBaseModel() *BaseModel {

	m := &BaseModel{}
	m.Init()

	return m
}

func (b *BaseModel) GetID() primitive.ObjectID {

	return b.ID
}

func (b *BaseModel) UnsetID() {

	b.ID = primitive.NilObjectID
}

func (b *BaseModel) GetCreatedAt() time.Time {

	return b.CreatedAt
}

func (b *BaseModel) GetUpdateAt() time.Time {

	return b.UpdatedAt
}

func (b *BaseModel) SetID() {

	b.ID = primitive.NewObjectID()
}

func (b *BaseModel) SetUpdatedAt() {

	b.UpdatedAt = time.Now()
}

func (b *BaseModel) SetCreatedAt() {

	b.CreatedAt = time.Now()
}

func (b *BaseModel) Init() {

	b.SetID()
	b.SetCreatedAt()
	b.SetUpdatedAt()
}
