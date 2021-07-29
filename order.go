package main

type Direction string

const (
	Ascending  Direction = ""
	Descending Direction = "-"
)

type Order struct {
	Name      string    `msgpack:"name"`
	Direction Direction `msgpack:"direction"`
}
