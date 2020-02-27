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

	"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/rdf"
)

// extractIndividuals returns all nodes with type http://www.w3.org/2002/07/owl#NamedIndividual
func extractIndividuals(g *rdf.Graph,
	classes map[string]*Class) (individuals map[string]*Individual, err error) {
	fmt.Println("\tExtract individuals")
	individuals = make(map[string]*Individual)
	// detrmine all individuals
	for i := range g.Nodes {
		for j := range g.Nodes[i].Edge {
			if g.Nodes[i].Edge[j].Pred.String() ==
				"http://www.w3.org/1999/02/22-rdf-syntax-ns#type" {
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
func GetBaseType(individual []string, individuals map[string]*Individual,
	classes map[string]*Class) (base *Class, err error) {
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
