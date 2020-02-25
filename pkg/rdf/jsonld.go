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

package rdf

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/piprate/json-gold/ld"
)

// DecodeJSONLD decodes a jsonld input to rdf triples
func DecodeJSONLD(input io.Reader) (trip []Triple, err error) {
	jsonDec := json.NewDecoder(input)
	var doc interface{}
	err = jsonDec.Decode(&doc)
	if err != nil {
		return
	}

	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Format = "application/n-quads"
	tripleString, err := proc.ToRDF(doc, options)
	if err != nil {
		return
	}
	var file *os.File
	file, err = os.Create("triples.ttl")
	file.Write([]byte(tripleString.(string)))
	file.Close()
	r := strings.NewReader(tripleString.(string))
	trip, err = DecodeTTL(r)
	return
}

// EncodeJSONLD encodes triples in json ld format
func EncodeJSONLD(triple []Triple, output io.Writer) (err error) {
	prefix := make(map[string]string)
	tripleString := ""
	for i := range triple {
		tripleString += triple[i].SerializeTTL(prefix) + "\n"
	}
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	doc, err := proc.FromRDF(tripleString, options)
	if err != nil {
		return
	}
	jsonEnc := json.NewEncoder(output)
	err = jsonEnc.Encode(doc)
	return
}
