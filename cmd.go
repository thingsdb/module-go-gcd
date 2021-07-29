package main

type Cmd string

const (
	// InsertEntity -
	InsertEntityCmd Cmd = "InsertEntity"

	// InsertEntities -
	InsertEntitiesCmd Cmd = "InsertEntities"

	// UpdateEntity -
	UpdateEntityCmd Cmd = "UpdateEntity"

	// UpdateEntities -
	UpdateEntitiesCmd Cmd = "UpdateEntities"

	// GetEntity -
	GetEntityCmd Cmd = "GetEntity"

	// GetEntities -
	GetEntitiesCmd Cmd = "GetEntities"

	// DeleteEntity -
	DeleteEntityCmd Cmd = "DeleteEntity"

	// DeleteEntities -
	DeleteEntitiesCmd Cmd = "DeleteEntities"
)
