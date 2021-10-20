package main

import (
	"fmt"
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

func errorMsg(msg string) error {
	return fmt.Errorf("Error: " + msg)
}

func retMsg(num int, title string) string {
	label := "entities"
	if num < 2 {
		label = "entity"
	}
	return fmt.Sprintf("%s %d %s", title, num, label)
}
