package rdf

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// EncodeTTL serializes triples in ttl format
func EncodeTTL(triple []Triple, output io.Writer) (err error) {
	prefix := make(map[string]string)
	singleOccurrence := make(map[string]interface{})
	prCounter := 0
	for i := range triple {
		if triple[i].Sub.Type() == TermIRI {
			iri := triple[i].Sub.(IRI).name
			prCounter = checkPrefix(iri, prefix, singleOccurrence, prCounter)
		}
		if triple[i].Pred.Type() == TermIRI {
			iri := triple[i].Pred.(IRI).name
			prCounter = checkPrefix(iri, prefix, singleOccurrence, prCounter)
		}
	}

	fmt.Println(prefix)
	for i := range prefix {
		output.Write([]byte("@prefix " + prefix[i] + ": " + i + " .\n"))
	}

	for i := range triple {
		_, err = output.Write([]byte(triple[i].SerializeTTL() + "\n"))
		if err != nil {
			return
		}
	}
	return
}

// SerializeTTL serializes IRI in ttl format
func (iri IRI) SerializeTTL() (ret string) {
	ret = "<" + iri.name + ">"
	return
}

// SerializeTTL serializes Literal in ttl format
func (lit Literal) SerializeTTL() (ret string) {
	ret = "\"" + lit.str + "\""
	if lit.langTag != "" {
		ret += "@" + lit.langTag
	}
	if lit.typeIRI != "" {
		ret += "^^<" + lit.typeIRI + ">"
	}
	return
}

// SerializeTTL serializes blank node in ttl format
func (blank BlankNode) SerializeTTL() (ret string) {
	ret = "_:" + blank.name + ">"
	return
}

// SerializeTTL serializes a single Triple in ttl format
func (trip Triple) SerializeTTL() (ret string) {
	ret = trip.Sub.SerializeTTL() + " " + trip.Pred.SerializeTTL() + " " + trip.Obj.SerializeTTL() + " ."
	return
}

// checkPrefix checks if iri starts with same string as other iris and adds a prefix
func checkPrefix(iri string, prefix map[string]string, single map[string]interface{}, prCounter int) (retCounter int) {
	retCounter = prCounter

	prefTemp := getPrefix(iri)
	if prefTemp == "" {
		return
	}

	if _, ok := prefix[prefTemp]; ok {
		// prefix already stored
		return
	}

	if _, ok := single[prefTemp]; ok {
		// prefix has been determined before
		prefix[prefTemp] = "pr" + strconv.Itoa(retCounter)
		delete(single, prefTemp)
		retCounter++
		fmt.Println(prefTemp)
	} else {
		single[prefTemp] = nil
	}

	return
}

// getPrefix determines the prefix of an iri
func getPrefix(iri string) (prefix string) {
	sep := "#"
	temp := strings.Split(iri, sep)
	if len(temp) == 1 {
		sep = "/"
		temp = strings.Split(iri, sep)
	}
	if len(temp) == 1 {
		return
	}
	for i := 0; i < len(temp)-1; i++ {
		prefix += temp[i] + sep
	}

	return
}
