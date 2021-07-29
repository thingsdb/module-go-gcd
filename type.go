package main

import (
	"reflect"
)

var intList = []reflect.Kind{reflect.Int8, reflect.Int16, reflect.Int32}

func convInvalidInts(value interface{}) interface{} {
	rv := reflect.ValueOf(value)

	for _, intType := range intList {
		if intType == rv.Kind() {
			return rv.Int()
		}
	}

	return value
}
