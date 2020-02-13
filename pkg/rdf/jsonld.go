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
	if err != nil {
		return
	}
	return
}
