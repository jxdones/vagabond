package commands

import "fmt"

const VERSION = "0.0.1"

func ShowVersion(_ []string) error {
	fmt.Println("vagabond version", VERSION)
	return nil
}
