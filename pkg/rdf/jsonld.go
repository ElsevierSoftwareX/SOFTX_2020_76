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
	var replaced map[string]string
	triple, replaced = addHTTP(triple)
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
	var js []byte
	js, err = json.Marshal(doc)
	if err != nil {
		return
	}
	outString := string(js)
	for i, r := range replaced {
		outString = strings.Replace(outString, i, r, -1)
	}

	_, err = output.Write([]byte(outString))
	// jsonEnc := json.NewEncoder(output)
	// err = jsonEnc.Encode(doc)
	// if err != nil {
	// 	return
	// }
	return
}

// addHTTP adds "http://" to avoid errors of json-gold package
func addHTTP(triple []Triple) (ret []Triple, replaced map[string]string) {
	replaced = make(map[string]string)
	for i := range triple {
		if triple[i].Sub.Type() == TermIRI {
			if !strings.HasPrefix(triple[i].Sub.String(), "http") {
				if _, ok := replaced["http://"+triple[i].Sub.String()]; !ok {
					replaced["http://"+triple[i].Sub.String()] = triple[i].Sub.String()
				}
				iri := NewIRI("http://" + triple[i].Sub.String())
				triple[i].Sub = iri
			}
		}
		if triple[i].Obj.Type() == TermIRI {
			if !strings.HasPrefix(triple[i].Obj.String(), "http") {
				if _, ok := replaced["http://"+triple[i].Obj.String()]; !ok {
					replaced["http://"+triple[i].Obj.String()] = triple[i].Obj.String()
				}
				iri := NewIRI("http://" + triple[i].Obj.String())
				triple[i].Obj = iri
			}
		}
	}
	ret = triple
	return
}