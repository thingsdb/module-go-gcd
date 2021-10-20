package main

import (
	"context"

	"cloud.google.com/go/datastore"
)

type Transaction struct {
	Delete *Delete      `msgpack:"delete"`
	Get    *Get         `msgpack:"get"`
	Next   *Transaction `msgpack:"next"`
	Upsert *Upsert      `msgpack:"upsert"`
}

func (transaction *Transaction) run(ctx context.Context, client *datastore.Client) (interface{}, error) {
	var ret interface{}
	var err error

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var err error
		ret, err = transaction.execTransactionQuery(tx)
		if err != nil {
			return err
		}
		return nil
	})

	return ret, err
}

func (transaction *Transaction) execTransactionQuery(tx *datastore.Transaction) (interface{}, error) {
	var fn func(tx *datastore.Transaction) (interface{}, error)
	no := 0
	key := ""
	ret := make(map[string]interface{})
	if transaction.Get != nil {
		fn = transaction.Get.runInTransaction
		no++
		key = "get"
	}

	if transaction.Delete != nil {
		fn = transaction.Delete.runInTransaction
		no++
		key = "delete"
	}

	if transaction.Upsert != nil {
		fn = transaction.Upsert.runInTransaction
		no++
		key = "upsert"
	}

	if no == 0 {
		return ret, errorMsg("GCD transaction requires either `get`, `delete`, `upsert`")
	}

	if no > 1 {
		return ret, errorMsg("GCD transaction requires either `get`, `delete`, `upsert`, not more then one")
	}

	res, err := fn(tx)
	if err != nil {
		return nil, err
	}
	ret[key] = res

	if transaction.Next != nil {
		next, err := transaction.Next.execTransactionQuery(tx)
		if err != nil {
			return nil, err
		}

		ret["next"] = next
	}

	return ret, nil
}
