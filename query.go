package main

import (
	"context"
	"fmt"
	"reflect"

	"cloud.google.com/go/datastore"
)

var intList = []reflect.Kind{reflect.Int8, reflect.Int16, reflect.Int32}

type ReturnEntity struct {
	Key          datastore.Key
	PropertyList datastore.PropertyList
}

// InsertEntities inserts an new entity to the datastore,
// returning the key of the newly created entity.
func InsertEntities(ctx context.Context, client *datastore.Client, req *reqMySQL) ([]*datastore.Key, error) {
	if len(req.Entities) < 1 {
		return nil, fmt.Errorf("GCD InsertEntity requires `Entities`")
	}

	cap := len(req.Entities)
	keys := make([]*datastore.Key, 0, cap)
	props := make([]datastore.PropertyList, 0, cap)
	for _, entity := range req.Entities {
		var key *datastore.Key

		if entity.Key.Kind == "" {
			return nil, fmt.Errorf("GCD InsertEntity requires `Kind`")
		}

		if entity.Key.ID != 0 {
			key = datastore.IDKey(entity.Key.Kind, entity.Key.ID, entity.Key.Parent)
		} else if entity.Key.Name != "" {
			key = datastore.NameKey(entity.Key.Kind, entity.Key.Name, entity.Key.Parent)
		} else {
			key = datastore.IncompleteKey(entity.Key.Kind, entity.Key.Parent)
		}
		key.Namespace = entity.Key.Namespace

		var propertyList datastore.PropertyList
		propertySlice := make([]datastore.Property, 0, len(entity.Properties))
		for _, prop := range entity.Properties {
			rv := reflect.ValueOf(prop.Value)

			isInvalidInt := false
			for _, intType := range intList {
				if intType == rv.Kind() {
					isInvalidInt = true
					break
				}
			}

			if isInvalidInt {
				prop.Value = rv.Int()
				propertySlice = append(propertySlice, datastore.Property(prop))
			} else {
				propertySlice = append(propertySlice, datastore.Property(prop))
			}
		}
		propertyList.Load(propertySlice)

		keys = append(keys, key)
		props = append(props, propertyList)
	}

	return client.PutMulti(ctx, keys, props)
}

// GetEntities gets all the entities from the datastore.
func GetEntities(ctx context.Context, client *datastore.Client, req *reqMySQL) (interface{}, error) {
	if len(req.Entities) > 0 {
		propertyList, err := queryByKeys(ctx, client, req)
		if err != nil {
			return nil, err
		}

		return propertyList, nil
	} else {
		keys, propertyList, err := query(ctx, client, req)
		if err != nil {
			return nil, err
		}

		switch fetch := req.Fetch; fetch {
		case Keys:
			return keys, nil
		default:
			returnEntity := make([]ReturnEntity, 0, len(keys))
			for i, key := range keys {
				returnEntity = append(returnEntity, ReturnEntity{
					Key:          *key,
					PropertyList: propertyList[i],
				})
			}
			return returnEntity, nil
		}
	}
}

func execReq(ctx context.Context, client *datastore.Client, req *reqMySQL) (interface{}, error) {
	switch cmd := req.Cmd; cmd {
	case InsertEntitiesCmd:
		return InsertEntities(ctx, client, req)
	case UpdateEntityCmd:
		return nil, fmt.Errorf("UpdateEntityCmd not implemented yet")
	case UpdateEntitiesCmd:
		return nil, fmt.Errorf("UpdateEntitiesCmd not implemented yet")
	case GetEntitiesCmd:
		return GetEntities(ctx, client, req)
	case DeleteEntityCmd:
		return nil, fmt.Errorf("DeleteEntityCmd not implemented yet")
	case DeleteEntitiesCmd:
		return nil, fmt.Errorf("DeleteEntitiesCmd not implemented yet")
	default:
		return nil, fmt.Errorf("Cmd parameter unknown; valid options are `InsertEntity`, `InsertEntities`, `UpdateEntity`, `UpdateEntities`, `GetEntity`, `GetEntities`, `DeleteEntity`, or `DeleteEntities`")
	}
}

func query(ctx context.Context, client *datastore.Client, req *reqMySQL) ([]*datastore.Key, []datastore.PropertyList, error) {
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

	return keys, propertyList, nil
}

func queryByKeys(ctx context.Context, client *datastore.Client, req *reqMySQL) ([]datastore.PropertyList, error) {
	keys := make([]*datastore.Key, 0, len(req.Entities))
	for _, entity := range req.Entities {
		keys = append(keys, entity.Key)
	}

	propertyList := make([]datastore.PropertyList, len(req.Entities)) // Need len otherwise hit `return errors.New("datastore: keys and dst slices have different length")``
	if err := client.GetMulti(ctx, keys, propertyList); err != nil {
		return nil, err
	}

	return propertyList, nil
}

// type Filter struct {
// 	filterStr string      `msgpack:"filterStr"`
// 	value     interface{} `msgpack:"value"`
// }

// type Query struct {
// 	Ancestor            string   `msgpack:"ancestor"`
// 	Distinct            bool     `msgpack:"distinct"`
// 	DistinctOn          []string `msgpack:"distinctOn"`
// 	End                 string   `msgpack:"end"`
// 	EventualConsistency bool     `msgpack:"eventualConsistency"`
// 	Filter              Filter   `msgpack:"filter"`
// 	KeysOnly            bool     `msgpack:"keysOnly"`
// 	Limit               int      `msgpack:"limit"`
// 	Namespace           string   `msgpack:"namespace"`
// 	Offset              int      `msgpack:"offset"`
// 	Order               string   `msgpack:"order"`
// 	Project             []string `msgpack:"project"`
// 	Start               string   `msgpack:"Start"`
// 	Transaction         bool     `msgpack:"transaction"` // TODO ????????
// 	Value               string   `msgpack:"value"`
// }
