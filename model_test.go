package friendlymongo_test

import (
	"testing"
	"time"

	"github.com/pmatteo/friendlymongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBaseModel_OnCreate(t *testing.T) {
	t.Parallel()

	model := &friendlymongo.BaseModel{}

	assert.True(t, model.ID.IsZero())
	assert.Empty(t, model.CreatedAt)
	assert.Empty(t, model.UpdatedAt)

	model.OnCreate()

	assert.False(t, model.ID.IsZero())
	assert.NotEmpty(t, model.CreatedAt)
	assert.NotEmpty(t, model.UpdatedAt)
}

func TestBaseModel_OnUpdate(t *testing.T) {
	t.Parallel()

	tm := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	id := primitive.NewObjectID()
	model := &friendlymongo.BaseModel{
		ID:        id,
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	assert.Equal(t, id, model.ID)
	assert.Equal(t, tm, model.CreatedAt)
	assert.Equal(t, tm, model.UpdatedAt)

	model.OnUpdate()

	assert.Equal(t, primitive.NilObjectID, model.ID)
	assert.Equal(t, tm, model.CreatedAt)
	assert.NotEqual(t, tm, model.UpdatedAt)
}

func TestBaseModel_OnReplace(t *testing.T) {
	t.Parallel()

	tm := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	id := primitive.NewObjectID()
	model := &friendlymongo.BaseModel{
		ID:        id,
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	assert.Equal(t, id, model.ID)
	assert.Equal(t, tm, model.CreatedAt)
	assert.Equal(t, tm, model.UpdatedAt)

	model.OnReplace()

	assert.Equal(t, primitive.NilObjectID, model.ID)
	assert.Equal(t, tm, model.CreatedAt)
	assert.NotEqual(t, tm, model.UpdatedAt)
}

type address struct {
	Street string `bson:"street"`
	Number int    `bson:"number"`
	City   string `bson:"city"`
}
type customModel struct {
	friendlymongo.BaseModel `bson:",inline"`

	Name           string     `bson:"name"`
	Email          string     `bson:"email"`
	Active         bool       `bson:"active"`
	ActivationDate *time.Time `bson:"activationDate"`
	Address        *address   `bson:"address"`
}

type otherModel struct {
	friendlymongo.BaseModel

	FK primitive.ObjectID `bson:"fk"`

	Group       string   `bson:"group"`
	Permissions []string `bson:"permissions"`
}

func newCustomModel(name, email string, active bool, a *address) *customModel {
	m := &customModel{
		Name:    name,
		Email:   email,
		Active:  active,
		Address: a,
	}
	if active {
		now := time.Now()
		m.ActivationDate = &now
	}

	return m
}

func TestCustomModel_OnCreate(t *testing.T) {

	t.Parallel()

	model := newCustomModel("John Doe", "test@test", true, basicAddress)

	assert.True(t, model.ID.IsZero())
	assert.Empty(t, model.CreatedAt)
	assert.Empty(t, model.UpdatedAt)

	model.OnCreate()

	assert.NotEqual(t, primitive.NilObjectID, model.ID)
	assert.NotEmpty(t, model.CreatedAt)
	assert.NotEmpty(t, model.UpdatedAt)
}

func TestCustomModel_OnUpdate(t *testing.T) {

	t.Parallel()

	model := newCustomModel("John Doe", "test@test", true, basicAddress)
	model.OnCreate()

	createdAt := model.CreatedAt
	updatedAt := model.UpdatedAt

	model.OnUpdate()

	assert.True(t, model.ID.IsZero())
	assert.Equal(t, createdAt, model.CreatedAt)
	assert.NotEqual(t, updatedAt, model.UpdatedAt)
}

func TestCustomModel_OnReplace(t *testing.T) {

	t.Parallel()

	model := newCustomModel("John Doe", "test@test", true, basicAddress)
	model.OnCreate()

	createdAt := model.CreatedAt
	updatedAt := model.UpdatedAt

	model.OnReplace()

	assert.True(t, model.ID.IsZero())
	assert.Equal(t, createdAt, model.CreatedAt)
	assert.NotEqual(t, updatedAt, model.UpdatedAt)
}
