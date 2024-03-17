package friendlymongo_test

import (
	"testing"
	"time"

	"github.com/pmatteo/friendlymongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBaseModel_Init(t *testing.T) {
	t.Parallel()

	model := &friendlymongo.BaseModel{}

	assert.True(t, model.ID.IsZero())
	assert.Empty(t, model.CreatedAt)
	assert.Empty(t, model.UpdatedAt)

	model.Init()

	assert.False(t, model.ID.IsZero())
	assert.NotEmpty(t, model.CreatedAt)
	assert.NotEmpty(t, model.UpdatedAt)
}

func TestBaseModel_Setters(t *testing.T) {
	t.Parallel()

	model := &friendlymongo.BaseModel{}

	assert.True(t, model.ID.IsZero())
	assert.Empty(t, model.CreatedAt)
	assert.Empty(t, model.UpdatedAt)

	model.SetID()

	assert.False(t, model.ID.IsZero())

	model.SetCreatedAt()
	assert.NotEmpty(t, model.CreatedAt)

	model.SetUpdatedAt()
	assert.NotEmpty(t, model.UpdatedAt)
}

func TestBaseModel_Getters(t *testing.T) {
	t.Parallel()

	model := friendlymongo.NewBaseModel()

	require.NotNil(t, model.GetID())
	assert.NotEqual(t, primitive.NilObjectID, model.GetID())
	assert.NotEmpty(t, model.GetUpdateAt())
	assert.NotEmpty(t, model.GetCreatedAt())
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

func NewCustomModel(name, email string, active bool, a *address) *customModel {
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

// TestCustomModel_Init tests that the Init method sets the BaseModel fields
// if they are not already set.
func TestCustomModel_Init(t *testing.T) {

	t.Parallel()

	now := time.Now()

	model := &customModel{
		Name:           "John Doe",
		Email:          "johndoe@test.com",
		Active:         true,
		ActivationDate: &now,
		Address:        basicAddress,
	}

	assert.True(t, model.ID.IsZero())
	assert.Empty(t, model.CreatedAt)
	assert.Empty(t, model.UpdatedAt)

	model.Init()

	assert.NotEqual(t, primitive.NilObjectID, model.ID)
	assert.NotEmpty(t, model.CreatedAt)
	assert.NotEmpty(t, model.UpdatedAt)
}

// TestCustomModel_Setters tests that the Init method does not
// override the BaseModel fields if they are already set.
func TestCustomModel_Setters(t *testing.T) {

	t.Parallel()

	now := time.Now()

	model := &customModel{
		Name:           "John Doe",
		Email:          "johndoe@test.com",
		Active:         true,
		ActivationDate: &now,
		Address:        basicAddress,
	}

	assert.True(t, model.ID.IsZero())
	assert.Empty(t, model.CreatedAt)
	assert.Empty(t, model.UpdatedAt)

	model.SetID()

	assert.False(t, model.ID.IsZero())

	model.SetCreatedAt()
	assert.NotEmpty(t, model.CreatedAt)

	model.SetUpdatedAt()
	assert.NotEmpty(t, model.UpdatedAt)
}
