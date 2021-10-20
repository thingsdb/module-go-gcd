package main

type Operator string

const (
	Equal              Operator = "="
	LessThan           Operator = "<"
	LessThanOrEqual    Operator = "<="
	GreaterThan        Operator = ">"
	GreaterThanOrEqual Operator = ">="
)

type PropertyFilter struct {
	Name     string      `msgpack:"name"`
	Operator string      `msgpack:"operator"`
	Value    interface{} `msgpack:"value"`
}

type Filter struct {
	Properties []PropertyFilter `msgpack:"properties"`
}
