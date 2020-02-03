package owl

import (
	"errors"
	"strconv"

	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/rdf"
)

// extractProperties returns returns all nodes with type owl:DatatypeProperty and owl:ObjectProperty
func extractProperties(g *rdf.Graph) (properties map[string]*Property, err error) {
	properties = make(map[string]*Property)
	// detrmine all properties
	for i := range g.Nodes {
		for j := range g.Nodes[i].Edge {
			if g.Nodes[i].Edge[j].Pred.String() == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" &&
				(g.Nodes[i].Edge[j].Object.Term.String() ==
					"http://www.w3.org/2002/07/owl#ObjectProperty" ||
					g.Nodes[i].Edge[j].Object.Term.String() ==
						"http://www.w3.org/2002/07/owl#DatatypeProperty") {
				prop := Property{
					Node:    g.Nodes[i],
					Name:    g.Nodes[i].Term.String(),
					Comment: getComment(g.Nodes[i]),
					Range:   getRange(g.Nodes[i]),
					Type:    g.Nodes[i].Edge[j].Object.Term.String(),
				}
				err = prop.extractPropertyCharacteristics()
				if err != nil {
					err = errors.New(err.Error() + " property " + prop.Name)
					return
				}
				properties[g.Nodes[i].Term.String()] = &prop
				break
			}
		}
	}
	return
}

// getRange returns a range if it exists
func getRange(node *rdf.Node) (ret []string) {
	for i := range node.Edge {
		if node.Edge[i].Pred.String() == "http://www.w3.org/2000/01/rdf-schema#range" {
			if node.Edge[i].Object.Term.Type() == rdf.TermBlankNode {
				for j := range node.Edge[i].Object.Edge {
					if node.Edge[i].Object.Edge[j].Pred.String() == "http://www.w3.org/2002/07/owl#oneOf" {
						enumNodes := getUnionValues(node.Edge[i].Object.Edge[j].Object)
						for k := range enumNodes {
							ret = append(ret, enumNodes[k].Term.String())
						}
						break
					}
				}
			} else {
				ret = append(ret, node.Edge[i].Object.Term.String())
			}
		}
	}
	return
}

// extractPropertyCharacteristics extracts // owl:functionalProperty, // owlInverseFunctionalProperty,
// owl:TransistiveProperty, owl:SysmmetricProperty
func (prop *Property) extractPropertyCharacteristics() (err error) {
	prop.IsFunctional = false
	prop.IsInverseFunctional = false
	prop.IsTransitive = false
	prop.IsSymmetric = false
	for i := range prop.Node.Edge {
		if prop.Node.Edge[i].Pred.String() == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" {
			if prop.Node.Edge[i].Object.Term.String() ==
				"http://www.w3.org/2002/07/owl#FunctionalProperty" {
				prop.IsFunctional = true
			} else if prop.Node.Edge[i].Object.Term.String() ==
				"http://www.w3.org/2002/07/owl#InverseFunctionalProperty" {
				prop.IsInverseFunctional = true
			} else if prop.Node.Edge[i].Object.Term.String() ==
				"http://www.w3.org/2002/07/owl#TransitiveProperty" {
				prop.IsTransitive = true
			} else if prop.Node.Edge[i].Object.Term.String() ==
				"http://www.w3.org/2002/07/owl#SymmetricProperty" {
				prop.IsSymmetric = true
			}
		}
	}
	return
}

// postProcessProperties extracts inverseOf, domain and subPropertyOf
func (on *Ontology) postProcessProperties() (err error) {
	for i := range on.Property {
		propNode := on.Property[i].Node
		for j := range propNode.Edge {
			pred := propNode.Edge[j]
			if pred.Pred.String() == "http://www.w3.org/2002/07/owl#inverseOf" {
				// extract owl:inverseOf
				if inv, ok := on.Property[pred.Object.Term.String()]; ok {
					on.Property[i].Inverse = inv
					inv.Inverse = on.Property[i]
				} else {
					err = errors.New("Property " + on.Property[i].Name + " unknown inverse: " +
						pred.Object.Term.String())
					return
				}
			} else if pred.Pred.String() == "http://www.w3.org/2000/01/rdf-schema#domain" {
				// extract rdf:domain
				if class, ok := on.Class[pred.Object.Term.String()]; ok {
					on.Property[i].Domain = append(on.Property[i].Domain, class)
				} else {
					err = errors.New("Property " + on.Property[i].Name + " unknown domain: " +
						pred.Object.Term.String())
					return
				}
			} else if pred.Pred.String() == "http://www.w3.org/2000/01/rdf-schema#subPropertyOf" {
				// extract rdf:subPropertyOf
				if sup, ok := on.Property[pred.Object.Term.String()]; ok {
					on.Property[i].SubPropertyOf = append(on.Property[i].SubPropertyOf, sup)
				} else {
					err = errors.New("Property " + on.Property[i].Name + " unknown property: " +
						pred.Object.Term.String())
					return
				}
			}
		}
	}
	return
}

// applyPropertyDomain adds restrictions to classes according to property domains
func (on *Ontology) addPropertyDomain() (err error) {
	for i := range on.Property {
		for j := range on.Property[i].Domain {
			rest := Restriction{
				Node:     on.Property[i].Node,
				Property: on.Property[i],
			}
			rest.Value = append(rest.Value, on.Property[i].Range...)
			//if gon.Property[i].IsFunctional {
			rest.CardinalityConstraint = "http://www.w3.org/2002/07/owl#cardinality"
			rest.Multiplicity = 1
			// } else {
			// 	rest.ValueConstraint = "http://www.w3.org/2002/07/owl#allValuesFrom"
			// }
			restExist := false
			for k := range on.Property[i].Domain[j].Restriction {
				if on.Property[i].Domain[j].Restriction[k].Property.Name == rest.Property.Name {
					// err = gon.Property[i].Domain[j].Restriction[k].mergeRestriction(&rest)
					// if err != nil {
					// 	return
					// }
					restExist = true
					break
				}
			}
			if !restExist {
				on.Property[i].Domain[j].Restriction = append(on.Property[i].Domain[j].Restriction, &rest)
				// fmt.Println("added property " + gon.Property[i].Name + " to class " + gon.Property[i].Domain[j].Name)
			}
		}
	}
	return
}

// String prints a property
func (prop *Property) String() (ret string) {
	ret = prop.Node.Term.String() + ": " + prop.Comment + "\n"
	ret += "\tType: " + prop.Type + "\n"
	ret += "\tRange: "
	for i := range prop.Range {
		ret += prop.Range[i] + ", "
	}
	ret += "\n\tIsFunctional: " + strconv.FormatBool(prop.IsFunctional) + "\n"
	ret += "\tIsInverseFunctional: " + strconv.FormatBool(prop.IsInverseFunctional) + "\n"
	ret += "\tIsTransitive: " + strconv.FormatBool(prop.IsTransitive) + "\n"
	ret += "\tIsSymmetric: " + strconv.FormatBool(prop.IsSymmetric) + "\n"
	ret += "\tInverse: "
	if prop.Inverse != nil {
		ret += prop.Inverse.Name
	}
	ret += "\n\tDomain: "
	for i := range prop.Domain {
		ret += prop.Domain[i].Name + ", "
	}
	ret += "\n\tSubPropertyOf: "
	for i := range prop.SubPropertyOf {
		ret += prop.SubPropertyOf[i].Name + ", "
	}
	return
}
