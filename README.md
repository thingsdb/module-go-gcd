# module-go-gcd
ThingsDB module for communication with the Google Cloud Datastore

```ti
new_module(
    "GCD",
    "module-go-gcd",
    {
        datastore_project_id: "id",
        datastore_emulator_host: "localhost:8085"
    }
);

future({
    module: 'GCD',
    upsert: {
        entities: [
            {
                key: {
                    ID: 3,
                    Namespace: 'test',
                    Kind: 'Test',
                    Parent: {
                        Kind: 'Parent',
                        ID: 3,
                        Namespace: 'test'
                    }
                },
                properties: [
                    {
                        Name: "age",
                        Value: 6
                    }, {
                        Name: "kind",
                        Value: "dog"
                    }
                ]
            }
        ],
        kind: 'Test',
        namespace: 'test',
    },
    deep: 5

}).then(|res| res);

future({
    module: 'GCD',
    get: {
        entities: [
            {
                key: {
                    ID: 3,
                    Namespace: "test",
                    Kind: 'Test',
                    Parent: {
                        Kind: 'Parent',
                        ID: 3,
                        Namespace: 'test'
                    }
                }
            }
        ],
        kind: 'Test',
        namespace: 'test',
    },
    deep: 5
}).then(|res| res);

future({
    module: 'GCD',
    delete: {
        entities: [
            {
                key: {
                    ID: 3,
                    Namespace: "test",
                    Kind: 'Test',
                    Parent: {
                        Kind: 'Parent',
                        ID: 3,
                        Namespace: 'test'
                    }
                }
            }
        ],
    },
    deep: 5
}).then(|res| res);
```
