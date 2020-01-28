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
	TermBlank TermType = iota
	TermLiteral
	TermIRI
)

// Term represents an RDF term. possible types: Blank node, Literal and IRI.
type Term interface {
	Type() (typ TermType)
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

// Literal is a possible RDF term
type Literal struct {
	value   string
	typeIRI string
}

// Type denotes the term type
func (lit Literal) Type() (typ TermType) {
	typ = TermLiteral
	return
}
