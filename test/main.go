// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"

// 	"cloud.google.com/go/datastore"
// )

// func main() {
// 	projID := os.Getenv("DATASTORE_PROJECT_ID")
// 	if projID == "" {
// 		log.Fatal(`You need to set the environment variable "DATASTORE_PROJECT_ID"`)
// 	}
// 	ctx := context.Background()
// 	client, err := datastore.NewClient(ctx, projID)
// 	if err != nil {
// 		log.Fatalf("Could not create datastore client: %v", err)
// 	}

// 	// Increment a counter.
// 	// See https://cloud.google.com/appengine/articles/sharding_counters for
// 	// a more scalable solution.
// 	type Counter struct {
// 		Count int
// 	}

// 	var count int
// 	key := datastore.NameKey("Counter", "singleton", nil)
// 	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
// 		var x Counter
// 		if err := tx.Get(key, &x); err != nil && err != datastore.ErrNoSuchEntity {
// 			return err
// 		}
// 		x.Count++
// 		if _, err := tx.Put(key, &x); err != nil {
// 			return err
// 		}
// 		count = x.Count
// 		return nil
// 	})
// 	if err != nil {
// 		log.Fatalf("Error: %v", err)
// 	}
// 	// The value of count is only valid once the transaction is successful
// 	// (RunInTransaction has returned nil).
// 	fmt.Printf("Count=%d\n", count)
// }

// package main

// import (
// 	"bufio"
// 	"context"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"text/tabwriter"
// 	"time"

// 	"cloud.google.com/go/datastore"
// )

// func main() {
// 	projID := os.Getenv("DATASTORE_PROJECT_ID")
// 	if projID == "" {
// 		log.Fatal(`You need to set the environment variable "DATASTORE_PROJECT_ID"`)
// 	}
// 	// [START datastore_build_service]
// 	ctx := context.Background()
// 	client, err := datastore.NewClient(ctx, projID)
// 	// [END datastore_build_service]
// 	if err != nil {
// 		log.Fatalf("Could not create datastore client: %v", err)
// 	}
// 	defer client.Close()

// 	// Print welcome message.
// 	fmt.Println("Cloud Datastore Task List")
// 	fmt.Println()
// 	usage()

// 	// Read commands from stdin.
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Print("> ")

// 	for scanner.Scan() {
// 		cmd, args, n := parseCmd(scanner.Text())
// 		switch cmd {
// 		case "insert":
// 			if args == "" {
// 				log.Printf("Missing description in %q command", cmd)
// 				usage()
// 				break
// 			}
// 			key, err := AddTask(ctx, client, args)
// 			if err != nil {
// 				log.Printf("Failed to create task: %v", err)
// 				break
// 			}
// 			fmt.Printf("Created new task with ID %d\n", key.ID)

// 		case "update":
// 			if n == 0 {
// 				log.Printf("Missing numerical task ID in %q command", cmd)
// 				usage()
// 				break
// 			}
// 			if err := MarkDone(ctx, client, n); err != nil {
// 				log.Printf("Failed to mark task done: %v", err)
// 				break
// 			}
// 			fmt.Printf("Task %d marked done\n", n)

// 		case "get":
// 			tasks, err := ListTasks(ctx, client)
// 			if err != nil {
// 				log.Printf("Failed to fetch task list: %v", err)
// 				break
// 			}
// 			PrintTasks(os.Stdout, tasks)

// 		case "delete":
// 			if n == 0 {
// 				log.Printf("Missing numerical task ID in %q command", cmd)
// 				usage()
// 				break
// 			}
// 			if err := DeleteTask(ctx, client, n); err != nil {
// 				log.Printf("Failed to delete task: %v", err)
// 				break
// 			}
// 			fmt.Printf("Task %d deleted\n", n)

// 		default:
// 			log.Printf("Unknown command %q", cmd)
// 			usage()
// 		}

