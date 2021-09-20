# module-go-gcd
ThingsDB module for communication with the Google Cloud Datastore

## Building the module

To build the GCD module, make sure Go is installed and configured.

First go to the module path and run the following command to install the module dependencies:

```
go mod tidy
```

Next, you are ready to build the module:

```
go build
```

Copy the created binary file to the ThingsDB module path.

## Configure the module

The GCD module must be configured before it can be used.

In this example we will name the module `GCD`. This name is arbitrary and can be anything you like.

Run the following code in the `@thingsdb` scope:

```ti
// The values MUST be changed according to your situation
new_module(
    "GCD",
    "module-go-gcd",
    {
        datastore_project_id: "id",
        datastore_emulator_host: "host:port",
        google_app_cred_path: "path/to/keyfile.json"
    }
);
```

### Arguments

Argument | Type | Description
-------- | ---- | -----------
`datastore_project_id` | `string` | Sets the project ID.
`datastore_emulator_host` | `string` (optional) | Sets the DATASTORE_EMULATOR_HOST environment variable. The client will connect to a locally-running datastore emulator when this value is set.
`google_app_cred_path` | `string` (optional) | Sets the path to the JSON key file for authorization.

```ti
new_module(
    "GCD",
    "module-go-gcd",
    {
        datastore_project_id: "id",
        datastore_emulator_host: "localhost:8085"
    }
);
```

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
    deep: 10

}).then(|res| res);
```

Argument | Type | Description
-------- | ---- | -----------
`module` | `string`| The module name.
`query` | `Query` | Object with the query properties, see [Query](#Query).
`transaction` | `boolean` (optional) | Indicates if the query needs to be wrapped in a transaction or not.
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
    transaction: false,
    timeout: 10,
    deep: 5

}).then(|res| res);
```

Argument | Type | Description
-------- | ---- | -----------
`cmd` | `string`| The action name, which can be either `upsert`, `get` or `delete`.
`upsert` | `Upsert` (optional) | Object with the `upsert` properties, see [Upsert](#Upsert).
`get` | `Get` (optional) | Object with the `get` properties, see [Get](#Get).
`delete` | `Delete` (optional) | Object with the `delete` properties, see [Delete](#Delete).
`next` | `Query` (optional) | The next query object.

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
    deep: 10

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
    deep: 10

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
    deep: 5
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
