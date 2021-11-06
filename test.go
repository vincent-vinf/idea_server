package main

import "idea_server/util"

func main() {

	// or error handling
	//u2 := uuid.NewV4()
	//
	//fmt.Printf("UUIDv4: %s\n", u2)
	//
	//// Parsing UUID from string input
	//u2, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	util.UploadFile("f.png", "f.png")
}
