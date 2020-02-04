package owl

import (
	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/rdf"
)

// Thing is the common base class of all types (owl:Thing)
type Thing interface {
	IRI() string
	String() string
	InitFromNode(*rdf.Node) error
	ToGraph(*rdf.Graph)
	RemoveObject(Thing, string)
}
