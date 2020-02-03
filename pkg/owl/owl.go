package owl

import (
	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/rdf"
)

// Ontology holds all information of one ontology
type Ontology struct {
	Class       map[string]*Class      // all classes (key = iri)
	Property    map[string]*Property   // all properties (key = iri)
	Individual  map[string]*Individual // all individuals (key = iri)
	Imports     map[string][]string
	Description map[string]string // comment about Ontology
	Content     map[string][]byte // Ontology specification in ttl format
	graph       *rdf.Graph
}

// Class is one ontology class
type Class struct {
	Node         *rdf.Node      // graph Node of class
	Parent       []*Class       // all direct parent classes of class
	Child        []*Class       // all direct children classes of class
	Enumeration  []*Individual  // enumerations in owl:OneOf
	Restriction  []*Restriction // restrictions of class in owl:Restriction
	Intersection []*Class       // intersections in owl:IntersectionOf
	Union        []*Class       // unions in owl:UnionOf
	Complement   []*Class       // complements in owl:Complement
	Name         string         // class name (IRI)
	Comment      string         // comment
}

// Restriction is a restriction of a class property
type Restriction struct {
	Node                  *rdf.Node // graph node of restriction
	Property              *Property // property
	ValueConstraint       string    // owl:allValuesFrom, owl:someValuesFrom or owl:hasValue
	CardinalityConstraint string    // owl:maxCardinality, owl:minCardinality or owl:cardinality
	Multiplicity          int       // if restriction has cardinality constraint this shows the multiplicity
	Value                 []string  // value type
}

// Property is one ontology property
type Property struct {
	Node                *rdf.Node   // graph node of property
	Name                string      // name of property (IRI)
	Comment             string      // comment
	Domain              []*Class    // rdfs:domain
	Range               []string    // rdfs:range
	Equivalent          *Property   // owl:equivalentProperty
	Inverse             *Property   // owl:inverseOf
	SubPropertyOf       []*Property //owl:subPropertyOf
	Type                string      // object or datatype property
	IsFunctional        bool        // owl:functionalProperty
	IsInverseFunctional bool        // owlInverseFunctionalProperty
	IsTransitive        bool        // owl:TransistiveProperty
	IsSymmetric         bool        // owl:SysmmetricProperty
}

// Individual is one ontology individual
type Individual struct {
	Node          *rdf.Node     // graph node of individual
	Comment       string        // comment
	Name          string        // name
	Type          *Class        // value type
	SameAs        *Individual   // owl:sameAs
	DifferentFrom *Individual   // owl:differentFrom
	AllDifferent  []*Individual // owl:AllDifferent
}
