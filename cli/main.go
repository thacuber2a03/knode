package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/thacuber2a03/knode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <file>", os.Args[0])
		os.Exit(-1)
	}

	buf, e := os.ReadFile(os.Args[1])
	if e != nil {
		log.Fatal(e)
	}

	node, e := knode.ParseFromSlice(buf)
	if e.(*knode.ParseError) != nil {
		log.Fatal(e)
	}

	buf, e = json.MarshalIndent(node, "", "\t")
	if e != nil {
		log.Fatal(e)
	}

	fmt.Printf("%s\n", buf)
}
