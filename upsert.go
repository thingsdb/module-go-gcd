package main

import (
	"context"

	"cloud.google.com/go/datastore"
)

type Upsert struct {
	Entities []Entity `msgpack:"entities"`
}

// upsert inserts an entities if they do not exist or updates them if they do.
// Returns the keys.
func (upsert *Upsert) run(ctx context.Context, client *datastore.Client) (interface{}, error) {
	keys, props, err := upsert.prepare()
	if err != nil {
		return "", err
	}

	_, err = client.PutMulti(ctx, keys, props)
	if err != nil {
		return "", err
	}

	num := len(keys)
	return retMsg(num, "Upserted"), nil
}

// transactionUpsert inserts an entities if they do not exist or updates them if they do.
// Returns pending keys.
func (upsert *Upsert) runInTransaction(tx *datastore.Transaction) (interface{}, error) {
	keys, props, err := upsert.prepare()
	if err != nil {
		return "", err
	}

	_, err = tx.PutMulti(keys, props)
	if err != nil {
		return "", err
	}

	num := len(keys)
	return retMsg(num, "Upserted"), nil
}

// prepare prepares the entities before upsert.
func (upsert *Upsert) prepare() ([]*datastore.Key, []datastore.PropertyList, error) {
	if len(upsert.Entities) < 1 {
		return nil, nil, errorMsg("`upsert` requires `entities`")
	}

	cap := len(upsert.Entities)
	keys := make([]*datastore.Key, 0, cap)
	props := make([]datastore.PropertyList, 0, cap)
	for _, entity := range upsert.Entities {
		if entity.Key.Kind == "" {
			return nil, nil, errorMsg("`upsert` requires `kind`")
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
