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

package owl

import (
	"errors"
	"fmt"
	"strconv"

	"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/rdf"
)

// extractClasses extracts all classes from a graph and fills them with basic information
func extractClasses(g *rdf.Graph) (classes map[string]*Class, err error) {
	fmt.Println("\tExtract classes")
	classes = make(map[string]*Class)
	// detrmine all classes
	for i := range g.Nodes {
		for j := range g.Nodes[i].Edge {
			if g.Nodes[i].Edge[j].Pred.String() ==
				"http://www.w3.org/1999/02/22-rdf-syntax-ns#type" &&
				g.Nodes[i].Edge[j].Object.Term.String() == "http://www.w3.org/2002/07/owl#Class" {
				isDeprecated := false
				for k := range g.Nodes[i].Edge {
					if g.Nodes[i].Edge[k].Pred.String() ==
						"http://www.w3.org/1999/02/22-rdf-syntax-ns#type" &&
						g.Nodes[i].Edge[k].Object.Term.String() ==
							"http://www.w3.org/2002/07/owl#DeprecatedClass" {
						isDeprecated = true
						break
					}
				}
				if isDeprecated {
					break
				}
				class := Class{
					Node:    g.Nodes[i],
					Name:    g.Nodes[i].Term.String(),
					Comment: getComment(g.Nodes[i]),
				}
				classes[g.Nodes[i].Term.String()] = &class
				break
			}
		}
	}
	return
}

// postProcessClasses fills parents, children and restrictions of classes
func (on *Ontology) postProcessClasses() (err error) {
	fmt.Println("\tPostprocess classes")
	for i := range on.Class {
		err = on.Class[i].extractInheritance(on)
		if err != nil {
			err = errors.New(err.Error() + " class " + on.Class[i].Name)
			return
		}
		err = on.Class[i].extractEnumeration(on)
		if err != nil {
			err = errors.New(err.Error() + " class " + on.Class[i].Name)
			return
		}
		err = on.Class[i].extractRestrictions(on)
		if err != nil {
			err = errors.New(err.Error() + " class " + on.Class[i].Name)
			return
		}
		err = on.Class[i].extractUnionOf(on)
		if err != nil {
			err = errors.New(err.Error() + " class " + on.Class[i].Name)
			return
		}
	}
	for i := range on.Class {
		on.Class[i].fillEmptyRestrictionValues()
	}
	return
}

// extractInheritance extracts all direct parent classes and adds itself to children of these
// classes
func (class *Class) extractInheritance(on *Ontology) (err error) {
	for i := range class.Node.Edge {
		if class.Node.Edge[i].Pred.String() == "http://www.w3.org/2000/01/rdf-schema#subClassOf" {
			if parent, ok := on.Class[class.Node.Edge[i].Object.Term.String()]; ok &&
				class.Node.Edge[i].Object.Term.Type() != rdf.TermBlankNode {
				class.Parent = append(class.Parent, parent)
				parent.Child = append(parent.Child, class)
			}
		}
	}
	return
}

// extractEnumeration extracts all Individuals the class consists of
func (class *Class) extractEnumeration(on *Ontology) (err error) {
	for i := range class.Node.Edge {
		if class.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#oneOf" {
			enumNodes := getUnionValues(class.Node.Edge[i].Object)
			for j := range enumNodes {
				if ind, ok := on.Individual[enumNodes[j].Term.String()]; ok {
					class.Enumeration = append(class.Enumeration, ind)
				} else {
					err = errors.New("unknown individual " + enumNodes[j].Term.String())
					fmt.Println(class.Name)
					return
				}
			}
		}
	}
	return
}

