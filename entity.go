package main

import "cloud.google.com/go/datastore"

type Entity struct {
	Key        *datastore.Key       `msgpack:"key"`
	Properties []datastore.Property `msgpack:"properties"`
}
