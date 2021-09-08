package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

func execReq(ctx context.Context, client *datastore.Client, req *reqMySQL) (interface{}, error) {
	switch cmd := req.Cmd; cmd {
	case InsertUpdateCmd:
		return insertEntities(ctx, client, req)
	case GetCmd:
		return getEntities(ctx, client, req)
	case DeleteCmd:
		return nil, deleteEntities(ctx, client, req)
	default:
		return nil, fmt.Errorf("Cmd parameter unknown; valid options are `InsertEntity`, `InsertEntities`, `UpdateEntity`, `UpdateEntities`, `GetEntity`, `GetEntities`, `DeleteEntity`, or `DeleteEntities`")
	}
}

// insertEntities inserts an new entity to the datastore,
// returning the key of the newly created entity.
func insertEntities(ctx context.Context, client *datastore.Client, req *reqMySQL) ([]*datastore.Key, error) {
	if len(req.Entities) < 1 {
		return nil, fmt.Errorf("GCD InsertEntities requires `Entities`")
	}

	cap := len(req.Entities)
	keys := make([]*datastore.Key, 0, cap)
	props := make([]datastore.PropertyList, 0, cap)
	for _, entity := range req.Entities {
		if entity.Key.Kind == "" {
			return nil, fmt.Errorf("GCD InsertEntity requires `Kind`")
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

	return client.PutMulti(ctx, keys, props) // updates existing entities... ?
}

// getEntities gets entities from the datastore.
func getEntities(ctx context.Context, client *datastore.Client, req *reqMySQL) (interface{}, error) {
	var keys []*datastore.Key
	var entities []Entity
	var err error
	if len(req.Entities) > 0 {
		keys, entities, err = queryByKeys(ctx, client, req)
		if err != nil {
			return nil, err
		}
	} else {
		keys, entities, err = query(ctx, client, req)
		if err != nil {
			return nil, err
		}
	}

	switch fetch := req.Fetch; fetch {
	case Keys:
		return keys, nil
	default:
		return entities, nil
	}
}

// deleteEntities deletes entities from the datastore.
func deleteEntities(ctx context.Context, client *datastore.Client, req *reqMySQL) error {
	if len(req.Entities) < 1 {
		return fmt.Errorf("GCD DeleteEntities requires `Entities`")
	}

	cap := len(req.Entities)
	keys := make([]*datastore.Key, 0, cap)
	for _, entity := range req.Entities {
		if entity.Key.Kind == "" {
			return fmt.Errorf("GCD DeleteEntities requires `Kind`")
		}

		keys = append(keys, entity.Key)
	}

	return client.DeleteMulti(ctx, keys)
}

func query(ctx context.Context, client *datastore.Client, req *reqMySQL) ([]*datastore.Key, []Entity, error) {
	if req.Kind == "" {
		return nil, nil, fmt.Errorf("GCD InsertEntity requires `Kind`")
	}

	query := datastore.NewQuery(req.Kind)

	if req.Namespace != "" {
		query = query.Namespace(req.Namespace)
	}

	if req.Filter.Ancestor != nil {
		query = query.Ancestor(req.Filter.Ancestor)
	}

	for _, filter := range req.Filter.Properties {
		query = query.Filter(fmt.Sprintf("%s %s", filter.Name, filter.Operator), filter.Value)
	}

	if req.Limit != 0 {
		query = query.Limit(req.Limit)
	}

	if req.Cursor != "" {
		cursor, err := datastore.DecodeCursor(req.Cursor)
		if err != nil {
			return nil, nil, err
		}
		query = query.Start(cursor)
	}

	if req.Order.Name != "" {
		query = query.Order(fmt.Sprintf("%s%s", req.Order.Direction, req.Order.Name))
	}

	var propertyList []datastore.PropertyList
	keys, err := client.GetAll(ctx, query, &propertyList)
	if err != nil {
		return nil, nil, err
	}

	entities, err := returnEntities(keys, propertyList)
	if err != nil {
		return nil, nil, err
	}

	return keys, entities, nil
}

func queryByKeys(ctx context.Context, client *datastore.Client, req *reqMySQL) ([]*datastore.Key, []Entity, error) {
	keys := make([]*datastore.Key, 0, len(req.Entities))
	for _, entity := range req.Entities {
		keys = append(keys, entity.Key)
	}

	propertyList := make([]datastore.PropertyList, len(req.Entities)) // Need len otherwise hit `return errors.New("datastore: keys and dst slices have different length")``
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
