package rdf

// Triple is one rdf triple
type Triple struct {
	Sub  Subject
	Pred Predicate
	Obj  Object
}

// TermType is the type of a RDF term
type TermType int

// possible RDF term types
const (
	TermBlankNode TermType = iota
	TermLiteral
	TermIRI
)

// Term represents an RDF term. possible types: Blank node, Literal and IRI.
type Term interface {
	Type() (typ TermType)
	String() (str string)
	SerializeTTL() (str string)
}

// Subject reprsents the subject of a rdf triple
type Subject interface {
	Term
}

// Predicate represents the predicate of a rdf triple
type Predicate interface {
	Term
}

// Object represents the object of a rdf triple
type Object interface {
	Term
}

// IRI is a possible RDF term
type IRI struct {
	name string
}

// Type denotes the term type
func (iri IRI) Type() (typ TermType) {
	typ = TermIRI
	return
}

// String prints the IRI
func (iri IRI) String() (str string) {
	str = iri.name
	return
}

// Literal is a possible RDF term
type Literal struct {
	str     string
	typeIRI string
	langTag string
	value   interface{}
}

// Type denotes the term type
func (lit Literal) Type() (typ TermType) {
	typ = TermLiteral
	return
}

// String prints the literal string
func (lit Literal) String() (str string) {
	str = lit.str
	return
}

// BlankNode is a possible RDF term
type BlankNode struct {
	name string
}

// Type denotes the term type
func (blank BlankNode) Type() (typ TermType) {
	typ = TermBlankNode
	return
}

// String prints the name
func (blank BlankNode) String() (str string) {
	str = blank.name
	return
}
