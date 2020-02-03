package rdf

import "io"

// EncodeTTL serializes triples in ttl format
func EncodeTTL(triple []Triple, output io.Writer) (err error) {
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
		ret += "^^" + lit.typeIRI
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
	ret = trip.Sub.SerializeTTL() + " " + trip.Pred.SerializeTTL() + " " + trip.Obj.SerializeTTL() + "."
	return
}
