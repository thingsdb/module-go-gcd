package main

import (
	"context"

	"cloud.google.com/go/datastore"
)

type Delete struct {
	Entities []Entity `msgpack:"entities"`
}

// delete deletes entities from the datastore.
func (delete *Delete) run(ctx context.Context, client *datastore.Client) (interface{}, error) {
	keys, err := delete.prepare()
	if err != nil {
		return "", err
	}

	err = client.DeleteMulti(ctx, keys)
	if err != nil {
		return "", err
	}

	num := len(keys)
	return retMsg(num, "Deleted"), nil
}

// transactionDelete entities from the datastore.
func (delete *Delete) runInTransaction(tx *datastore.Transaction) (interface{}, error) {
	keys, err := delete.prepare()
	if err != nil {
		return "", err
	}

	err = tx.DeleteMulti(keys)
	if err != nil {
		return "", err
	}

	num := len(keys)
	return retMsg(num, "Deleted"), nil
}

// prepare prepares the entities before delete.
func (delete *Delete) prepare() ([]*datastore.Key, error) {
	if len(delete.Entities) < 1 {
		return nil, errorMsg("`delete` requires `entities`")
	}

	cap := len(delete.Entities)
	keys := make([]*datastore.Key, 0, cap)
	for _, entity := range delete.Entities {
		if entity.Key.Kind == "" {
			return nil, errorMsg("`delete` requires `kind`")
		}

		keys = append(keys, entity.Key)
	}

	return keys, nil
}
