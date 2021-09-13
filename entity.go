package main

import (
	"cloud.google.com/go/datastore"
	"github.com/vmihailenco/msgpack/v4"
)

type key datastore.Key
type property datastore.Property

type entity struct {
	Key        *key       `msgpack:"key"`
	Properties []property `msgpack:"properties"`
}

type Entity struct {
	Key        *datastore.Key       `msgpack:"key"`
	Properties []datastore.Property `msgpack:"properties"`
}

func (e *Entity) UnmarshalMsgpack(data []byte) error {
	var ret entity
	_ = msgpack.Unmarshal(data, &ret)
	e.Key = (*datastore.Key)(ret.Key)
	e.Properties = make([]datastore.Property, len(ret.Properties))
	for i, p := range ret.Properties {
		e.Properties[i] = datastore.Property(p)
	}

	return nil
}

func (k *key) UnmarshalMsgpack(data []byte) error {
	var ret map[string]interface{}
	_ = msgpack.Unmarshal(data, &ret)

	ki, ok := ret["kind"].(string)
	if ok {
		k.Kind = ki
	}
	i, ok := ret["id"].(int64)
	if ok {
		k.ID = i
	}
	n, ok := ret["name"].(string)
	if ok {
		k.Name = n
	}
	p, ok := ret["parent"].(*datastore.Key)
	if ok {
		k.Parent = p
	}
	ns, ok := ret["namespace"].(string)
	if ok {
		k.Namespace = ns
	}

	return nil
}

func (p *property) UnmarshalMsgpack(data []byte) error {
	var ret map[string]interface{}
	_ = msgpack.Unmarshal(data, &ret)
	n, ok := ret["name"].(string)
	if ok {
		p.Name = n
	}
	p.Value = ret["value"]
	ni, ok := ret["no_index"].(bool)
	if ok {
		p.NoIndex = ni
	}
	return nil
}
