package friendlymongo

import (
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type stages map[string]*stage

func (s stages) append(id string, filters interface{}) {

	s[id].append(filters)
}

func (s stages) stages() []*stage {

	v := make([]*stage, 0, len(s))
	for _, value := range s {
		v = append(v, value)
	}
	return v
}

func (s stages) sorted() []*stage {

	v := s.stages()
	sort.Slice(v, func(i, j int) bool {
		return v[i].Priority < v[j].Priority
	})
	return v
}

type StageBuilder struct {
	stages   stages
	priority int8
}

func NewStageBuilder() *StageBuilder {

	return &StageBuilder{
		stages:   make(map[string]*stage),
		priority: 0,
	}
}

func (pb *StageBuilder) Project(id string, filters interface{}) *StageBuilder {

	return pb.AddStage(id, "$project", filters)
}

func (pb *StageBuilder) Match(id string, filters interface{}) *StageBuilder {

	return pb.AddStage(id, "$match", filters)
}

// Bucket categorizes incoming documents into groups, called buckets, based on a specified expression and bucket
// boundaries and outputs a document per each bucket. Each output document contains an _id field whose value specifies
// the inclusive lower bound of the bucket.
// The output option specifies the fields included in each output document.
//
// $bucket only produces output documents for buckets that contain at least one input document.
//
// [More info]: https://www.mongodb.com/docs/manual/reference/operator/aggregation/bucket
func (pb *StageBuilder) Bucket(id string, groupBy string, boundaries []any, def any, output bson.M) *StageBuilder {

	return pb.AddStage(id, "$bucket", bson.M{
		"groupBy":    groupBy,
		"boundaries": boundaries,
		"default":    def,
		"output":     output,
	})
}

func (pb *StageBuilder) Group(id string, filters interface{}) *StageBuilder {

	return pb.AddStage(id, "$group", filters)
}

func (pb *StageBuilder) Sort(id string, filters interface{}) *StageBuilder {

	return pb.AddStage(id, "$sort", filters)
}

func (pb *StageBuilder) Limit(id string, filters interface{}) *StageBuilder {

	return pb.AddStage(id, "$limit", filters)
}

func (pb *StageBuilder) Skip(id string, filters interface{}) *StageBuilder {

	return pb.AddStage(id, "$skip", filters)
}

func (pb *StageBuilder) Lookup(id, from, localField, foreignField, as string) *StageBuilder {

	return pb.AddStage(id, "$lookup", bson.M{
		"from":         from,
		"localField":   localField,
		"foreignField": foreignField,
		"as":           as,
	})
}

// TODO DOC
type unwindOpts struct {
	includeIndex         *string
	preserveNullAndEmpty *bool
}

type unwindOptsFunc func(*unwindOpts)

func IncludeIndex(index string) unwindOptsFunc {

	return func(opts *unwindOpts) {
		opts.includeIndex = &index
	}
}
func PreserveNullEmpty(preserve bool) unwindOptsFunc {

	return func(opts *unwindOpts) {
		opts.preserveNullAndEmpty = &preserve
	}
}

func (pb *StageBuilder) Unwind(id, fieldPath string, opts ...unwindOptsFunc) *StageBuilder {

	unwind := bson.M{"path": fieldPath}

	unwindOpts := unwindOpts{}
	for _, opt := range opts {
		opt(&unwindOpts)
	}

	if unwindOpts.includeIndex != nil {
		unwind["includeArrayIndex"] = *unwindOpts.includeIndex
	}
	if unwindOpts.preserveNullAndEmpty != nil {
		unwind["preserveNullAndEmptyArrays"] = *unwindOpts.preserveNullAndEmpty
	}

	return pb.AddStage(id, "$unwind", unwind)
}

func (pb *StageBuilder) SortByCount(id string, expr interface{}) *StageBuilder {

	return pb.AddStage(id, "$sortByCount", expr)
}

func (pb *StageBuilder) Facet(id string, aggregations map[string]*StageBuilder) *StageBuilder {

	var aggrBuild primitive.D
	for outputField, stages := range aggregations {
		aggrBuild = append(aggrBuild, bson.E{Key: outputField, Value: stages.Build()})
	}

	return pb.AddStage(id, "$facet", aggrBuild)
}

func (pb *StageBuilder) Count(id string, filter string) *StageBuilder {

	return pb.AddStage(id, "$count", filter)
}

func (pb *StageBuilder) AddStage(id string, stageType string, filters interface{}) *StageBuilder {

	if pb.stages[id] != nil {
		panic("stage already exists")
	}

	pb.stages[id] = newStage(stageType, pb.priority, filters)
	pb.priority++
	return pb
}

func (pb *StageBuilder) Append(id string, filters interface{}) *StageBuilder {

	if pb.stages[id] == nil {
		panic("stage does not exist")
	}
	pb.stages.append(id, filters)
	return pb
}

func (pb *StageBuilder) Build() mongo.Pipeline {

	var pipeline mongo.Pipeline

	for _, stage := range pb.stages.sorted() {
		pipeline = append(pipeline, stage.toBsonD())
	}
	return pipeline
}

func (pb *StageBuilder) Stages() []*stage {

	return pb.stages.stages()
}
