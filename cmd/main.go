package main

import (
	"fmt"
	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/rdf"
	"os"
)

func main() {
	var err error
	var file *os.File
	file, err = os.Open("test.ttl")
	if err != nil {
		fmt.Println(err)
	}
	var triples []rdf.Triple
	triples, err = rdf.DecodeTTL(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(triples)
}
