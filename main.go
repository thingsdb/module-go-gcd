// Demo is a ThingsDB module which may be used as a template to build modules.
//
// This module simply extract a given `message` property from a request and
// returns this message.
//
// For example:
//
//     // Create the module (@thingsdb scope)
//     new_module('DEMO', 'demo', nil, nil);
//
//     // When the module is loaded, use the module in a future
//     future({
//       module: 'DEMO',
//       message: 'Hi ThingsDB module!',
//     }).then(|msg| {
//	      `Got the message back: {msg}`
//     });
//
package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/datastore"
	timod "github.com/thingsdb/go-timod"
	"github.com/vmihailenco/msgpack/v4"
	"google.golang.org/api/option"
)

var client *datastore.Client = nil
var mux sync.Mutex

type confMySQL struct {
	DatastoreProjectId    string `msgpack:"datastore_project_id"`
	DatastoreEmulatorHost string `msgpack:"datastore_emulator_host"`
	GoogleAppCredPath     string `msgpack:"google_app_cred_path"`
}

type reqMySQL struct {
	Query   *Query `msgpack:"query"`
	Timeout int    `msgpack:"timeout"`
}

func handleConf(config *confMySQL) {
	mux.Lock()
	defer mux.Unlock()

	if client != nil {
		client.Close()
	}

	if config.DatastoreProjectId == "" {
		timod.WriteConfErr()
	}

	var opts []option.ClientOption
	if config.DatastoreEmulatorHost != "" {
		os.Setenv("DATASTORE_EMULATOR_HOST", config.DatastoreEmulatorHost)
	} else if config.GoogleAppCredPath != "" {
		opts = []option.ClientOption{
			option.WithCredentialsFile(config.GoogleAppCredPath),
		}
	}

	ctx := context.Background()

	var err error
	client, err = datastore.NewClient(ctx, config.DatastoreProjectId, opts...)
	if err != nil {
		log.Println("Error: Failed to configure", err)
		timod.WriteConfErr()
		return
	}

	timod.WriteConfOk()
}

func onModuleReq(pkg *timod.Pkg) {
	mux.Lock()
	defer mux.Unlock()

	if client == nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExOperation,
			"Error: GCD is not connected; please check the module configuration")
		return
	}

	var req reqMySQL
	err := msgpack.Unmarshal(pkg.Data, &req)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Failed to unpack GCD request")
		return
	}

	if req.Timeout == 0 {
		req.Timeout = 10
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Duration(req.Timeout)*time.Second)
	defer cancelfunc()

	if req.Query == nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExOperation,
			"Error: Query parameter is required")
		return
	}

	ret, err := req.Query.query(ctx, client)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExOperation,
			err.Error())
		return
	}

	timod.WriteResponse(pkg.Pid, ret)
}

func handler(buf *timod.Buffer, quit chan bool) {
	for {
		select {
		case pkg := <-buf.PkgCh:
			switch timod.Proto(pkg.Tp) {
			case timod.ProtoModuleConf:
				var conf confMySQL
				err := msgpack.Unmarshal(pkg.Data, &conf)
				if err == nil {
					handleConf(&conf)
				} else {
					log.Println("Error: Failed to unpack MySQL configuration")
					timod.WriteConfErr()
				}

			case timod.ProtoModuleReq:
				onModuleReq(pkg)

			default:
				log.Printf("Error: Unexpected package type: %d", pkg.Tp)
			}
		case err := <-buf.ErrCh:
			// In case of an error you probably want to quit the module.
			// ThingsDB will try to restart the module a few times if this
			// happens.
			log.Printf("Error: %s", err)
			quit <- true
		}
	}
}

func main() {
	// Starts the module
	timod.StartModule("gcd", handler)

	if client != nil {
		client.Close()
	}
}
