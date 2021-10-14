package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

type Query struct {
	Cmd         Cmd     `msgpack:"cmd"`
	Delete      *Delete `msgpack:"delete"`
	Get         *Get    `msgpack:"get"`
	Next        *Query  `msgpack:"next"`
	Transaction bool    `msgpack:"transaction"`
	Upsert      *Upsert `msgpack:"upsert"`
}

func (query *Query) query(ctx context.Context, client *datastore.Client) (interface{}, error) {
	var ret interface{}
	var err error

	if query.Transaction {
		_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			var err error
			ret, err = query.execTransactionQuery(tx)
			if err != nil {
				return err
			}
			return nil
		})
	} else {
		ret, err = query.execQuery(ctx, client)
	}

	return ret, err
}

func (query *Query) execQuery(ctx context.Context, client *datastore.Client) (interface{}, error) {
	ret := make(map[string]interface{})
	switch cmd := query.Cmd; cmd {
	case UpsertCmd:
		if query.Upsert == nil {
			return nil, fmt.Errorf("Error: Upsert parameter is required")
		}
		upsertRet, err := query.Upsert.upsert(ctx, client)
		if err != nil {
			return nil, err
		}
		ret["upsert"] = upsertRet
	case GetCmd:
		if query.Get == nil {
			return nil, fmt.Errorf("Error: Get parameter is required")
		}
		getRet, err := query.Get.get(ctx, client)
		if err != nil {
			return nil, err
		}
		ret["get"] = getRet
	case DeleteCmd:
		if query.Delete == nil {
			return nil, fmt.Errorf("Error: Delete parameter is required")
		}
		deleteRet, err := query.Delete.delete(ctx, client)
		if err != nil {
			return nil, err
		}
		ret["delete"] = deleteRet
	default:
		return ret, fmt.Errorf("Error: Cmd parameter unknown; valid options are `upsert`, `get` or `delete`")
	}

	if query.Next != nil {
		next, err := query.Next.execQuery(ctx, client)
		if err != nil {
			return nil, err
		}

		ret["next"] = next
	}

	return ret, nil
}

func (query *Query) execTransactionQuery(tx *datastore.Transaction) (interface{}, error) {
	ret := make(map[string]interface{})
	switch cmd := query.Cmd; cmd {
	case UpsertCmd:
		if query.Upsert == nil {
			return nil, fmt.Errorf("Error: Upsert parameter is required")
		}
		upsertRet, err := query.Upsert.transactionUpsert(tx)
		if err != nil {
			return nil, err
		}
		ret["upsert"] = upsertRet
	case GetCmd:
		if query.Get == nil {
			return nil, fmt.Errorf("Error: Get parameter is required")
		}
		getRet, err := query.Get.transactionGet(tx)
		if err != nil {
			return nil, err
		}
		ret["get"] = getRet
	case DeleteCmd:
		if query.Delete == nil {
			return nil, fmt.Errorf("Error: Delete parameter is required")
		}
		deleteRet, err := query.Delete.transactionDelete(tx)
		if err != nil {
			return nil, err
		}
		ret["delete"] = deleteRet
	default:
		return ret, fmt.Errorf("Error: Cmd parameter unknown; valid options are `upsert`, `get` or `delete`")
	}

	if query.Next != nil {
		next, err := query.Next.execTransactionQuery(tx)
		if err != nil {
			return nil, err
		}

		ret["next"] = next
	}

	return ret, nil
}
