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

Property                | Type            | Description
----------------------- | --------------- | -----------
datastore_project_id    | str (required)  | Project ID.
datastore_emulator_host | str (optional)  | Host of locally-running datastore emulator.
google_app_cred_path    | str (optional)  | Path to the JSON key file for authorization.

Example configuration:

```javascript
set_module_conf('gcd', {
    datastore_project_id: 'id',
    datastore_emulator_host: 'host:port',
    google_app_cred_path: 'path/to/keyfile.json'
});
```

## Exposed functions

Name                        | Description
--------------------------- | -----------
[upsert](#upsert)           | Upsert entities into GCD.
[get](#get)                 | Get entities from GCD.
[delete](#delete)           | Delete entities from GCD.
[transaction](#transaction) | Make a transaction request; provide a `delete`, `get` or `upsert` request and a `next` `delete`, `get` or `upsert` request.

### upsert

#### Arguments

Argument | Type                | Description
---------|---------------------| -----------
`upsert` | `Upsert` (required) | Thing with `upsert` properties, see [Upsert](#Upsert).
`deep`   | `int` (optional)    | Deep value of the thing with `upsert` properties.

#### Example:

```javascript
gcd.upsert({
    entities: [
        {
            key: {
                id: 4,
                namespace: 'test',
                kind: 'Test',
                parent: {
                    kind: 'Parent',
                    id: 2,
                    namespace: 'test',
                }
            },
            properties: [
                {name: 'foo', value: 'FOO', no_index: false},
                {name: 'bar', value: 'BAR', no_index: false},
            ]
        }
    ],
}).then(|res| {
    res; // just return the response.
});
```

### get

#### Arguments

Argument | Type             | Description
---------|------------------| -----------
`get`    | `Get` (required) | Thing with `get` properties, see [Get](#Get).
`deep`   | `int` (optional) | Deep value of the thing with `get` properties.

#### Example:

```javascript
gcd.get({
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
}).then(|res| {
    res; // just return the response.
});
```

### delete

#### Arguments

Argument | Type                | Description
---------|---------------------| -----------
`delete` | `Delete` (required) | Thing with `delete` properties, see [Delete](#Delete).
`deep`   | `int` (optional)    | Deep value of the thing with `delete` properties.

#### Example:

```javascript
gcd.delete({
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
}).then(|res| {
    res; // just return the response.
});
```

### transaction

#### Arguments

Argument      | Type                     | Description
--------------|--------------------------| -----------
`transaction` | `Transaction` (required) | Thing with `transaction` properties, see [Transaction](#Transaction).
`deep`        | `int` (optional)         | Deep value of the thing with `transaction` properties.

#### Example:

```javascript
gcd.transaction({
    get: {
        entities: [
            {
                key: {
                    id: 4,
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
    },
    next: {
        upsert: {
            entities: [
                {
                    key: {
                        id: 4,
                        namespace: 'test',
                        kind: 'Test',
                        parent: {
                            kind: 'Parent',
                            id: 2,
                            namespace: 'test',
                        }
                    },
                    properties: [
                        {name: 'foo', value: 'FOO', no_index: false},
                        {name: 'bar', value: 'BAR', no_index: false},
                    ]
                }
            ],
        },
    }
}).then(|res| {
    res; // just return the response.
});
```

### Types

#### Transaction

Argument      | Type                | Description
--------------|---------------------| -----------
`delete`      | `Delete` (optional) | Thing with `delete` properties, see [Delete](#Delete).
`get`         | `Get` (optional)    | Thing with `get` properties, see [Get](#Get).
`next`        | `Query` (optional)  | The next query Thing.
`upsert`      | `Upsert` (optional) | Thing with the `upsert` properties, see [Upsert](#Upsert).

#### Upsert

Argument   | Type                | Description
---------- | ------------------- | -----------
`entities` | `list with entities`| A list containing entities that should be either inserted or updated, see [Entity](#Entity).

#### Get

Argument    | Type                            | Description
----------- | ------------------------------- | -----------
`ancestor`  | `Key` (required)                | Object with the `key` properties, see [Key](#Key).
`cursor`    | `str` (optional)                | The start cursor.
`entities`  | `list with entities` (optional) | A list containing entities that should be returned, see [Entity](#Entity).
`fetch`     | `str` (optional)                | The way the result will be returned, options are: `Keys` or `Entities`,
`filter`    | `Filter` (optional)             | Object with the `filter` properties, see [Delete](#Delete).
`kind`      | `str` (optional)                | A specific entity kind.
`limit`     | `integer` (optional)            | The maximum number of results that are returned.
`namespace` | `str` (optional)                | A specific namespace.
`order`     | `Order` (optional)              | Object with the `order` properties, see [Order](#Order).

#### Delete

Argument   | Type                            | Description
---------- | ------------------------------- | -----------
`entities` | `list with entities` (required) | A list containing entities that should be deleted, see [Entity](#Entity).

#### Entity

Argument     | Type                              | Description
------------ | --------------------------------- | -----------
`key`        | `Key` (required)                  | Object with the `key` properties, see [Key](#Key).
`properties` | `list with properties` (optional) | A list containing properties, see [Property](#Property).

#### Key

Argument    | Type                 | Description
----------- | -------------------- | -----------
`kind`      | `str` (required)     | A specific entity kind.
`id`        | `integer` (optional) | The id of a key. Either `id` or `name` must be zero for the Key to be valid. If both are zero, the `Key` is incomplete.
`name`      | `str` (optional)     | The name of a key. Either `id` or `name` must be zero for the Key to be valid. If both are zero, the `Key` is incomplete.
`parent`    | `Key` (optional)     | The parent key.
`namespace` | `str` (optional)     | A specific namespace.

#### Property

Argument   | Type              | Description
---------- | ----------------- | -----------
`name`     | `str` (required)  | The property name.
`value`    | `any` (required)  | The property value.
`no_index` | `bool` (optional) | Whether the datastore cannot index this property.

#### Filter

Argument     | Type                                  | Description
------------ | ------------------------------------- | -----------
`properties` | `list of property filters` (required) | A list containing property filters, see [PropertyFilter](#PropertyFilter).

#### PropertyFilter

Argument   | Type             | Description
---------- | ---------------- | -----------
`name`     | `str` (required) | The property name.
`operator` | `str` (required) | The operator which can be either `=`, `<`, `<=`, `>` or `>=`.
`value`    | `any` (required) | The property value to compare.

#### Order

Argument    | Type             | Description
----------- | ---------------- | -----------
`name`      | `str` (required) | The property name.
`direction` | `str` (optional) | The direction which can be either an empty string for ascending order or `-` meaning descending order.
