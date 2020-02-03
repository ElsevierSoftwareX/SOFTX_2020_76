package owl

import (
	"errors"

	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/rdf"
)

// extractIndividuals returns all nodes with type http://www.w3.org/2002/07/owl#NamedIndividual
func extractIndividuals(g *rdf.Graph, classes map[string]*Class) (individuals map[string]*Individual,
	err error) {
	individuals = make(map[string]*Individual)
	// detrmine all individuals
	for i := range g.Nodes {
		for j := range g.Nodes[i].Edge {
			if g.Nodes[i].Edge[j].Pred.String() == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" {
				if class, ok := classes[g.Nodes[i].Edge[j].Object.Term.String()]; ok {
					ind := Individual{
						Node:    g.Nodes[i],
						Name:    g.Nodes[i].Term.String(),
						Comment: getComment(g.Nodes[i]),
						Type:    class,
					}
					individuals[g.Nodes[i].Term.String()] = &ind
					break
				}
			}
		}
	}
	return
}

// GetBaseType returns the base type of specified individuals
func GetBaseType(individual []string, individuals map[string]*Individual, classes map[string]*Class) (base *Class, err error) {
	indClass := make([]string, len(individual))
	for i := range individual {
		if ind, ok := individuals[individual[i]]; ok {
			indClass[i] = ind.Type.Name
		} else {
			err = errors.New("unkbown individual: " + individual[i])
		}
	}
	base, err = GetBaseClass(indClass, classes)
	return
}

// String prints a individual
func (ind *Individual) String() (ret string) {
	ret = ind.Node.Term.String() + ": " + ind.Comment + " Type: " + ind.Type.Name
	return
}
