/*
Copyright 2020 Institute for Automation of Complex Power Systems,
E.ON Energy Research Center, RWTH Aachen University

This project is licensed under either of
- Apache License, Version 2.0
- MIT License
at your option.

Apache License, Version 2.0:

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

MIT License:

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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
	var file *os.File
	// on, err = owl.ExtractOntology("https://data.nasa.gov/ontologies/atmonto/general.ttl")
	on, err = owl.ExtractOntologyLink("https://w3id.org/saref")
	// on, err = owl.ExtractOntologyLink("https://w3id.org/saref4ener")
	// file, _ = os.Open("test.ttl")
	// on, err = owl.ExtractOntology(file)
	if err != nil {
		fmt.Println(err)
		return
	}

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
	// mod, err = owl.MapModel(&on, "git-ce.rwth-aachen.de/acs/private/research/agent/saref.git")
	mod, err = owl.MapModel(&on, "git-ce.rwth-aachen.de/acs/private/research/agent/test.git")
	if err != nil {
		fmt.Println(err)
	}

	err = codegen.GenerateGoCode(mod, "../../test")
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
