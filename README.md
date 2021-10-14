# GCD ThingsDB Module (Go)
GCD module written using the [Go language](https://golang.org). This module can be used to communicate with the Google Cloud Datastore

## Installation

Install the module by running the following command in the `@thingsdb` scope:

```javascript
new_module('gcd', 'github.com/thingsdb/module-go-gcd');
```

Optionally, you can choose a specific version by adding a `@` followed with the release tag. For example: `@v0.1.0`.

## Configuration

The GCD module requires a configuration with the following properties:

Property | Type            | Description
-------- | --------------- | -----------
datastore_project_id    | str (required)  | Project ID.
datastore_emulator_host | str (optional)  | Host of locally-running datastore emulator.
google_app_cred_path    | str (optional)  | Path to the JSON key file for authorization.

Example configuration:

```javascript
set_module_conf('siridb', {
    datastore_project_id: 'id',
    datastore_emulator_host: 'host:port',
    google_app_cred_path: 'path/to/keyfile.json'
});
```

## Exposed functions

Name              | Description
----------------- | -----------
[query](#query)   | Run a GCD query.

### Query

#### Arguments

Argument      | Type                | Description
--------------|---------------------| -----------
`cmd`         | `str` (required)    | The command which can be either `upsert`, `get` or `delete`.
`delete`      | `Delete` (optional) | Thing with `delete` properties, see [Delete](#Delete).
`get`         | `Get` (optional)    | Thing with `get` properties, see [Get](#Get).
`next`        | `Query` (optional)  | The next query Thing.
`transaction` | `bool` (optional)   | Indicates if the query needs to be wrapped in a transaction or not.
`upsert`      | `Upsert` (optional) | Thing with the `upsert` properties, see [Upsert](#Upsert).

#### Example:



## Using the module

### Request

```ti
future({
    module: 'GCD',
    query: {
        cmd: "get",
        get: {
            entities: [
                {
                    key: {
                        id: 3,
                        namespace: 'test',
                        kind: 'Test',
                        parent: {
                            kind: 'Parent',
                            id: 2,
                            namespace: 'test',
                        }
                    },
                }
            ],
            kind: 'Test',
            namespace: 'test',
        },
    },
    deep: 6

}).then(|res| res);
```

Argument | Type | Description
-------- | ---- | -----------
`module` | `string`| The module name.
`query` | `Query` | Object with the query properties, see [Query](#Query).
`timeout` | `integer` (optional) | Provide a custom timeout in seconds (Default: 10 seconds).
`deep` | `integer` | The depth of the deepest object. Every object raises the depth one level. In the examples above the `parent` object is the deepest object and the deep should be 6.

### Query

```ti
future({
    module: 'GCD',
    query: {
        cmd: 'upsert',
        upsert: {
            entities: [
                {
                    key: {
                        id: 3,
                        namespace: 'test',
                        kind: 'Test',
                        parent: {
                            kind: 'Parent',
                            id: 2,
                            namespace: 'test'
                        }
                    },
                    properties: [
                        {
                            name: "age",
                            value: 6
                        }, {
                            name: "kind",
                            value: "dog"
                        }
                    ]
                }
            ],
        },
        next: {
            cmd: 'get',
            get: {
                entities: [
                    {
                        key: {
                            id: 3,
                            namespace: 'test',
                            kind: 'Test',
                            parent: {
                                kind: 'Parent',
                                id: 2,
                                namespace: 'test'
                            }
                        },
                        properties: [
                            {
                                name: "age",
                                value: 6
                            }, {
                                name: "kind",
                                value: "dog"
                            }
                        ]
                    }
                ],
                kind: 'Test',
                namespace: 'test',
            }
        }
    },
    timeout: 10,
    deep: 7

}).then(|res| res);
```

Argument | Type | Description
-------- | ---- | -----------
`cmd` | `string`| The action name, which can be either `upsert`, `get` or `delete`.
`upsert` | `Upsert` (optional) | Object with the `upsert` properties, see [Upsert](#Upsert).
`get` | `Get` (optional) | Object with the `get` properties, see [Get](#Get).
`delete` | `Delete` (optional) | Object with the `delete` properties, see [Delete](#Delete).
`next` | `Query` (optional) | The next query object.
`transaction` | `boolean` (optional) | Indicates if the query needs to be wrapped in a transaction or not.

### Upsert

```ti
future({
    module: 'GCD',
    query: {
        cmd: "upsert",
        upsert: {
            entities: [
                {
                    key: {
                        id: 3,
                        namespace: 'test',
                        kind: 'Test',
                        parent: {
                            kind: 'Parent',
                            id: 2,
                            namespace: 'test',
                        }
                    },
                    properties: [
                        {
                            name: "age",
                            value: 6
                        }, {
                            name: "kind",
                            value: "dog"
                        }
                    ]
                }
            ]
        }
    },
    deep: 6

}).then(|res| res);
```

Argument | Type | Description
-------- | ---- | -----------
`entities` | `list with entities`| A list containing entities that should be either inserted or updated, see [Entity](#Entity).

### Get

```ti
future({
    module: 'GCD',
    query: {
        cmd: "get",
        get: {
            entities: [
                {
                    key: {
                        id: 3,
                        namespace: 'test',
                        kind: 'Test',
                        parent: {
                            kind: 'Parent',
                            id: 2,
                            namespace: 'test',
                        }
                    },
                }
            ],
            kind: 'Test',
            namespace: 'test',
        },
    },
    deep: 6

}).then(|res| res);
```

Argument | Type | Description
-------- | ---- | -----------
`entities` | `list with entities`| A list containing entities that should be returned, see [Entity](#Entity).
`fetch` | `string` (optional) | The way the result will be returned, options are: `Keys` or `Entities`,
`filter` | `Filter` (optional) | Object with the `filter` properties, see [Delete](#Delete).
`kind` | `string` (optional) | A specific entity kind.
`limit` | `integer` (optional) | The maximum number of results that are returned.
`namespace` | `string` (optional) | A specific namespace.
`order` | `Order` (optional) | Object with the `order` properties, see [Order](#Order).

### Delete

```ti
future({
    module: 'GCD',
    query: {
        cmd: "delete",
        delete: {
            entities: [
                {
                    key: {
                        id: 3,
                        namespace: "test",
                        kind: 'Test',
                        parent: {
                            kind: 'Parent',
                            id: 2,
                            namespace: 'test'
                        }
                    }
                }
            ],
        },
    },
    deep: 6
}).then(|res| res);
```

Argument | Type | Description
-------- | ---- | -----------
`entities` | `list with entities`| A list containing entities that should be deleted, see [Entity](#Entity).

### Entity

Argument | Type | Description
-------- | ---- | -----------
`key` | `Key` | Object with the `key` properties, see [Key](#Key).
`properties` | `list with properties` | A list containing properties, see [Property](#Property).

### Key

Argument | Type | Description
-------- | ---- | -----------
`kind` | `string` | A specific entity kind.
`id` | `integer` (optional) | The id of a key. Either `id` or `name` must be zero for the Key to be valid. If both are zero, the `Key` is incomplete.
`name` | `string` (optional) | The name of a key. Either `id` or `name` must be zero for the Key to be valid. If both are zero, the `Key` is incomplete.
`parent` | `Key` (optional) | The parent key.
`namespace` | `string` (optional) | A specific namespace.

### Property

Argument | Type | Description
-------- | ---- | -----------
`name` | `string` | The property name.
`value` | `any` | The property value.
`no_index` | `boolean` | Whether the datastore cannot index this property.

### Filter

Argument | Type | Description
-------- | ---- | -----------
`ancestor` | `Key` | Object with the `key` properties, see [Key](#Key).
`properties` | `list of property filters` | A list containing property filters, see [PropertyFilter](#PropertyFilter).

### PropertyFilter

Argument | Type | Description
-------- | ---- | -----------
`name` | `string` | The property name.
`operator` | `string` | The operator which can be either `=`, `<`, `<=`, `>` or `>=`.
`value` | `any` | The property value to compare.

### Order

Argument | Type | Description
-------- | ---- | -----------
`name` | `string` | The property name.
`direction` | `string` | The direction which can be either an empty string for ascending order or `-` meaning descending order.
