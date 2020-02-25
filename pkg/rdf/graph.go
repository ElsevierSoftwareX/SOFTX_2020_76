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

package rdf

import (
	"io"
	"strconv"
	"strings"
)

// Node is a node (subject and/or object) in a rdf graph
type Node struct {
	Term        Term
	Edge        []*Edge
	InverseEdge []*Edge
}

// Edge is a edge (predicate) in a rdf graph
type Edge struct {
	Pred    Predicate
	Subject *Node
	Object  *Node
}

// Graph is a rdf grapgh containing nodes and edges
type Graph struct {
	Nodes map[string]*Node
	Edges []*Edge
}

// NewGraph creates a graph from an rdf triple slice
func NewGraph(triple []Triple) (graph Graph, err error) {
	graph.Nodes = make(map[string]*Node)
	for i := range triple {
		// object
		obj, ok := graph.Nodes[triple[i].Obj.String()]
		if !ok {
			obj = &Node{
				Term: triple[i].Obj,
			}
			graph.Nodes[obj.Term.String()] = obj
		}

		// subject
		subj, ok := graph.Nodes[triple[i].Sub.String()]
		if !ok {
			subj = &Node{
				Term: triple[i].Sub,
			}
			graph.Nodes[subj.Term.String()] = subj
		}

		// predicate
		edge := &Edge{
			Pred:    triple[i].Pred,
			Subject: subj,
			Object:  obj,
		}
		subj.Edge = append(subj.Edge, edge)
		obj.InverseEdge = append(obj.InverseEdge, edge)
		graph.Edges = append(graph.Edges, edge)
	}
	err = nil
	return
}

// ToTriples extracts triples from a graph
func (graph *Graph) ToTriples() (ret []Triple) {
	for i := range graph.Edges {
		trip := Triple{
			Sub:  graph.Edges[i].Subject.Term,
			Pred: graph.Edges[i].Pred,
			Obj:  graph.Edges[i].Object.Term,
		}
		ret = append(ret, trip)
	}
	return
}

// SubGraph returns a graph containing the specified nodes and all transitive objects
func (graph *Graph) SubGraph(nodes ...*Node) (g Graph) {
	sub := make(map[string]*Node)
	for i := range nodes {
		nodes[i].addDependentNodes(sub)
	}
	g.Nodes = make(map[string]*Node)
	for i := range sub {
		newNode := &Node{Term: sub[i].Term}
		g.Nodes[newNode.Term.String()] = newNode
	}
	for i := range sub {
		subj, ok := g.Nodes[i]
		if ok {
			for j := range sub[i].Edge {
				obj, ok := g.Nodes[sub[i].Edge[j].Object.Term.String()]
				if ok {
					pred := &Edge{
						Pred:    sub[i].Edge[j].Pred,
						Subject: subj,
						Object:  obj,
					}
					subj.Edge = append(subj.Edge, pred)
					obj.InverseEdge = append(obj.InverseEdge, pred)
					g.Edges = append(g.Edges, pred)
				}
			}
		}
	}
	return
}

// addDependentNodes adds all nodes that are connected to node via an edge
func (node *Node) addDependentNodes(nodes map[string]*Node) {
	for i := range node.Edge {
		if _, ok := nodes[node.Edge[i].Object.Term.String()]; !ok {
			obj := node.Edge[i].Object
			nodes[obj.Term.String()] = obj
			obj.addDependentNodes(nodes)
		}
	}
	// for i := range node.InversePredicates {
	// 	if _, ok := nodes[node.InversePredicates[i].Subject.Name]; !ok {
	// 		subj := node.InversePredicates[i].Subject
	// 		nodes[subj.Name] = subj
	// 		subj.addDependentNodes(nodes)
	// 	}
	// }
	return
}

// Merge merges gIn into graph (nodes and edges are copied)
func (graph *Graph) Merge(gIn *Graph) (err error) {
	blankID := 0
	for i := range graph.Nodes {
		if graph.Nodes[i].Term.Type() == TermBlankNode {
			temp := strings.Split(i, "bn")
			if len(temp) > 1 {
				id, err := strconv.Atoi(temp[1])
				if err != nil {
					return err
				}
				if id > blankID {
					blankID = id
				}
			}
		}
	}
	blankID++
	for i := range gIn.Nodes {
		if gIn.Nodes[i].Term.Type() == TermBlankNode {
			gIn.Nodes[i].Term = BlankNode{name: "bn" + strconv.Itoa(blankID)}
			blankID++
		}
	}

	for i := range gIn.Nodes {
		if _, ok := graph.Nodes[gIn.Nodes[i].Term.String()]; !ok {
			n := &Node{Term: gIn.Nodes[i].Term}
			graph.Nodes[gIn.Nodes[i].Term.String()] = n
		}
	}
	for i := range gIn.Edges {
		if subj, ok := graph.Nodes[gIn.Edges[i].Subject.Term.String()]; ok {
			predExist := false
			for j := range subj.Edge {
				if subj.Edge[j].Pred.String() == gIn.Edges[i].Pred.String() &&
					subj.Edge[j].Object.Term.String() == gIn.Edges[i].Object.Term.String() {
					predExist = true
					break
				}
			}
			if !predExist {
				if obj, ok := graph.Nodes[gIn.Edges[i].Object.Term.String()]; ok {
					pred := &Edge{
						Pred:    gIn.Edges[i].Pred,
						Subject: subj,
						Object:  obj,
					}
					subj.Edge = append(subj.Edge, pred)
					obj.InverseEdge = append(obj.InverseEdge, pred)
					graph.Edges = append(graph.Edges, pred)
				}
			}
		}
	}
	return
}

// ToGraphvizDot exports a graph to the graphviz dot format
func (graph *Graph) ToGraphvizDot(output io.Writer, replace map[string]string,
	nodeShape map[string]string) (err error) {
	labelIndex := make(map[string]int)
	Index := 0
	dot := "digraph model\n{\n"
	for i := range graph.Nodes {
		label := i
		shape := ""
		// if graph.Nodes[i].Literal != nil {
		// 	label = graph.Nodes[i].Literal.String()
		// }
		for j := range replace {
			if strings.HasPrefix(label, j) {
				label = replace[j] + strings.TrimPrefix(label, j)
				break
			}
		}
		for j := range nodeShape {
			if strings.HasPrefix(i, j) {
				shape = "shape=" + nodeShape[j] + ", "
			}
		}

		labelIndex[i] = Index
		dot += "n" + strconv.Itoa(Index) + " [" + shape + "label=<" + label + ">]\n"
		Index++
	}
	dot += "\n"
	for i := range graph.Edges {
		label := graph.Edges[i].Pred.String()
		for j := range replace {
			if strings.HasPrefix(label, j) {
				label = replace[j] + strings.TrimPrefix(label, j)
				break
			}
		}
		subj := "n" + strconv.Itoa(labelIndex[graph.Edges[i].Subject.Term.String()])
		obj := "n" + strconv.Itoa(labelIndex[graph.Edges[i].Object.Term.String()])
		dot += subj + " -> " + obj + " [label=<" + label + ">]\n"
	}
	dot += "}\n"
	output.Write([]byte(dot))
	return
}

// Print prints all nodes of the graph
func (graph *Graph) String() (ret string) {
	ret = ""
	for i := range graph.Nodes {
		ret += graph.Nodes[i].Term.String() + "[ "
		for j := range graph.Nodes[i].Edge {
			ret += graph.Nodes[i].Edge[j].Pred.String() + " " +
				graph.Nodes[i].Edge[j].Object.Term.String() + ", "
		}
		ret += "]\n"
	}
	return
}
