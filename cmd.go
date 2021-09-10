package main

type Cmd string

const (
	// Upsert -
	UpsertCmd Cmd = "upsert"

	// Get -
	GetCmd Cmd = "get"

	// Delete -
	DeleteCmd Cmd = "delete"
)