// extractRestrictions extracts all restrictions of a class and unions of restrictions
func (class *Class) extractRestrictions(on *Ontology) (err error) {
	for i := range class.Node.Edge {
		if class.Node.Edge[i].Pred.String() == "http://www.w3.org/2000/01/rdf-schema#subClassOf" {
			parent := class.Node.Edge[i].Object
			for j := range parent.Edge {
				// single restriction
				if parent.Edge[j].Pred.String() ==
					"http://www.w3.org/1999/02/22-rdf-syntax-ns#type" &&
					parent.Edge[j].Object.Term.String() ==
						"http://www.w3.org/2002/07/owl#Restriction" {
					rest := Restriction{
						Node: parent,
					}
					err = rest.extractRestriction(on)
					if err != nil {
						return
					}
					if rest.Property == nil {
						err = errors.New("invalid restriction: " + class.Name)
						return
					}
					restExist := false
					for k := range class.Restriction {
						if class.Restriction[k].Property.Name == rest.Property.Name {
							class.Restriction[k].mergeRestriction(&rest)
							restExist = true
							break
						}
					}
					if !restExist {
						class.Restriction = append(class.Restriction, &rest)
					}
					// union of restrictions
				} else if parent.Edge[j].Pred.String() == "http://www.w3.org/2002/07/owl#unionOf" {
					unionNodes := getUnionValues(parent.Edge[j].Object)
					for k := range unionNodes {
						for l := range unionNodes[k].Edge {
							if unionNodes[k].Edge[l].Pred.String() ==
								"http://www.w3.org/1999/02/22-rdf-syntax-ns#type" &&
								unionNodes[k].Edge[l].Object.Term.String() ==
									"http://www.w3.org/2002/07/owl#Restriction" {
								rest := Restriction{
									Node: unionNodes[k],
								}
								rest.extractRestriction(on)
								restExist := false
								for m := range class.Restriction {
									if class.Restriction[m].Property.Name == rest.Property.Name {
										class.Restriction[m].mergeRestriction(&rest)
										restExist = true
										break
									}
								}
								if rest.Property == nil {
									err = errors.New("invalid restriction: ")
									return
								}
								if !restExist {
									class.Restriction = append(class.Restriction, &rest)
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

// extractRestriction extracts information of a restriction
func (rest *Restriction) extractRestriction(on *Ontology) (err error) {
	err = rest.extractValueConstraint(on)
	if err != nil {
		return
	}
	if rest.ValueConstraint == "" {
		err = rest.extractCardinalityConstraint(on)
	}
	if err != nil {
		return
	}
	if len(rest.Value) == 0 && len(rest.Property.Range) > 0 {
		rest.Value = append(rest.Value, rest.Property.Range...)
	}
	return
}

// extractValueConstraint extracts the value constraint of a restriction
func (rest *Restriction) extractValueConstraint(on *Ontology) (err error) {
	property := on.Property
	rest.ValueConstraint = ""
	for i := range rest.Node.Edge {
		if rest.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#onProperty" {
			if prop, ok := property[rest.Node.Edge[i].Object.Term.String()]; ok {
				rest.Property = prop
			} else {
				err = errors.New("unknown property " + rest.Node.Edge[i].Object.Term.String())
				return
			}
		} else if rest.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#allValuesFrom" {
			rest.ValueConstraint = "http://www.w3.org/2002/07/owl#allValuesFrom"
			rest.Value = append(rest.Value, getRestrictionValues(rest.Node.Edge[i].Object)...)
		} else if rest.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#hasValue" {
			rest.ValueConstraint = "http://www.w3.org/2002/07/owl#hasValue"
			rest.Value = append(rest.Value, getRestrictionValues(rest.Node.Edge[i].Object)...)
		} else if rest.Node.Edge[i].Pred.String() ==
			"http://www.w3.org/2002/07/owl#someValuesFrom" {
			rest.ValueConstraint = "http://www.w3.org/2002/07/owl#someValuesFrom"
			rest.Value = append(rest.Value, getRestrictionValues(rest.Node.Edge[i].Object)...)
		}
	}
	return
}

// extractCardinalityConstraint extracts the cardinality constraint of a restriction
func (rest *Restriction) extractCardinalityConstraint(on *Ontology) (err error) {
	property := on.Property
	rest.CardinalityConstraint = ""
	for i := range rest.Node.Edge {
		if rest.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#onProperty" {
			if prop, ok := property[rest.Node.Edge[i].Object.Term.String()]; ok {
				rest.Property = prop
			} else {
				err = errors.New("unknown property " + rest.Node.Edge[i].Object.Term.String())
				return
			}
		} else if rest.Node.Edge[i].Pred.String() ==
			"http://www.w3.org/2002/07/owl#minQualifiedCardinality" ||
			rest.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#minCardinality" ||
			rest.Node.Edge[i].Pred.String() ==
				"http://www.w3.org/2002/07/owl#qualifiedCardinality" ||
			rest.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#cardinality" ||
			rest.Node.Edge[i].Pred.String() ==
				"http://www.w3.org/2002/07/owl#maxQualifiedCardinality" ||
			rest.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#maxCardinality" {
			rest.CardinalityConstraint = rest.Node.Edge[i].Pred.String()
			num, err := strconv.Atoi(rest.Node.Edge[i].Object.Term.String())
			if err == nil {
				rest.Multiplicity = num
			} else {
				rest.Multiplicity = 1
			}
		} else if rest.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#onClass" ||
			rest.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#onDataRange" {
			rest.Value = append(rest.Value, getRestrictionValues(rest.Node.Edge[i].Object)...)
		}
	}
	return
}

// mergeRestriction merges values of the mRest into rest if both only differ in value
func (rest *Restriction) mergeRestriction(mRest *Restriction) (err error) {
	if rest.Property.Name == mRest.Property.Name && rest.ValueConstraint == mRest.ValueConstraint &&
		rest.CardinalityConstraint == mRest.CardinalityConstraint {
		for i := range mRest.Value {
			inRest := false
			for j := range rest.Value {
				if mRest.Value[i] == rest.Value[j] {
					inRest = true
					break
				}
			}
			if !inRest {
				rest.Value = append(rest.Value, mRest.Value[i])
			}
		}
	} else {
		err = errors.New("cannot merge restriction " + rest.Property.Name)
	}
	return
}

// fillEmptyRestrictionValues sets empty restrictions values to child's values
func (class *Class) fillEmptyRestrictionValues() (err error) {
	for i := range class.Restriction {
		if len(class.Restriction[i].Value) == 0 {
			var val []string
			val, err = class.getChildRestrictionValue(class.Restriction[i])
			if len(val) > 0 && err == nil {
				class.Restriction[i].Value = val
			}
		}
	}
	return
}

// getChildRestrictionValue returns restriction values of a child if only one child exists
func (class *Class) getChildRestrictionValue(rest *Restriction) (value []string, err error) {
	restExist := false
	if len(class.Child) == 1 {
		for i := range class.Child[0].Restriction {
			if rest.Property.Name == class.Child[0].Restriction[i].Property.Name {
				if len(class.Child[0].Restriction[i].Value) > 0 {
					value = append(value, class.Child[0].Restriction[i].Value...)
				} else {
					value, err = class.Child[0].getChildRestrictionValue(rest)
				}
				restExist = true
				break
			}
		}
	}
	if !restExist {
		err = errors.New("No child restriction " + rest.Property.Name + " available")
	}
	return
}

// getRestrictionValues returns the values of a restriction
func getRestrictionValues(node *rdf.Node) (value []string) {
	isUnion := false
	// isClass := false
	// for i := range node.Predicates {
	// 	if node.Predicates[i].Name == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" &&
	// 		node.Predicates[i].Object.Name == "http://www.w3.org/2002/07/owl#Class" {
	// 		isClass = true
	// 		break
	// 	}
	// }
	if node.Term.Type() == rdf.TermBlankNode {
		for i := range node.Edge {
			if node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#unionOf" {
				isUnion = true
				temp := getUnionValues(node.Edge[i].Object)
				for j := range temp {
					value = append(value, temp[j].Term.String())
				}
				break
			}
		}
	}
	if !isUnion {
		value = append(value, node.Term.String())
	}
	return
}

// extractUnionOf extracts unionOf
func (class *Class) extractUnionOf(on *Ontology) (err error) {
	for i := range class.Node.Edge {
		if class.Node.Edge[i].Pred.String() == "http://www.w3.org/2002/07/owl#unionOf" {
			unionNodes := getUnionValues(class.Node.Edge[i].Object)
			for j := range unionNodes {
				if unionClass, ok := on.Class[unionNodes[j].Term.String()]; ok {
					class.Union = append(class.Union, unionClass)
				}
			}
		}
	}
	return
}

// GetBaseClass returns the base class of all values
func GetBaseClass(values []string, classes map[string]*Class) (base *Class, err error) {
	base = nil
	if len(values) == 1 {
		if _, ok := classes[values[0]]; !ok {
			err = errors.New("unknown class: " + values[0])
			return
		}
		base, _ = classes[values[0]]
	} else if len(values) > 1 {
		if _, ok := classes[values[0]]; !ok {
			err = errors.New("unknown class: " + values[0])
			return
		}
		if _, ok := classes[values[1]]; !ok {
			err = errors.New("unknown class: " + values[1])
			return
		}

		// get parents of first two classes determine their base class and
		// call function again with base class and rest
		par0 := classes[values[0]].GetAllParents()
		par1 := classes[values[1]].GetAllParents()

		for i := range par0 {
			for j := range par1 {
				if par0[i].Node.Term.String() == par1[j].Node.Term.String() {
					var temp []string
					temp = append(temp, par0[i].Name)
					for k := range values {
						if k == 0 || k == 1 {
							continue
						}
						temp = append(temp, values[k])
					}
					base, err = GetBaseClass(temp, classes)
					return
				}
			}
		}
	}
	return
}

// GetAllParents returns all transitive parents
func (class *Class) GetAllParents() (parents []*Class) {
	parents = append(parents, class)
	parents = append(parents, class.Parent...)
	for i := range class.Parent {
		temp := class.Parent[i].GetAllParents()
		for j := range temp {
			parExist := false
			for k := range parents {
				if temp[j].Node.Term.String() == parents[k].Node.Term.String() {
					parExist = true
				}
			}
			if !parExist {
				parents = append(parents, temp[j])
			}
		}
	}
	return
}

// GetRestrictions returns all restrictions of a class and its base classes
func (class *Class) GetRestrictions() (ret []*Restriction) {
	for i := range class.Restriction {
		restExist := false
		for j := range ret {
			if class.Restriction[i].Property.Name == ret[j].Property.Name {
				restExist = true
			}
		}
		if !restExist {
			ret = append(ret, class.Restriction[i])
		}
	}
	for i := range class.Parent {
		temp := class.Parent[i].GetRestrictions()
		for j := range temp {
			restExist := false
			for k := range ret {
				if temp[j].Property.Name == ret[k].Property.Name {
					restExist = true
					break
				}
			}
			if !restExist {
				ret = append(ret, temp[j])
			}
		}
	}
	return
}

// GetRestrictionsInverse returns all restrictions of a class and its base classes prioritizing base
// class restrictions
func (class *Class) GetRestrictionsInverse() (ret []*Restriction) {
	for i := range class.Parent {
		temp := class.Parent[i].GetRestrictionsInverse()
		for j := range temp {
			restExist := false
			for k := range ret {
				if temp[j].Property.Name == ret[k].Property.Name {
					restExist = true
					break
				}
			}
			if !restExist {
				ret = append(ret, temp[j])
			}
		}
	}
	for i := range class.Restriction {
		restExist := false
		for j := range ret {
			if class.Restriction[i].Property.Name == ret[j].Property.Name {
				restExist = true
			}
		}
		if !restExist {
			ret = append(ret, class.Restriction[i])
		}
	}
	return
}

// String prints a class
func (class *Class) String() (ret string) {
	ret = class.Node.Term.String() + ": " + class.Comment + "\n"
	ret += "\tParents: "
	for i := range class.Parent {
		ret += class.Parent[i].Node.Term.String() + ", "
	}
	ret += "\n\tChildren: "
	for i := range class.Child {
		ret += class.Child[i].Node.Term.String() + ", "
	}
	ret += "\n\tEnumerations: "
	for i := range class.Enumeration {
		ret += class.Enumeration[i].Name + ", "
	}
	ret += "\n\tRestrictions: "
	for i := range class.Restriction {
		ret += class.Restriction[i].Node.Term.String() + " on " +
			class.Restriction[i].Property.Name + " with ValueConstraint {" +
			class.Restriction[i].ValueConstraint + "} " + " with CardinalityConstraint {" +
			class.Restriction[i].CardinalityConstraint + " " +
			strconv.Itoa(class.Restriction[i].Multiplicity) + "} and value ["
		for j := range class.Restriction[i].Value {
			ret += class.Restriction[i].Value[j] + ", "
		}
		ret += "], "
	}
	ret += "\n\tUnionOf: "
	for i := range class.Union {
		ret += class.Union[i].Name + ", "
	}
	return
}
