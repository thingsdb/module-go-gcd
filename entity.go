package main

import "cloud.google.com/go/datastore"

// type Key struct {
// 	Kind      string `msgpack:"kind"`
// 	ID        int64  `msgpack:"id"`
// 	Name      string `msgpack:"name"`
// 	Parent    *Key   `msgpack:"parent"`
// 	Namespace string `msgpack:"namespace"`
// }

// type Property struct {
// 	Name    string      `msgpack:"name"`
// 	Value   interface{} `msgpack:"value"`
// 	NoIndex bool        `msgpack:"no_index"`
// }

type Entity struct {
	Key        *datastore.Key       `msgpack:"key"`
	Properties []datastore.Property `msgpack:"properties"`
}
