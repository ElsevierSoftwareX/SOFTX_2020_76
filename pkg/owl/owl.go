/*
Copyright 2020 Institute for Automation of Complex Power Systems,
E.ON Energy Research Center, RWTH Aachen University

This project is licensed under either of
- Apache License, Version 2.0
- MIT License
at your option.

Apache License, Version 2.0:

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

MIT License:

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package owl extracts classes and properties from an owl ontology and maps them to a Go model.
package owl

import (
	"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/rdf"
)

// Ontology holds all information of one ontology
type Ontology struct {
	IRI         string
	Class       map[string]*Class      // all classes (key = iri)
	Property    map[string]*Property   // all properties (key = iri)
	Individual  map[string]*Individual // all individuals (key = iri)
	Imports     map[string][]string    // imported ontologies
	Description map[string]string      // comment about Ontology
	Content     map[string][]byte      // Ontology specification in ttl format
	graph       *rdf.Graph             // graph of parsed owl document
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
