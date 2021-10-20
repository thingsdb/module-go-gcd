package main

import (
	"cloud.google.com/go/datastore"
	"github.com/vmihailenco/msgpack/v4"
)

type key datastore.Key
type property datastore.Property

type subKey struct {
	Kind      string `msgpack:"kind"`
	ID        int64  `msgpack:"id"`
	Name      string `msgpack:"name"`
	Parent    *key   `msgpack:"parent"`
	Namespace string `msgpack:"namespace"`
}

type subProperty struct {
	Name    string      `msgpack:"name"`
	Value   interface{} `msgpack:"value"`
	NoIndex bool        `msgpack:"no_index"`
}

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
	err := msgpack.Unmarshal(data, &ret)
	if err != nil {
		return err
	}

	e.Key = (*datastore.Key)(ret.Key)
	e.Properties = make([]datastore.Property, len(ret.Properties))
	for i, p := range ret.Properties {
		e.Properties[i] = datastore.Property(p)
	}
	return nil
}

// func (e *Entity) MarshalMsgpack() ([]byte, error) {
// 	var ret entity
// 	ret.Key = (*key)(e.Key)

// 	ret.Properties = make([]property, len(e.Properties))
// 	for i, p := range e.Properties {
// 		ret.Properties[i] = property(p)
// 	}

// 	return msgpack.Marshal(&ret)
// }

func (k *key) UnmarshalMsgpack(data []byte) error {
	var ret subKey
	err := msgpack.Unmarshal(data, &ret)
	if err != nil {
		return err
	}

	k.Kind = ret.Kind
	k.ID = ret.ID
	k.Name = ret.Name
	k.Parent = (*datastore.Key)(ret.Parent)
	k.Namespace = ret.Namespace
	return nil
}

// func (k *key) MarshalMsgpack() ([]byte, error) {
// 	var ret subKey
// 	ret.Kind = k.Kind
// 	ret.ID = k.ID
// 	ret.Name = k.Name
// 	ret.Parent = (*key)(k.Parent)
// 	ret.Namespace = k.Namespace
// 	return msgpack.Marshal(&ret)
// }

func (p *property) UnmarshalMsgpack(data []byte) error {
	var ret subProperty
	err := msgpack.Unmarshal(data, &ret)
	if err != nil {
		return err
	}

	p.Name = ret.Name
	p.Value = ret.Value
	p.NoIndex = ret.NoIndex
	return nil
}

// func (p *property) MarshalMsgpack() ([]byte, error) {
// 	var ret subProperty
// 	ret.Name = p.Name
// 	ret.Value = p.Value
// 	ret.NoIndex = p.NoIndex
// 	return msgpack.Marshal(&ret)
// }
