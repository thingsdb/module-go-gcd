package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

type Delete struct {
	Entities []Entity `msgpack:"entities"`
}

// delete deletes entities from the datastore.
func (delete *Delete) delete(ctx context.Context, client *datastore.Client) (string, error) {
	keys, err := delete.prepare()
	if err != nil {
		return "", err
	}

	err = client.DeleteMulti(ctx, keys)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Removed %d entities", len(keys)), nil
}

// transactionDelete entities from the datastore.
func (delete *Delete) transactionDelete(tx *datastore.Transaction) (string, error) {
	keys, err := delete.prepare()
	if err != nil {
		return "", err
	}

	err = tx.DeleteMulti(keys)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Removed %d entities", len(keys)), nil
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
