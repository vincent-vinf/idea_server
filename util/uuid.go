package util

import (
	"fmt"
	"github.com/satori/go.uuid"
)

func test() {
	u2 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u2)

	// Parsing UUID from string input
	u2, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	fmt.Printf("Successfully parsed: %s", u2)
}