// 		fmt.Print("> ")
// 	}

// 	if err := scanner.Err(); err != nil {
// 		log.Fatalf("Failed reading stdin: %v", err)
// 	}
// }

// // Task is the model used to store tasks in the datastore.
// type Task struct {
// 	Desc    string    `datastore:"description"`
// 	Created time.Time `datastore:"created"`
// 	Done    bool      `datastore:"done"`
// 	id      int64     // The integer ID used in the datastore.
// }

// // AddTask adds a task with the given description to the datastore,
// // returning the key of the newly created entity.
// func AddTask(ctx context.Context, client *datastore.Client, desc string) (*datastore.Key, error) {
// 	task := &Task{
// 		Desc:    desc,
// 		Created: time.Now(),
// 	}
// 	key := datastore.IncompleteKey("Task", nil)
// 	return client.Put(ctx, key, task)
// }

// // MarkDone marks the task done with the given ID.
// func MarkDone(ctx context.Context, client *datastore.Client, taskID int64) error {
// 	// Create a key using the given integer ID.
// 	key := datastore.IDKey("Task", taskID, nil)

// 	// In a transaction load each task, set done to true and store.
// 	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
// 		var task Task
// 		if err := tx.Get(key, &task); err != nil {
// 			return err
// 		}
// 		task.Done = true
// 		_, err := tx.Put(key, &task)
// 		return err
// 	})
// 	return err
// }

// // ListTasks returns all the tasks in ascending order of creation time.
// func ListTasks(ctx context.Context, client *datastore.Client) ([]*Task, error) {
// 	var tasks []*Task

// 	// Create a query to fetch all Task entities, ordered by "created".
// 	query := datastore.NewQuery("Task").Order("created")
// 	keys, err := client.GetAll(ctx, query, &tasks)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Set the id field on each Task from the corresponding key.
// 	for i, key := range keys {
// 		tasks[i].id = key.ID
// 	}

// 	return tasks, nil
// }

// // DeleteTask deletes the task with the given ID.
// func DeleteTask(ctx context.Context, client *datastore.Client, taskID int64) error {
// 	return client.Delete(ctx, datastore.IDKey("Task", taskID, nil))
// }

// // PrintTasks prints the tasks to the given writer.
// func PrintTasks(w io.Writer, tasks []*Task) {
// 	// Use a tab writer to help make results pretty.
// 	tw := tabwriter.NewWriter(w, 8, 8, 1, ' ', 0) // Min cell size of 8.
// 	fmt.Fprintf(tw, "ID\tDescription\tStatus\n")
// 	for _, t := range tasks {
// 		if t.Done {
// 			fmt.Fprintf(tw, "%d\t%s\tdone\n", t.id, t.Desc)
// 		} else {
// 			fmt.Fprintf(tw, "%d\t%s\tcreated %v\n", t.id, t.Desc, t.Created)
// 		}
// 	}
// 	tw.Flush()
// }

// func usage() {
// 	fmt.Print(`Usage:
//   insert_entity <kind> <properties> 							Adds an entity
//   insert_entities <kind> <list of properties> 					Adds multiple entities
//   update_entity <kind> <entity>     							Updates an entity
//   update_entities <kind> <list of entities> 					Updates multiple entities
//   get_entity <kind> <id>  									Gets an entity
//   get_entities 													Gets multiple entities
//   	<kind> <filter (ancestor/property> <order> <limit:cursor>
//   delete_entity <kind> <id>   			  					Deletes an entity
//   delete_entities <kind> <list of ids>  					Deletes multiple entities
// `)
// }

// // parseCmd splits a line into a command and optional extra args.
// // n will be set if the extra args can be parsed as an int64.
// func parseCmd(line string) (cmd, args string, n int64) {
// 	if f := strings.Fields(line); len(f) > 0 {
// 		cmd = f[0]
// 		args = strings.Join(f[1:], " ")
// 	}
// 	if i, err := strconv.ParseInt(args, 10, 64); err == nil {
// 		n = i
// 	}
// 	return cmd, args, n
// }

