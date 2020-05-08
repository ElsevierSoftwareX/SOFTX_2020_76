/*
This file is auto-generated by OWL2Go (https://git.rwth-aachen.de/acs/public/ontology/owl/owl2go).

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
	"strconv"
	"time"

	"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/codegen"
	"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/owl"
	"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/test/input"
)

func main() {
	var err error
	err = generateInputFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = measurePerformance(1000)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func generateInputFiles() (err error) {
	var file *os.File
	file, err = os.Create("input10.ttl")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = input.GenerateInputFile(10, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err = os.Create("input100.ttl")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = input.GenerateInputFile(100, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err = os.Create("input1000.ttl")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = input.GenerateInputFile(1000, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func measurePerformance(size int) (err error) {
	var file *os.File

	var on owl.Ontology

	tstart := time.Now()
	file, err = os.Open("input" + strconv.Itoa(size) + ".ttl")
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
	fmt.Println("Time for extraction: ", time.Since(tstart))

	tstart = time.Now()
	var mod []owl.GoModel
	mod, err = owl.MapModel(&on, "test")
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	fmt.Println("Time for mapping: ", time.Since(tstart))

	tstart = time.Now()
	err = codegen.GenerateGoCode(mod, "output"+strconv.Itoa(size))
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	fmt.Println("Time for generation: ", time.Since(tstart))
	return
}
