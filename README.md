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

Run the following in the `@thingsdb` scope:

```ti
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

### Query

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
`module` | `string`| The module name.
`query` | `Query` |
`transaction` | `boolean` (optional) | Indicates if the query needs to be wrapped in a transaction or not.
`timeout` | `integer` (optional) | Provide a custom timeout in seconds (Default: 10 seconds).

### Fetch types