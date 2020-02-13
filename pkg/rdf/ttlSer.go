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
		if triple[i].Obj.Type() == TermIRI {
			iri := triple[i].Obj.(IRI).name
			prCounter = checkPrefix(iri, prefix, singleOccurrence, prCounter)
		} else if triple[i].Obj.Type() == TermLiteral && triple[i].Obj.(Literal).typeIRI != "" {
			iri := triple[i].Obj.(Literal).typeIRI
			prCounter = checkPrefix(iri, prefix, singleOccurrence, prCounter)
		}
	}

	for i := range prefix {
		output.Write([]byte("@prefix " + prefix[i] + ": <" + i + "> .\n"))
	}

	for i := range triple {
		_, err = output.Write([]byte(triple[i].SerializeTTL(prefix) + "\n"))
		if err != nil {
			return
		}
	}
	return
}

// SerializeTTL serializes IRI in ttl format
func (iri IRI) SerializeTTL(prefix map[string]string) (ret string) {
	pr := getPrefix(iri.name)
	if temp, ok := prefix[pr]; ok {
		rest := strings.Split(iri.name, pr)
		ret = temp + ":" + rest[len(rest)-1]
	} else {
		ret = "<" + iri.name + ">"
	}
	return
}

// SerializeTTL serializes Literal in ttl format
func (lit Literal) SerializeTTL(prefix map[string]string) (ret string) {
	ret = "\"" + lit.str + "\""
	if lit.langTag != "" {
		ret += "@" + lit.langTag
	}
	if lit.typeIRI != "" {
		ret += "^^"
		pr := getPrefix(lit.typeIRI)
		if temp, ok := prefix[pr]; ok {
			rest := strings.Split(lit.typeIRI, pr)
			ret += temp + ":" + rest[len(rest)-1]
		} else {
			ret += "<" + lit.typeIRI + ">"
		}
	}
	return
}

// SerializeTTL serializes blank node in ttl format
func (blank BlankNode) SerializeTTL(prefix map[string]string) (ret string) {
	ret = "_:" + blank.name
	return
}

// SerializeTTL serializes a single Triple in ttl format
func (trip Triple) SerializeTTL(prefix map[string]string) (ret string) {
	ret = trip.Sub.SerializeTTL(prefix) + " " + trip.Pred.SerializeTTL(prefix) + " " + trip.Obj.SerializeTTL(prefix) + " ."
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
		if st := standardPrefix(prefTemp); st == "" {
			prefix[prefTemp] = "pr" + strconv.Itoa(retCounter)
			retCounter++
		} else {
			prefix[prefTemp] = st
		}
		delete(single, prefTemp)
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

// standardPrefix checks some standard prefixes
func standardPrefix(name string) (pref string) {
	switch name {
	case "http://www.wurvoc.org/vocabularies/om-1.8/":
		pref = "om"
	case "http://www.w3.org/2002/07/owl#":
		pref = "owl"
	case "http://www.w3.org/1999/02/22-rdf-syntax-ns#":
		pref = "rdf"
	case "http://www.w3.org/2001/XMLSchema#":
		pref = "xsd"
	case "http://www.w3.org/2000/01/rdf-schema#":
		pref = "rdfs"
	default:
		pref = ""
	}
	return
}
