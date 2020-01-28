package rdf

// Triple is one rdf triple
type Triple struct {
	Sub  Subject
	Pred Predicate
	Obj  Object
}

// Term represents an RDF term. possible types: Blank node, Literal and IRI.
type Term interface {
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
