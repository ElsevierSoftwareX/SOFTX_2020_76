package rdf

import (
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

// New creates a graph from an rdf triple slice
func New(triple []Triple) (graph Graph, err error) {
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

// addDependentNodes
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
			temp := strings.Split(i, "b")
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
			gIn.Nodes[i].Term = BlankNode{name: "b" + strconv.Itoa(blankID)}
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
