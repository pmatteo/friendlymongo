package friendlymongo_test

import (
	"testing"

	"github.com/pmatteo/friendlymongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestAggregationBuilder_AddStage tests the AddStage method
// of the AggregationBuilder type.
func TestAggregationBuilder_AddStage(t *testing.T) {
	t.Parallel()

	builder := friendlymongo.
		NewStageBuilder().
		AddStage("1", "$match", bson.M{"test": bson.M{"$eq": 1}})

	require.Len(t, builder.Stages(), 1)

	expected := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"test": bson.M{"$eq": 1}}}},
	}
	assert.Equal(t, expected, builder.Build())
}

// TestAggregationBuilder_AddCount tests the AddCount method of the AggregationBuilder type.
func TestAggregationBuilder_AddCount(t *testing.T) {
	t.Parallel()

	builder := friendlymongo.
		NewStageBuilder().
		Count("1", "test")

	require.Len(t, builder.Stages(), 1)

	expected := mongo.Pipeline{
		{{Key: "$count", Value: "test"}},
	}
	assert.Equal(t, expected, builder.Build())
}

// TestAggregationBuilder_AddStageWithPriority tests the AddStageWithPriority method
// of the AggregationBuilder type.
func TestAggregationBuilder_AddStageWithPriority(t *testing.T) {
	t.Parallel()

	builder := friendlymongo.
		NewStageBuilder().
		AddStage("a", "$match", bson.M{"test": bson.M{"$eq": 1}})

	expected := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"test": bson.M{"$eq": 1}}}},
	}
	assert.Equal(t, expected, builder.Build())
}

// TestAggregationBuilder_AddStageWithPriority tests the AddStageWithPriority method
// of the AggregationBuilder type.
func TestAggregationBuilder_AddStageFilter(t *testing.T) {
	t.Parallel()

	builder := friendlymongo.
		NewStageBuilder().
		AddStage("a", "$match", bson.M{"test": bson.M{"$eq": 1}}).
		Append("a", bson.M{"test2": bson.M{"$eq": 2}})

	stages := builder.Stages()
	require.Len(t, stages, 1)

	expected := mongo.Pipeline{{
		{
			Key: "$match",
			Value: bson.M{
				"test":  bson.M{"$eq": 1},
				"test2": bson.M{"$eq": 2},
			},
		},
	}}
	assert.Equal(t, expected, builder.Build())
}

// TestAggregationBuilder_Build tests the Build method of the AggregationBuilder type
// with 2 stages added, no priority.
func TestAggregationBuilder_Build(t *testing.T) {
	t.Parallel()

	builder := friendlymongo.
		NewStageBuilder().
		Match("a", bson.M{"test": bson.M{"$eq": 1}}).
		Append("a", bson.M{"test2": bson.M{"$lt": 2}}).
		Lookup("b", bson.M{
			"from":         "otherCollection",
			"localField":   "id",
			"foreignField": "fkId",
			"as":           "other",
		})

	build := builder.Build()
	require.Len(t, build, 2)

	expected := mongo.Pipeline{
		{{
			Key: "$match",
			Value: bson.M{
				"test":  bson.M{"$eq": 1},
				"test2": bson.M{"$lt": 2},
			},
		}},
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         "otherCollection",
				"localField":   "id",
				"foreignField": "fkId",
				"as":           "other",
			},
		}},
	}

	assert.Equal(t, expected, builder.Build())
}
