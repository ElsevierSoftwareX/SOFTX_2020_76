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

	"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/codegen"
	"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/owl"
)

func main() {
	var err error
	if len(os.Args) != 5 {
		fmt.Println("Error: Wrong number of command line arguments")
		return
	}
	f := false
	if os.Args[1] == "-f" {
		f = true
	} else if os.Args[1] == "-l" {

	} else {
		fmt.Println("Error: Wrong ontology location (-f or -l)")
		return
	}

	ontLoc := os.Args[2]
	module := os.Args[3]
	path := os.Args[4]

	var on owl.Ontology

	if f {
		var file *os.File
		file, err = os.Open(ontLoc)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		on, err = owl.ExtractOntology(file)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		file.Close()
	} else {
		on, err = owl.ExtractOntologyLink(ontLoc)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
	}

	fmt.Println(on.Imports)

	var file *os.File
	// classes
	file, err = os.Create("classes.dat")
	for i := range on.Class {
		fmt.Fprintln(file, on.Class[i].String())
	}
	file.Close()

	// // properties
	// file, err = os.Create("properties.dat")
	// for i := range on.Property {
	// 	fmt.Fprintln(file, on.Property[i].String())
	// }
	// file.Close()

	// // individuals
	// file, err = os.Create("individuals.dat")
	// for i := range on.Individual {
	// 	fmt.Fprintln(file, on.Individual[i].String())
	// }
	// file.Close()

	var mod owl.GoModel
	mod, err = owl.MapModel(&on, module)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	err = codegen.GenerateGoCode(mod, path)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
}
