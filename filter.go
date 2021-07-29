package main

import "cloud.google.com/go/datastore"

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
	Ancestor   *datastore.Key   `msgpack:"ancestor"`
	Properties []PropertyFilter `msgpack:"properties"`
}
