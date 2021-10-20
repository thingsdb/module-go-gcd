package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type Get struct {
	Ancestor  *datastore.Key `msgpack:"ancestor"`
	Cursor    string         `msgpack:"cursor"`
	Fetch     Fetch          `msgpack:"fetch"`
	Filter    Filter         `msgpack:"filter"`
	Entities  []Entity       `msgpack:"entities"`
	Kind      string         `msgpack:"kind"`
	Limit     int            `msgpack:"limit"`
	Namespace string         `msgpack:"namespace"`
	Order     Order          `msgpack:"order"`
}

// get gets entities from the datastore.
func (get *Get) run(ctx context.Context, client *datastore.Client) (interface{}, error) {
	var propertyList []datastore.PropertyList
	var keys []*datastore.Key
	var cursor datastore.Cursor

	if len(get.Entities) > 0 {
		keys, propertyList = get.prepare()

		if err := client.GetMulti(ctx, keys, propertyList); err != nil {
			return nil, err
		}
	} else {
		query, err := get.query()
		if err != nil {
			return nil, err
		}

		it := client.Run(ctx, query)
		for {
			var p datastore.PropertyList
			key, err := it.Next(&p)

			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}

			propertyList = append(propertyList, p)
			keys = append(keys, key)

			cursor, err = it.Cursor()
			if err != nil {
				return nil, err
			}
		}
	}

	entities, err := returnEntities(keys, propertyList)
	if err != nil {
		return nil, err
	}

	var c = cursor.String()
	switch fetch := get.Fetch; fetch {
	case Keys:
		return []interface{}{keys, c}, nil
	default:
		return []interface{}{entities, c}, nil
	}
}

// get gets entities from the datastore.
func (get *Get) runInTransaction(tx *datastore.Transaction) (interface{}, error) {
	if len(get.Entities) < 1 {
		return nil, errorMsg("`get` in transaction requires `entities`")
	}

	keys, propertyList := get.prepare()

	if err := tx.GetMulti(keys, propertyList); err != nil {
		return nil, err
	}

	entities, err := returnEntities(keys, propertyList)
	if err != nil {
		return nil, err
	}

	switch fetch := get.Fetch; fetch {
	case Keys:
		return keys, nil
	default:
		return entities, nil
	}
}

// prepare prepares the entities.
func (get *Get) prepare() ([]*datastore.Key, []datastore.PropertyList) {
	keys := make([]*datastore.Key, 0, len(get.Entities))
	propertyList := make([]datastore.PropertyList, len(get.Entities)) // Need len otherwise hit `return errors.New("datastore: keys and dst slices have different length")``

	for _, entity := range get.Entities {
		keys = append(keys, entity.Key)
	}

	return keys, propertyList
}

func (get Get) query() (*datastore.Query, error) {
	if get.Kind == "" {
		return nil, errorMsg("`get` requires `kind`")
	}

	query := datastore.NewQuery(get.Kind)

	if get.Namespace != "" {
		query = query.Namespace(get.Namespace)
	}

	if get.Ancestor != nil {
		query = query.Ancestor(get.Ancestor)
	}

	for _, filter := range get.Filter.Properties {
		query = query.Filter(fmt.Sprintf("%s %s", filter.Name, filter.Operator), filter.Value)
	}

	if get.Limit != 0 {
		query = query.Limit(get.Limit)
	}

	if get.Cursor != "" {
		cursor, err := datastore.DecodeCursor(get.Cursor)
		if err != nil {
			return nil, err
		}
		query = query.Start(cursor)
	}

	if get.Order.Name != "" {
		query = query.Order(fmt.Sprintf("%s%s", get.Order.Direction, get.Order.Name))
	}

	return query, nil
}

func returnEntities(keys []*datastore.Key, propertyList []datastore.PropertyList) ([]Entity, error) {
	entities := make([]Entity, 0, len(keys))
	for i, key := range keys {
		properties, err := propertyList[i].Save()
		if err != nil {
			return nil, err
		}

		entities = append(entities, Entity{
			Key:        key,
			Properties: properties,
		})
	}

	return entities, nil
}
