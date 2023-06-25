package main

import (
	"CometDB"
	"fmt"
)

func main() {

	opts := CometDB.DefaultOptions
	opts.DirPath = "/tmp/CometDB"
	db, err := CometDB.Open(opts)
	if err != nil {
		panic(err)
	}

	err = db.Put([]byte("name"), []byte("CometDB"))
	if err != nil {
		panic(err)
	}
	val, err := db.Get([]byte("name"))
	if err != nil {
		panic(err)
	}
	fmt.Println("val = ", string(val))

	err = db.Delete([]byte("name"))
	if err != nil {
		panic(err)
	}
}
