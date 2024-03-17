package friendlymongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type stage struct {
	Priority int8

	typ     string
	filters interface{}
}

func newStage(stageType string, priority int8, filters interface{}) *stage {
	return &stage{
		Priority: priority,
		typ:      stageType,
		filters:  filters,
	}
}

func (s *stage) append(filters interface{}) *stage {
	switch s.filters.(type) {

	case bson.M:
		s.filters = mergeBSON(s.filters.(bson.M), filters.(bson.M))

	case bson.D:
		s.filters = mergeBSOND(s.filters.(bson.D), filters.(bson.D))

	case string:
		panic(fmt.Errorf("cannot append to a string filter of stage of type %s", s.typ))
	}

	return s
}

func (s *stage) toBsonD() bson.D {
	return bson.D{{
		Key:   s.typ,
		Value: s.filters,
	}}
}

// mergeBSOND merges two BSON documents represented as bson.D
func mergeBSOND(d1, d2 bson.D) bson.D {
	// Create a map to store unique keys
	mergedKeys := make(map[string]struct{})

	// Create a new bson.D slice to store merged elements
	merged := append(bson.D{}, d1...)

	// Iterate over the elements of the second document
	for _, e := range d2 {
		// Check if the key is already present in the merged document
		if _, ok := mergedKeys[e.Key]; !ok {
			// If not present, append the element to the merged document
			merged = append(merged, e)
			mergedKeys[e.Key] = struct{}{}
		} else {
			// If present, update the value in the merged document
			for i, existing := range merged {
				if existing.Key == e.Key {
					merged[i] = e
					break
				}
			}
		}
	}

	return merged
}

// mergeBSON merges two BSON maps
func mergeBSON(m1, m2 bson.M) bson.M {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
