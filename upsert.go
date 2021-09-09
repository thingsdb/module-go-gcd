package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

type Upsert struct {
	Entities []Entity `msgpack:"entities"`
}

// upsert inserts an entities if they do not exist or updates them if they do.
// Returns the keys.
func (upsert *Upsert) upsert(ctx context.Context, client *datastore.Client) ([]*datastore.Key, error) {
	keys, props, err := upsert.prepare()
	if err != nil {
		return nil, err
	}

	return client.PutMulti(ctx, keys, props)
}

// transactionUpsert inserts an entities if they do not exist or updates them if they do.
// Returns pending keys.
func (upsert *Upsert) transactionUpsert(tx *datastore.Transaction) ([]*datastore.PendingKey, error) {
	keys, props, err := upsert.prepare()
	if err != nil {
		return nil, err
	}

	return tx.PutMulti(keys, props)
}

// prepare prepares the entities before upsert.
func (upsert *Upsert) prepare() ([]*datastore.Key, []datastore.PropertyList, error) {
	if len(upsert.Entities) < 1 {
		return nil, nil, fmt.Errorf("GCD upsert requires `Entities`")
	}

	cap := len(upsert.Entities)
	keys := make([]*datastore.Key, 0, cap)
	props := make([]datastore.PropertyList, 0, cap)
	for _, entity := range upsert.Entities {
		if entity.Key.Kind == "" {
			return nil, nil, fmt.Errorf("GCD upsert requires `Kind`")
		}

		var propertyList datastore.PropertyList
		propertySlice := make([]datastore.Property, 0, len(entity.Properties))
		for _, prop := range entity.Properties {
			prop.Value = convInvalidInts(prop.Value)
			propertySlice = append(propertySlice, prop)
		}
		propertyList.Load(propertySlice)

		keys = append(keys, entity.Key)
		props = append(props, propertyList)
	}

	return keys, props, nil
}
