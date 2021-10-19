package main

import "fmt"

func errorMsg(msg string) error {
	return fmt.Errorf("Error: " + msg)
}