// package main

// import (
// 	"fmt"

// 	"cloud.google.com/go/datastore"
// 	"github.com/vmihailenco/msgpack/v4"
// )

// func init() {
// 	msgpack.RegisterExt(1, (*key)(nil))
// 	msgpack.RegisterExt(2, (*property)(nil))
// }

// var _ msgpack.Unmarshaler = (*key)(nil)
// var _ msgpack.Unmarshaler = (*property)(nil)

// type key datastore.Key
// type property datastore.Property

// type entity struct {
// 	Key        *key       `msgpack:"key"`
// 	Properties []property `msgpack:"properties"`
// }

// type Entity struct {
// 	Key        *datastore.Key       `msgpack:"key"`
// 	Properties []datastore.Property `msgpack:"properties"`
// }

// func (e *Entity) UnmarshalMsgpack(data []byte) error {
// 	var ret entity
// 	_ = msgpack.Unmarshal(data, &ret)
// 	e.Key = (*datastore.Key)(ret.Key)
// 	e.Properties = make([]datastore.Property, len(ret.Properties))
// 	for i, p := range ret.Properties {
// 		e.Properties[i] = datastore.Property(p)
// 	}

// 	return nil
// }

// func (k *key) UnmarshalMsgpack(data []byte) error {
// 	var ret map[string]interface{}
// 	_ = msgpack.Unmarshal(data, &ret)

// 	ki, ok := ret["kind"].(string)
// 	if ok {
// 		k.Kind = ki
// 	}
// 	i, ok := ret["id"].(int64)
// 	if ok {
// 		k.ID = i
// 	}
// 	n, ok := ret["name"].(string)
// 	if ok {
// 		k.Name = n
// 	}
// 	p, ok := ret["parent"].(*datastore.Key)
// 	if ok {
// 		k.Parent = p
// 	}
// 	ns, ok := ret["namespace"].(string)
// 	if ok {
// 		k.Namespace = ns
// 	}

// 	return nil
// }

// func (p *property) UnmarshalMsgpack(data []byte) error {
// 	var ret map[string]interface{}
// 	_ = msgpack.Unmarshal(data, &ret)
// 	n, ok := ret["name"].(string)
// 	if ok {
// 		p.Name = n
// 	}
// 	p.Value = ret["value"]
// 	ni, ok := ret["no_index"].(bool)
// 	if ok {
// 		p.NoIndex = ni
// 	}
// 	return nil
// }

// type reqMySQL struct {
// 	Cmd       string   `msgpack:"cmd"`
// 	Entities  []Entity `msgpack:"entities"`
// 	Kind      string   `msgpack:"kind"`
// 	Namespace string   `msgpack:"namespace"`
// }

// func main() {
// 	e := map[string]interface{}{
// 		"name":     "testname",
// 		"value":    "testvalue",
// 		"no_index": false,
// 	}
// 	var s []map[string]interface{}
// 	s = append(s, e)
// 	s = append(s, e)
// 	s = append(s, e)

// 	k := map[string]interface{}{
// 		"name":      "name",
// 		"kind":      "testKind",
// 		"id":        1,
// 		"namespace": "testnamespace",
// 	}

// 	entity := map[string]interface{}{
// 		"properties": s,
// 		"key":        k,
// 	}

// 	var slice []map[string]interface{}
// 	slice = append(slice, entity)
// 	newmap := map[string]interface{}{
// 		"cmd":       "InsertEntity",
// 		"kind":      "Test",
// 		"namespace": "test",
// 		"entities":  slice,
// 	}

// 	b, err := msgpack.Marshal(newmap)
// 	_ = err

// 	var out reqMySQL
// 	err = msgpack.Unmarshal(b, &out)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(out.Entities)

// }
