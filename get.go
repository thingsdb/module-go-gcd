package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type Get struct {
	Cursor    string   `msgpack:"cursor"`
	Fetch     Fetch    `msgpack:"fetch"`
	Filter    Filter   `msgpack:"filter"`
	Entities  []Entity `msgpack:"entities"`
	Kind      string   `msgpack:"kind"`
	Limit     int      `msgpack:"limit"`
	Namespace string   `msgpack:"namespace"`
	Order     Order    `msgpack:"order"`
}

// get gets entities from the datastore.
func (get *Get) get(ctx context.Context, client *datastore.Client) (interface{}, string, error) {
	var keys []*datastore.Key
	var entities []Entity
	var err error
	var cursor string
	if len(get.Entities) > 0 {
		keys, entities, err = get.queryByKeys(ctx, client)
		if err != nil {
			return nil, "", err
		}
	} else {
		keys, entities, cursor, err = get.query(ctx, client)
		if err != nil {
			return nil, "", err
		}
	}

	switch fetch := get.Fetch; fetch {
	case Keys:
		return keys, cursor, nil
	default:
		return entities, cursor, nil
	}
}

// get gets entities from the datastore.
func (get *Get) transactionGet(tx *datastore.Transaction) (interface{}, error) {
	if len(get.Entities) == 0 {
		return nil, errorMsg("`get` in transaction requires `entities`")
	}
	keys := make([]*datastore.Key, 0, len(get.Entities))
	for _, entity := range get.Entities {
		keys = append(keys, entity.Key)
	}

	propertyList := make([]datastore.PropertyList, len(get.Entities)) // Need len otherwise hit `return errors.New("datastore: keys and dst slices have different length")``
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

func (get Get) query(ctx context.Context, client *datastore.Client) ([]*datastore.Key, []Entity, string, error) {
	if get.Kind == "" {
		return nil, nil, "", errorMsg("`get` requires `kind`")
	}

	query := datastore.NewQuery(get.Kind)

	if get.Namespace != "" {
		query = query.Namespace(get.Namespace)
	}

	if get.Filter.Ancestor != nil {
		query = query.Ancestor(get.Filter.Ancestor)
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
			return nil, nil, "", err
		}
		query = query.Start(cursor)
	}

	if get.Order.Name != "" {
		query = query.Order(fmt.Sprintf("%s%s", get.Order.Direction, get.Order.Name))
	}

	var propertyList []datastore.PropertyList
	var keys []*datastore.Key
	var cursor datastore.Cursor
	it := client.Run(ctx, query)
	for {
		var p datastore.PropertyList
		key, err := it.Next(&p)

		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}

		propertyList = append(propertyList, p)
		keys = append(keys, key)

		cursor, err := it.Cursor()
		if err != nil {
			break
		}
		// When printed, a cursor will display as a string that can be passed
		// to datastore.DecodeCursor.
		fmt.Printf("to resume with this post, use cursor %s\n", cursor)
	}

	// keys, err := client.GetAll(ctx, query, &propertyList)
	// if err != nil {
	// 	return nil, nil, "", err
	// }

	entities, err := returnEntities(keys, propertyList)
	if err != nil {
		return nil, nil, "", err
	}

	return keys, entities, cursor.String(), nil
}

func (get Get) queryByKeys(ctx context.Context, client *datastore.Client) ([]*datastore.Key, []Entity, error) {
	keys := make([]*datastore.Key, 0, len(get.Entities))
	for _, entity := range get.Entities {
		keys = append(keys, entity.Key)
	}

	propertyList := make([]datastore.PropertyList, len(get.Entities)) // Need len otherwise hit `return errors.New("datastore: keys and dst slices have different length")``
	if err := client.GetMulti(ctx, keys, propertyList); err != nil {
		return nil, nil, err
	}

	entities, err := returnEntities(keys, propertyList)
	if err != nil {
		return nil, nil, err
	}

	return keys, entities, nil
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
