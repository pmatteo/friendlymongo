package friendlymongo

import (
	"sort"

	"go.mongodb.org/mongo-driver/bson"
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

func (pb *StageBuilder) Lookup(id string, filters bson.M) *StageBuilder {

	return pb.AddStage(id, "$lookup", filters)
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
