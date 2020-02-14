package main

import (
	"fmt"
	"os"

	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/codegen"
	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/owl"
	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/rdf"
	// "git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/codegen"
)

func main() {
	// testjsonld()
	//
	// for i := range triples {
	// 	fmt.Fprintln(file, triples[i])
	// }
	// fmt.Println(triples)

	var err error
	var on owl.Ontology
	// on, err = owl.ExtractOntology("https://data.nasa.gov/ontologies/atmonto/general.ttl")
	on, err = owl.ExtractOntology("https://w3id.org/saref")
	if err != nil {
		fmt.Println(err)
		return
	}

	var file *os.File
	// classes
	file, err = os.Create("classes.dat")
	for i := range on.Class {
		fmt.Fprintln(file, on.Class[i].String())
	}
	file.Close()

	// properties
	file, err = os.Create("properties.dat")
	for i := range on.Property {
		fmt.Fprintln(file, on.Property[i].String())
	}
	file.Close()

	// individuals
	file, err = os.Create("individuals.dat")
	for i := range on.Individual {
		fmt.Fprintln(file, on.Individual[i].String())
	}
	file.Close()

	var mod []owl.GoModel
	mod, err = owl.MapModel(&on, "git-ce.rwth-aachen.de/acs/private/research/agent/saref.git")
	if err != nil {
		fmt.Println(err)
	}

	err = codegen.GenerateGoCode(mod, "../../saref")
	if err != nil {
		fmt.Println(err)
	}
}

func testjsonld() {
	var err error
	var file *os.File
	file, err = os.Open("time.ttl")
	if err != nil {
		fmt.Println(err)
	}
	var triples []rdf.Triple
	triples, err = rdf.DecodeTTL(file)
	if err != nil {
		fmt.Println(err)
	}
	file, err = os.Create("triples.json")
	err = rdf.EncodeJSONLD(triples, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	file.Close()

	file, err = os.Open("triples.json")
	if err != nil {
		fmt.Println(err)
	}
	triples, err = rdf.DecodeJSONLD(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(triples)
	return
}
