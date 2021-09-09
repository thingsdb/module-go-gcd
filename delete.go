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
func (delete *Delete) delete(ctx context.Context, client *datastore.Client) error {
	keys, err := delete.prepare()
	if err != nil {
		return err
	}

	return client.DeleteMulti(ctx, keys)
}

// transactionDelete entities from the datastore.
func (delete *Delete) transactionDelete(tx *datastore.Transaction) error {
	keys, err := delete.prepare()
	if err != nil {
		return err
	}

	return tx.DeleteMulti(keys)
}

// prepare prepares the entities before delete.
func (delete *Delete) prepare() ([]*datastore.Key, error) {
	if len(delete.Entities) < 1 {
		return nil, fmt.Errorf("GCD delete requires `Entities`")
	}

	cap := len(delete.Entities)
	keys := make([]*datastore.Key, 0, cap)
	for _, entity := range delete.Entities {
		if entity.Key.Kind == "" {
			return nil, fmt.Errorf("GCD delete requires `Kind`")
		}

		keys = append(keys, entity.Key)
	}

	return keys, nil
}
